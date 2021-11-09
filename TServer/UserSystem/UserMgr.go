package UserSystem

import (
	"TServerGo/TServer/NotifySystem"
	"TServerGo/TServer/Sessionx"
	"TServerGo/TServer/dbproxy"
	"fmt"
	"log"
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
	Sess        *Sessionx.Session
}

var (
	PlayerOpenIdMap = make(map[string]*Player)
	PlayerRemoteMap = make(map[string]*Player)
)

func init() {
	NotifySystem.NotifyRegister(NotifySystem.NotifyTypeRoleLogout, PlayerLogout)
}

func PlayerLogin(u *Player) {
	if oldUser, ok := PlayerOpenIdMap[u.OpenId]; ok {
		// 已有登陆角色直接顶号替换
		PlayerOpenIdMap[oldUser.OpenId] = u
		PlayerRemoteMap[u.RemoteAddr] = u
		delete(PlayerRemoteMap, oldUser.RemoteAddr)
		oldUser.Sess.Close()
	} else {
		PlayerOpenIdMap[u.OpenId] = u
		PlayerRemoteMap[u.RemoteAddr] = u
	}

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

	log.Printf("%v login success, OpenId:%v RemoteAddr:%v", u.NickName, u.OpenId, u.RemoteAddr)

	NotifySystem.NotifyExec(NotifySystem.NotifyTypeRoleLoginIn, NotifySystem.NotifyRoleLoginParam{
		OpenId:     u.OpenId,
		RemoteAddr: u.RemoteAddr,
	})
}

func PlayerLogout(params ...interface{}) {
	param := params[0].(NotifySystem.NotifyRoleLogoutParam)
	if tmp, ok := PlayerRemoteMap[param.RemoteAddr]; ok {
		delete(PlayerRemoteMap, tmp.RemoteAddr)
		delete(PlayerOpenIdMap, tmp.OpenId)
		tmp.Sess.Close()
		log.Printf("Player closed. OpenId:%v RemoteAddr:%v NickName:%v", tmp.OpenId, tmp.RemoteAddr, tmp.NickName)
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
