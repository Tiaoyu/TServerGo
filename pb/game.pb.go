// Code generated by protoc-gen-go. DO NOT EDIT.
// source: game.proto

/*
Package gamepb is a generated protocol buffer package.

It is generated from these files:
	game.proto
	error.proto
	enum.proto

It has these top-level messages:
	Error
	Point
	ChessStep
	GobangInfo
	C2SPing
	S2CPing
	C2SLogin
	S2CLogin
	C2SMatch
	S2CMatch
	C2SStep
	S2CStep
	S2CGameResult
*/
package gamepb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Error struct {
	ErrorCode ErrorType `protobuf:"varint,1,opt,name=errorCode,enum=gamepb.ErrorType" json:"errorCode,omitempty"`
	ErrorMsg  string    `protobuf:"bytes,2,opt,name=errorMsg" json:"errorMsg,omitempty"`
}

func (m *Error) Reset()                    { *m = Error{} }
func (m *Error) String() string            { return proto.CompactTextString(m) }
func (*Error) ProtoMessage()               {}
func (*Error) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Error) GetErrorCode() ErrorType {
	if m != nil {
		return m.ErrorCode
	}
	return ErrorType_SUCCESS
}

func (m *Error) GetErrorMsg() string {
	if m != nil {
		return m.ErrorMsg
	}
	return ""
}

// 位置 及 阵营
type Point struct {
	X    int32 `protobuf:"varint,1,opt,name=x" json:"x,omitempty"`
	Y    int32 `protobuf:"varint,2,opt,name=y" json:"y,omitempty"`
	Camp int32 `protobuf:"varint,3,opt,name=Camp" json:"Camp,omitempty"`
}

func (m *Point) Reset()                    { *m = Point{} }
func (m *Point) String() string            { return proto.CompactTextString(m) }
func (*Point) ProtoMessage()               {}
func (*Point) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Point) GetX() int32 {
	if m != nil {
		return m.X
	}
	return 0
}

func (m *Point) GetY() int32 {
	if m != nil {
		return m.Y
	}
	return 0
}

func (m *Point) GetCamp() int32 {
	if m != nil {
		return m.Camp
	}
	return 0
}

// 位置 阵营 及 次序
type ChessStep struct {
	Point *Point `protobuf:"bytes,1,opt,name=point" json:"point,omitempty"`
	Index int32  `protobuf:"varint,2,opt,name=index" json:"index,omitempty"`
}

func (m *ChessStep) Reset()                    { *m = ChessStep{} }
func (m *ChessStep) String() string            { return proto.CompactTextString(m) }
func (*ChessStep) ProtoMessage()               {}
func (*ChessStep) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *ChessStep) GetPoint() *Point {
	if m != nil {
		return m.Point
	}
	return nil
}

func (m *ChessStep) GetIndex() int32 {
	if m != nil {
		return m.Index
	}
	return 0
}

// 当前对局
type GobangInfo struct {
	ChessSteps []*ChessStep `protobuf:"bytes,1,rep,name=chessSteps" json:"chessSteps,omitempty"`
}

func (m *GobangInfo) Reset()                    { *m = GobangInfo{} }
func (m *GobangInfo) String() string            { return proto.CompactTextString(m) }
func (*GobangInfo) ProtoMessage()               {}
func (*GobangInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *GobangInfo) GetChessSteps() []*ChessStep {
	if m != nil {
		return m.ChessSteps
	}
	return nil
}

// ============PROTOCOL===========//
// ping
type C2SPing struct {
	Timestamp int64 `protobuf:"varint,1,opt,name=timestamp" json:"timestamp,omitempty"`
}

func (m *C2SPing) Reset()                    { *m = C2SPing{} }
func (m *C2SPing) String() string            { return proto.CompactTextString(m) }
func (*C2SPing) ProtoMessage()               {}
func (*C2SPing) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *C2SPing) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

type S2CPing struct {
	Error     *Error `protobuf:"bytes,1,opt,name=error" json:"error,omitempty"`
	Timestamp int64  `protobuf:"varint,2,opt,name=timestamp" json:"timestamp,omitempty"`
}

func (m *S2CPing) Reset()                    { *m = S2CPing{} }
func (m *S2CPing) String() string            { return proto.CompactTextString(m) }
func (*S2CPing) ProtoMessage()               {}
func (*S2CPing) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *S2CPing) GetError() *Error {
	if m != nil {
		return m.Error
	}
	return nil
}

