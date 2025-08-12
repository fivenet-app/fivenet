package common

func (x *IDMapping) Merge(in *IDMapping) *IDMapping {
	if in.GetId() != 0 && in.GetId() != x.GetId() {
		x.Id = in.GetId()
	}

	return x
}
