package protoutils

import (
	"github.com/cespare/xxhash/v2"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// UnmarshalPartialJSON unmarshals a JSON representation (can be partial) into a proto message.
func UnmarshalPartialJSON(b []byte, m proto.Message) error {
	return protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}.Unmarshal(b, m)
}

// MarshalToJSON marshals a proto message to a JSON representation.
// It uses the default protojson options.
func MarshalToJSON(m proto.Message) ([]byte, error) {
	return protojson.MarshalOptions{}.Marshal(m)
}

// MarshalToPrettyJSON marshals a proto message to a pretty-printed JSON representation.
func MarshalToPrettyJSON(m proto.Message) ([]byte, error) {
	return protojson.MarshalOptions{
		Multiline: true,
		Indent:    "    ",
	}.Marshal(m)
}

// JSONAndHash is a helper that returns the JSON representation of a proto message along with its xxhash hash.
func JSONAndHash(m proto.Message) ([]byte, uint64, error) {
	jsonBytes, err := MarshalToJSON(m)
	if err != nil {
		return nil, 0, err
	}
	hash := xxhash.Sum64(jsonBytes)
	return jsonBytes, hash, nil
}
