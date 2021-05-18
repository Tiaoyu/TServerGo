package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"github.com/golang/protobuf/proto"
	"log"
	"net"
	"os"
	gamepb "pb"
	"reflect"
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

	go handler()

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		//line := input.Text()
		protocol := &gamepb.C2SGobangStep{Point: &gamepb.Point{
			X: 1, Y: 2,
		}}
		if out, err := proto.Marshal(protocol); err != nil {

		} else {
			conn.Write(out)
		}
	}
}

func handler() {
	for {
		var tmp = make([]byte, 8192)
		if len, err := conn.Read(tmp); len > 0 && err == nil {
			ping := &gamepb.S2CPing{}
			if err := proto.Unmarshal(tmp, ping); err != nil {

			}
			log.Println(ping, ". Ping:", time.Now().Unix()-ping.Timestamp)
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
		ping := &gamepb.C2SPing{Timestamp: time.Now().Unix()}
		out, err := proto.Marshal(ping)
		if err != nil {
			log.Fatalln("Failed to encode ping:", err)
		}

		prefix := make([]byte, 8)
		binary.LittleEndian.PutUint32(prefix, uint32(len(out)+4+4))
		binary.LittleEndian.PutUint32(prefix[4:], uint32(gamepb.ProtocolTypeMap[reflect.TypeOf(ping)]))
		prefix = append(prefix, out...)
		len, err := conn.Write(prefix)
		if err != nil {
			log.Println("cannot connect to server, error:", err, len)
			go reconnect()
			break
		} else {
			log.Println("ping to server...")
		}
		time.Sleep(time.Millisecond * 5000)
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
