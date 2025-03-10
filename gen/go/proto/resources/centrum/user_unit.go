package centrum

func (x *UserUnitMapping) Merge(in *UserUnitMapping) *UserUnitMapping {
	if x.UnitId != in.UnitId {
		x.UnitId = in.UnitId
	}

	if x.UserId != in.UserId {
		x.UserId = in.UserId
	}

	return x
}
