package centrum

import (
	"database/sql/driver"

	"github.com/fivenet-app/fivenet/pkg/utils/protoutils"
	"google.golang.org/protobuf/encoding/protojson"
)

// Scan implements driver.Valuer for protobuf PredefinedStatus.
func (x *PredefinedStatus) Scan(value any) error {
	switch t := value.(type) {
	case string:
		return protojson.Unmarshal([]byte(t), x)
	case []byte:
		return protojson.Unmarshal(t, x)
	}
	return nil
}

// Value marshals the value into driver.Valuer.
func (x *PredefinedStatus) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}

	out, err := protoutils.Marshal(x)
	return string(out), err
}

// Scan implements driver.Valuer for protobuf Timings.
func (x *Timings) Scan(value any) error {
	switch t := value.(type) {
	case string:
		return protojson.Unmarshal([]byte(t), x)
	case []byte:
		return protojson.Unmarshal(t, x)
	}
	return nil
}

// Value marshals the value into driver.Valuer.
func (x *Timings) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}

	out, err := protoutils.Marshal(x)
	return string(out), err
}

func (x *Settings) Default(job string) {
	x.Job = job

	if x.Mode <= CentrumMode_CENTRUM_MODE_UNSPECIFIED {
		x.Mode = CentrumMode_CENTRUM_MODE_MANUAL
	}

	if x.FallbackMode <= CentrumMode_CENTRUM_MODE_UNSPECIFIED {
		x.FallbackMode = CentrumMode_CENTRUM_MODE_MANUAL
	}

	if x.PredefinedStatus == nil {
		x.PredefinedStatus = &PredefinedStatus{}
	}

	if x.Timings == nil {
		x.Timings = &Timings{}
	}
	if x.Timings.DispatchMaxWait == 0 {
		x.Timings.DispatchMaxWait = 900
	}
	if x.Timings.RequireUnitReminderSeconds == 0 {
		x.Timings.RequireUnitReminderSeconds = 180
	}
}