func (m *S2CPing) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

// 登陆
type C2SLogin struct {
	NickName  string `protobuf:"bytes,1,opt,name=NickName" json:"NickName,omitempty"`
	AvatarUrl string `protobuf:"bytes,2,opt,name=AvatarUrl" json:"AvatarUrl,omitempty"`
}

func (m *C2SLogin) Reset()                    { *m = C2SLogin{} }
func (m *C2SLogin) String() string            { return proto.CompactTextString(m) }
func (*C2SLogin) ProtoMessage()               {}
func (*C2SLogin) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *C2SLogin) GetNickName() string {
	if m != nil {
		return m.NickName
	}
	return ""
}

func (m *C2SLogin) GetAvatarUrl() string {
	if m != nil {
		return m.AvatarUrl
	}
	return ""
}

type S2CLogin struct {
	ErrorCode string `protobuf:"bytes,1,opt,name=ErrorCode" json:"ErrorCode,omitempty"`
}

func (m *S2CLogin) Reset()                    { *m = S2CLogin{} }
func (m *S2CLogin) String() string            { return proto.CompactTextString(m) }
func (*S2CLogin) ProtoMessage()               {}
func (*S2CLogin) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *S2CLogin) GetErrorCode() string {
	if m != nil {
		return m.ErrorCode
	}
	return ""
}

// 匹配
type C2SMatch struct {
	MatchType MatchType `protobuf:"varint,1,opt,name=match_type,json=matchType,enum=gamepb.MatchType" json:"match_type,omitempty"`
}

func (m *C2SMatch) Reset()                    { *m = C2SMatch{} }
func (m *C2SMatch) String() string            { return proto.CompactTextString(m) }
func (*C2SMatch) ProtoMessage()               {}
func (*C2SMatch) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *C2SMatch) GetMatchType() MatchType {
	if m != nil {
		return m.MatchType
	}
	return MatchType_MatchTypeDefault
}

type S2CMatch struct {
	EnemyName      string      `protobuf:"bytes,1,opt,name=enemyName" json:"enemyName,omitempty"`
	EnemyAvatarUrl string      `protobuf:"bytes,2,opt,name=enemyAvatarUrl" json:"enemyAvatarUrl,omitempty"`
	Color          ColorType   `protobuf:"varint,3,opt,name=color,enum=gamepb.ColorType" json:"color,omitempty"`
	Result         MatchResult `protobuf:"varint,4,opt,name=result,enum=gamepb.MatchResult" json:"result,omitempty"`
}

func (m *S2CMatch) Reset()                    { *m = S2CMatch{} }
func (m *S2CMatch) String() string            { return proto.CompactTextString(m) }
func (*S2CMatch) ProtoMessage()               {}
func (*S2CMatch) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *S2CMatch) GetEnemyName() string {
	if m != nil {
		return m.EnemyName
	}
	return ""
}

func (m *S2CMatch) GetEnemyAvatarUrl() string {
	if m != nil {
		return m.EnemyAvatarUrl
	}
	return ""
}

func (m *S2CMatch) GetColor() ColorType {
	if m != nil {
		return m.Color
	}
	return ColorType_ColorTypeDefault
}

func (m *S2CMatch) GetResult() MatchResult {
	if m != nil {
		return m.Result
	}
	return MatchResult_MatResultDefault
}

// 走一步
type C2SStep struct {
	Point *Point `protobuf:"bytes,1,opt,name=point" json:"point,omitempty"`
}

func (m *C2SStep) Reset()                    { *m = C2SStep{} }
func (m *C2SStep) String() string            { return proto.CompactTextString(m) }
func (*C2SStep) ProtoMessage()               {}
func (*C2SStep) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *C2SStep) GetPoint() *Point {
	if m != nil {
		return m.Point
	}
	return nil
}

type S2CStep struct {
	Error      *Error      `protobuf:"bytes,1,opt,name=error" json:"error,omitempty"`
	GobangInfo *GobangInfo `protobuf:"bytes,2,opt,name=gobangInfo" json:"gobangInfo,omitempty"`
}

func (m *S2CStep) Reset()                    { *m = S2CStep{} }
func (m *S2CStep) String() string            { return proto.CompactTextString(m) }
func (*S2CStep) ProtoMessage()               {}
func (*S2CStep) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

func (m *S2CStep) GetError() *Error {
	if m != nil {
		return m.Error
	}
	return nil
}

