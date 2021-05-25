package main

import (
	"TServer/MatchSystem"
	"TServer/PB"
	"TServer/RoomSystem"
	"TServer/UserSystem"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

const (
	SECRET = "5fb51181444d801fcc9aa42b44f86b57"
	APP_ID = "wx76bf9a66a06b39c3"
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
				c.Logger().Error("client closed. RemoteAddr:", ws.RemoteAddr().String())
				break
			}
		}
		handlerJson(ws, msg)
		log.Printf("Recv %s\n", msg)
	}
	return nil
}

func handlerJson(ws *websocket.Conn, msg []byte) ([]byte, error) {
	m := make(map[string]int)
	json.Unmarshal(msg, &m)
	switch m["id"] {
	case 1001: // ping
		req := &PB.Pong{}
		json.Unmarshal(msg, &req)
		res, err := json.Marshal(&PB.Pong{Id: 1002, Timestamp: req.Timestamp})
		UserSystem.GetPlayerByAddr(ws.RemoteAddr().String()).SendChannel <- res
		return res, err
	case 1101: // 登陆
		req := &PB.LoginReq{}
		err := json.Unmarshal(msg, &req)
		if err != nil {
			return nil, nil
		}
		wxLogin := handlerGetWXLogin(req.Token)
		res, err := json.Marshal(&PB.LoginAck{
			Id:        1102,
			ErrorCode: "SUCCESS",
			OpenId:    wxLogin.Openid,
		})
		if wxLogin.Openid == "" {
			uuid, _ := uuid.NewUUID()
			wxLogin.Openid = uuid.String()
		}
		player := &UserSystem.Player{
			OpenId:      wxLogin.Openid,
			NickName:    req.NickName,
			AvatarUrl:   req.AvatarUrl,
			RemoteAddr:  ws.RemoteAddr().String(),
			SessionKey:  wxLogin.Session_key,
			SendChannel: make(chan []byte),
			Conn:        ws,
		}
		UserSystem.PlayerLogin(player)
		player.SendChannel <- res
		return res, err
	case 1201: // 匹配
		req := &PB.MatchReq{}
		err := json.Unmarshal(msg, &req)
		if err != nil {
			return nil, nil
		}
		player, ok := UserSystem.PlayerRemoteMap[ws.RemoteAddr().String()]
		if !ok {
			return nil, nil
		}
		MatchSystem.JoinMatch(player)

		res, err := json.Marshal(&PB.MatchAck{
			Id:        1202,
			ErrorCode: "MATCHING",
		})
		player.SendChannel <- res
	case 1301: // 走棋
		req := &PB.ChessStepReq{}
		err := json.Unmarshal(msg, &req)
		if err != nil {
			return nil, nil
		}
		player, ok := UserSystem.PlayerRemoteMap[ws.RemoteAddr().String()]
		if !ok {
			log.Println("player is nil, RemoteAddr:", ws.RemoteAddr().String())
			return nil, nil
		}
		room, ok := RoomSystem.RoomOpenIdMap[player.OpenId]
		if !ok {
			log.Println("room is nil, OpenId:", player.OpenId)
		}
		if room.TurnId != player.OpenId {
			res, _ := json.Marshal(&PB.ChessStepAck{
				Id:        0,
				ErrorCode: "SUCCESS",
				Steps:     room.ChessStepList,
			})
			return res, nil
		}
		room.MsgChannel <- req.Step
	}

	return nil, nil
}

func handlerGetWXLogin(token string) *PB.WXLoginAck {
	res, err := http.Get("https://api.weixin.qq.com/sns/jscode2session?" +
		"appid=" + APP_ID + "&secret=" + SECRET + "&js_code=" + token + "&grant_type=authorization_code")
	if err != nil {
		print(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		print(err)
	}
	ack := &PB.WXLoginAck{}
	json.Unmarshal(body, ack)
	return ack
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/ws", hello)
	e.Logger.Fatal(e.Start(":8081"))
}
