package main

import (
	"TServerGo/log"
	"TServerGo/pb"
	"sync"
	"time"
)

var (
	matchPool  = make(chan *MatchItem) // 匹配池子
	cancelPool = make(chan *MatchItem) // 取消匹配池子
	matchMap   = new(sync.Map)         // 在匹配中的
)

type MatchItem struct {
	OpenId     string
	RemoteAddr string
}

func initMatch() {
	NotifyRegister(NotifyTypeRoleLogout, onMatchPlayerLogout)
	// 匹配线程
	go func() {
		pair := make([]*MatchItem, 0)
		for {
			select {
			case item := <-matchPool:
				pair = append(pair, item)
				log.Debugf("Join match queue success! OpenId:%v", item.OpenId)
			case item := <-cancelPool:
				for i, p := range pair {
					if p.OpenId == item.OpenId {
						pair = append(pair[:i], pair[i+1:]...)
						break
					}
				}
				matchMap.Delete(item.OpenId)
				// 返回取消匹配成功
				if player := GetPlayerByOpenId(item.OpenId); player != nil {
					tmp := SendMsg(&pb.S2CMatch{
						Result: pb.MatchResult_MatResultCancel,
					}, pb.ProtocolType_ES2CMatch)
					player.Sess.SendChannel <- tmp
				}
			}

			// 只要匹配到两个就进行创建房间逻辑
			if len(pair) >= 2 {
				room := &Room{
					RedId:         pair[0].OpenId,
					BlackId:       pair[1].OpenId,
					CreateTime:    time.Now().Unix(),
					ChessStepList: make([]*pb.ChessStep, 0),
					GobangInfo:    [15][15]int32{},
					TurnId:        pair[0].OpenId,
					GoBangTemp:    [15][15]*Piece{},
				}
				matchMap.Delete(pair[0].OpenId)
				matchMap.Delete(pair[1].OpenId)
				RoomLogic(room)
				pair = make([]*MatchItem, 0)
				log.Debugf("Match success! %v vs %v", room.RedId, room.BlackId)
			}
		}
	}()

}

// JoinMatch 加入匹配
func JoinMatch(player *Player) {
	if _, ok := matchMap.Load(player.OpenId); ok {
		log.Warnf("%v already in mathing, so join match failed!", player.OpenId)
		return
	}

	item := &MatchItem{
		OpenId: player.OpenId,
	}

	matchPool <- item
	log.Debugf(player.OpenId, " join to match.")
}

// CancelMatch 取消匹配
func CancelMatch(player *Player) {
	item := &MatchItem{
		OpenId: player.OpenId,
	}

	cancelPool <- item
	log.Debugf(player.OpenId, " cancel match.")
}

func CancelMatchById(openId, remoteAddr string) {
	item := &MatchItem{
		OpenId:     openId,
		RemoteAddr: remoteAddr,
	}
	cancelPool <- item
	log.Debugf(openId, " cancel match.")
}

func onMatchPlayerLogout(params ...interface{}) {
	param := params[0].(*NotifyRoleLogoutParam)
	if _, ok := matchMap.Load(param.OpenId); ok {
		CancelMatchById(param.OpenId, param.RemoteAddr)
	}
}