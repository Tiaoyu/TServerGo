package pbhandler

import (
	"github.com/gorilla/websocket"
)

type ConnectInfo struct {
	WS *websocket.Conn
}
