// Code generated by protoc-gen-go. DO NOT EDIT.
// source: enum.proto

package gamepb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// 协议类型
type ProtocolType int32

const (
	ProtocolType_EC2SPing       ProtocolType = 0
	ProtocolType_ES2CPing       ProtocolType = 1
	ProtocolType_EC2SLogin      ProtocolType = 2
	ProtocolType_ES2CLogin      ProtocolType = 3
	ProtocolType_EC2SMatch      ProtocolType = 4
	ProtocolType_ES2CMatch      ProtocolType = 5
	ProtocolType_EC2SStep       ProtocolType = 6
	ProtocolType_ES2CStep       ProtocolType = 7
	ProtocolType_ES2CGameResult ProtocolType = 8
)

var ProtocolType_name = map[int32]string{
	0: "EC2SPing",
	1: "ES2CPing",
	2: "EC2SLogin",
	3: "ES2CLogin",
	4: "EC2SMatch",
	5: "ES2CMatch",
	6: "EC2SStep",
	7: "ES2CStep",
	8: "ES2CGameResult",
}
var ProtocolType_value = map[string]int32{
	"EC2SPing":       0,
	"ES2CPing":       1,
	"EC2SLogin":      2,
	"ES2CLogin":      3,
	"EC2SMatch":      4,
	"ES2CMatch":      5,
	"EC2SStep":       6,
	"ES2CStep":       7,
	"ES2CGameResult": 8,
}

func (x ProtocolType) String() string {
	return proto.EnumName(ProtocolType_name, int32(x))
}
func (ProtocolType) EnumDescriptor() ([]byte, []int) { return fileDescriptor2, []int{0} }

// 匹配类型
type MatchType int32

const (
	MatchType_MatchTypeDefault MatchType = 0
	MatchType_MatchTypeMatch   MatchType = 1
	MatchType_MatchTypeCancel  MatchType = 2
)

var MatchType_name = map[int32]string{
	0: "MatchTypeDefault",
	1: "MatchTypeMatch",
	2: "MatchTypeCancel",
}
var MatchType_value = map[string]int32{
	"MatchTypeDefault": 0,
	"MatchTypeMatch":   1,
	"MatchTypeCancel":  2,
}

func (x MatchType) String() string {
	return proto.EnumName(MatchType_name, int32(x))
}
func (MatchType) EnumDescriptor() ([]byte, []int) { return fileDescriptor2, []int{1} }

// 匹配结果
type MatchResult int32

const (
	MatchResult_MatResultDefault  MatchResult = 0
	MatchResult_MatResultSuccess  MatchResult = 1
	MatchResult_MatResultFailed   MatchResult = 2
	MatchResult_MatResultMatching MatchResult = 3
	MatchResult_MatResultCancel   MatchResult = 4
)

var MatchResult_name = map[int32]string{
	0: "MatResultDefault",
	1: "MatResultSuccess",
	2: "MatResultFailed",
	3: "MatResultMatching",
	4: "MatResultCancel",
}
var MatchResult_value = map[string]int32{
	"MatResultDefault":  0,
	"MatResultSuccess":  1,
	"MatResultFailed":   2,
	"MatResultMatching": 3,
	"MatResultCancel":   4,
}

func (x MatchResult) String() string {
	return proto.EnumName(MatchResult_name, int32(x))
}
func (MatchResult) EnumDescriptor() ([]byte, []int) { return fileDescriptor2, []int{2} }

// 对局双方类型
type ColorType int32

const (
	ColorType_ColorTypeDefault ColorType = 0
	ColorType_ColorTypeRed     ColorType = 1
	ColorType_ColorTypeBlack   ColorType = 2
)

var ColorType_name = map[int32]string{
	0: "ColorTypeDefault",
	1: "ColorTypeRed",
	2: "ColorTypeBlack",
}
var ColorType_value = map[string]int32{
	"ColorTypeDefault": 0,
	"ColorTypeRed":     1,
	"ColorTypeBlack":   2,
}

func (x ColorType) String() string {
	return proto.EnumName(ColorType_name, int32(x))
}
func (ColorType) EnumDescriptor() ([]byte, []int) { return fileDescriptor2, []int{3} }

// 对局结果
type GameResult int32

const (
	GameResult_GameResultDefault GameResult = 0
	GameResult_GameResultWin     GameResult = 1
	GameResult_GameResultFail    GameResult = 2
	GameResult_GameResultDraw    GameResult = 3
)

var GameResult_name = map[int32]string{
	0: "GameResultDefault",
	1: "GameResultWin",
	2: "GameResultFail",
	3: "GameResultDraw",
}
var GameResult_value = map[string]int32{
	"GameResultDefault": 0,
	"GameResultWin":     1,
	"GameResultFail":    2,
	"GameResultDraw":    3,
}

func (x GameResult) String() string {
	return proto.EnumName(GameResult_name, int32(x))
}
func (GameResult) EnumDescriptor() ([]byte, []int) { return fileDescriptor2, []int{4} }

