package main

import (
	"DBProxy"
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"log"
	"net"
	"net/http"
	"os"
	gamepb "pb"
)

var (
	conns    []net.Conn
	addr     = flag.String("addr", "localhost:9999", "http service address")
	upgrader = websocket.Upgrader{}
)

const (
	maxMessageLength = 8192
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/ws", hello)
	e.Logger.Fatal(e.Start(":8081"))
	log.Fatal(http.ListenAndServe(*addr, nil))

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
func hello(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), c.Request().Header)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		// Write
		err := ws.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))
		if err != nil {
			c.Logger().Error(err)
		}

		// Read
		_, msg, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
		}
		fmt.Printf("%s\n", msg)
	}
}
