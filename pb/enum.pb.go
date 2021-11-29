// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.14.0
// source: enum.proto

package gamepb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

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
	ProtocolType_ES2CPushMsg    ProtocolType = 9
)

// Enum value maps for ProtocolType.
var (
	ProtocolType_name = map[int32]string{
		0: "EC2SPing",
		1: "ES2CPing",
		2: "EC2SLogin",
		3: "ES2CLogin",
		4: "EC2SMatch",
		5: "ES2CMatch",
		6: "EC2SStep",
		7: "ES2CStep",
		8: "ES2CGameResult",
		9: "ES2CPushMsg",
	}
	ProtocolType_value = map[string]int32{
		"EC2SPing":       0,
		"ES2CPing":       1,
		"EC2SLogin":      2,
		"ES2CLogin":      3,
		"EC2SMatch":      4,
		"ES2CMatch":      5,
		"EC2SStep":       6,
		"ES2CStep":       7,
		"ES2CGameResult": 8,
		"ES2CPushMsg":    9,
	}
)

func (x ProtocolType) Enum() *ProtocolType {
	p := new(ProtocolType)
	*p = x
	return p
}

func (x ProtocolType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ProtocolType) Descriptor() protoreflect.EnumDescriptor {
	return file_enum_proto_enumTypes[0].Descriptor()
}

func (ProtocolType) Type() protoreflect.EnumType {
	return &file_enum_proto_enumTypes[0]
}