func init() {
	proto.RegisterEnum("gamepb.ProtocolType", ProtocolType_name, ProtocolType_value)
	proto.RegisterEnum("gamepb.MatchType", MatchType_name, MatchType_value)
	proto.RegisterEnum("gamepb.MatchResult", MatchResult_name, MatchResult_value)
	proto.RegisterEnum("gamepb.ColorType", ColorType_name, ColorType_value)
	proto.RegisterEnum("gamepb.GameResult", GameResult_name, GameResult_value)
}

func init() { proto.RegisterFile("enum.proto", fileDescriptor2) }

var fileDescriptor2 = []byte{
	// 323 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x92, 0xed, 0x4e, 0xf2, 0x30,
	0x14, 0xc7, 0xd9, 0xe0, 0xe1, 0x81, 0x23, 0xe8, 0xa1, 0xbe, 0x5c, 0x44, 0x3f, 0x40, 0x82, 0x77,
	0xc0, 0x50, 0x12, 0xa3, 0x09, 0x61, 0x26, 0x26, 0x7e, 0x30, 0x29, 0xa5, 0xe2, 0x62, 0xb7, 0x2e,
	0x7b, 0x89, 0x99, 0x17, 0xe2, 0xf5, 0x9a, 0xb6, 0x5b, 0x2b, 0x9f, 0x9a, 0xdf, 0xef, 0xb4, 0xff,
	0xd3, 0xd3, 0x14, 0x40, 0x64, 0x75, 0x3a, 0xcf, 0x0b, 0x55, 0x29, 0x32, 0x3c, 0xb2, 0x54, 0xe4,
	0x7b, 0xfa, 0x13, 0xc0, 0x64, 0xab, 0x0d, 0x57, 0xf2, 0xb9, 0xc9, 0x05, 0x99, 0xc0, 0xe8, 0x2e,
	0x5a, 0xc6, 0xdb, 0x24, 0x3b, 0x62, 0xcf, 0x50, 0xbc, 0x8c, 0x0c, 0x05, 0x64, 0x0a, 0x63, 0x5d,
	0x7b, 0x54, 0xc7, 0x24, 0xc3, 0xd0, 0x60, 0xbc, 0x8c, 0x2c, 0xf6, 0xbb, 0xea, 0x13, 0xab, 0xf8,
	0x07, 0x0e, 0xba, 0xaa, 0xc5, 0x7f, 0x5d, 0x6e, 0x5c, 0x89, 0x1c, 0x87, 0x5d, 0xae, 0xa1, 0xff,
	0x84, 0xc0, 0xb9, 0xa6, 0x0d, 0x4b, 0xc5, 0x4e, 0x94, 0xb5, 0xac, 0x70, 0x44, 0x1f, 0x60, 0x6c,
	0x8e, 0x9a, 0x4b, 0x5d, 0x01, 0x3a, 0x58, 0x8b, 0x77, 0xa6, 0xb7, 0xf4, 0xf4, 0x31, 0x67, 0x6d,
	0x9b, 0x80, 0x5c, 0xc2, 0x85, 0x73, 0x11, 0xcb, 0xb8, 0x90, 0x18, 0xd2, 0x6f, 0x38, 0x33, 0xd2,
	0x86, 0xb7, 0x69, 0x16, 0x7c, 0xda, 0x5f, 0x1b, 0xd7, 0x9c, 0x8b, 0xb2, 0x74, 0x79, 0xd6, 0xde,
	0xb3, 0x44, 0x8a, 0x03, 0x86, 0xe4, 0x1a, 0x66, 0x4e, 0x9a, 0x60, 0xfd, 0x3c, 0xfd, 0x93, 0xbd,
	0x6d, 0xef, 0x01, 0xdd, 0xc0, 0x38, 0x52, 0x52, 0x15, 0xdd, 0x1c, 0x0e, 0x7c, 0x67, 0x84, 0x89,
	0xb3, 0x3b, 0x71, 0xc0, 0x40, 0x4f, 0xe6, 0xcc, 0x4a, 0x32, 0xfe, 0x89, 0x21, 0x7d, 0x03, 0xf0,
	0x0f, 0xa4, 0xaf, 0xe0, 0xc9, 0x47, 0xcd, 0x60, 0xea, 0xf5, 0x4b, 0x92, 0xd9, 0x2c, 0xaf, 0xf4,
	0x08, 0x18, 0x9e, 0xba, 0x75, 0xc1, 0xbe, 0xb0, 0xbf, 0xa2, 0xd0, 0xfe, 0x89, 0x55, 0xbb, 0x6e,
	0x83, 0xd7, 0x9b, 0x2a, 0x61, 0xaa, 0xa9, 0x9b, 0x7a, 0xce, 0x55, 0xba, 0xc8, 0xf7, 0x0b, 0x5b,
	0xd9, 0x0f, 0xcd, 0x27, 0xba, 0xfd, 0x0d, 0x00, 0x00, 0xff, 0xff, 0xf1, 0xc2, 0xe1, 0xaa, 0x52,
	0x02, 0x00, 0x00,
}
