package tests

func (x *SimpleObject) Merge(in *SimpleObject) *SimpleObject {
	if in == nil {
		return x
	}

	x.SetField1(in.GetField1())
	x.SetField2(in.GetField2())

	return x
}
