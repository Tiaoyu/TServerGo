//角色匹配后加入匹配池子
//匹配线程每秒从池子中取若干角色
//有足够角色后进行创建房间
//创建完房间后开始进入房间线程

package RoomSystem

import (
	"TServerGo/TServer/NotifySystem"
	"TServerGo/TServer/UserSystem"
	"TServerGo/TServer/dbproxy"
	"TServerGo/TServer/pbsend"
	gamepb "TServerGo/pb"
	"encoding/json"
	"log"
	"sync"
	"time"

	"xorm.io/xorm"
)

type Room struct {
	RedId         string
	BlackId       string
	CreateTime    int64
	ChessStepList []*gamepb.ChessStep
	GobangInfo    [15][15]int32
	TurnId        string // 当前手
	MsgChannel    chan *gamepb.ChessStep

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
	RoomOpenIdMap = new(sync.Map)
)

func init() {
	NotifySystem.NotifyRegister(NotifySystem.NotifyTypeRoleLoginIn, PlayerLogin)
}

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
	res := pbsend.SendMsg(&gamepb.S2CMatch{
		EnemyName:      blackPlayer.NickName,
		EnemyAvatarUrl: blackPlayer.AvatarUrl,
		Color:          gamepb.ColorType_ColorTypeBlack,
		Result:         gamepb.MatchResult_MatResultSuccess,
	}, gamepb.ProtocolType_ES2CMatch)
	redPlayer.Sess.SendChannel <- res
	res = pbsend.SendMsg(&gamepb.S2CMatch{
		EnemyName:      redPlayer.NickName,
		EnemyAvatarUrl: redPlayer.AvatarUrl,
		Color:          gamepb.ColorType_ColorTypeBlack,
		Result:         gamepb.MatchResult_MatResultSuccess,
	}, gamepb.ProtocolType_ES2CMatch)
	blackPlayer.Sess.SendChannel <- res

	// 加入房间管理
	RoomOpenIdMap.Store(room.RedId, room)
	RoomOpenIdMap.Store(room.BlackId, room)
	room.MsgChannel = make(chan *gamepb.ChessStep, 0)
	room.CreateTime = time.Now().Unix()
	go func() {
		d := time.Duration(time.Second * 2)
		t := time.NewTimer(d)
		defer t.Stop()
		var finished = false
		for {
			select {
			case <-t.C:
				t.Reset(time.Second * 2)
				if time.Now().Unix()-room.CreateTime > 3600 {
					log.Printf("Room is time out, so it will be destroyed! Names:%v-%v", redPlayer.NickName, blackPlayer.NickName)
					finished = true
				}
				break
			case step := <-room.MsgChannel:
				{
					if !isPosValid(room, step.Point) {
						log.Printf("%v step to an wrong pos (%v)\n", step.Point.Camp, step.Point)
						continue
					}
					log.Println(step.Point.Camp, " step to ", step.Point)
					room.GobangInfo[step.Point.X][step.Point.Y] = step.Point.Camp
					room.ChessStepList = append(room.ChessStepList, step)
					// 当前位置没人下过则创建一步棋
					if temp := room.GoBangTemp[step.Point.X][step.Point.Y]; temp == nil {
						temp = &Piece{
							horizontal: 0,
							vertical:   0,
							lOblique:   0,
							rOblique:   0,
						}
						if step.Point.Camp == int32(gamepb.ColorType_ColorTypeRed) {
							temp.openId = room.RedId
							room.TurnId = room.BlackId
						} else if step.Point.Camp == int32(gamepb.ColorType_ColorTypeBlack) {
							temp.openId = room.BlackId
							room.TurnId = room.RedId
						}
						room.GoBangTemp[step.Point.X][step.Point.Y] = temp
					}

					// 更新棋盘数据
					updateGobangTemp(room, step.Point.X, step.Point.Y)
					res := pbsend.SendMsg(&gamepb.S2CStep{
						Error: nil,
						GobangInfo: &gamepb.GobangInfo{
							ChessSteps: room.ChessStepList,
						},
					}, gamepb.ProtocolType_ES2CStep)
					if user := UserSystem.GetPlayerByOpenId(room.RedId); user != nil {
						user.Sess.SendChannel <- res
						redPlayer = user
					}
					if user := UserSystem.GetPlayerByOpenId(room.BlackId); user != nil {
						user.Sess.SendChannel <- res
						blackPlayer = user
					}

					// 判断胜负
					winId, isWin := WhoWin(room)
					if isWin {
						wRes := pbsend.SendMsg(&gamepb.S2CGameResult{
							Result: gamepb.GameResult_GameResultWin,
						}, gamepb.ProtocolType_ES2CGameResult)
						lRes := pbsend.SendMsg(&gamepb.S2CGameResult{
							Result: gamepb.GameResult_GameResultFail,
						}, gamepb.ProtocolType_ES2CGameResult)

						if winId == redPlayer.OpenId {
							redPlayer.Sess.SendChannel <- wRes
							blackPlayer.Sess.SendChannel <- lRes
						} else if winId == blackPlayer.OpenId {
							blackPlayer.Sess.SendChannel <- wRes
							redPlayer.Sess.SendChannel <- lRes
						}
						// 存储胜负数据
						dbproxy.Instance().Transaction(func(session *xorm.Session) (interface{}, error) {
							rUser := &dbproxy.User{}
							bUser := &dbproxy.User{}
							session.Where("open_id = ?", redPlayer.OpenId).Get(rUser)
							session.Where("open_id = ?", blackPlayer.OpenId).Get(bUser)
							if winId == redPlayer.OpenId {
								rUser.WinCount++
								rUser.Score++
								bUser.Score--
								bUser.FailedCount++
							} else if winId == blackPlayer.OpenId {
								bUser.WinCount++
								bUser.Score++
								rUser.Score--
								rUser.FailedCount++
							}
							session.Where("open_id = ?", redPlayer.OpenId).Update(rUser)
							session.Where("open_id = ?", blackPlayer.OpenId).Update(bUser)
							race := &dbproxy.Race{
								RedOpenId:   redPlayer.OpenId,
								BlackOpenId: blackPlayer.OpenId,
								WinnerId:    winId,
							}

							gobangInfo, _ := json.Marshal(room.ChessStepList)
							race.GobangInfo = string(gobangInfo)
							session.Insert(race)
							return nil, nil
						})
						finished = true
					}
				}
				break
			}
			if finished {
				RoomOpenIdMap.Delete(room.RedId)
				RoomOpenIdMap.Delete(room.BlackId)
				log.Printf("room destroyed! Red:%v Black:%v\n", redPlayer.NickName, blackPlayer.NickName)
				break
			}
		}
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
func isPosValid(room *Room, pos *gamepb.Point) bool {
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

func PlayerLogin(params ...interface{}) {
	param := params[0].(NotifySystem.NotifyRoleLoginParam)
	if room, ok := RoomOpenIdMap.Load(param.OpenId); ok {
		if user := UserSystem.GetPlayerByOpenId(param.OpenId); user != nil {
			res := pbsend.SendMsg(&gamepb.S2CStep{
				GobangInfo: &gamepb.GobangInfo{ChessSteps: room.(*Room).ChessStepList},
			}, gamepb.ProtocolType_ES2CStep)
			user.Sess.SendChannel <- res
		}
	}
}
