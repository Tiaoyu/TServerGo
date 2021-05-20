package main

type Ping struct {
	Id        int
	Timestamp int
}

type Pong struct {
	Id        int
	Timestamp int
}

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
