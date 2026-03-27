package citizenslabels

func (x *Label) Equal(a *Label) bool {
	return x.GetName() == a.GetName()
}
