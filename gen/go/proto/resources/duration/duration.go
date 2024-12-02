package duration

import (
	"database/sql/driver"
	"errors"
	"time"

	"google.golang.org/protobuf/types/known/durationpb"
)

// New constructs a new Duration from the provided time.Duration.
func New(t time.Duration) *Duration {
	return &Duration{
		Duration: durationpb.New(t),
	}
}

// AsDuration converts x to a time.Duration.
func (x *Duration) AsDuration() time.Duration {
	return x.GetDuration().AsDuration()
}

// Scan implements sql.Scanner for protobuf Duration.
func (ts *Duration) Scan(value any) error {
	switch t := value.(type) {
	case time.Duration:
		ts.Duration = durationpb.New(t) // google proto version
	case string:
		dur, err := time.ParseDuration(t)
		if err != nil {
			return err
		}
		ts.Duration = durationpb.New(dur)
	default:
		return errors.New("not a protobuf Duration")
	}
	return nil
}

// Value marshals the value into driver.Valuer.
func (ts *Duration) Value() (driver.Value, error) {
	if ts == nil {
		return nil, nil
	}
	return ts.Duration.AsDuration(), nil
}
