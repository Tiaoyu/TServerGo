package main

import (
	logger "TServerGo/Log"
	pbhandler "TServerGo/TServer/PBHandler"
	"TServerGo/TServer/constants"
	"TServerGo/TServer/dbproxy"
	"encoding/binary"
	"errors"
	"flag"
	"io"
	"net"
)

var (
	Mysql     = flag.String("MYSQL", "", "please set mysql")
	MysqlHost = flag.String("MYSQL_HOST", "", "please set mysql host")
	connMap   = make(map[string]*pbhandler.ConnectInfo)
)

func main() {
	logger.Init("TServer", logger.LogLevelDEBUG|logger.LogLevelINFO|logger.LogLevelWARN|logger.LogLevelERROR)
	flag.Parse()

	// 数据库初始化
	db := dbproxy.Instance()
	db.Init(*Mysql, *MysqlHost)
	db.Sync()

	// tcp服务初始化
	logger.Debug("init tcp...")
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8081")
	if err != nil {
		logger.Errorf("resolve tcp addr error, err:%v", err)
	}

	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		logger.Errorf("net error, err: %v", err)
		return
	}

	// socket accept
	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			logger.Errorf("net accept error, err: %v", err)
		}
		conn.SetReadBuffer(constants.SocketReadBufferSize)
		conn.SetWriteBuffer(1024)
		conn.SetNoDelay(true)
		go handlerConnect(conn)
	}
}

func handlerConnect(conn net.Conn) {
	defer conn.Close()
	logger.Debugf("Connected, Addr:%v", conn.RemoteAddr())
	connectInfo, ok := connMap[conn.RemoteAddr().String()]
	if !ok {
		connectInfo = &pbhandler.ConnectInfo{
			SOCKET: conn,
		}
	}
	handler := &pbhandler.HandlerProtobuf{}
	pLen := make([]byte, 4)
	for {
		_, err := io.ReadFull(conn, pLen)
		if err != nil {
			logger.Errorf("net error, err:%v", err)
			break
		}
		var len uint32
		len = binary.BigEndian.Uint32(pLen)
		if len < 4 {
			logger.Errorf("net error, err:%v", errors.New("protocol len is invalid"))
			break
		}
		msg := make([]byte, len)
		_, err = io.ReadFull(conn, msg)
		if err != nil {
			logger.Errorf("net error, err:%v", err)
			break
		}
		handler.ParsePB(connectInfo, msg)
		//
		//msg = make([]byte, 1024)
		//len, err = conn.Read(msg)
		//if err != nil || len == 0 {
		//	logger.Errorf("net error, err:%v", err)
		//	break
		//}
		//logger.Debugf("Recv msg, len:%v msg:%v", len, msg[:len])
		//handler.HandlerPB(connectInfo, msg[:len])
	}
}
