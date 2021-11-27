package main

import (
	logger "TServerGo/log"
	gamepb "TServerGo/pb"
	"encoding/binary"
	"errors"
	"log"

	"google.golang.org/protobuf/proto"
)

var (
	pbMap map[gamepb.ProtocolType]func(sess *Session, msg []byte) ([]byte, uint32, error)
)

func init() {
	pbMap = make(map[gamepb.ProtocolType]func(sess *Session, msg []byte) ([]byte, uint32, error))
	pbMap[gamepb.ProtocolType_EC2SPing] = OnPing
	pbMap[gamepb.ProtocolType_EC2SLogin] = OnLogin
	pbMap[gamepb.ProtocolType_EC2SMatch] = OnMatch
	pbMap[gamepb.ProtocolType_EC2SStep] = OnStep

}

type HandlerProtobuf struct {
	sess *Session
}

func (h *HandlerProtobuf) HandlerPB(ws *ConnectInfo, msg []byte) ([]byte, error) {
	defer ws.Clear()

	// 获取协议长度 4字节
	if len(ws.MsgHead) < 4 {
		if len(ws.MsgLastBytes) > 0 {
			//将上次剩余字节和新收到的字节合并
			msg = append(ws.MsgLastBytes, msg...)
			ws.MsgLastBytes = ws.MsgLastBytes[0:0]
		}
		remaindLen := 4 - len(ws.MsgHead)
		if len(msg) >= remaindLen {
			ws.MsgHead = append(ws.MsgHead, msg[0:remaindLen]...)
			msg = msg[remaindLen:]
			protoSize := binary.BigEndian.Uint32(ws.MsgHead)
			if protoSize <= 4 {
				//协议长度必须大于4 因为协议号会占4字节
				logger.Error("Get msg size error, msg size must greater than 4")
				ws.Clear()
				return nil, errors.New("msg size error")
			}
			ws.MsgSize = int32(protoSize)
			logger.Debugf("Recv msg size:%v", protoSize)
		} else {
			ws.MsgHead = append(ws.MsgHead, msg...)
			logger.Debugf("Recv msg head is not enough, len:%v", len(ws.MsgHead))
			return nil, nil
		}
	}
	// 获取协议内容 ws.MsgSize字节
	if ws.MsgSize > 0 && len(ws.MsgContent) < int(ws.MsgSize) {
		remaindLen := int(ws.MsgSize) - len(ws.MsgContent)
		if len(msg) >= remaindLen {
			ws.MsgContent = append(ws.MsgContent, msg[:remaindLen]...)
			//将剩余字节保留到下次
			ws.MsgLastBytes = msg[remaindLen:]
		} else {
			ws.MsgContent = append(ws.MsgContent, msg...)
			logger.Debugf("Recv msg content is not enough, len:%v", len(ws.MsgContent))
			return nil, nil
		}
	}
	// 反序列化协议
	h.ParsePB(ws, ws.MsgContent)
	return nil, nil
}

// ParsePB 反序列化协议
func (h *HandlerProtobuf) ParsePB(connectInfo *ConnectInfo, msg []byte) (error, error) {
	// 反序列化协议
	{
		protoId := binary.BigEndian.Uint32(msg[:4])
		call, ok := pbMap[gamepb.ProtocolType(protoId)]
		if !ok {
			return nil, nil
		}

		// 登陆
		var sess *Session
		if protoId == uint32(gamepb.ProtocolType_EC2SLogin) {
			if h.sess != nil {
				logger.Errorf("repeat login error")
				return nil, nil
			}
			sess = &Session{
				Conn:        connectInfo.SOCKET,
				RemoteAttr:  connectInfo.SOCKET.RemoteAddr().String(),
				SendChannel: make(chan []byte),
			}
			h.sess = sess
			go SendLoop(sess)
		}

		if h.sess == nil {
			return nil, nil
		}

		ack, ackPId, _ := call(h.sess, msg[4:])
		if ack == nil {
			return nil, nil
		}
		var bufHead = make([]byte, 4)
		var bufPId = make([]byte, 4)
		binary.BigEndian.PutUint32(bufPId, ackPId)
		binary.BigEndian.PutUint32(bufHead, uint32(len(ack)+4))
		bufHead = append(bufHead, bufPId...)
		bufHead = append(bufHead, ack...)
		h.sess.SendChannel <- bufHead
	}
	return nil, nil
}

