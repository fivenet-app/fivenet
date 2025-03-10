package qualifications

func (x *QualificationJobAccess) SetJob(job string) {
	x.Job = job
}

func (x *QualificationJobAccess) SetJobLabel(label string) {
	x.JobLabel = &label
}

func (x *QualificationJobAccess) GetJobGrade() int32 {
	return x.MinimumGrade
}

func (x *QualificationJobAccess) SetJobGrade(grade int32) {
	x.MinimumGrade = grade
}

func (x *QualificationJobAccess) SetJobGradeLabel(label string) {
	x.JobGradeLabel = &label
}
