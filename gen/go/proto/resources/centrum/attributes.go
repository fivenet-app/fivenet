package centrum

import (
	"slices"
)

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

func (x *DispatchAttributes) Has(attribute DispatchAttribute) bool {
	if len(x.GetList()) == 0 {
		return false
	}

	return slices.Contains(x.GetList(), attribute)
}

func (x *DispatchAttributes) Add(attribute DispatchAttribute) bool {
	if x.Has(attribute) {
		return false
	}

	if x.List == nil {
		x.List = []DispatchAttribute{attribute}
	} else {
		x.List = append(x.List, attribute)
	}

	return true
}

func (x *DispatchAttributes) Remove(attribute DispatchAttribute) bool {
	if !x.Has(attribute) {
		return false
	}

	x.List = slices.DeleteFunc(x.GetList(), func(item DispatchAttribute) bool {
		return item == attribute
	})

	return true
}
