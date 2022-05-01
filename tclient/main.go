package main

import "C"
import (
	"TServerGo/pb"
	"bufio"
	"encoding/binary"
	"fmt"
	"google.golang.org/protobuf/proto"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	// 收包队列
	var s2cChan = make(chan proto.Message, 10)
	var receiveQ Queue
	var sendQ Queue
	receiveQ.Init()
	sendQ.Init()

	// TCP连接
	client, err := net.Dial("tcp", "192.168.0.107:8081")
	if err != nil {
		fmt.Println(err)
		return
	}

	// 处理stdin协程
	go func() {
		var desc = `
			login [nick] [avatar] [userid]
			match [type:1-match 2-cancel]
			step []
		`
		var descOutput = func(l int, p []string, f func()) {
			if l != len(p) {
				fmt.Println(desc)
			} else {
				f()
			}
		}
		input := bufio.NewReader(os.Stdin)
		for {
			text, _ := input.ReadString('\n')
			textParams := strings.Split(strings.TrimSpace(text), " ")
			textLen := len(textParams)
			if textLen == 0 {
				fmt.Println(desc)
			}
			switch textParams[0] {
			case "login":
				descOutput(4, textParams, func() {
					sendQ.PushLast(msgToBytes(&pb.C2SLogin{
						NickName:  textParams[1],
						AvatarUrl: textParams[2],
						UserId:    textParams[3],
					}, uint32(pb.ProtocolType_EC2SLogin)))
				})
			case "match":
				descOutput(2, textParams, func() {
					t, _ := strconv.Atoi(textParams[1])
					sendQ.PushLast(msgToBytes(&pb.C2SMatch{
						MatchType: pb.MatchType(t),
					}, uint32(pb.ProtocolType_EC2SMatch)))
				})
			case "step":
				descOutput(3, textParams, func() {
					x, _ := strconv.Atoi(textParams[1])
					y, _ := strconv.Atoi(textParams[2])
					sendQ.PushLast(msgToBytes(&pb.C2SStep{
						Point: &pb.Point{
							X: int32(x),
							Y: int32(y),
						},
					}, uint32(pb.ProtocolType_EC2SStep)))
				})
			}
		}
	}()

	// ping 协程
	go func() {
		defer fmt.Println("ping goroutine down!!!")
		C := time.Tick(time.Second * 1)
	L:
		for {
			select {
			case _, ok := <-C:
				{
					if !ok {
						break L
					}
					sendQ.PushLast(msgToBytes(&pb.C2SPing{Timestamp: time.Now().UnixNano() / 1e6}, uint32(pb.ProtocolType_EC2SPing)))
				}
			}
		}
	}()

	// 发包协程
	go func() {
		for {
			if sendQ.Len() == 0 {
				continue
			}
			msg := sendQ.PopFront()
			client.Write(msg)
		}
	}()

	// 解包协程
	go func() {
		for {
			if receiveQ.Len() == 0 {
				continue
			}
			msg := receiveQ.PopFront()
			pbId := binary.BigEndian.Uint32(msg[:4])
			switch pbId {
			case uint32(pb.ProtocolType_ES2CLogin):
				res := &pb.S2CLogin{}
				proto.Unmarshal(msg[4:], res)
				fmt.Println("login: " + res.Error.ErrorMsg)
				s2cChan <- res
			case uint32(pb.ProtocolType_ES2CMatch):
				res := &pb.S2CMatch{}
				proto.Unmarshal(msg[4:], res)
				s2cChan <- res
			case uint32(pb.ProtocolType_ES2CPing):
				res := &pb.S2CPing{}
				proto.Unmarshal(msg[4:], res)
				s2cChan <- res
			case uint32(pb.ProtocolType_ES2CGameResult):
				res := &pb.S2CGameResult{}
				proto.Unmarshal(msg[4:], res)
				s2cChan <- res
			case uint32(pb.ProtocolType_ES2CStep):
				res := &pb.S2CStep{}
				proto.Unmarshal(msg[4:], res)
				s2cChan <- res
			}
		}
	}()

	// 逻辑协程
	go func() {
		for {
			b := <-s2cChan
			if res, ok := b.(*pb.S2CLogin); ok {
				fmt.Printf("login: %v\n", res)
			} else if res, ok := b.(*pb.S2CPing); ok {
				fmt.Printf("ping: %v %v\n", res, time.Now().UnixNano()/1e6-res.Timestamp)
			} else if res, ok := b.(*pb.S2CMatch); ok {
				fmt.Printf("match: %v\n", res)
			} else if res, ok := b.(*pb.S2CStep); ok {
				fmt.Printf("step: %v\n", res)
			}
		}
	}()

	// 主协程 接收网络包
	pLen := make([]byte, 4)
	for {
		_, err := io.ReadFull(client, pLen)
		if err != nil {
			fmt.Printf("connection closed!!! err:%v\n", err)
			break
		}
		l := binary.BigEndian.Uint32(pLen)
		if l < 4 || l > 2048 {
			fmt.Println("pLen too long!!!")
			break
		}
		msg := make([]byte, l)
		_, _ = io.ReadFull(client, msg)
		go receiveQ.PushLast(msg) // 异步不阻塞
	}
}

func msgToBytes(message proto.Message, pbId uint32) []byte {
	msg, _ := proto.Marshal(message)
	var bufHead = make([]byte, 4)
	var bufPId = make([]byte, 4)
	binary.BigEndian.PutUint32(bufPId, pbId)
	binary.BigEndian.PutUint32(bufHead, uint32(len(msg)+4))
	bufHead = append(bufHead, bufPId...)
	bufHead = append(bufHead, msg...)
	return bufHead
}

type Queue struct {
	sync.RWMutex
	data [][]byte
}

func (q *Queue) Init() {
	q.Lock()
	defer q.Unlock()
	q.data = make([][]byte, 0, 10)
}
func (q *Queue) PushFront(e []byte) {
	q.Lock()
	defer q.Unlock()
	q.data = append([][]byte{e}, q.data...)
}

func (q *Queue) PushLast(e []byte) {
	q.Lock()
	defer q.Unlock()
	q.data = append(q.data, e)
}

func (q *Queue) PopFront() []byte {
	q.Lock()
	defer q.Unlock()
	e := q.data[0]
	q.data = q.data[1:]
	return e
}

func (q *Queue) PopLast() []byte {
	q.Lock()
	defer q.Unlock()
	e := q.data[len(q.data)-1]
	q.data = q.data[0 : len(q.data)-1]
	return e
}

func (q *Queue) Len() int {
	q.RLock()
	defer q.RUnlock()
	return len(q.data)
}
