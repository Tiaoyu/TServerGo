package MatchSystem

import (
	"TServer/PB"
	"TServer/RoomSystem"
	"time"
)

var (
	matchPool = make(chan *MatchItem, 2)
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
				RoomSystem.RoomLogic(room)
				pair = make([]*MatchItem, 0)
			}
		}
	}()
}

func JoinMatch(openId, remoteAddr string) {
	item := &MatchItem{
		OpenId:     openId,
		RemoteAddr: remoteAddr,
	}
	matchPool <- item
}
