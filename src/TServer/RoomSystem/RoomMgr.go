//角色匹配后加入匹配池子
//匹配线程每秒从池子中取若干角色
//有足够角色后进行创建房间
//创建完房间后开始进入房间线程

package RoomSystem

import (
	"TServer/PB"
	"TServer/UserSystem"
	"encoding/json"
	"log"
)

type Room struct {
	RedId         string
	BlackId       string
	CreateTime    int64
	ChessStepList []PB.ChessStep
	GobangInfo    [15][15]string
	TurnId        string // 当前手
	MsgChannel    chan PB.ChessStep
}

var (
	RoomOpenIdMap = make(map[string]*Room)
)

func RoomLogic(room *Room) error {
	redPlayer, ok := UserSystem.PlayerOpenIdMap[room.RedId]
	if !ok {
		return nil
	}
	blackPlayer, ok := UserSystem.PlayerOpenIdMap[room.BlackId]
	if !ok {
		return nil
	}
	res, _ := json.Marshal(&PB.MatchAck{
		Id:             1302,
		ErrorCode:      "SUCCESS",
		EnemyName:      blackPlayer.NickName,
		EnemyAvatarUrl: blackPlayer.AvatarUrl,
		Color:          "RED",
	})
	redPlayer.SendChannel <- res
	res, _ = json.Marshal(&PB.MatchAck{
		Id:             1302,
		ErrorCode:      "SUCCESS",
		EnemyName:      redPlayer.NickName,
		EnemyAvatarUrl: redPlayer.AvatarUrl,
		Color:          "BLACK",
	})
	blackPlayer.SendChannel <- res

	RoomOpenIdMap[room.RedId] = room
	RoomOpenIdMap[room.BlackId] = room
	room.MsgChannel = make(chan PB.ChessStep, 0)

	go func() {
		for {
			select {
			case step := <-room.MsgChannel:
				log.Println(step.Color, " step to ", step.Pos)
				room.GobangInfo[step.Pos.X][step.Pos.Y] = step.Color
				room.ChessStepList = append(room.ChessStepList, step)
				res, _ := json.Marshal(&PB.ChessStepAck{
					Id:        1302,
					ErrorCode: "SUCCESS",
					Steps:     room.ChessStepList,
				})
				log.Println("step list:", string(res))
				UserSystem.GetPlayerByOpenId(room.RedId).SendChannel <- res
				UserSystem.GetPlayerByOpenId(room.BlackId).SendChannel <- res
			}
			// 判断胜负
		}
		delete(RoomOpenIdMap, room.RedId)
		delete(RoomOpenIdMap, room.BlackId)
		log.Println("room destroyed!")
	}()
	log.Println("create room success, RedId:", room.RedId, " BlackId:", room.BlackId)
	return nil
}
