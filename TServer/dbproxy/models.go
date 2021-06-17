package dbproxy

// User 账号
type User struct {
	Id          int64
	OpenId      string
	UserName    string
	Score       uint32
	WinCount    uint32
	FailedCount uint32
	Created     int64 `xorm:"created"`
	Updated     int64 `xorm:"updated"`
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
