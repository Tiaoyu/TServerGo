package PB

import "encoding/json"

// ToJsonBytes 协议转json字节流
func ToJsonBytes(v interface{}) []byte {
	res, _ := json.Marshal(v)
	return res
}

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

// 匹配类型
const (
	MatchTypeMatch  = 1
	MatchTypeCancel = 2
)

// 对局双方类型
const (
	ColorTypeRed   = 1 // 红方
	ColorTypeBlack = 2 // 黑方
)

type MatchAck struct {
	Id             int
	ErrorCode      string
	EnemyName      string // 敌方阵营
	EnemyAvatarUrl string // 敌方头像
	Color          int    // 阵营（红方、黑方）
}

// ChessStepReq 走棋
type ChessStepReq struct {
	Id    int
	Step  ChessStep
	Color int
}
type ChessStepAck struct {
	Id        int
	ErrorCode string
	Steps     []ChessStep
}

// ChessStep 一步棋
type ChessStep struct {
	Pos   Pos
	Color int
}

// Pos 位置
type Pos struct {
	X, Y int32
}

// GameResultAck 对局结果
type GameResultAck struct {
	Id         int
	ErrorCode  string
	GameResult string
}
