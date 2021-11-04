package gamepb

import (
	"github.com/golang/protobuf/proto"
	"log"
	"net"
	"reflect"
)

var (
	ProtocolTypeMap    map[reflect.Type]ProtocolType
	ProtocolHandlerMap map[int32]func([]byte, *net.Conn)
)

func init() {
	ProtocolTypeMap = make(map[reflect.Type]ProtocolType, 0)
	ProtocolTypeMap[reflect.TypeOf(C2SPing{})] = ProtocolType_EC2SPing
	ProtocolTypeMap[reflect.TypeOf(S2CPing{})] = ProtocolType_ES2CPing
	ProtocolTypeMap[reflect.TypeOf(C2SGobangStep{})] = ProtocolType_EC2SGobangStep
	ProtocolTypeMap[reflect.TypeOf(S2CGobangStep{})] = ProtocolType_ES2CGobangStep

	ProtocolHandlerMap = make(map[int32]func([]byte, *net.Conn), 0)
	ProtocolHandlerMap[int32(ProtocolType_EC2SPing)] = DoC2SPing
	ProtocolHandlerMap[int32(ProtocolType_EC2SGobangStep)] = DoC2SGobangStep
}

// DoC2SPing ping
func DoC2SPing(bytes []byte, conn *net.Conn) {
	protocol := &C2SPing{}
	if err := proto.Unmarshal(bytes, protocol); err != nil {
		log.Fatalln("Failed to parse C2SPing:", err)
	} else {
		// log.Println(ping)
		out, err := proto.Marshal(&S2CPing{Timestamp: protocol.Timestamp,
			Error: &Error{ErrorMsg: "success", ErrorCode: ErrorType_SUCCESS}})
		if err != nil {
			log.Println("Failed to encode S2CPing:", err)
		}
		(*conn).Write(out)
	}
}

// DoC2SGobangStep 一步
func DoC2SGobangStep(bytes []byte, conn *net.Conn) {
	protocol := &C2SGobangStep{}
	if err := proto.Unmarshal(bytes, protocol); err != nil {
		log.Fatalln("Failed to parse C2SGobangStep:", err)
	} else {
		// TODO 使用channel 找到对应的房间 处理相应的逻辑
	}
}
