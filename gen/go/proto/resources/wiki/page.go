package wiki

func (x *Page) SetJobLabel(label string) {
	x.JobLabel = &label
}

func (x *PageShort) SetJobLabel(label string) {
	x.JobLabel = &label
}
