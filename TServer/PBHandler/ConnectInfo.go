package pbhandler

import (
	"net"

	"github.com/gorilla/websocket"
)

type ConnectInfo struct {
	WS     *websocket.Conn
	SOCKET net.Conn
}
