package mailer

func (x *EmailShort) GetJobGrade() int32 {
	return 0
}

func (x *EmailShort) SetJob(job string) {
	x.Job = &job
}

func (x *EmailShort) SetJobLabel(label string) {}

func (x *EmailShort) SetJobGrade(grade int32) {}

func (x *EmailShort) SetJobGradeLabel(label string) {}
