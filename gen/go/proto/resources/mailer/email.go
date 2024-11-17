package mailer

func (x *Email) ToEmailShort() *EmailShort {
	return &EmailShort{
		Id:        x.Id,
		CreatedAt: x.CreatedAt,
		UpdatedAt: x.UpdatedAt,
		DeletedAt: x.DeletedAt,
		Disabled:  x.Disabled,
		Job:       x.Job,
		UserId:    x.UserId,
		User:      x.User,
		Email:     x.Email,
		Label:     x.Label,
		Internal:  x.Internal,
	}
}

func (x *EmailShort) GetJobGrade() int32 {
	return 0
}

func (x *EmailShort) SetJob(job string) {
	x.Job = &job
}

func (x *EmailShort) SetJobLabel(label string) {}

func (x *EmailShort) SetJobGrade(grade int32) {}

func (x *EmailShort) SetJobGradeLabel(label string) {}
