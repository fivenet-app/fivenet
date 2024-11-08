package protoutils

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type ProtoMessageWithMerge[T any] interface {
	*T
	proto.Message

	Merge(in *T) *T
}

type ProtoMessage[T any] interface {
	*T
	proto.Message
}

type ProtoEnum interface {
	Number() protoreflect.EnumNumber
}