func (m *S2CStep) GetGobangInfo() *GobangInfo {
	if m != nil {
		return m.GobangInfo
	}
	return nil
}

// 对局结果
type S2CGameResult struct {
	Result GameResult `protobuf:"varint,1,opt,name=result,enum=gamepb.GameResult" json:"result,omitempty"`
}

func (m *S2CGameResult) Reset()                    { *m = S2CGameResult{} }
func (m *S2CGameResult) String() string            { return proto.CompactTextString(m) }
func (*S2CGameResult) ProtoMessage()               {}
func (*S2CGameResult) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

func (m *S2CGameResult) GetResult() GameResult {
	if m != nil {
		return m.Result
	}
	return GameResult_GameResultDefault
}

func init() {
	proto.RegisterType((*Error)(nil), "gamepb.Error")
	proto.RegisterType((*Point)(nil), "gamepb.Point")
	proto.RegisterType((*ChessStep)(nil), "gamepb.ChessStep")
	proto.RegisterType((*GobangInfo)(nil), "gamepb.GobangInfo")
	proto.RegisterType((*C2SPing)(nil), "gamepb.C2SPing")
	proto.RegisterType((*S2CPing)(nil), "gamepb.S2CPing")
	proto.RegisterType((*C2SLogin)(nil), "gamepb.C2SLogin")
	proto.RegisterType((*S2CLogin)(nil), "gamepb.S2CLogin")
	proto.RegisterType((*C2SMatch)(nil), "gamepb.C2SMatch")
	proto.RegisterType((*S2CMatch)(nil), "gamepb.S2CMatch")
	proto.RegisterType((*C2SStep)(nil), "gamepb.C2SStep")
	proto.RegisterType((*S2CStep)(nil), "gamepb.S2CStep")
	proto.RegisterType((*S2CGameResult)(nil), "gamepb.S2CGameResult")
}

