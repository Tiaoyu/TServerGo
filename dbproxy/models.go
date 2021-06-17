package dbproxy

// User 账号
type User struct {
	Id       int64
	OpenId   string
	UserName string
	Score    uint32
	Created  int64 `xorm:"created"`
	Updated  int64 `xorm:"updated"`
}

// Role 角色
type Role struct {
	Id       int64
	RoleName string `xorm:"varchar(25) notnull unique 'user_name' comment('角色名') pk"`
	Created  int64  `xorm:"created"`
	Updated  int64  `xorm:"updated"`
}

// Race 赛局
type Race struct {
	Id          int64
	RedRoleId   int64 `xorm:"pk"`
	BlackRoleId int64 `xorm:"pk"`
	WinnerId    int64
	GobangInfo  string `xorm:"longtext notnull comment('当局棋盘信息')"`
	Created     int64  `xorm:"created"`
	Updated     int64  `xorm:"updated"`
}
