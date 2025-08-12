package tracker

func (x *UserMapping) Merge(in *UserMapping) *UserMapping {
	x.UserId = in.GetUserId()

	if in.UnitId != nil {
		x.UnitId = in.UnitId
	} else {
		x.UnitId = nil
	}

	if x.GetCreatedAt() == nil {
		x.CreatedAt = in.GetCreatedAt()
	}

	return x
}
