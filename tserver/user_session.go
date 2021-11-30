package main

import "net"

type UserSession struct {
	Conn        net.Conn
	OpenId      string
	SendChannel chan []byte //消息发送chan
}

func (s *UserSession) Send(msg []byte) {
	s.Conn.Write(msg)
}
func (s *UserSession) Close() {
	s.Conn.Close()
}
