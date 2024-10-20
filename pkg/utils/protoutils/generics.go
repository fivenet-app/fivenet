package protoutils

import "google.golang.org/protobuf/proto"

type ProtoMessage[T any] interface {
	*T
	proto.Message

	Merge(in *T) *T
}
