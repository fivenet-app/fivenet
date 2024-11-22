package common

import (
	"database/sql/driver"

	"github.com/fivenet-app/fivenet/pkg/utils/protoutils"
	"google.golang.org/protobuf/encoding/protojson"
)

// Scan implements driver.Valuer for protobuf TranslateItem.
func (x *TranslateItem) Scan(value any) error {
	switch t := value.(type) {
	case string:
		return protojson.Unmarshal([]byte(t), x)
	case []byte:
		return protojson.Unmarshal(t, x)
	}
	return nil
}

// Value marshals the value into driver.Valuer.
func (x *TranslateItem) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}

	out, err := protoutils.Marshal(x)
	return string(out), err
}

func NewTranslateItem(key string) *TranslateItem {
	return &TranslateItem{
		Key: key,
	}
}

func NewTranslateItemWithParams(key string, params map[string]string) *TranslateItem {
	return &TranslateItem{
		Key:        key,
		Parameters: params,
	}
}
