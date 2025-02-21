package tests

func (x *SimpleObject) Merge(in *SimpleObject) *SimpleObject {
	if in == nil {
		return x
	}

	x.Field1 = in.Field1

	x.Field2 = in.Field2

	return x
}
