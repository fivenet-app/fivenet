package mailer

func (x *Email) GetJobGrade() int32 {
	return 0
}

func (x *Email) SetJob(job string) {
	x.Job = &job
}

func (x *Email) SetJobLabel(label string) {}

func (x *Email) SetJobGrade(grade int32) {}

func (x *Email) SetJobGradeLabel(label string) {}

func (x *Access) IsEmpty() bool {
	return len(x.GetJobs()) == 0 && len(x.GetUsers()) == 0 && len(x.GetQualifications()) == 0
}