func init() { proto.RegisterFile("game.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 509 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x53, 0x4d, 0x6f, 0xd3, 0x40,
	0x10, 0x95, 0x9b, 0x3a, 0x8d, 0x27, 0xb4, 0x12, 0x0b, 0x42, 0x56, 0x94, 0x43, 0xb4, 0x95, 0x68,
	0x54, 0x24, 0x07, 0xcc, 0x81, 0x03, 0x48, 0x88, 0x9a, 0x52, 0x21, 0xb5, 0x55, 0xb4, 0x2e, 0x17,
	0x2e, 0xc8, 0x76, 0x17, 0xd7, 0x22, 0xeb, 0xb5, 0xec, 0x0d, 0x8a, 0xff, 0x10, 0xbf, 0x13, 0xed,
	0xac, 0x3f, 0x92, 0x70, 0xa0, 0x27, 0xef, 0xbc, 0x79, 0x7e, 0xf3, 0x76, 0x66, 0x16, 0x20, 0x8d,
	0x04, 0xf7, 0x8a, 0x52, 0x2a, 0x49, 0x86, 0xfa, 0x5c, 0xc4, 0x13, 0xe0, 0xf9, 0x5a, 0x18, 0x6c,
	0x32, 0xe6, 0x65, 0x29, 0x4b, 0x13, 0xd0, 0x3b, 0xb0, 0x2f, 0x75, 0x48, 0x16, 0xe0, 0x20, 0x1e,
	0xc8, 0x7b, 0xee, 0x5a, 0x33, 0x6b, 0x7e, 0xe2, 0x3f, 0xf5, 0xcc, 0xdf, 0x1e, 0x32, 0xee, 0xea,
	0x82, 0xb3, 0x9e, 0x43, 0x26, 0x30, 0xc2, 0xe0, 0xa6, 0x4a, 0xdd, 0x83, 0x99, 0x35, 0x77, 0x58,
	0x17, 0xd3, 0x77, 0x60, 0x2f, 0x65, 0x96, 0x2b, 0xf2, 0x04, 0xac, 0x0d, 0xaa, 0xd9, 0xcc, 0xda,
	0xe8, 0xa8, 0x46, 0xae, 0xcd, 0xac, 0x9a, 0x10, 0x38, 0x0c, 0x22, 0x51, 0xb8, 0x03, 0x04, 0xf0,
	0x4c, 0xbf, 0x80, 0x13, 0x3c, 0xf0, 0xaa, 0x0a, 0x15, 0x2f, 0xc8, 0x29, 0xd8, 0x85, 0x56, 0x41,
	0x81, 0xb1, 0x7f, 0xdc, 0xda, 0x41, 0x69, 0x66, 0x72, 0xe4, 0x39, 0xd8, 0x59, 0x7e, 0xcf, 0x37,
	0x8d, 0xae, 0x09, 0xe8, 0x47, 0x80, 0x2b, 0x19, 0x47, 0x79, 0xfa, 0x35, 0xff, 0x29, 0xc9, 0x1b,
	0x80, 0xa4, 0x55, 0xad, 0x5c, 0x6b, 0x36, 0x98, 0x8f, 0xfb, 0xcb, 0x75, 0xf5, 0xd8, 0x16, 0x89,
	0x9e, 0xc1, 0x51, 0xe0, 0x87, 0xcb, 0x2c, 0x4f, 0xc9, 0x14, 0x1c, 0x95, 0x09, 0x5e, 0x29, 0x6d,
	0x56, 0x5b, 0x19, 0xb0, 0x1e, 0xa0, 0xd7, 0x70, 0x14, 0xfa, 0x01, 0x12, 0x4f, 0xc1, 0xc6, 0x0e,
	0xec, 0xfb, 0xc5, 0xf6, 0x31, 0x93, 0xdb, 0x55, 0x3b, 0xd8, 0x57, 0xfb, 0x0c, 0xa3, 0xc0, 0x0f,
	0xaf, 0x65, 0x9a, 0xe5, 0xba, 0xc1, 0xb7, 0x59, 0xf2, 0xeb, 0x36, 0x12, 0x66, 0x20, 0x0e, 0xeb,
	0x62, 0xad, 0xf2, 0xe9, 0x77, 0xa4, 0xa2, 0xf2, 0x5b, 0xb9, 0x6a, 0xba, 0xdf, 0x03, 0x74, 0x0e,
	0xa3, 0xd0, 0x0f, 0x8c, 0xca, 0x14, 0x9c, 0xcb, 0x9d, 0xb9, 0x3a, 0xac, 0x07, 0xe8, 0x07, 0xac,
	0x77, 0x13, 0xa9, 0xe4, 0x81, 0xbc, 0x06, 0x10, 0xfa, 0xf0, 0x43, 0xd5, 0xc5, 0x3f, 0x2b, 0x80,
	0x14, 0xb3, 0x02, 0xa2, 0x3d, 0xd2, 0x3f, 0x16, 0x16, 0x32, 0xbf, 0x4f, 0xc1, 0xe1, 0x39, 0x17,
	0xf5, 0x96, 0xdf, 0x1e, 0x20, 0x2f, 0xe1, 0x04, 0x83, 0x7d, 0xd7, 0x7b, 0x28, 0x39, 0x03, 0x3b,
	0x91, 0x2b, 0x59, 0xe2, 0x56, 0x6c, 0xd5, 0x0f, 0x34, 0x88, 0xf5, 0x4d, 0x9e, 0xbc, 0x82, 0x61,
	0xc9, 0xab, 0xf5, 0x4a, 0xb9, 0x87, 0xc8, 0x7c, 0xb6, 0xe3, 0x94, 0x61, 0x8a, 0x35, 0x14, 0xea,
	0xe1, 0x34, 0x1f, 0xbd, 0x54, 0x34, 0xc6, 0xa1, 0xb6, 0xfc, 0xff, 0x0f, 0xd5, 0x07, 0x48, 0xbb,
	0x75, 0xc3, 0x9b, 0x8d, 0x7d, 0xd2, 0x32, 0xfb, 0x45, 0x64, 0x5b, 0x2c, 0xfa, 0x1e, 0x8e, 0x43,
	0x3f, 0xb8, 0x8a, 0x04, 0x37, 0x66, 0xc9, 0x79, 0x77, 0x23, 0xd3, 0xfb, 0x5e, 0xa0, 0xe3, 0xb4,
	0x17, 0xba, 0x38, 0x87, 0xe6, 0x65, 0x5f, 0x34, 0xdf, 0xa5, 0xf5, 0xfd, 0x85, 0xca, 0x22, 0x59,
	0xaf, 0xeb, 0xb5, 0x97, 0x48, 0xb1, 0x28, 0xe2, 0x85, 0xc9, 0xc4, 0x43, 0x7c, 0xe9, 0x6f, 0xff,
	0x06, 0x00, 0x00, 0xff, 0xff, 0x62, 0x74, 0x55, 0xb0, 0x18, 0x04, 0x00, 0x00,
}
