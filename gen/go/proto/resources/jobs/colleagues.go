package jobs

func (x *Colleague) SetJob(job string) {
	x.Job = job
}

func (x *Colleague) SetJobLabel(label string) {
	x.JobLabel = &label
}

func (x *Colleague) SetJobGrade(grade int32) {
	x.JobGrade = grade
}

func (x *Colleague) SetJobGradeLabel(label string) {
	x.JobGradeLabel = &label
}
