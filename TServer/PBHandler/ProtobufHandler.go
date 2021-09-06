package pbhandler

import (
	logger "TServerGo/Log"
	"TServerGo/TServer/PB"
	"log"

	"google.golang.org/protobuf/proto"
)

type HandlerProtobuf struct {
}

func (h *HandlerProtobuf) HandlerPB(ws *ConnectInfo, msg []byte) ([]byte, error) {
	p := &PB.C2SPing{}
	if err := proto.Unmarshal(msg, p); err != nil {
		log.Fatalln("Failed to parse address book:", err)
	}
	logger.Debugf("Recv msg:%v", p)
	ack, _ := proto.Marshal(&PB.S2CPong{Time: p.Time})
	ws.SOCKET.Write(ack) // todo 放到写协程中处理

	return nil, nil
}
