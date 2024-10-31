package livemap

import (
	"database/sql/driver"

	"github.com/paulmach/orb"
	"google.golang.org/protobuf/proto"
)

// Scan implements driver.Valuer for protobuf MarkerData.
func (x *MarkerData) Scan(value any) error {
	switch t := value.(type) {
	case string:
		return proto.Unmarshal([]byte(t), x)
	case []byte:
		return proto.Unmarshal(t, x)
	}
	return nil
}

// Value marshals the value into driver.Valuer.
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
		x.Unit = nil
	} else {
		x.UnitId = in.UnitId
		x.Unit = in.Unit
	}

	if x.UserId != in.UserId {
		x.UserId = in.UserId
	}

	if in.User != nil {
		if x.User == nil {
			x.User = in.User
		} else {
			proto.Merge(x.User, in.User)
		}
	}

	if in.Info != nil {
		proto.Merge(x.Info, in.Info)
	}

	x.Hidden = in.Hidden

	return x
}
