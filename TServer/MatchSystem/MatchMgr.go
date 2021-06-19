package MatchSystem

import (
	"TServerGo/TServer/NotifySystem"
	"TServerGo/TServer/PB"
	"TServerGo/TServer/RoomSystem"
	"TServerGo/TServer/UserSystem"
	"log"
	"sync"
	"time"
)

var (
	matchPool  = make(chan *MatchItem, 0) // 匹配池子
	cancelPool = make(chan *MatchItem, 0) // 取消匹配池子
	matchMap   = new(sync.Map)            // 在匹配中的
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
				log.Printf("Join match queue success! OpenId:%v\n", item.OpenId)
			case item := <-cancelPool:
				for i, p := range pair {
					if p.OpenId == item.OpenId {
						pair = append(pair[:i], pair[i+1:]...)
						break
					}
				}
				matchMap.Delete(item.OpenId)
				// 返回取消匹配成功
				if player := UserSystem.GetPlayerByOpenId(item.OpenId); player != nil {
					player.SendChannel <- PB.ToJsonBytes(&PB.MatchAck{
						Id:        1202,
						ErrorCode: "CANCEL",
					})
				}
			}

			// 只要匹配到两个就进行创建房间逻辑
			if len(pair) >= 2 {
				room := &RoomSystem.Room{
					RedId:         pair[0].OpenId,
					BlackId:       pair[1].OpenId,
					CreateTime:    time.Now().Unix(),
					ChessStepList: make([]PB.ChessStep, 0),
					GobangInfo:    [15][15]int{},
					TurnId:        pair[0].OpenId,
					GoBangTemp:    [15][15]*RoomSystem.Piece{},
				}
				matchMap.Delete(pair[0].OpenId)
				matchMap.Delete(pair[1].OpenId)
				RoomSystem.RoomLogic(room)
				pair = make([]*MatchItem, 0)
				log.Printf("Match success! %v vs %v\n", room.RedId, room.BlackId)
			}
		}
	}()

	NotifySystem.NotifyRegister(NotifySystem.NotifyTypeRoleLogout, PlayerLogout)
}

// JoinMatch 加入匹配
func JoinMatch(player *UserSystem.Player) {
	if _, ok := matchMap.Load(player.OpenId); ok {
		log.Printf("%v already in mathing, so join match failed!", player.OpenId)
		return
	}

	item := &MatchItem{
		OpenId:     player.OpenId,
		RemoteAddr: player.RemoteAddr,
	}

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
	if _, ok := matchMap.Load(param.OpenId); ok {
		CancelMatchById(param.OpenId, param.RemoteAddr)
	}
}
