package MatchSystem

import (
	"TServer/PB"
	"TServer/RoomSystem"
	"TServer/UserSystem"
	"log"
	"time"
)

var (
	matchPool = make(chan *MatchItem, 0)
	matchMap  = make(map[string]struct{}, 0)
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
			}
			// 只要匹配到两个就进行创建房间逻辑
			if len(pair) == 2 {
				room := &RoomSystem.Room{
					RedId:         pair[0].OpenId,
					BlackId:       pair[1].OpenId,
					CreateTime:    time.Now().Unix(),
					ChessStepList: make([]PB.ChessStep, 0),
					GobangInfo:    [15][15]string{},
					TurnId:        pair[0].OpenId,
				}
				delete(matchMap, pair[0].OpenId)
				delete(matchMap, pair[1].OpenId)
				RoomSystem.RoomLogic(room)
				pair = make([]*MatchItem, 0)
			}
		}
	}()
}

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
