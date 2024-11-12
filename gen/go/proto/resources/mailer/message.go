package mailer

import (
	"database/sql/driver"

	"google.golang.org/protobuf/encoding/protojson"
)

// Scan implements driver.Valuer for protobuf MessageData.
func (x *MessageData) Scan(value any) error {
	switch t := value.(type) {
	case string:
		return protojson.Unmarshal([]byte(t), x)
	case []byte:
		return protojson.Unmarshal(t, x)
	}
	return nil
}

// Value marshals the value into driver.Valuer.
func (x *MessageData) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}

	out, err := protojson.Marshal(x)
	return out, err
}
