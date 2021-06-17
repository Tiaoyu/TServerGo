package pbhandler

import (
	"github.com/gorilla/websocket"
)

type PBHandler interface {
	HandlerPB(ws *websocket.Conn, msg []byte) ([]byte, error)
}

func GetHandler(name string) PBHandler {
	switch name {
	case "json":
		return new(HandlerJson)
	default:
		return new(HandlerJson)
	}
}
