package main

import (
	"TServerGo/log"
	"sync"

	"xorm.io/xorm"
)

type PlayerState uint8

const (
	PlayerStateOnline   PlayerState = 1
	PlayerStateMatching PlayerState = 2
	PlayerStateInRoom   PlayerState = 3
	PlayerStateLogout   PlayerState = 4
)

// Player 玩家数据
type Player struct {
	OpenId string
	Sess   *UserSession
	State  PlayerState
}

var (
	PlayerOpenIdMap = make(map[string]*Player)
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
		u.State = oldUser.State
		// 已有登陆角色直接顶号替换
		PlayerOpenIdMap[oldUser.OpenId] = u
		oldUser.Sess.Close()
	} else {
		u.State = PlayerStateOnline
		PlayerOpenIdMap[u.OpenId] = u
	}

	fun := func(session *xorm.Session) (interface{}, error) {
		user := &User{}
		has, err := session.Where("open_id=?", u.OpenId).Get(user)
		if err != nil || !has {
			user.OpenId = u.OpenId
			res, err := session.Insert(user)
			return res, err
		} else if has {
			user.OpenId = u.OpenId
			res, err := session.Update(user)
			return res, err
		}
		return nil, err
	}
	// 登陆成功后 需要更新数据库
	_, err := dbProxy.Transaction(fun)
	if err != nil {
		log.Errorf("Login failed! OpenId:%v Error:%v", u.OpenId, err)
		return err
	}

	log.Debugf("Login success, OpenId:%v", u.OpenId)

	NotifyExec(NotifyTypeRoleLoginIn, &NotifyRoleLoginParam{
		OpenId: u.OpenId,
	})
	return nil
}

func onPlayerLogout(params ...interface{}) {
	param := params[0].(*NotifyRoleLogoutParam)
	PlayerMapLock.Lock()
	defer PlayerMapLock.Unlock()
	if tmp, ok := PlayerOpenIdMap[param.OpenId]; ok {
		delete(PlayerOpenIdMap, tmp.OpenId)
		tmp.Sess.Close()
		log.Debugf("Player closed. OpenId:%v", tmp.OpenId)
	}
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
