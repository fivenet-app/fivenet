package centrum

func (x *JobAccessEntry) SetJobLabel(label string) {
	if x == nil {
		return
	}
	x.JobLabel = &label
}
