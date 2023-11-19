package centrum

import users "github.com/galexrt/fivenet/gen/go/proto/resources/users"

func (x *Disponents) Merge(in *Disponents) *Disponents {
	if len(in.Disponents) == 0 {
		x.Disponents = []*users.UserShort{}
	} else {
		x.Disponents = in.Disponents
	}

	return x
}

func (x *UserUnitMapping) Merge(in *UserUnitMapping) *UserUnitMapping {
	if x.UnitId != in.UnitId {
		x.UnitId = in.UnitId
	}

	if x.UserId != in.UserId {
		x.UserId = in.UserId
	}

	return x
}
