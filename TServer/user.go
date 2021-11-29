package main

import (
	logger "TServerGo/log"
	"fmt"
	"sync"

	"xorm.io/xorm"
)

// Player 玩家数据
type Player struct {
	OpenId     string
	NickName   string
	AvatarUrl  string
	RemoteAddr string
	SessionKey string
	Sess       *UserSession
}

var (
	PlayerOpenIdMap = make(map[string]*Player)
	PlayerRemoteMap = make(map[string]*Player)
	PlayerMapLock   *sync.RWMutex
)

func initUser() {
	PlayerMapLock = &sync.RWMutex{}
	NotifyRegister(NotifyTypeRoleLogout, onPlayerLogout)
}

func PlayerLogin(u *Player) error {
	PlayerMapLock.Lock()
	defer PlayerMapLock.Unlock()
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
		user := &User{}
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
	_, err := dbProxy.Transaction(fun)
	if err != nil {
		logger.Errorf("Login failed! OpenId:%v NickName:%v Error:%v\n", u.OpenId, u.NickName, err)
		return err
	}

	logger.Debugf("%v login success, OpenId:%v RemoteAddr:%v", u.NickName, u.OpenId, u.RemoteAddr)

	NotifyExec(NotifyTypeRoleLoginIn, &NotifyRoleLoginParam{
		OpenId:     u.OpenId,
		RemoteAddr: u.RemoteAddr,
	})
	return nil
}

func onPlayerLogout(params ...interface{}) {
	param := params[0].(*NotifyRoleLogoutParam)
	PlayerMapLock.Lock()
	defer PlayerMapLock.Unlock()
	if tmp, ok := PlayerRemoteMap[param.RemoteAddr]; ok {
		delete(PlayerRemoteMap, tmp.RemoteAddr)
		delete(PlayerOpenIdMap, tmp.OpenId)
		tmp.Sess.Close()
		logger.Debugf("Player closed. OpenId:%v RemoteAddr:%v NickName:%v", tmp.OpenId, tmp.RemoteAddr, tmp.NickName)
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
	PlayerMapLock.RLock()
	defer PlayerMapLock.RUnlock()
	p, ok := PlayerOpenIdMap[openId]
	if !ok {
		return nil
	}
	return p
}
