package usershort

func (x *UserShort) SetJob(job string) {
	x.Job = job
}

func (x *UserShort) SetJobLabel(label string) {
	x.JobLabel = &label
}

func (x *UserShort) SetJobGrade(grade int32) {
	x.JobGrade = grade
}

func (x *UserShort) SetJobGradeLabel(label string) {
	x.JobGradeLabel = &label
}
