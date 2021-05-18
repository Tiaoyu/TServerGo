package main

// User 账号
type User struct {
	UserName  string `xorm:"varchar(25) notnull unique 'user_name' comment('账号名') pk"`
	CreatedAt int64  `xorm:"created"`
}

// Role 角色
type Role struct {
	RoleName  string `xorm:"varchar(25) notnull unique 'user_name' comment('角色名') pk"`
	CreatedAt int64  `xorm:"created"`
}

// Race 赛局
type Race struct {
	RedRoleId   int64 `xorm:"pk"`
	BlackRoleId int64 `xorm:"pk"`
	WinnerId    int64
	GobangInfo  string `xorm:"longtext notnull comment('当局棋盘信息')"`
	CreatedAt   int64  `xorm:"created"`
}
