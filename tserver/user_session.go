package main

import (
	"TServerGo/log"
	"net"
)

type UserSession struct {
	Conn        net.Conn
	OpenId      string
	SendChannel chan []byte //消息发送chan
}

func (s *UserSession) Send(msg []byte) {
	if _, err := s.Conn.Write(msg); err != nil {
		log.Errorf("write to socket failed!")
	}
}
func (s *UserSession) Close() {
	err := s.Conn.Close()
	if err != nil {
		log.Errorf("Conn close failed! Error:%v", err)
	}
}
