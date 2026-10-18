package protoutils

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// ProtoMessageWithMerge is a generic interface for protobuf messages with the
// application-specific merge semantics required by the cached KV store.
type ProtoMessageWithMerge[T any] interface {
	*T
	proto.Message

	Merge(in *T) *T
}

// ProtoMessage is a generic interface for protobuf messages.
// T is the underlying message type.
type ProtoMessage[T any] interface {
	*T
	proto.Message
}

// ProtoEnum is a generic interface for protobuf enums, requiring a Number method.
type ProtoEnum interface {
	Number() protoreflect.EnumNumber
	String() string
}
