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

	fun := func(session *xorm.Session) (interface{}, error) {
		user := &dbproxy.User{}
		has, err := session.Where("open_id=?", u.OpenId).Get(user)
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
		return nil, err
	}
	// 登陆成功后 需要更新数据库
	_, err := dbproxy.Instance().Transaction(fun)
	if err != nil {
		log.Printf("Login failed! OpenId:%v NickName:%v Error:%v\n", u.OpenId, u.NickName, err)
	}

	go func() {
		for {
			select {
			case res := <-u.SendChannel:
				u.Conn.WriteMessage(websocket.TextMessage, res)
			}
		}
	}()
	log.Printf("%v login success, OpenId:%v RemoteAddr:%v", u.NickName, u.OpenId, u.RemoteAddr)
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
