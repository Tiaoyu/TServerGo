package pbhandler

import (
	configs "TServerGo/TServer/Configs"
	"TServerGo/TServer/MatchSystem"
	"TServerGo/TServer/PB"
	"TServerGo/TServer/RoomSystem"
	"TServerGo/TServer/UserSystem"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type HandlerJson struct {
}

func (h *HandlerJson) HandlerPB(conn *ConnectInfo, msg []byte) ([]byte, error) {
	m := make(map[string]int)
	json.Unmarshal(msg, &m)
	switch m["id"] {
	case 1001: // ping
		req := &PB.Pong{}
		json.Unmarshal(msg, &req)
		res, err := json.Marshal(&PB.Pong{Id: 1002, Timestamp: req.Timestamp})
		UserSystem.GetPlayerByAddr(conn.WS.RemoteAddr().String()).SendChannel <- res
		return res, err
	case 1101: // 登陆
		req := &PB.LoginReq{}
		err := json.Unmarshal(msg, &req)
		if err != nil {
			return nil, nil
		}
		wxLogin := handlerGetWXLogin(req.Token)
		res, err := json.Marshal(&PB.LoginAck{
			Id:        1102,
			ErrorCode: "SUCCESS",
			OpenId:    wxLogin.Openid,
		})
		if wxLogin.Openid == "" {
			wxLogin.Openid = req.Token
		}
		player := &UserSystem.Player{
			OpenId:      wxLogin.Openid,
			NickName:    req.NickName,
			AvatarUrl:   req.AvatarUrl,
			RemoteAddr:  conn.WS.RemoteAddr().String(),
			SessionKey:  wxLogin.Session_key,
			SendChannel: make(chan []byte),
			Conn:        conn.WS,
		}
		UserSystem.PlayerLogin(player)
		player.SendChannel <- res
		return res, err
	case 1201: // 匹配
		req := &PB.MatchReq{}
		err := json.Unmarshal(msg, &req)
		if err != nil {
			return nil, nil
		}
		player, ok := UserSystem.PlayerRemoteMap[conn.WS.RemoteAddr().String()]
		if !ok {
			return nil, nil
		}
		if req.MatchType == PB.MatchTypeMatch {
			MatchSystem.JoinMatch(player)
		} else if req.MatchType == PB.MatchTypeCancel {
			MatchSystem.CancelMatch(player)
		}

		res, err := json.Marshal(&PB.MatchAck{
			Id:        1202,
			ErrorCode: "MATCHING",
		})
		player.SendChannel <- res
	case 1301: // 走棋
		req := &PB.ChessStepReq{}
		err := json.Unmarshal(msg, &req)
		if err != nil {
			return nil, nil
		}
		player, ok := UserSystem.PlayerRemoteMap[conn.WS.RemoteAddr().String()]
		if !ok {
			log.Println("player is nil, RemoteAddr:", conn.WS.RemoteAddr().String())
			return nil, nil
		}
		value, ok := RoomSystem.RoomOpenIdMap.Load(player.OpenId)
		if !ok {
			log.Println("room is nil, OpenId:", player.OpenId)
		}
		room := value.(*RoomSystem.Room)
		if room.TurnId != player.OpenId {
			res, _ := json.Marshal(&PB.ChessStepAck{
				Id:        0,
				ErrorCode: "NOT YOUR TURN",
				Steps:     room.ChessStepList,
			})
			return res, nil
		}
		if room.BlackId == player.OpenId {
			req.Step.Color = PB.ColorTypeBlack
		} else if room.RedId == player.OpenId {
			req.Step.Color = PB.ColorTypeRed
		}
		room.MsgChannel <- req.Step
	}

	return nil, nil
}

// 微信平台账号验证
func handlerGetWXLogin(token string) *PB.WXLoginAck {
	res, err := http.Get("https://api.weixin.qq.com/sns/jscode2session?" +
		"appid=" + configs.AppId + "&secret=" + configs.Secret + "&js_code=" + token + "&grant_type=authorization_code")
	if err != nil {
		print(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		print(err)
	}
	ack := &PB.WXLoginAck{}
	json.Unmarshal(body, ack)
	return ack
}
