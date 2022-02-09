package main

import (
	"TServerGo/log"
	"encoding/binary"
	"errors"
	"flag"
	"io"
	"net"
	"net/http"
	_ "net/http/pprof"
)

const (
	maxPackageSize = 1024
)

var (
	Mysql     = flag.String("MYSQL", "", "please set mysql")
	MysqlHost = flag.String("MYSQL_HOST", "", "please set mysql host")
	connMap   = make(map[string]*ConnectInfo)
)

func main() {
	log.Init("TServer", log.LogLevelDEBUG|log.LogLevelINFO|log.LogLevelWARN|log.LogLevelERROR)
	flag.Parse()

	go func() {
		log.Debugf("http listen error: %v", http.ListenAndServe("localhost:6060", nil))
	}()

	// 数据库初始化
	db := dbProxy
	db.Init(*Mysql, *MysqlHost)
	db.Sync()

	// tcp服务初始化
	log.Info("init tcp...")
	addr, err := net.ResolveTCPAddr("tcp", ":8081")
	if err != nil {
		log.Errorf("resolve tcp addr error, err:%v", err)
		return
	}
	ln, err := net.ListenTCP("tcp", addr)

	if err != nil {
		log.Errorf("net error, err: %v", err)
		return
	}

	// 各个系统初始化
	initNotify()
	initMatch()
	initRoom()
	initUser()

	// socket accept
	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			log.Errorf("net AcceptTCP error, err: %v", err)
			continue
		}
		if err = conn.SetReadBuffer(SocketReadBufferSize); err != nil {
			log.Errorf("net SetReadBuffer error, err: %v", err)
			continue
		}
		if err = conn.SetWriteBuffer(SocketSendBufferSize); err != nil {
			log.Errorf("net SetWriteBuffer error, err: %v", err)
			continue
		}
		if err = conn.SetNoDelay(true); err != nil {
			log.Errorf("net SetNoDelay error, err: %v", err)
			continue
		}
		go handlerConnect(conn)
	}
}

func handlerConnect(conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Errorf("Conn close failed. Error:%v", err)
		}
	}()
	log.Debugf("Connected, Addr:%v", conn.RemoteAddr())
	connectInfo, ok := connMap[conn.RemoteAddr().String()]
	if !ok {
		connectInfo = &ConnectInfo{
			SOCKET: conn,
		}
	}
	handler := &HandlerProtobuf{}
	pLen := make([]byte, 4)
	for {
		_, err := io.ReadFull(conn, pLen)
		if err != nil {
			log.Errorf("net error on ReadFull pbLen, err:%v", err)
			handler.Error()
			break
		}
		l := binary.BigEndian.Uint32(pLen)
		if l < 4 || l > maxPackageSize {
			log.Errorf("net error on PBLen, err:%v", errors.New("protocol len is invalid"))
			break
		}
		msg := make([]byte, l)
		_, err = io.ReadFull(conn, msg)
		if err != nil {
			log.Errorf("net error on ReadFull msg, err:%v", err)
			break
		}
		handler.ParsePB(connectInfo, msg)
	}
}
