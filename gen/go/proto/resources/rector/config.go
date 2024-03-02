package rector

import (
	"database/sql/driver"

	"google.golang.org/protobuf/encoding/protojson"
)

func (x *AppConfig) Default() {
	// TODO
}

func (x *PluginConfig) Default() {
	// TODO
}

func (x *AppConfig) Scan(value any) error {
	switch t := value.(type) {
	case string:
		return protojson.Unmarshal([]byte(t), x)
	case []byte:
		return protojson.Unmarshal(t, x)
	}
	return nil
}

// Scan implements driver.Valuer for protobuf AppConfig.
func (x *AppConfig) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}

	out, err := protojson.Marshal(x)
	return string(out), err
}

func (x *PluginConfig) Scan(value any) error {
	switch t := value.(type) {
	case string:
		return protojson.Unmarshal([]byte(t), x)
	case []byte:
		return protojson.Unmarshal(t, x)
	}
	return nil
}

// Scan implements driver.Valuer for protobuf PluginConfig.
func (x *PluginConfig) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}

	out, err := protojson.Marshal(x)
	return string(out), err
}
