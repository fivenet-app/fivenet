package users

import (
	"database/sql/driver"

	"github.com/fivenet-app/fivenet/pkg/utils/protoutils"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	DefaultJobAbsencePastDays   = 7
	DefaultJobAbsenceFutureDays = 93 // ~3 months
)

func (x *JobSettings) Default() {
	if x.AbsencePastDays <= 0 {
		x.AbsencePastDays = DefaultJobAbsencePastDays
	}
	if x.AbsenceFutureDays <= 0 {
		x.AbsenceFutureDays = DefaultJobAbsenceFutureDays
	}
}

// Scan implements driver.Valuer for protobuf JobSettings.
func (x *JobSettings) Scan(value any) error {
	switch t := value.(type) {
	case string:
		return protojson.Unmarshal([]byte(t), x)
	case []byte:
		return protojson.Unmarshal(t, x)
	}
	return nil
}

// Value marshals the value into driver.Valuer.
func (x *JobSettings) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}

	out, err := protoutils.Marshal(x)
	return string(out), err
}
