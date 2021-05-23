package PB

// Ping PING
type Ping struct {
	Id        int
	Timestamp int
}

type Pong struct {
	Id        int
	Timestamp int
}

// LoginReq 登陆
type LoginReq struct {
	Id        int
	NickName  string
	AvatarUrl string
	Token     string
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
	Id int
}

type MatchAck struct {
	Id             int
	ErrorCode      string
	EnemyName      string
	EnemyAvatarUrl string
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
