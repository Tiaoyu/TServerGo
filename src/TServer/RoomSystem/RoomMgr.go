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
	GobangInfo    [15][15]int
	TurnId        string // 当前手
	MsgChannel    chan PB.ChessStep

	GoBangTemp [15][15]*Piece // 对局辅助信息 用来标志每个位置的四个方向的连珠数
}
type Piece struct {
	openId     string // 阵营
	horizontal int    // 横
	vertical   int    // 竖
	lOblique   int    // 左斜
	rOblique   int    // 右斜
}

var (
	RoomOpenIdMap = make(map[string]*Room)
)

// RoomLogic 房间逻辑
func RoomLogic(room *Room) error {
	//获取红方、黑方玩家
	redPlayer, ok := UserSystem.PlayerOpenIdMap[room.RedId]
	if !ok {
		return nil
	}
	blackPlayer, ok := UserSystem.PlayerOpenIdMap[room.BlackId]
	if !ok {
		return nil
	}

	//分别给红方、黑方发送对手 消息
	res, _ := json.Marshal(&PB.MatchAck{
		Id:             1202,
		ErrorCode:      "SUCCESS",
		EnemyName:      blackPlayer.NickName,
		EnemyAvatarUrl: blackPlayer.AvatarUrl,
		Color:          1,
	})
	redPlayer.SendChannel <- res
	res, _ = json.Marshal(&PB.MatchAck{
		Id:             1202,
		ErrorCode:      "SUCCESS",
		EnemyName:      redPlayer.NickName,
		EnemyAvatarUrl: redPlayer.AvatarUrl,
		Color:          2,
	})
	blackPlayer.SendChannel <- res

	// 加入房间管理
	RoomOpenIdMap[room.RedId] = room
	RoomOpenIdMap[room.BlackId] = room
	room.MsgChannel = make(chan PB.ChessStep, 0)

	go func() {
		for {
			select {
			case step := <-room.MsgChannel:
				{
					if !isPosValid(room, step.Pos) {
						log.Printf("%v step to an wrong pos (%v)\n", step.Color, step.Pos)
						continue
					}
					log.Println(step.Color, " step to ", step.Pos)
					room.GobangInfo[step.Pos.X][step.Pos.Y] = step.Color
					room.ChessStepList = append(room.ChessStepList, step)
					// 当前位置没人下过则创建一步棋
					if temp := room.GoBangTemp[step.Pos.X][step.Pos.Y]; temp == nil {
						temp = &Piece{
							horizontal: 0,
							vertical:   0,
							lOblique:   0,
							rOblique:   0,
						}
						if step.Color == PB.ColorTypeRed {
							temp.openId = room.RedId
							room.TurnId = room.BlackId
						} else if step.Color == PB.ColorTypeBlack {
							temp.openId = room.BlackId
							room.TurnId = room.RedId
						}
						room.GoBangTemp[step.Pos.X][step.Pos.Y] = temp
					}

					// 更新棋盘数据
					updateGobangTemp(room, step.Pos.X, step.Pos.Y)
					res, _ := json.Marshal(&PB.ChessStepAck{
						Id:        1302,
						ErrorCode: "SUCCESS",
						Steps:     room.ChessStepList,
					})
					UserSystem.GetPlayerByOpenId(room.RedId).SendChannel <- res
					UserSystem.GetPlayerByOpenId(room.BlackId).SendChannel <- res

					// 判断胜负
					winId, isWin := WhoWin(room)
					if isWin {
						wRes, _ := json.Marshal(&PB.GameResultAck{
							Id:         1402,
							ErrorCode:  "SUCCESS",
							GameResult: "WIN",
						})
						lRes, _ := json.Marshal(&PB.GameResultAck{
							Id:         1402,
							ErrorCode:  "SUCCESS",
							GameResult: "LOSE",
						})

						if winId == redPlayer.OpenId {
							redPlayer.SendChannel <- wRes
							blackPlayer.SendChannel <- lRes
						} else if winId == blackPlayer.OpenId {
							blackPlayer.SendChannel <- wRes
							redPlayer.SendChannel <- lRes
						}
						break
					}
				}
			}

		}
		delete(RoomOpenIdMap, room.RedId)
		delete(RoomOpenIdMap, room.BlackId)
		log.Println("room destroyed!")
	}()
	log.Println("create room success, RedId:", room.RedId, " BlackId:", room.BlackId)
	return nil
}

