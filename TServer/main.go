package main

import (
	logger "TServerGo/log"
	"encoding/binary"
	"errors"
	"flag"
	"io"
	"net"
)

var (
	Mysql     = flag.String("MYSQL", "", "please set mysql")
	MysqlHost = flag.String("MYSQL_HOST", "", "please set mysql host")
	connMap   = make(map[string]*ConnectInfo)
)

func main() {
	logger.Init("TServer", logger.LogLevelDEBUG|logger.LogLevelINFO|logger.LogLevelWARN|logger.LogLevelERROR)
	flag.Parse()

	// 数据库初始化
	db := dbProxy
	db.Init(*Mysql, *MysqlHost)
	db.Sync()

	// tcp服务初始化
	logger.Info("init tcp...")
	addr, err := net.ResolveTCPAddr("tcp", ":8081")
	if err != nil {
		logger.Errorf("resolve tcp addr error, err:%v", err)
		return
	}

	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		logger.Errorf("net error, err: %v", err)
		return
	}
	// init notify
	initNotify()
	initMatch()
	initRoom()
	initUser()

	// socket accept
	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			logger.Errorf("net accept error, err: %v", err)
		}
		conn.SetReadBuffer(SocketReadBufferSize)
		conn.SetWriteBuffer(SocketSendBufferSize)
		conn.SetNoDelay(true)
		go handlerConnect(conn)
	}
}

func handlerConnect(conn net.Conn) {
	defer conn.Close()
	logger.Debugf("Connected, Addr:%v", conn.RemoteAddr())
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
			logger.Errorf("net error on ReadFull pbLen, err:%v", err)
			handler.Error()
			break
		}
		len := binary.BigEndian.Uint32(pLen)
		if len < 4 {
			logger.Errorf("net error on PBLen, err:%v", errors.New("protocol len is invalid"))
			break
		}
		msg := make([]byte, len)
		_, err = io.ReadFull(conn, msg)
		if err != nil {
			logger.Errorf("net error on ReadFull msg, err:%v", err)
			break
		}
		handler.ParsePB(connectInfo, msg)
	}
}
