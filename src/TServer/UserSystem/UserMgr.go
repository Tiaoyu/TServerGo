package UserSystem

import (
	"github.com/gorilla/websocket"
	"log"
)

// Player 玩家数据
type Player struct {
	OpenId      string
	NickName    string
	AvatarUrl   string
	RemoteAddr  string
	SessionKey  string
	SendChannel chan []byte
	Conn        *websocket.Conn
}

var (
	PlayerOpenIdMap = make(map[string]*Player)
	PlayerRemoteMap = make(map[string]*Player)
)

func PlayerLogin(u *Player) {
	PlayerRemoteMap[u.RemoteAddr] = u
	PlayerOpenIdMap[u.OpenId] = u
	go func() {
		for {
			select {
			case res := <-u.SendChannel:
				u.Conn.WriteMessage(websocket.TextMessage, res)
			}
		}
	}()
	log.Println(u.NickName, " login success, OpenId:", u.OpenId)
}

func GetPlayerByAddr(addr string) *Player {
	p, ok := PlayerRemoteMap[addr]
	if !ok {
		return nil
	}
	return p
}

func GetPlayerByOpenId(openId string) *Player {
	p, ok := PlayerOpenIdMap[openId]
	if !ok {
		return nil
	}
	return p
}
