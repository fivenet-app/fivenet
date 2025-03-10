package users

func (x *CitizenLabel) Equal(a *CitizenLabel) bool {
	return x.Name == a.Name
}
