package common

func (x *IDMapping) Merge(in *IDMapping) *IDMapping {
	if in.Id != 0 && in.Id != x.Id {
		x.Id = in.Id
	}

	return x
}
