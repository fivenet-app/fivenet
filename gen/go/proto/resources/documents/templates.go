package documents

import (
	"database/sql/driver"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func (x *TemplateSchema) Scan(value any) error {
	switch t := value.(type) {
	case string:
		return json.UnmarshalFromString(t, x)
	case []byte:
		return json.Unmarshal(t, x)
	}
	return nil
}

// Scan implements driver.Valuer for protobuf TemplateSchema.
func (x *TemplateSchema) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}

	out, err := json.MarshalToString(x)
	return out, err
}
