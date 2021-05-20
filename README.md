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


##
1. 如何在Canvas里面放一张图片，
2. 在Canvas放置一枚棋子，位置怎么确定

## 功能

### PING

Req ID:1001


| -         | -   |            |
| --------- | --- | ---------- |
| id        | int | 1001       |
| timestamp | int | 时间戳(秒) |

Res ID:1002

| -         | -   |            |
| --------- | --- | ---------- |
| id        | int | 1002       |
| timestamp | int | 时间戳(秒) |


```json
{"id":1001,"timestamp":1621521841}
```

### 登陆

Req ID : 1101

| -         | -      | -       |
| --------- | ------ | ------- |
| id        | int    | 1101    |
| nickName  | string | 昵称    |
| avatarUrl | string | 头像    |
| token     | string | js_code |

Res ID : 1102

| -         | -      | -        |
| --------- | ------ | -------- |
| id        | int    | 1101     |
| errorCode | string | 错误信息 |
| openId    | string |          |

```json
{
    "id": 1101,
    "token":"031IAJ0w32ruoW2bJJ3w3eaBxX1IAJ0A",
    "nickName":"0.0",
    "avatarUrl":"https://thirdwx.qlogo.cn/mmopen/vi_32/Q0j4TwGTfTIcOPgba5had6WBXqu8V1ZFsdcVkJy83RfWbrQ5k0qQkg8F0AWJqACUmKq6Oib1ASdyibq8CPl283QA/132"
}
```


## 协议

加入t转成 byte数组后，数组长度为10
t这个协议号是 1001，1001转成byte数组后，数组长度为4
10 + 4 + 4 = 18
[长度][协议号][协议本身]



属性	类型	说明
openid	string	用户唯一标识
session_key	string	会话密钥
unionid	string
errcode	number	错误码
errmsg	string	错误信息
