package documents

func (x *SignatureTask) GetJobGrade() int32 {
	return x.GetMinimumGrade()
}

func (x *SignatureTask) SetJobGrade(grade int32) {
	x.MinimumGrade = &grade
}

func (x *SignatureTask) SetJobGradeLabel(label string) {
	x.JobGradeLabel = &label
}

func (x *SignatureTask) SetJob(job string) {
	x.Job = &job
}

func (x *SignatureTask) SetJobLabel(label string) {
	x.JobLabel = &label
}
