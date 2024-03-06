package protoutils

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func UnmarshalPartial(b []byte, m proto.Message) error {
	return protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}.Unmarshal(b, m)
}

func Marshal(m proto.Message) ([]byte, error) {
	return protojson.MarshalOptions{}.Marshal(m)
}

func MarshalPretty(m proto.Message) ([]byte, error) {
	return protojson.MarshalOptions{
		Multiline: true,
		Indent:    "    ",
	}.Marshal(m)
}
