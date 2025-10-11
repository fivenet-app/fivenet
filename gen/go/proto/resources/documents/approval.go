package documents

func (x *ApprovalTask) GetJobGrade() int32 {
	return x.GetMinimumGrade()
}

func (x *ApprovalTask) SetJobGrade(grade int32) {
	x.MinimumGrade = &grade
}

func (x *ApprovalTask) SetJobGradeLabel(label string) {
	x.JobGradeLabel = &label
}

func (x *ApprovalTask) SetJob(job string) {
	x.Job = &job
}

func (x *ApprovalTask) SetJobLabel(label string) {
	x.JobLabel = &label
}
