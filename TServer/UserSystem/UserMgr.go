package UserSystem

import (
	"TServerGo/TServer/NotifySystem"
	"TServerGo/TServer/dbproxy"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"xorm.io/xorm"
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
	// 登陆成功后 需要更新数据库
	user := &dbproxy.User{}
	has, err := dbproxy.Instance().Engine.Where("open_id=?", u.OpenId).Get(user)
	dbproxy.Instance().Transaction(func(session *xorm.Session) (interface{}, error) {
		if err != nil || !has {
			user.UserName = u.NickName
			user.OpenId = u.OpenId
			res, err := session.Insert(user)
			return res, err
		} else if has {
			fmt.Println(user.Created)
			user.UserName = u.NickName
			user.OpenId = u.OpenId
			res, err := session.Update(user)
			return res, err
		}
		return nil, nil
	})

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
