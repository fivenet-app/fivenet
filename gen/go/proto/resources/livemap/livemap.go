package livemap

import (
	"database/sql/driver"

	"github.com/paulmach/orb"
	"google.golang.org/protobuf/proto"
)

func (x *MarkerData) Scan(value any) error {
	switch t := value.(type) {
	case string:
		return proto.Unmarshal([]byte(t), x)
	case []byte:
		return proto.Unmarshal(t, x)
	}
	return nil
}

// Scan implements driver.Valuer for protobuf Data.
func (x *MarkerData) Value() (driver.Value, error) {
	if x == nil {
		return nil, nil
	}

	out, err := proto.Marshal(x)
	return out, err
}

func (x *UserMarker) Point() orb.Point {
	if x.Info == nil {
		return orb.Point{}
	}

	return orb.Point{x.Info.X, x.Info.Y}
}

func (x *MarkerInfo) SetJobLabel(label string) {
	x.JobLabel = label
}

func (x *UserMarker) Merge(in *UserMarker) *UserMarker {
	if in.UnitId == nil {
		x.UnitId = nil
	}
	if in.Unit == nil {
		x.Unit = nil
	}

	proto.Merge(x, in)

	return x
}
