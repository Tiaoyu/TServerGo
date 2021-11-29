package main

import "net"

type UserSession struct {
	Conn        net.Conn
	OpenId      string
	RemoteAttr  string
	SendChannel chan []byte
}

func (s *UserSession) Send(msg []byte) {
	s.Conn.Write(msg)
}
func (s *UserSession) Close() {
	s.Conn.Close()
}
