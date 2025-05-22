package centrum

func (x *UserUnitMapping) Merge(in *UserUnitMapping) *UserUnitMapping {
	if x.UnitId != in.UnitId {
		x.UnitId = in.UnitId
	}
	if x.UnitJob != in.UnitJob {
		x.UnitJob = in.UnitJob
	}

	if x.UserId != in.UserId {
		x.UserId = in.UserId
	}
	if x.UserJob != in.UserJob {
		x.UserJob = in.UserJob
	}

	return x
}
