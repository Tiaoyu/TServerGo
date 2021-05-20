# TServerGo


## 目标
想要做一个五子棋的小游戏，后端用go实现，前端用微信小游戏。因此通信协议应该就只能使用websocket了

## TDL

1. go实现websocket
2. websocket+json协议交互

## go labstack

```
go get github.com/labstack/echo/v4/middleware
go get github.com/labstack/echo/v4
```

## 问题

1. websocket server无法被连接上
：增加子域名game.tiaoyuyu.com, 并加上ssl证书 使用nginx将子域名代理到websocket服务端口上
