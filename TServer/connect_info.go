package main

import (
	"net"
)

type ConnectInfo struct {
	SOCKET       net.Conn
	MsgSize      int32
	MsgHead      []byte
	MsgContent   []byte
	MsgLastBytes []byte
}

func (c *ConnectInfo) Clear() {
	c.MsgSize = 0
	c.MsgHead = c.MsgHead[0:0]
	c.MsgContent = c.MsgContent[0:0]
}