func (x ProtocolType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ProtocolType.Descriptor instead.
func (ProtocolType) EnumDescriptor() ([]byte, []int) {
	return file_enum_proto_rawDescGZIP(), []int{0}
}

// 匹配类型
type MatchType int32

const (
	MatchType_MatchTypeDefault MatchType = 0
	MatchType_MatchTypeMatch   MatchType = 1
	MatchType_MatchTypeCancel  MatchType = 2
)

// Enum value maps for MatchType.
var (
	MatchType_name = map[int32]string{
		0: "MatchTypeDefault",
		1: "MatchTypeMatch",
		2: "MatchTypeCancel",
	}
	MatchType_value = map[string]int32{
		"MatchTypeDefault": 0,
		"MatchTypeMatch":   1,
		"MatchTypeCancel":  2,
	}
)

func (x MatchType) Enum() *MatchType {
	p := new(MatchType)
	*p = x
	return p
}

func (x MatchType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MatchType) Descriptor() protoreflect.EnumDescriptor {
	return file_enum_proto_enumTypes[1].Descriptor()
}

func (MatchType) Type() protoreflect.EnumType {
	return &file_enum_proto_enumTypes[1]
}

func (x MatchType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MatchType.Descriptor instead.
func (MatchType) EnumDescriptor() ([]byte, []int) {
	return file_enum_proto_rawDescGZIP(), []int{1}
}

// 匹配结果
type MatchResult int32

const (
	MatchResult_MatResultDefault  MatchResult = 0
	MatchResult_MatResultSuccess  MatchResult = 1
	MatchResult_MatResultFailed   MatchResult = 2
	MatchResult_MatResultMatching MatchResult = 3
	MatchResult_MatResultCancel   MatchResult = 4
)

// Enum value maps for MatchResult.
var (
	MatchResult_name = map[int32]string{
		0: "MatResultDefault",
		1: "MatResultSuccess",
		2: "MatResultFailed",
		3: "MatResultMatching",
		4: "MatResultCancel",
	}
	MatchResult_value = map[string]int32{
		"MatResultDefault":  0,
		"MatResultSuccess":  1,
		"MatResultFailed":   2,
		"MatResultMatching": 3,
		"MatResultCancel":   4,
	}
)

func (x MatchResult) Enum() *MatchResult {
	p := new(MatchResult)
	*p = x
	return p
}

func (x MatchResult) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MatchResult) Descriptor() protoreflect.EnumDescriptor {
	return file_enum_proto_enumTypes[2].Descriptor()
}

func (MatchResult) Type() protoreflect.EnumType {
	return &file_enum_proto_enumTypes[2]
}

func (x MatchResult) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MatchResult.Descriptor instead.
func (MatchResult) EnumDescriptor() ([]byte, []int) {
	return file_enum_proto_rawDescGZIP(), []int{2}
}

// 对局双方类型
type ColorType int32

const (
	ColorType_ColorTypeDefault ColorType = 0
	ColorType_ColorTypeRed     ColorType = 1 // 红方
	ColorType_ColorTypeBlack   ColorType = 2 // 黑方
)

// Enum value maps for ColorType.
var (
	ColorType_name = map[int32]string{
		0: "ColorTypeDefault",
		1: "ColorTypeRed",
		2: "ColorTypeBlack",
	}
	ColorType_value = map[string]int32{
		"ColorTypeDefault": 0,
		"ColorTypeRed":     1,
		"ColorTypeBlack":   2,
	}
)

func (x ColorType) Enum() *ColorType {
	p := new(ColorType)
	*p = x
	return p
}

func (x ColorType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ColorType) Descriptor() protoreflect.EnumDescriptor {
	return file_enum_proto_enumTypes[3].Descriptor()
}

func (ColorType) Type() protoreflect.EnumType {
	return &file_enum_proto_enumTypes[3]
}

func (x ColorType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ColorType.Descriptor instead.
func (ColorType) EnumDescriptor() ([]byte, []int) {
	return file_enum_proto_rawDescGZIP(), []int{3}
}

// 对局结果
type GameResult int32

const (
	GameResult_GameResultDefault GameResult = 0
	GameResult_GameResultWin     GameResult = 1 // 赢
	GameResult_GameResultFail    GameResult = 2 // 输
	GameResult_GameResultDraw    GameResult = 3 // 和
)

// Enum value maps for GameResult.
var (
	GameResult_name = map[int32]string{
		0: "GameResultDefault",
		1: "GameResultWin",
		2: "GameResultFail",
		3: "GameResultDraw",
	}
	GameResult_value = map[string]int32{
		"GameResultDefault": 0,
		"GameResultWin":     1,
		"GameResultFail":    2,
		"GameResultDraw":    3,
	}
)

func (x GameResult) Enum() *GameResult {
	p := new(GameResult)
	*p = x
	return p
}

func (x GameResult) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (GameResult) Descriptor() protoreflect.EnumDescriptor {
	return file_enum_proto_enumTypes[4].Descriptor()
}

func (GameResult) Type() protoreflect.EnumType {
	return &file_enum_proto_enumTypes[4]
}

func (x GameResult) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use GameResult.Descriptor instead.
func (GameResult) EnumDescriptor() ([]byte, []int) {
	return file_enum_proto_rawDescGZIP(), []int{4}
}

var File_enum_proto protoreflect.FileDescriptor

var file_enum_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x65, 0x6e, 0x75, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x67, 0x61,
	0x6d, 0x65, 0x70, 0x62, 0x2a, 0xa7, 0x01, 0x0a, 0x0c, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f,
	0x6c, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0c, 0x0a, 0x08, 0x45, 0x43, 0x32, 0x53, 0x50, 0x69, 0x6e,
	0x67, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x45, 0x53, 0x32, 0x43, 0x50, 0x69, 0x6e, 0x67, 0x10,
	0x01, 0x12, 0x0d, 0x0a, 0x09, 0x45, 0x43, 0x32, 0x53, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x10, 0x02,
	0x12, 0x0d, 0x0a, 0x09, 0x45, 0x53, 0x32, 0x43, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x10, 0x03, 0x12,
	0x0d, 0x0a, 0x09, 0x45, 0x43, 0x32, 0x53, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x10, 0x04, 0x12, 0x0d,
	0x0a, 0x09, 0x45, 0x53, 0x32, 0x43, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x10, 0x05, 0x12, 0x0c, 0x0a,
	0x08, 0x45, 0x43, 0x32, 0x53, 0x53, 0x74, 0x65, 0x70, 0x10, 0x06, 0x12, 0x0c, 0x0a, 0x08, 0x45,
	0x53, 0x32, 0x43, 0x53, 0x74, 0x65, 0x70, 0x10, 0x07, 0x12, 0x12, 0x0a, 0x0e, 0x45, 0x53, 0x32,
	0x43, 0x47, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x10, 0x08, 0x12, 0x0f, 0x0a,
	0x0b, 0x45, 0x53, 0x32, 0x43, 0x50, 0x75, 0x73, 0x68, 0x4d, 0x73, 0x67, 0x10, 0x09, 0x2a, 0x4a,
	0x0a, 0x09, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x54, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x10, 0x4d,
	0x61, 0x74, 0x63, 0x68, 0x54, 0x79, 0x70, 0x65, 0x44, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x10,
	0x00, 0x12, 0x12, 0x0a, 0x0e, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x54, 0x79, 0x70, 0x65, 0x4d, 0x61,
	0x74, 0x63, 0x68, 0x10, 0x01, 0x12, 0x13, 0x0a, 0x0f, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x54, 0x79,
	0x70, 0x65, 0x43, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x10, 0x02, 0x2a, 0x7a, 0x0a, 0x0b, 0x4d, 0x61,
	0x74, 0x63, 0x68, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x14, 0x0a, 0x10, 0x4d, 0x61, 0x74,
	0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x44, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x10, 0x00, 0x12,
	0x14, 0x0a, 0x10, 0x4d, 0x61, 0x74, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x53, 0x75, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x10, 0x01, 0x12, 0x13, 0x0a, 0x0f, 0x4d, 0x61, 0x74, 0x52, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x46, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x10, 0x02, 0x12, 0x15, 0x0a, 0x11, 0x4d, 0x61,
	0x74, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x69, 0x6e, 0x67, 0x10,
	0x03, 0x12, 0x13, 0x0a, 0x0f, 0x4d, 0x61, 0x74, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x43, 0x61,
	0x6e, 0x63, 0x65, 0x6c, 0x10, 0x04, 0x2a, 0x47, 0x0a, 0x09, 0x43, 0x6f, 0x6c, 0x6f, 0x72, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x10, 0x43, 0x6f, 0x6c, 0x6f, 0x72, 0x54, 0x79, 0x70, 0x65,
	0x44, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x10, 0x00, 0x12, 0x10, 0x0a, 0x0c, 0x43, 0x6f, 0x6c,
	0x6f, 0x72, 0x54, 0x79, 0x70, 0x65, 0x52, 0x65, 0x64, 0x10, 0x01, 0x12, 0x12, 0x0a, 0x0e, 0x43,
	0x6f, 0x6c, 0x6f, 0x72, 0x54, 0x79, 0x70, 0x65, 0x42, 0x6c, 0x61, 0x63, 0x6b, 0x10, 0x02, 0x2a,
	0x5e, 0x0a, 0x0a, 0x47, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x15, 0x0a,
	0x11, 0x47, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x44, 0x65, 0x66, 0x61, 0x75,
	0x6c, 0x74, 0x10, 0x00, 0x12, 0x11, 0x0a, 0x0d, 0x47, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x57, 0x69, 0x6e, 0x10, 0x01, 0x12, 0x12, 0x0a, 0x0e, 0x47, 0x61, 0x6d, 0x65, 0x52,
	0x65, 0x73, 0x75, 0x6c, 0x74, 0x46, 0x61, 0x69, 0x6c, 0x10, 0x02, 0x12, 0x12, 0x0a, 0x0e, 0x47,
	0x61, 0x6d, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x44, 0x72, 0x61, 0x77, 0x10, 0x03, 0x42,
	0x2a, 0x0a, 0x06, 0x67, 0x61, 0x6d, 0x65, 0x70, 0x62, 0x42, 0x06, 0x67, 0x61, 0x6d, 0x65, 0x70,
	0x62, 0x50, 0x01, 0x5a, 0x16, 0x74, 0x69, 0x61, 0x6f, 0x79, 0x75, 0x79, 0x75, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x70, 0x62, 0x2f, 0x67, 0x61, 0x6d, 0x65, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_enum_proto_rawDescOnce sync.Once
	file_enum_proto_rawDescData = file_enum_proto_rawDesc
)

func file_enum_proto_rawDescGZIP() []byte {
	file_enum_proto_rawDescOnce.Do(func() {
		file_enum_proto_rawDescData = protoimpl.X.CompressGZIP(file_enum_proto_rawDescData)
	})
	return file_enum_proto_rawDescData
}

var file_enum_proto_enumTypes = make([]protoimpl.EnumInfo, 5)
var file_enum_proto_goTypes = []interface{}{
	(ProtocolType)(0), // 0: gamepb.ProtocolType
	(MatchType)(0),    // 1: gamepb.MatchType
	(MatchResult)(0),  // 2: gamepb.MatchResult
	(ColorType)(0),    // 3: gamepb.ColorType
	(GameResult)(0),   // 4: gamepb.GameResult
}
var file_enum_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_enum_proto_init() }
func file_enum_proto_init() {
	if File_enum_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_enum_proto_rawDesc,
			NumEnums:      5,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_enum_proto_goTypes,
		DependencyIndexes: file_enum_proto_depIdxs,
		EnumInfos:         file_enum_proto_enumTypes,
	}.Build()
	File_enum_proto = out.File
	file_enum_proto_rawDesc = nil
	file_enum_proto_goTypes = nil
	file_enum_proto_depIdxs = nil
}
