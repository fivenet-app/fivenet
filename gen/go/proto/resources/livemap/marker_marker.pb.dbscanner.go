// Code generated by protoc-gen-customizer. DO NOT EDIT.
// source: resources/livemap/marker_marker.proto

package livemap

import (
	"database/sql/driver"

	"google.golang.org/protobuf/proto"
)

// Scan implements driver.Valuer for protobuf MarkerData.
func (x *MarkerData) Scan(value any) error {
	switch t := value.(type) {
	case string:
		return proto.Unmarshal([]byte(t), x)
	case []byte:
		return proto.Unmarshal(t, x)
	}
	return nil
}

// Value marshals the MarkerData value into driver.Valuer.
func (x *MarkerData) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}

	out, err := proto.Marshal(x)
	return string(out), err
}
