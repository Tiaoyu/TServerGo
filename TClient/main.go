package main

import (
	"flag"
	"log"
	"net"
	"time"
)

const (
	RECONNECT_COUNT = 10
)

var (
	port           = flag.String("port", "", "string类型参数")
	conn           net.Conn
	reconnectCount int
)

func main() {
	flag.Parse()

	if c, err := connect(); c != nil && err == nil {
		conn = c
		go ping()
	} else {
		log.Fatal("cannot connect to server!!!")
	}

	for {
		var tmp = make([]byte, 8192)
		if len, err := conn.Read(tmp); len > 0 && err == nil {
			log.Println(string(tmp[:len]))
		}
	}
}
func connect() (net.Conn, error) {
	c, err := net.Dial("tcp", "localhost:"+*port)
	if err != nil {
		log.Println("connect to server error:", err)
		return nil, err
	} else {
		log.Println("connect to server success!")
	}
	return c, nil
}
func ping() {
	for {
		len, err := conn.Write([]byte("ping"))
		if err != nil {
			log.Println("cannot connect to server, error:", err, len)
			go reconnect()
			break
		} else {
			log.Println("ping to server...")
		}
		time.Sleep(time.Millisecond * 1000)
	}
}

func reconnect() {
	for {
		if reconnectCount >= RECONNECT_COUNT {
			reconnectCount = 0
			log.Println("finally cannot connect to server after try ", RECONNECT_COUNT, " times.")
			break
		}
		if c, err := connect(); c != nil && err == nil {
			conn = c
			go ping()
			break
		}
		reconnectCount++
		log.Println("try to connect to server...", reconnectCount)
	}
}
