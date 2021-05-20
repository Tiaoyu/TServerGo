package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
		err := ws.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))
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
		res, err := handlerJson(msg)
		ws.WriteMessage(websocket.TextMessage, res)
		fmt.Printf("%s\n", msg)
	}
	return nil
}

func handlerJson(msg []byte) ([]byte, error) {
	m := make(map[string]int)
	json.Unmarshal(msg, &m)
	switch m["id"] {
	case 1001:
		req := &Pong{}
		json.Unmarshal(msg, &req)
		res, err := json.Marshal(&Pong{Id: 1002, Timestamp: req.Timestamp})
		return res, err
	case 1101:
		req := &LoginReq{}
		json.Unmarshal(msg, &req)
		wxLogin := handlerGetWXLogin(req.Token)
		res, err := json.Marshal(&LoginAck{
			Id:        1102,
			ErrorCode: "SUCCESS",
			OpenId:    wxLogin.Openid,
		})
		return res, err
	}

	return nil, nil
}

func handlerGetWXLogin(token string) *WXLoginAck {
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
	ack := &WXLoginAck{}
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
