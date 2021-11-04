package main

import (
	"net"
	"os"
)

func main() {
	l, err := net.Listen("tcp", "localhost:8888")
	if err != nil {
		println("Error listening:", err)
		os.Exit(1)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			println("Error accepting:", err)
			os.Exit(1)
		}
		println("Someone connected! " + conn.RemoteAddr().String())
	}
}
