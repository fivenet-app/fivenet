package documents

func (x *Category) Merge(in *Category) *Category {
	if in != nil {
		x = in
	}

	return x
}
