package permissions

func (x *Role) GetJobGrade() int32 {
	return x.GetGrade()
}

func (x *Role) SetJob(job string) {
	x.Job = job
}

func (x *Role) SetJobLabel(label string) {
	x.JobLabel = &label
}

func (x *Role) SetJobGrade(grade int32) {
	x.Grade = grade
}

func (x *Role) SetJobGradeLabel(label string) {
	x.JobGradeLabel = &label
}
