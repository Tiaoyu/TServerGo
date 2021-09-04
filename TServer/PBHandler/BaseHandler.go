package pbhandler

type PBHandler interface {
	HandlerPB(conn *ConnectInfo, msg []byte) ([]byte, error)
}

func GetHandler(name string) PBHandler {
	switch name {
	case "json":
		return new(HandlerJson)
	default:
		return new(HandlerJson)
	}
}
