package livemap

import (
	"github.com/paulmach/orb"
	"google.golang.org/protobuf/proto"
)

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
