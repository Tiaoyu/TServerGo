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
	OpenId string
}

func initMatch() {
	NotifyRegister(NotifyTypeRoleLogout, onMatchPlayerLogout)
	// 匹配线程
	go func() {
		pair := make([]*MatchItem, 0)
	L:
		for {
			select {
			case item, ok := <-matchPool:
				if !ok {
					log.Error("matchPool channel get error!!!")
					break L
				}
				pair = append(pair, item)
				log.Debugf("Join match queue success! OpenId:%v", item.OpenId)

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
			case item, ok := <-cancelPool:
				if !ok {
					log.Error("cancelPool channel get error!!!")
					break L
				}
				for i, p := range pair {
					if p.OpenId == item.OpenId {
						pair = append(pair[:i], pair[i+1:]...)
						log.Debugf("Cancel match success from match pair! OpenId:%v", item.OpenId)
						break
					}
				}
				if _, isPresent := matchMap.LoadAndDelete(item.OpenId); isPresent {
					log.Debugf("Cancel match success from match map! OpenId:%v", item.OpenId)
				}
				// 返回取消匹配成功
				if player := GetPlayerByOpenId(item.OpenId); player != nil {
					tmp := MsgToBytes(&pb.S2CMatch{
						Result: pb.MatchResult_MatResultCancel,
					}, pb.ProtocolType_ES2CMatch)
					player.Sess.SendChannel <- tmp
				}
			}
		}
	}()

}

// JoinMatch 加入匹配
func JoinMatch(player *Player) {
	if player.State == PlayerStateMatching {
		return
	}
	if _, ok := matchMap.LoadOrStore(player.OpenId, struct{}{}); ok {
		log.Warnf("%v already in matching, so join match failed!", player.OpenId)
		return
	}

	item := &MatchItem{
		OpenId: player.OpenId,
	}

	matchPool <- item
	player.State = PlayerStateMatching
	log.Debugf("Join to match. OpenId:%v", player.OpenId)
}

// CancelMatch 取消匹配
func CancelMatch(player *Player) {
	item := &MatchItem{
		OpenId: player.OpenId,
	}

	cancelPool <- item
	log.Debugf("Cancel match. OpenId:%v", player.OpenId)
}

func CancelMatchById(openId string) {
	item := &MatchItem{
		OpenId: openId,
	}
	cancelPool <- item
	log.Debugf("Cancel match. OpenId:%v", openId)
}

func onMatchPlayerLogout(params ...interface{}) {
	param := params[0].(*NotifyRoleLogoutParam)
	CancelMatchById(param.OpenId)
}
