package main

import (
	"TServerGo/pb"
	"encoding/binary"

	"google.golang.org/protobuf/proto"
)

func SendMsg(message proto.Message, t pb.ProtocolType) []byte {
	ack, err := proto.Marshal(message)
	if err != nil {
		return nil
	}

	var bufHead = make([]byte, 4)
	var bufPId = make([]byte, 4)
	binary.BigEndian.PutUint32(bufPId, uint32(t))
	binary.BigEndian.PutUint32(bufHead, uint32(len(ack)+4))
	bufHead = append(bufHead, bufPId...)
	bufHead = append(bufHead, ack...)

	return bufHead
}
