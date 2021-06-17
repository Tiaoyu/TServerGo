package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8081", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	sendQueue := make(chan []byte)
	go func() {
		sendQueue <- []byte("{\"id\": 1101,\"nickName\":\"条鱼鱼丶炕\",\"token\":\"0311lZ000MMFKL1qpI100iKyxo41lZ0s\",\"avatarUrl\":\"\"}")
		for {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter text: ")
			text, _, _ := reader.ReadLine()
			sendQueue <- text
		}
	}()

	for {
		select {
		case <-done:
			return
		//case t := <-ticker.C:
		//err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
		//if err != nil {
		//	log.Println("write:", err)
		//	return
		//}
		case t := <-sendQueue:
			err := c.WriteMessage(websocket.TextMessage, []byte(t))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
