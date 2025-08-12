package tests

func (x *SimpleObject) Merge(in *SimpleObject) *SimpleObject {
	if in == nil {
		return x
	}

	x.Field1 = in.GetField1()

	x.Field2 = in.GetField2()

	return x
}
