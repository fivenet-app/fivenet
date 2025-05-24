package users

func (x *Label) Equal(a *Label) bool {
	return x.Name == a.Name
}
