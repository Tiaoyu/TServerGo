package MatchSystem

import (
	"TServer/NotifySystem"
	"TServer/PB"
	"TServer/RoomSystem"
	"TServer/UserSystem"
	"log"
	"time"
)

var (
	matchPool  = make(chan *MatchItem, 0)     // 匹配池子
	cancelPool = make(chan *MatchItem, 0)     // 取消匹配池子
	matchMap   = make(map[string]struct{}, 0) // 在匹配中的
)

type MatchItem struct {
	OpenId     string
	RemoteAddr string
}

func init() {
	// 匹配线程
	go func() {
		pair := make([]*MatchItem, 0)
		for {
			select {
			case item := <-matchPool:
				pair = append(pair, item)
			case item := <-cancelPool:
				for i, p := range pair {
					if p.OpenId == item.OpenId {
						pair = append(pair[:i], pair[i+1:]...)
						break
					}
				}

				// 返回取消匹配成功
				if player := UserSystem.GetPlayerByOpenId(item.OpenId); player != nil {
					player.SendChannel <- PB.ToJsonBytes(&PB.MatchAck{
						Id:        1202,
						ErrorCode: "CANCEL",
					})
				}
			}

			// 只要匹配到两个就进行创建房间逻辑
			if len(pair) == 2 {
				room := &RoomSystem.Room{
					RedId:         pair[0].OpenId,
					BlackId:       pair[1].OpenId,
					CreateTime:    time.Now().Unix(),
					ChessStepList: make([]PB.ChessStep, 0),
					GobangInfo:    [15][15]int{},
					TurnId:        pair[0].OpenId,
					GoBangTemp:    [15][15]*RoomSystem.Piece{},
				}
				delete(matchMap, pair[0].OpenId)
				delete(matchMap, pair[1].OpenId)
				RoomSystem.RoomLogic(room)
				pair = make([]*MatchItem, 0)
			}
		}
	}()

	NotifySystem.NotifyRegister(NotifySystem.NotifyTypeRoleLogout, PlayerLogout)
}

// JoinMatch 加入匹配
func JoinMatch(player *UserSystem.Player) {
	if _, ok := matchMap[player.OpenId]; ok {
		return
	}

	item := &MatchItem{
		OpenId:     player.OpenId,
		RemoteAddr: player.RemoteAddr,
	}
	matchMap[player.OpenId] = struct{}{}
	matchPool <- item
	log.Println(player.OpenId, " join to match.")
}

// CancelMatch 取消匹配
func CancelMatch(player *UserSystem.Player) {
	item := &MatchItem{
		OpenId:     player.OpenId,
		RemoteAddr: player.RemoteAddr,
	}
	cancelPool <- item
	log.Println(player.OpenId, " cancel match.")
}

func CancelMatchById(openId, remoteAddr string) {
	item := &MatchItem{
		OpenId:     openId,
		RemoteAddr: remoteAddr,
	}
	cancelPool <- item
	log.Println(openId, " cancel match.")
}

func PlayerLogout(params ...interface{}) {
	param := params[0].(NotifySystem.NotifyRoleLogoutParam)
	CancelMatchById(param.OpenId, param.RemoteAddr)
}
