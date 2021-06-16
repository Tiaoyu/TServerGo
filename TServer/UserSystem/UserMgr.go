package UserSystem

import (
	"TServerGo/TServer/NotifySystem"
	"TServerGo/dbproxy"
	"fmt"
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

func init() {
	NotifySystem.NotifyRegister(NotifySystem.NotifyTypeRoleLogout, PlayerLogout)
}

func PlayerLogin(u *Player) {
	PlayerRemoteMap[u.RemoteAddr] = u
	PlayerOpenIdMap[u.OpenId] = u
	NotifySystem.NotifyExec(NotifySystem.NotifyTypeRoleLoginIn, &NotifySystem.NotifyRoleLoginParam{
		OpenId:     u.OpenId,
		RemoteAddr: u.RemoteAddr,
	})
	//TODO 登陆成功后 需要更新数据库
	user := &dbproxy.User{}
	u.OpenId = "0311lZ000MMFKL1qpI100iKyxo41lZ0s"
	has, err := dbproxy.Instance().Engine.Where("open_id=?", u.OpenId).Get(user)

	if err != nil || !has {
		user.UserName = u.NickName
		user.OpenId = u.OpenId
		dbproxy.Instance().Engine.Insert(user)
	} else if has {
		fmt.Println(user.Created)
		user.UserName = u.NickName
		user.OpenId = u.OpenId
		dbproxy.Instance().Engine.Update(user)
	}

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

func PlayerLogout(params ...interface{}) {
	param := params[0].(NotifySystem.NotifyRoleLogoutParam)
	if tmp, ok := PlayerRemoteMap[param.RemoteAddr]; ok {
		delete(PlayerRemoteMap, tmp.RemoteAddr)
	}
	if tmp, ok := PlayerRemoteMap[param.RemoteAddr]; ok {
		delete(PlayerOpenIdMap, tmp.OpenId)
	}
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
