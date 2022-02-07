package main

import (
	"TServerGo/log"
	"TServerGo/pb"
	"encoding/binary"
	"google.golang.org/protobuf/proto"
)

var (
	// 协议类型 对应的 协议处理方法
	pbMap map[pb.ProtocolType]func(sess *UserSession, msg []byte) ([]byte, uint32, error)
)

func init() {
	pbMap = make(map[pb.ProtocolType]func(sess *UserSession, msg []byte) ([]byte, uint32, error))
	pbMap[pb.ProtocolType_EC2SPing] = OnPing
	pbMap[pb.ProtocolType_EC2SLogin] = OnLogin
	pbMap[pb.ProtocolType_EC2SMatch] = OnMatch
	pbMap[pb.ProtocolType_EC2SStep] = OnStep

}

type HandlerProtobuf struct {
	sess *UserSession
}

func (h *HandlerProtobuf) Error() {
	if h.sess != nil {
		h.sess.Close()
		NotifyExec(NotifyTypeRoleLogout, &NotifyRoleLogoutParam{OpenId: h.sess.OpenId})
	}
}

// ParsePB 反序列化协议
func (h *HandlerProtobuf) ParsePB(connectInfo *ConnectInfo, msg []byte) (error, error) {
	// 反序列化协议
	{
		protoId := binary.BigEndian.Uint32(msg[:4])
		call, ok := pbMap[pb.ProtocolType(protoId)]
		if !ok {
			return nil, nil
		}

		// 登陆
		var sess *UserSession
		if protoId == uint32(pb.ProtocolType_EC2SLogin) {
			if h.sess != nil {
				log.Errorf("repeat login error")
				return nil, nil
			}
			sess = &UserSession{
				Conn:        connectInfo.SOCKET,
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

func SendLoop(sess *UserSession) {
	for {
		msg, ok := <-sess.SendChannel
		if !ok {
			log.Debugf("Send channel is closed! OpenId:%v", sess.OpenId)
			break
		}
		pid := binary.BigEndian.Uint32(msg[4:8])
		log.Debugf("Send OpenId:%v pid:%v", sess.OpenId, pid)
		sess.Send(msg)
	}
}

func OnLogin(sess *UserSession, msg []byte) ([]byte, uint32, error) {
	log.Debugf("Login Recv msg bytes:%v OpenId:%v", msg, sess.OpenId)
	req := &pb.C2SLogin{}
	if err := proto.Unmarshal(msg, req); err != nil {
		log.Errorf("Failed to parse C2SPing:%v", err)
		return nil, 0, err
	}
	log.Debugf("Recv msg:%v", req)

	player := &Player{
		OpenId: req.NickName,
		Sess:   sess,
	}
	err := PlayerLogin(player)
	if err != nil {
		ackTmp, _ := proto.Marshal(&pb.S2CLogin{
			Error: &pb.Error{
				ErrorCode: pb.ErrorType_FAILED,
				ErrorMsg:  "failed",
			},
		})

		return ackTmp, uint32(pb.ProtocolType_ES2CLogin), err
	}
	sess.OpenId = player.OpenId
	ack, _ := proto.Marshal(&pb.S2CLogin{Error: &pb.Error{
		ErrorCode: pb.ErrorType_SUCCESS,
		ErrorMsg:  "success",
	}})
	return ack, uint32(pb.ProtocolType_ES2CLogin), err
}

func OnPing(sess *UserSession, msg []byte) ([]byte, uint32, error) {
	log.Debugf("Recv msg bytes:%v", msg)
	req := &pb.C2SPing{}
	if err := proto.Unmarshal(msg, req); err != nil {
		log.Errorf("Failed to parse C2SPing:%v", err)
		return nil, 0, err
	}
	log.Debugf("Recv msg:%v", req)

	ack, _ := proto.Marshal(&pb.S2CPing{Timestamp: req.Timestamp})

	return ack, uint32(pb.ProtocolType_ES2CPing), nil
}

func OnMatch(sess *UserSession, msg []byte) ([]byte, uint32, error) {
	req := &pb.C2SMatch{}
	if err := proto.Unmarshal(msg, req); err != nil {
		log.Errorf("Failed to parse C2SMatch:%v", err)
	}
	log.Debugf("OnMatch Recv msg:%v", req)
	player, ok := PlayerOpenIdMap[sess.OpenId]
	if !ok {
		return nil, 0, nil
	}
	if req.MatchType == pb.MatchType_MatchTypeMatch {
		JoinMatch(player)
	} else if req.MatchType == pb.MatchType_MatchTypeCancel {
		CancelMatch(player)
	}

	res, err := proto.Marshal(&pb.S2CMatch{
		Result: pb.MatchResult_MatResultMatching,
	})

	return res, uint32(pb.ProtocolType_ES2CMatch), err
}

func OnStep(sess *UserSession, msg []byte) ([]byte, uint32, error) {
	req := &pb.C2SStep{}
	if err := proto.Unmarshal(msg, req); err != nil {
		log.Errorf("Failed to parse C2SStep:%v", err)
	}
	log.Debugf("Recv msg:%v", req)

	player, ok := PlayerOpenIdMap[sess.OpenId]
	if !ok {
		log.Errorf("player is nil, OpenId:%v", player.OpenId)
		return nil, 0, nil
	}
	value, ok := RoomOpenIdMap.Load(player.OpenId)
	if !ok {
		log.Errorf("room is nil, OpenId:%v", player.OpenId)
		return nil, 0, nil
	}
	room := value.(*Room)
	if room.TurnId != player.OpenId {
		res, _ := proto.Marshal(&pb.S2CStep{
			Error:      nil,
			GobangInfo: nil,
		})
		return res, uint32(pb.ProtocolType_EC2SStep), nil
	}
	if room.BlackId == player.OpenId {
		req.Point.Camp = int32(pb.ColorType_ColorTypeBlack)
	} else if room.RedId == player.OpenId {
		req.Point.Camp = int32(pb.ColorType_ColorTypeRed)
	}
	room.MsgChannel <- &pb.ChessStep{
		Point: req.Point,
	}

	return nil, uint32(pb.ProtocolType_ES2CStep), nil
}
