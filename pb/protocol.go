package pb

import (
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
	ProtocolTypeMap[reflect.TypeOf(C2SStep{})] = ProtocolType_EC2SStep
	ProtocolTypeMap[reflect.TypeOf(S2CStep{})] = ProtocolType_ES2CStep

	ProtocolHandlerMap = make(map[int32]func([]byte, *net.Conn), 0)
}
