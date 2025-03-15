package centrum

import (
	"slices"

	"github.com/paulmach/orb"
	"google.golang.org/protobuf/proto"
)

func DispatchPointMatchFn(dspId uint64) func(p orb.Pointer) bool {
	return func(p orb.Pointer) bool {
		return p.(*Dispatch).Id == dspId
	}
}

func (x *Dispatch) Merge(in *Dispatch) *Dispatch {
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

	if in.Status != nil {
		// Only update status if it is newer (higher ID)
		if x.Status == nil || x.Status.Id < in.Status.Id {
			x.Status = in.Status
		}
	}

	if x.Message != in.Message {
		x.Message = in.Message
	}

	if in.Description != nil && (x.Description == nil || *x.Description != *in.Description) {
		x.Description = in.Description
	}

	if in.Attributes != nil {
		if x.Attributes == nil {
			x.Attributes = in.Attributes
		} else {
			x.Attributes.List = in.Attributes.List
		}
	}

	if x.X != in.X {
		x.X = in.X
	}

	if x.Y != in.Y {
		x.Y = in.Y
	}

	if in.Postal != nil && (x.Postal == nil || *x.Postal != *in.Postal) {
		x.Postal = in.Postal
	}

	if x.Anon != in.Anon {
		x.Anon = in.Anon
	}

	if in.CreatorId != nil && (x.CreatorId == nil || *x.CreatorId != *in.CreatorId) {
		x.CreatorId = in.CreatorId
	}

	if in.Creator != nil {
		if x.Creator == nil {
			x.Creator = in.Creator
		} else {
			proto.Merge(x.Creator, in.Creator)
		}
	}

	x.Units = in.Units

	if in.References != nil {
		if x.References == nil {
			x.References = in.References
		} else {
			x.References.References = in.References.References
		}
	}

	return x
}

func (x *Dispatch) Point() orb.Point {
	return orb.Point{x.X, x.Y}
}

func (x *DispatchStatus) Point() orb.Point {
	if x.X == nil || x.Y == nil {
		return orb.Point{}
	}

	return orb.Point{*x.X, *x.Y}
}

func (x *DispatchReferences) Has(dspId uint64) bool {
	if len(x.References) == 0 {
		return false
	}

	return slices.ContainsFunc(x.References, func(r *DispatchReference) bool {
		return r.TargetDispatchId == dspId
	})
}

func (x *DispatchReferences) Add(ref *DispatchReference) bool {
	if x.Has(ref.TargetDispatchId) {
		return false
	}

	if x.References == nil {
		x.References = []*DispatchReference{ref}
	} else {
		x.References = append(x.References, ref)
	}

	return true
}

func (x *DispatchReferences) Remove(dspId uint64) bool {
	if !x.Has(dspId) {
		return false
	}

	x.References = slices.DeleteFunc(x.References, func(item *DispatchReference) bool {
		return item.TargetDispatchId == dspId
	})

	return true
}
