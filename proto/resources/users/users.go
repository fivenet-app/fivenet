package users

type IJobInfo interface {
	GetJob() string
	GetJobGrade() int32
	SetJobLabel(label string)
	SetJobGradeLabel(label string)
}

func (x *User) SetJobLabel(label string) {
	x.JobLabel = label
}

func (x *User) SetJobGradeLabel(label string) {
	x.JobGradeLabel = label
}

func (x *UserShort) SetJobLabel(label string) {
	x.JobLabel = label
}

func (x *UserShort) SetJobGradeLabel(label string) {
	x.JobGradeLabel = label
}
