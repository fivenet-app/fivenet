package centrum

import (
	"slices"

	"github.com/paulmach/orb"
	"google.golang.org/protobuf/proto"
)

func DispatchPointMatchFn(dspId uint64) func(p orb.Pointer) bool {
	return func(p orb.Pointer) bool {
		//nolint:forcetypeassert // Value type is guaranteed to be Dispatch due to generics type
		return p.(*Dispatch).GetId() == dspId
	}
}

func (x *Dispatch) Merge(in *Dispatch) *Dispatch {
	if x.GetId() != in.GetId() {
		return x
	}

	if in.GetCreatedAt() != nil {
		x.CreatedAt = in.GetCreatedAt()
	}

	if in.GetUpdatedAt() != nil {
		x.UpdatedAt = in.GetUpdatedAt()
	}

	x.Jobs = in.GetJobs()

	if in.GetStatus() != nil {
		// Only update status if it is newer (higher ID)
		if x.GetStatus() == nil || x.GetStatus().GetId() < in.GetStatus().GetId() {
			x.Status = in.GetStatus()
		}
	}

	if x.GetMessage() != in.GetMessage() {
		x.Message = in.GetMessage()
	}

	if in.Description != nil &&
		(x.Description == nil || x.GetDescription() != in.GetDescription()) {
		x.Description = in.Description
	}

	if in.GetAttributes() != nil {
		if x.GetAttributes() == nil {
			x.Attributes = in.GetAttributes()
		} else {
			x.Attributes.List = in.GetAttributes().GetList()
		}
	}

	if x.GetX() != in.GetX() {
		x.X = in.GetX()
	}

	if x.GetY() != in.GetY() {
		x.Y = in.GetY()
	}

	if in.Postal != nil && (x.Postal == nil || x.GetPostal() != in.GetPostal()) {
		x.Postal = in.Postal
	}

	if x.GetAnon() != in.GetAnon() {
		x.Anon = in.GetAnon()
	}

	if in.CreatorId != nil && (x.CreatorId == nil || x.GetCreatorId() != in.GetCreatorId()) {
		x.CreatorId = in.CreatorId
	}

	if in.GetCreator() != nil {
		if x.GetCreator() == nil {
			x.Creator = in.GetCreator()
		} else {
			proto.Merge(x.GetCreator(), in.GetCreator())
		}
	}

	x.Units = in.GetUnits()

	if in.GetReferences() != nil {
		if x.GetReferences() == nil {
			x.References = in.GetReferences()
		} else {
			x.References.References = in.GetReferences().GetReferences()
		}
	}

	return x
}

func (x *Dispatch) Point() orb.Point {
	return orb.Point{x.GetX(), x.GetY()}
}

func (x *DispatchStatus) Point() orb.Point {
	if x.X == nil || x.Y == nil {
		return orb.Point{}
	}

	return orb.Point{x.GetX(), x.GetY()}
}

func (x *DispatchReferences) Has(dspId uint64) bool {
	if len(x.GetReferences()) == 0 {
		return false
	}

	return slices.ContainsFunc(x.GetReferences(), func(r *DispatchReference) bool {
		return r.GetTargetDispatchId() == dspId
	})
}

func (x *DispatchReferences) Add(ref *DispatchReference) bool {
	if x.Has(ref.GetTargetDispatchId()) {
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

	x.References = slices.DeleteFunc(x.GetReferences(), func(item *DispatchReference) bool {
		return item.GetTargetDispatchId() == dspId
	})

	return true
}
