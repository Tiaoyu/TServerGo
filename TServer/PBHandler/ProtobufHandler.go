package pbhandler

import (
	logger "TServerGo/Log"
	"TServerGo/TServer/MatchSystem"
	"TServerGo/TServer/PB"
	"TServerGo/TServer/Sessionx"
	"TServerGo/TServer/UserSystem"
	"TServerGo/pb"
	"encoding/binary"
	"google.golang.org/protobuf/proto"
	"log"
)

var (
	pbMap map[gamepb.ProtocolType]func(sess *Sessionx.Session, msg []byte) ([]byte, uint32, error)
)

func init() {
	pbMap = make(map[gamepb.ProtocolType]func(sess *Sessionx.Session, msg []byte) ([]byte, uint32, error))
	pbMap[gamepb.ProtocolType_EC2SPing] = OnPing
	pbMap[gamepb.ProtocolType_EC2SLogin] = OnLogin
	pbMap[gamepb.ProtocolType_EC2SMatch] = OnMatch
	pbMap[gamepb.ProtocolType_EC2SStep] = OnStep

}

type HandlerProtobuf struct {
	sess *Sessionx.Session
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
				return nil, nil
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
	{
		protoId := binary.BigEndian.Uint32(ws.MsgContent[:4])
		call, ok := pbMap[gamepb.ProtocolType(protoId)]
		if !ok {
			return nil, nil
		}

		// 登陆
		var sess *Sessionx.Session
		if protoId == uint32(gamepb.ProtocolType_EC2SLogin) {
			sess = &Sessionx.Session{
				Conn:       ws.SOCKET,
				RemoteAttr: ws.SOCKET.RemoteAddr().String(),
				SendBuffer: make(chan []byte),
			}
			h.sess = sess
			go SendLoop(sess)
		}

		ack, ackPId, _ := call(sess, ws.MsgContent[4:])

		var bufHead = make([]byte, 4)
		var bufPId = make([]byte, 4)
		binary.BigEndian.PutUint32(bufPId, ackPId)
		binary.BigEndian.PutUint32(bufHead, uint32(len(ack)+4))
		bufHead = append(bufHead, bufPId...)
		bufHead = append(bufHead, ack...)
		if protoId == uint32(gamepb.ProtocolType_EC2SPing) {
			ws.SOCKET.Write(bufHead)
		} else {
			sess.SendBuffer <- bufHead
		}
	}
	return nil, nil
}

func SendLoop(sess *Sessionx.Session) {
	for {
		select {
		case msg, ok := <-sess.SendBuffer:
			if !ok {
				continue
			}
			sess.Send(msg)
			break
		}
	}
}

func OnLogin(sess *Sessionx.Session, msg []byte) ([]byte, uint32, error) {
	logger.Debugf("Recv msg bytes:%v", msg)
	req := &gamepb.C2SLogin{}
	if err := proto.Unmarshal(msg, req); err != nil {
		log.Fatalln("Failed to parse C2SPing:", err)
	}
	logger.Debugf("Recv msg:%v", req)

	player := &UserSystem.Player{
		OpenId:      "",
		NickName:    req.NickName,
		AvatarUrl:   req.AvatarUrl,
		RemoteAddr:  sess.RemoteAttr,
		SendChannel: make(chan []byte),
		Sess:        sess,
	}
	UserSystem.PlayerLogin(player)

	ack, _ := proto.Marshal(&gamepb.S2CLogin{ErrorCode: "success"})
	return ack, uint32(gamepb.ProtocolType_ES2CPing), nil
}

func OnPing(sess *Sessionx.Session, msg []byte) ([]byte, uint32, error) {
	logger.Debugf("Recv msg bytes:%v", msg)
	req := &gamepb.C2SPing{}
	if err := proto.Unmarshal(msg, req); err != nil {
		log.Fatalln("Failed to parse C2SPing:", err)
	}
	logger.Debugf("Recv msg:%v", req)

	ack, _ := proto.Marshal(&gamepb.S2CPing{Timestamp: req.Timestamp})

	return ack, uint32(gamepb.ProtocolType_ES2CPing), nil
}

func OnMatch(sess *Sessionx.Session, msg []byte) ([]byte, uint32, error) {
	req := &gamepb.C2SMatch{}
	if err := proto.Unmarshal(msg, req); err != nil {
		log.Fatalln("Failed to parse C2SMatch:", err)
	}
	logger.Debugf("Recv msg:%v", req)
	player, ok := UserSystem.PlayerRemoteMap[sess.RemoteAttr]
	if !ok {
		return nil, 0, nil
	}
	if req.MatchType == PB.MatchTypeMatch {
		MatchSystem.JoinMatch(player)
	} else if req.MatchType == PB.MatchTypeCancel {
		MatchSystem.CancelMatch(player)
	}

	res, err := proto.Marshal(&gamepb.S2CMatch{
		Result: gamepb.MatchResult_MatResultMatching,
	})

	return res, uint32(gamepb.ProtocolType_ES2CMatch), err
}

func OnStep(sess *Sessionx.Session, msg []byte) ([]byte, uint32, error) {
	req := &gamepb.C2SStep{}
	if err := proto.Unmarshal(msg, req); err != nil {
		log.Fatalln("Failed to parse C2SStep:", err)
	}
	logger.Debugf("Recv msg:%v", req)

	ack, _ := proto.Marshal(&gamepb.S2CStep{
		Error:      nil,
		GobangInfo: nil,
	})
	return ack, uint32(gamepb.ProtocolType_ES2CStep), nil
}
