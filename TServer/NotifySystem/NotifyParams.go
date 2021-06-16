package NotifySystem

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
