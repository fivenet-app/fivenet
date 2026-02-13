package centrumunits

import (
	"slices"

	"github.com/paulmach/orb"
	"google.golang.org/protobuf/proto"
)

const DefaultUnitIcon = "MapMarkerIcon"

func (x *Unit) Merge(in *Unit) *Unit {
	if x.GetId() != in.GetId() {
		return x
	}

	if in.GetCreatedAt() != nil {
		if x.GetCreatedAt() == nil {
			x.CreatedAt = in.GetCreatedAt()
		} else {
			proto.Merge(x.GetCreatedAt(), in.GetCreatedAt())
		}
	}

	if in.GetUpdatedAt() != nil {
		if x.GetUpdatedAt() == nil {
			x.UpdatedAt = in.GetUpdatedAt()
		} else {
			proto.Merge(x.GetUpdatedAt(), in.GetUpdatedAt())
		}
	}

	if x.GetJob() != in.GetJob() {
		x.Job = in.GetJob()
	}

	if x.GetName() != in.GetName() {
		x.Name = in.GetName()
	}

	if x.GetInitials() != in.GetInitials() {
		x.Initials = in.GetInitials()
	}

	if x.GetColor() != in.GetColor() {
		x.Color = in.GetColor()
	}

	if in.Icon != nil && in.GetIcon() != "" {
		x.Icon = in.Icon
	} else if x.Icon == nil || x.GetIcon() == "" {
		def := DefaultUnitIcon
		x.Icon = &def
	}

	if in.Description != nil &&
		(x.Description == nil || x.GetDescription() != in.GetDescription()) {
		x.Description = in.Description
	}

	if in.GetStatus() != nil {
		// Only update status if it is newer (higher ID)
		if x.GetStatus() == nil || x.GetStatus().GetId() < in.GetStatus().GetId() {
			x.Status = in.GetStatus()
		}
	}

	if len(in.GetUsers()) == 0 {
		x.Users = []*UnitAssignment{}
	} else {
		x.Users = in.GetUsers()
	}

	if in.GetAttributes() != nil {
		if x.GetAttributes() == nil {
			x.Attributes = in.GetAttributes()
		} else {
			x.Attributes.List = in.GetAttributes().GetList()
		}
	}

	x.HomePostal = in.HomePostal

	if in.GetAccess() != nil {
		x.Access = in.GetAccess()
	}

	return x
}

func (x *UnitStatus) Point() orb.Point {
	if x.X == nil || x.Y == nil {
		return orb.Point{}
	}

	return orb.Point{x.GetX(), x.GetY()}
}

func (x *UnitAttributes) Has(attribute UnitAttribute) bool {
	if len(x.GetList()) == 0 {
		return false
	}

	return slices.Contains(x.GetList(), attribute)
}

func (x *UnitAttributes) Add(attribute UnitAttribute) bool {
	if x.Has(attribute) {
		return false
	}

	if x.List == nil {
		x.List = []UnitAttribute{attribute}
	} else {
		x.List = append(x.List, attribute)
	}

	return true
}

func (x *UnitAttributes) Remove(attribute UnitAttribute) bool {
	if !x.Has(attribute) {
		return false
	}

	x.List = slices.DeleteFunc(x.GetList(), func(item UnitAttribute) bool {
		return item == attribute
	})

	return true
}
