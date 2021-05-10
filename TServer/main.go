package main

import (
	"bufio"
	"log"
	"net"
	"os"
)

var (
	conns []net.Conn
)

func handler(conn net.Conn) {
	for {
		var tmp = make([]byte, 8192)
		if len, err := conn.Read(tmp); len > 0 && err == nil {
			log.Println(string(tmp[:len]))
		}
	}
}

func main() {
	l, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatalln("listen error:", err)
		return
	}
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				log.Fatalln("accept error:", err)
			}
			conns = append(conns, conn)
			go handler(conn)
		}

	}()

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
