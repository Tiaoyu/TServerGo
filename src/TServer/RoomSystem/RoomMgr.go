//角色匹配后加入匹配池子
//匹配线程每秒从池子中取若干角色
//有足够角色后进行创建房间
//创建完房间后开始进入房间线程

package RoomSystem

import (
	"log"
	gamepb "pb"
	"time"
)

type Room struct {
	RedId         int32
	BlackId       int32
	CreateTime    int64
	ChessStepList []gamepb.ChessStep
	GobangInfo    [15][15]int32
	TurnId        int32 // 当前手
}

type Role struct {
	RoleId int32
	Score  int32
}

var (
	MatchPool   map[int32]struct{}
	RoomUserMap map[int32]*Room
)

// JoinMatch 加入匹配
func JoinMatch(roleId int32) {
	MatchPool[roleId] = struct{}{}
}

// Matching 匹配
func Matching(roleId int32) {
	for {
		var roleList []int32
		if len(MatchPool) > 2 {
			for k, _ := range MatchPool {
				roleList = append(roleList, k)
				delete(MatchPool, k)
				if len(roleList) == 2 {
					break
				}
			}

			room := &Room{RedId: roleList[0], BlackId: roleList[1], CreateTime: time.Now().Unix(),
				ChessStepList: make([]gamepb.ChessStep, 0), GobangInfo: [15][15]int32{}, TurnId: roleList[0]}
			go RoomLogic(room)
		}
	}
}

func RoomLogic(room *Room) {
	log.Println("Match success. RedId:", room.RedId, " BlackId:", room.BlackId)
	RoomUserMap[room.RedId] = room
	RoomUserMap[room.BlackId] = room
	//gobang := [15][15]int32{
	//	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	//	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	//	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	//	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	//	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	//	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	//	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	//	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	//	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	//	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	//	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	//	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	//	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	//	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	//	{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	//}
	// 游戏主循环
	for {
		// TODO
		log.Println(time.Now().Unix())
	}
}
