package main

import (
	logger "TServerGo/Log"
	configs "TServerGo/TServer/Configs"
	"TServerGo/TServer/NotifySystem"
	pbhandler "TServerGo/TServer/PBHandler"
	"TServerGo/TServer/dbproxy"
	"flag"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

var (
	Secret      = flag.String("SECRET", "", "please set SECRET")
	AppId       = flag.String("APP_ID", "", "please set APP_ID")
	Mysql       = flag.String("MYSQL", "", "please set mysql")
	MysqlHost   = flag.String("MYSQL_HOST", "", "please set mysql host")
	IsWebsocket = flag.Bool("IS_WEBSOCKET", true, "please set net type")
)

func hello(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		// Write
		if err != nil {
			c.Logger().Error(err)
		}

		// Read
		msgType, msg, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
			if msgType < 0 || err.(*websocket.CloseError).Code == 1005 {
				NotifySystem.NotifyExec(NotifySystem.NotifyTypeRoleLogout, NotifySystem.NotifyRoleLogoutParam{
					OpenId:     "",
					RemoteAddr: ws.RemoteAddr().String(),
				})
				c.Logger().Error("client closed. RemoteAddr:", ws.RemoteAddr().String())
				break
			}
		}
		pbhandler.GetHandler("json").HandlerPB(&pbhandler.ConnectInfo{WS: ws}, msg)
		// handlerJson(ws, msg)
		log.Printf("Recv %s\n", msg)
	}
	return nil
}

func main() {
	logger.Init("TServer", logger.LogLevelDEBUG|logger.LogLevelINFO|logger.LogLevelWARN|logger.LogLevelERROR)
	flag.Parse()
	configs.Secret = *Secret
	configs.AppId = *AppId
	// 数据库初始化
	db := dbproxy.Instance()
	db.Init(*Mysql, *MysqlHost)
	db.Sync()

	// http服务初始化
	if *IsWebsocket {
		e := echo.New()
		e.Use(middleware.Logger())
		e.Use(middleware.Recover())
		e.GET("/ws", hello)
		e.Logger.Fatal(e.Start(":8081"))
	}
	// tcp服务初始化
	logger.Debug("init tcp...")
	ln, err := net.Listen("tcp", "127.0.0.1:8081")
	if err != nil {
		logger.Errorf("net error, err: %v", err)
		return
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			logger.Errorf("net error, err: %v", err)
		}
		go handlerConnect(conn)
	}
}

func handlerConnect(conn net.Conn) {
	defer conn.Close()
	logger.Debugf("Connected, Addr:%v", conn.RemoteAddr())
	connMap := make(map[string]*pbhandler.ConnectInfo)
	connectInfo, ok := connMap[conn.RemoteAddr().String()]
	if !ok {
		connectInfo = &pbhandler.ConnectInfo{
			SOCKET: conn,
		}
	}
	for {
		var msg = make([]byte, 1024)
		len, err := conn.Read(msg)
		if err != nil || len == 0 {
			logger.Errorf("net error, err:%v", err)
			break
		}
		logger.Debugf("Recv msg, len:%v msg:%v", len, msg[:len])
		pbhandler.GetHandler("pb").HandlerPB(connectInfo, msg[:len])
	}
}
