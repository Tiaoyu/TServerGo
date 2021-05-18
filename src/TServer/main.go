package main

import (
	"DBProxy"
	"bufio"
	"encoding/binary"
	"log"
	"net"
	"os"
	gamepb "pb"
)

var (
	conns []net.Conn
)

const (
	maxMessageLength = 8192
)

func main() {
	// 初始化数据库
	DBProxy.GetInstance().Init()

	// 监听socket
	l, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		log.Fatalln("listen error:", err)
		return
	}

	go accept(l)

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		for i, conn := range conns {
			_, err := conn.Write([]byte(line))
			if err != nil {
				conns = append(conns[:i], conns[i+1:]...)
			}
		}
	}
}
func accept(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalln("accept error:", err)
		}
		conns = append(conns, conn)
		go handler(conn)
	}
}

var (
	isFinishReadMap  = map[string]bool{}
	bufferMessageMap = map[string][]byte{}
	headLengthMap    = map[string]int32{}
)

// 处理 消息
func handler(conn net.Conn) {
	connkey := conn.RemoteAddr().String()
	for {
		tmp, ok := bufferMessageMap[connkey]
		if !ok {
			tmp = make([]byte, maxMessageLength)
			bufferMessageMap[connkey] = tmp
		}

		is, ok := isFinishReadMap[connkey]
		if !ok {
			isFinishReadMap[connkey] = true
		}

		if len, err := conn.Read(tmp); len > 0 && err == nil {
			if is {
				headLen := binary.LittleEndian.Uint32(tmp[0:4])
				if headLen > maxMessageLength {
					log.Println("Failed to parse protocol length:", headLen)
					conn.Write([]byte("error"))
					break
				}
			}

			protoId := binary.LittleEndian.Uint32(tmp[4:8])
			if pbHandler, ok := gamepb.ProtocolHandlerMap[int32(protoId)]; !ok {
				pbHandler(tmp[8:len], &conn)
			}
		}
	}
}