// WhoWin 谁赢了 是否有输赢
func WhoWin(room *Room) (string, bool) {
	for _, row := range room.GoBangTemp {
		for _, col := range row {
			if col == nil {
				continue
			}
			if col.horizontal >= 5 || col.vertical >= 5 || col.lOblique >= 5 || col.rOblique >= 5 {
				return col.openId, true
			}
		}
	}
	return "", false
}

// 位置是否合法
func isPosValid(room *Room, pos PB.Pos) bool {
	res := true
	if pos.X < 0 || pos.X >= 15 {
		res = false
	} else if pos.Y < 0 || pos.Y >= 15 {
		res = false
	} else if room.GobangInfo[pos.X][pos.Y] != 0 {
		res = false
	}

	return res
}

// 更新对局辅助信息
func updateGobangTemp(room *Room, x, y int32) (isWin bool) {
	curPiece := room.GoBangTemp[x][y]
	isWin = false

	// 寻找当前位置的四个方向的所有piece 并更新每个方向的连珠数
	arrPiece := make([]*Piece, 0)
	// 横
	{
		arrPiece = append(arrPiece, curPiece)
		// 1 相同的棋子全部取出来
		for i := y - 1; i >= 0; i-- {
			if t := room.GoBangTemp[x][i]; t != nil && t.openId == curPiece.openId {
				arrPiece = append(arrPiece, t)
			} else {
				break
			}
		}
		for i := y + 1; i < 15; i++ {
			if t := room.GoBangTemp[x][i]; t != nil && t.openId == curPiece.openId {
				arrPiece = append(arrPiece, t)
			} else {
				break
			}
		}

		// 2 更新每个piece
		for _, piece := range arrPiece {
			piece.horizontal = len(arrPiece)
		}

	}

	if len(arrPiece) >= 5 {
		isWin = true
	}

	// 竖
	{
		arrPiece = make([]*Piece, 0)
		arrPiece = append(arrPiece, curPiece)
		// 1 相同的棋子全部取出来
		for i := x - 1; i >= 0; i-- {
			if t := room.GoBangTemp[i][y]; t != nil && t.openId == curPiece.openId {
				arrPiece = append(arrPiece, t)
			} else {
				break
			}
		}
		for i := x + 1; i < 15; i++ {
			if t := room.GoBangTemp[i][y]; t != nil && t.openId == curPiece.openId {
				arrPiece = append(arrPiece, t)
			} else {
				break
			}
		}
		// 2 更新每个piece
		for _, piece := range arrPiece {
			piece.vertical = len(arrPiece)
		}
	}
	if len(arrPiece) >= 5 {
		isWin = true
	}
	// 左斜
	{
		arrPiece = make([]*Piece, 0)
		arrPiece = append(arrPiece, curPiece)
		// 1 相同的棋子全部取出来
		for i, j := x-1, y-1; i >= 0 && j >= 0; i, j = i-1, j-1 {
			if t := room.GoBangTemp[i][j]; t != nil && t.openId == curPiece.openId {
				arrPiece = append(arrPiece, t)
			} else {
				break
			}
		}
		for i, j := x+1, y+1; i < 15 && j < 15; i, j = i+1, j+1 {
			if t := room.GoBangTemp[i][j]; t != nil && t.openId == curPiece.openId {
				arrPiece = append(arrPiece, t)
			} else {
				break
			}
		}
		// 2 更新每个piece
		for _, piece := range arrPiece {
			piece.lOblique = len(arrPiece)
		}
	}
	if len(arrPiece) >= 5 {
		isWin = true
	}
	// 右斜
	{
		arrPiece = make([]*Piece, 0)
		arrPiece = append(arrPiece, curPiece)
		// 1 相同的棋子全部取出来
		for i, j := x+1, y-1; i < 15 && j >= 0; i, j = i+1, j-1 {
			if t := room.GoBangTemp[i][j]; t != nil && t.openId == curPiece.openId {
				arrPiece = append(arrPiece, t)
			} else {
				break
			}
		}
		for i, j := x-1, y+1; i >= 0 && j < 15; i, j = i-1, j+1 {
			if t := room.GoBangTemp[i][j]; t != nil && t.openId == curPiece.openId {
				arrPiece = append(arrPiece, t)
			} else {
				break
			}
		}
		// 2 更新每个piece
		for _, piece := range arrPiece {
			piece.rOblique = len(arrPiece)
		}
	}
	if len(arrPiece) >= 5 {
		isWin = true
	}

	return isWin
}
