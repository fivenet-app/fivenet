package livemap

import (
	"github.com/paulmach/orb"
	"google.golang.org/protobuf/proto"
)

func (x *UserMarker) Point() orb.Point {
	return orb.Point{x.X, x.Y}
}

func (x *UserMarker) SetJobLabel(label string) {
	x.JobLabel = label
}

func (x *UserMarker) Merge(in *UserMarker) *UserMarker {
	if x.UserId != in.UserId {
		x.UserId = in.UserId
	}

	x.X = in.X
	x.Y = in.Y

	if in.UpdatedAt != nil {
		x.UpdatedAt = in.UpdatedAt
	}

	if in.Postal == nil {
		x.Postal = nil
	} else {
		x.Postal = in.Postal
	}

	if in.Color == nil {
		x.Color = nil
	} else {
		x.Color = in.Color
	}

	x.Job = in.Job
	x.JobLabel = in.JobLabel

	if in.User != nil {
		if x.User == nil {
			x.User = in.User
		} else {
			proto.Merge(x.User, in.User)
		}
	}

	if in.UnitId == nil {
		x.UnitId = nil
		x.Unit = nil
	} else {
		x.UnitId = in.UnitId
		x.Unit = in.Unit
	}

	x.Hidden = in.Hidden

	return x
}

func (x *MarkerMarker) SetJobLabel(label string) {
	x.JobLabel = label
}
