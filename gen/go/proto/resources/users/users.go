package users

func (x *User) SetJobLabel(label string) {
	x.JobLabel = &label
}

func (x *User) SetJobGradeLabel(label string) {
	x.JobGradeLabel = &label
}

func (x *UserShort) SetJobLabel(label string) {
	x.JobLabel = &label
}

func (x *UserShort) SetJobGradeLabel(label string) {
	x.JobGradeLabel = &label
}

func (x *User) UserShort() *UserShort {
	return &UserShort{
		UserId:     x.UserId,
		Identifier: x.Identifier,
		Job:        x.Job,
		JobGrade:   x.JobGrade,
		Firstname:  x.Firstname,
		Lastname:   x.Lastname,
	}
}
