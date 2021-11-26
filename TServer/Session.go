package main

import "net"

type Session struct {
	Conn        net.Conn
	RemoteAttr  string
	SendChannel chan []byte
}

func (s *Session) Send(msg []byte) {
	s.Conn.Write(msg)
}
func (s *Session) Close() {
	s.Conn.Close()
}
