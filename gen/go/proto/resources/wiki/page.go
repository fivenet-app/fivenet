package wiki

func (x *Page) SetJobLabel(label string) {
	x.JobLabel = &label
}

func (x *PageShort) SetJobLabel(label string) {
	x.JobLabel = &label
}

func (x *PageAccess) IsEmpty() bool {
	return len(x.Jobs) == 0 && len(x.Users) == 0
}
