package main

import (
	"fmt"
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
		fmt.Printf("%s\n", msg)
	}
	return nil
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/ws", hello)
	e.Logger.Fatal(e.Start(":8081"))
}
