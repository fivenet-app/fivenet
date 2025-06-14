package tracker

func (x *UserMapping) Merge(in *UserMapping) *UserMapping {
	if x.UserId != in.UserId {
		x.UserId = in.UserId
	}

	if x.UnitId != in.UnitId {
		x.UnitId = in.UnitId
	}

	if x.CreatedAt == nil {
		x.CreatedAt = in.CreatedAt
	}

	return x
}
