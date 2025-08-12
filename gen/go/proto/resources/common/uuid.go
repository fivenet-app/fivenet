package common

import (
	"database/sql/driver"

	"github.com/google/uuid"
)

// Scan implements driver.Valuer for protobuf UUID.
func (x *UUID) Scan(value any) error {
	switch t := value.(type) {
	case string:
		uid, err := uuid.FromBytes([]byte(t))
		x.Uuid = uid.String()
		return err
	case []byte:
		uid, err := uuid.FromBytes(t)
		x.Uuid = uid.String()
		return err
	}
	return nil
}

// Value marshals the value into driver.Valuer.
func (x *UUID) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}

	return x.GetUuid(), nil
}