func SendLoop(sess *Session) {
	for {
		msg, ok := <-sess.SendChannel
		if !ok {
			continue
		}
		sess.Send(msg)
	}
}

func OnLogin(sess *Session, msg []byte) ([]byte, uint32, error) {
	logger.Debugf("Recv msg bytes:%v", msg)
	req := &gamepb.C2SLogin{}
	if err := proto.Unmarshal(msg, req); err != nil {
		log.Fatalln("Failed to parse C2SPing:", err)
	}
	logger.Debugf("Recv msg:%v", req)

	player := &Player{
		OpenId:     req.NickName,
		NickName:   req.NickName,
		AvatarUrl:  req.AvatarUrl,
		RemoteAddr: sess.RemoteAttr,
		Sess:       sess,
	}
	PlayerLogin(player)

	ack, _ := proto.Marshal(&gamepb.S2CLogin{ErrorCode: "success"})
	return ack, uint32(gamepb.ProtocolType_ES2CLogin), nil
}

func OnPing(sess *Session, msg []byte) ([]byte, uint32, error) {
	logger.Debugf("Recv msg bytes:%v", msg)
	req := &gamepb.C2SPing{}
	if err := proto.Unmarshal(msg, req); err != nil {
		log.Fatalln("Failed to parse C2SPing:", err)
	}
	logger.Debugf("Recv msg:%v", req)

	ack, _ := proto.Marshal(&gamepb.S2CPing{Timestamp: req.Timestamp})

	return ack, uint32(gamepb.ProtocolType_ES2CPing), nil
}

func OnMatch(sess *Session, msg []byte) ([]byte, uint32, error) {
	req := &gamepb.C2SMatch{}
	if err := proto.Unmarshal(msg, req); err != nil {
		log.Fatalln("Failed to parse C2SMatch:", err)
	}
	logger.Debugf("Recv msg:%v", req)
	player, ok := PlayerRemoteMap[sess.RemoteAttr]
	if !ok {
		return nil, 0, nil
	}
	if req.MatchType == gamepb.MatchType_MatchTypeMatch {
		JoinMatch(player)
	} else if req.MatchType == gamepb.MatchType_MatchTypeCancel {
		CancelMatch(player)
	}

	res, err := proto.Marshal(&gamepb.S2CMatch{
		Result: gamepb.MatchResult_MatResultMatching,
	})

	return res, uint32(gamepb.ProtocolType_ES2CMatch), err
}

func OnStep(sess *Session, msg []byte) ([]byte, uint32, error) {
	req := &gamepb.C2SStep{}
	if err := proto.Unmarshal(msg, req); err != nil {
		log.Fatalln("Failed to parse C2SStep:", err)
	}
	logger.Debugf("Recv msg:%v", req)

	player, ok := PlayerRemoteMap[sess.RemoteAttr]
	if !ok {
		log.Println("player is nil, RemoteAddr:", sess.RemoteAttr)
		return nil, 0, nil
	}
	value, ok := RoomOpenIdMap.Load(player.OpenId)
	if !ok {
		log.Println("room is nil, OpenId:", player.OpenId)
		return nil, 0, nil
	}
	room := value.(*Room)
	if room.TurnId != player.OpenId {
		res, _ := proto.Marshal(&gamepb.S2CStep{
			Error:      nil,
			GobangInfo: nil,
		})
		return res, uint32(gamepb.ProtocolType_EC2SStep), nil
	}
	if room.BlackId == player.OpenId {
		req.Point.Camp = int32(gamepb.ColorType_ColorTypeBlack)
	} else if room.RedId == player.OpenId {
		req.Point.Camp = int32(gamepb.ColorType_ColorTypeRed)
	}
	room.MsgChannel <- &gamepb.ChessStep{
		Point: req.Point,
	}

	// ack, _ := proto.Marshal(&gamepb.S2CStep{
	// 	Error:      nil,
	// 	GobangInfo: nil,
	// })
	return nil, uint32(gamepb.ProtocolType_ES2CStep), nil
}
