package pbhandler

type PBHandler interface {
	HandlerPB(conn *ConnectInfo, msg []byte) ([]byte, error)
}

func GetHandler(name string) PBHandler {
	switch name {
	case "json":
		return new(HandlerJson)
	case "pb":
		return new(HandlerProtobuf)
	default:
		return new(HandlerJson)
	}
}
