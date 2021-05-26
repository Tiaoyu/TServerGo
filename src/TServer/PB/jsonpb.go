package PB

// Ping PING
type Ping struct {
	Id        int
	Timestamp int // 时间戳
}

type Pong struct {
	Id        int
	Timestamp int // 时间戳
}

// LoginReq 登陆
type LoginReq struct {
	Id        int
	NickName  string // 昵称
	AvatarUrl string // 头像
	Token     string // jc_code
}

type LoginAck struct {
	Id        int
	ErrorCode string
	OpenId    string
}

type WXLoginAck struct {
	Openid      string
	Session_key string
	Unionid     string
	Errcode     int
	Errmsg      string
}

// MatchReq 匹配
type MatchReq struct {
	Id        int
	MatchType int
}

const (
	MatchTypeMatch  = 1
	MatchTypeCancel = 2
)

type MatchAck struct {
	Id             int
	ErrorCode      string
	EnemyName      string // 敌方阵营
	EnemyAvatarUrl string // 敌方头像
	Color          string // 阵营（红方、黑方）
}

// ChessStepReq 走棋
type ChessStepReq struct {
	Id    int
	Step  ChessStep
	Color string
}
type ChessStepAck struct {
	Id        int
	ErrorCode string
	Steps     []ChessStep
}

// ChessStep 一步棋
type ChessStep struct {
	Pos   Pos
	Color string
}

// Pos 位置
type Pos struct {
	X, Y int32
}
