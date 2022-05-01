package main

// User 账号
type User struct {
	OpenId      string `xorm:"pk"`
	UserName    string
	Score       int32  `xorm:"int(11) comment('分数')"`
	WinCount    uint32 `xorm:"int(11) comment('胜场')"`
	FailedCount uint32 `xorm:"int(11) comment('胜场')"`
	Created     int64  `xorm:"created"`
	Updated     int64  `xorm:"updated"`
}

// Race 赛局
type Race struct {
	Id          int64
	RedOpenId   string
	BlackOpenId string
	WinnerId    string
	GobangInfo  string `xorm:"longtext notnull comment('当局棋盘信息')"`
	Created     int64  `xorm:"created"`
	Updated     int64  `xorm:"updated"`
}
