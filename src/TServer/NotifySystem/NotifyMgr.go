package NotifySystem

import "log"

// NotifyType 通知类型
type NotifyType uint8

const (
	// NotifyTypeRoleLogout 角色登出
	NotifyTypeRoleLogout NotifyType = 0
	// NotifyTypeRoleLoginIn 角色登录
	NotifyTypeRoleLoginIn NotifyType = 1
)

var (
	NotifyMap map[NotifyType][]func(...interface{})
)

func init() {
	NotifyMap = make(map[NotifyType][]func(...interface{}))
	NotifyMap[NotifyTypeRoleLogout] = []func(...interface{}){}
	NotifyMap[NotifyTypeRoleLoginIn] = []func(...interface{}){}
}

// NotifyRegister 注册事件
func NotifyRegister(t NotifyType, f func(...interface{})) {
	NotifyMap[t] = append(NotifyMap[t], f)
}

// NotifyExec 执行事件
func NotifyExec(t NotifyType, param ...interface{}) {
	queue, ok := NotifyMap[t]
	if !ok {
		log.Printf("NotifyMap is nil, type:%v\n", t)
	}

	switch t {
	case NotifyTypeRoleLogout:
		for _, event := range queue {
			event(param[0].(NotifyRoleLogoutParam))
		}
	case NotifyTypeRoleLoginIn:
		for _, event := range queue {
			event(param[0].(NotifyRoleLoginParam))
		}
	}
}
