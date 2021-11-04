package pbhandler

import (
	logger "TServerGo/Log"
	"TServerGo/pb"
	"encoding/binary"
	"log"

	"google.golang.org/protobuf/proto"
)

var (
	pbMap map[gamepb.ProtocolType]func(msg []byte) ([]byte, uint32, error)
)

func init() {
	pbMap = make(map[gamepb.ProtocolType]func(msg []byte) ([]byte, uint32, error))
	pbMap[gamepb.ProtocolType_EC2SPing] = OnPing
}

type HandlerProtobuf struct {
}

func (h *HandlerProtobuf) HandlerPB(ws *ConnectInfo, msg []byte) ([]byte, error) {
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
			ws.Clear()
			return nil, nil
		}
		ack, ackPId, _ := call(ws.MsgContent[4:])

		ws.Clear()

		var bufHead = make([]byte, 4)
		var bufPId = make([]byte, 4)
		binary.BigEndian.PutUint32(bufPId, ackPId)
		binary.BigEndian.PutUint32(bufHead, uint32(len(ack)+4))
		bufHead = append(bufHead, bufPId...)
		bufHead = append(bufHead, ack...)
		ws.SOCKET.Write(bufHead) // todo 放到写协程中处理
	}
	return nil, nil
}
func OnPing(msg []byte) ([]byte, uint32, error) {
	logger.Debugf("Recv msg bytes:%v", msg)
	p := &gamepb.C2SPing{}
	if err := proto.Unmarshal(msg, p); err != nil {
		log.Fatalln("Failed to parse C2SPing:", err)
	}
	logger.Debugf("Recv msg:%v", p)

	ack, _ := proto.Marshal(&gamepb.S2CPing{Timestamp: p.Timestamp})

	return ack, uint32(gamepb.ProtocolType_ES2CPing), nil
}
