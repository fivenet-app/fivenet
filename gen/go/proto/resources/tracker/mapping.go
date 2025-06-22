package tracker

func (x *UserMapping) Merge(in *UserMapping) *UserMapping {
	x.UserId = in.UserId

	if in.UnitId != nil {
		x.UnitId = in.UnitId
	} else {
		x.UnitId = nil
	}

	if x.CreatedAt == nil {
		x.CreatedAt = in.CreatedAt
	}

	return x
}
