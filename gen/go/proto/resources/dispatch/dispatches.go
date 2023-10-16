package dispatch

import (
	"database/sql/driver"

	"dario.cat/mergo"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func (x *Attributes) Scan(value any) error {
	switch t := value.(type) {
	case string:
		return json.UnmarshalFromString(t, x)
	case []byte:
		return json.Unmarshal(t, x)
	}
	return nil
}

// Scan implements driver.Valuer for protobuf Attributes.
func (x *Attributes) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}

	out, err := json.MarshalToString(x)
	return out, err
}

func (x *Dispatch) Update(in *Dispatch) {
	if x.Id != in.Id {
		return
	}

	err := mergo.Merge(x, in)
	if err != nil {
		return
	}
}
