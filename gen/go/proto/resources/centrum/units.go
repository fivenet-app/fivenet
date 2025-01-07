package centrum

import (
	"github.com/paulmach/orb"
	"google.golang.org/protobuf/proto"
)

func (x *Unit) Merge(in *Unit) *Unit {
	if x.Id != in.Id {
		return x
	}

	if in.CreatedAt != nil {
		if x.CreatedAt == nil {
			x.CreatedAt = in.CreatedAt
		} else {
			proto.Merge(x.CreatedAt, in.CreatedAt)
		}
	}

	if in.UpdatedAt != nil {
		if x.UpdatedAt == nil {
			x.UpdatedAt = in.UpdatedAt
		} else {
			proto.Merge(x.UpdatedAt, in.UpdatedAt)
		}
	}

	if x.Job != in.Job {
		x.Job = in.Job
	}

	if x.Name != in.Name {
		x.Name = in.Name
	}

	if x.Initials != in.Initials {
		x.Initials = in.Initials
	}

	if x.Color != in.Color {
		x.Color = in.Color
	}

	if in.Description != nil && (x.Description == nil || x.Description != in.Description) {
		x.Description = in.Description
	}

	if in.Status != nil {
		// Only update status if it is newer (higher ID)
		if x.Status == nil || x.Status.Id < in.Status.Id {
			x.Status = in.Status
		}
	}

	if len(in.Users) == 0 {
		x.Users = []*UnitAssignment{}
	} else {
		x.Users = in.Users
	}

	if in.Attributes != nil {
		if x.Attributes == nil {
			x.Attributes = in.Attributes
		} else {
			x.Attributes.List = in.Attributes.List
		}
	}

	x.HomePostal = in.HomePostal

	if in.Access != nil {
		x.Access = in.Access
	}

	return x
}

func (x *UnitStatus) Point() orb.Point {
	if x.X == nil || x.Y == nil {
		return orb.Point{}
	}

	return orb.Point{*x.X, *x.Y}
}
