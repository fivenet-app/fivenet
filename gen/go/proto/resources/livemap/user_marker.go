package livemap

import (
	timestamp "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	"github.com/paulmach/orb"
	"google.golang.org/protobuf/proto"
)

func (x *UserMarker) Point() orb.Point {
	return orb.Point{x.GetX(), x.GetY()}
}

func (x *UserMarker) SetJobLabel(label string) {
	x.JobLabel = label
}

func (x *UserMarker) Merge(in *UserMarker) *UserMarker {
	if x.GetUserId() != in.GetUserId() {
		x.UserId = in.GetUserId()
	}

	x.X = in.GetX()
	x.Y = in.GetY()

	if in.GetUpdatedAt() != nil {
		//nolint:forcetypeassert // Value type is guaranteed to be timestamp.Timestamp
		x.UpdatedAt = proto.Clone(in.GetUpdatedAt()).(*timestamp.Timestamp)
	}

	if in.Postal == nil {
		x.Postal = nil
	} else {
		*x.Postal = in.GetPostal()
	}

	if in.Color == nil {
		x.Color = nil
	} else {
		*x.Color = in.GetColor()
	}

	x.Job = in.GetJob()
	x.JobLabel = in.GetJobLabel()
	if in.JobGrade == nil {
		x.JobGrade = nil
	} else {
		val := in.GetJobGrade()
		x.JobGrade = &val
	}

	if in.GetUser() != nil {
		if x.GetUser() == nil {
			x.User = in.GetUser()
		} else {
			proto.Merge(x.GetUser(), in.GetUser())
		}
	}

	if in.UnitId == nil {
		x.UnitId = nil
		x.Unit = nil
	} else {
		x.UnitId = in.UnitId
		x.Unit = in.GetUnit()
	}

	x.Hidden = in.GetHidden()

	return x
}
