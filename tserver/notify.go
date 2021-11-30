package main

import "TServerGo/log"

// NotifyType 通知类型
type NotifyType uint8

const (
	NotifyTypeRoleLogout  NotifyType = 0
	NotifyTypeRoleLoginIn NotifyType = 1
)

var (
	NotifyMap map[NotifyType][]func(...interface{})
)

func initNotify() {
	NotifyMap = make(map[NotifyType][]func(...interface{}))
	NotifyMap[NotifyTypeRoleLogout] = []func(...interface{}){}
	NotifyMap[NotifyTypeRoleLoginIn] = []func(...interface{}){}
}

// NotifyRegister 注册事件
func NotifyRegister(t NotifyType, f func(...interface{})) {
	NotifyMap[t] = append(NotifyMap[t], f)
}

// NotifyExec 执行事件
func NotifyExec(notifyType NotifyType, param ...interface{}) {
	queue, ok := NotifyMap[notifyType]
	if !ok {
		log.Debugf("NotifyMap is nil, type:%v", notifyType)
	}

	switch notifyType {
	case NotifyTypeRoleLogout:
		for _, event := range queue {
			event(param[0].(*NotifyRoleLogoutParam))
		}
	case NotifyTypeRoleLoginIn:
		for _, event := range queue {
			event(param[0].(*NotifyRoleLoginParam))
		}
	}
}

// NotifyRoleLogoutParam 角色登出事件参数
type NotifyRoleLogoutParam struct {
	OpenId     string
	RemoteAddr string
}

// NotifyRoleLoginParam 角色登陆事件参数
type NotifyRoleLoginParam struct {
	OpenId     string
	RemoteAddr string
}
