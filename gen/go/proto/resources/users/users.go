package users

import "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/jobs"

func (x *User) SetJob(job string) {
	x.Job = job
}

func (x *User) SetJobLabel(label string) {
	x.JobLabel = &label
}

func (x *User) SetJobGrade(grade int32) {
	x.JobGrade = grade
}

func (x *User) SetJobGradeLabel(label string) {
	x.JobGradeLabel = &label
}

func (x *UserShort) SetJob(job string) {
	x.Job = job
}

func (x *UserShort) SetJobLabel(label string) {
	x.JobLabel = &label
}

func (x *UserShort) SetJobGrade(grade int32) {
	x.JobGrade = grade
}

func (x *UserShort) SetJobGradeLabel(label string) {
	x.JobGradeLabel = &label
}

func (x *User) UserShort() *UserShort {
	return &UserShort{
		UserId:        x.UserId,
		Identifier:    x.Identifier,
		Job:           x.Job,
		JobGrade:      x.JobGrade,
		Firstname:     x.Firstname,
		Lastname:      x.Lastname,
		Dateofbirth:   x.Dateofbirth,
		PhoneNumber:   x.PhoneNumber,
		JobLabel:      x.JobLabel,
		JobGradeLabel: x.JobGradeLabel,
	}
}

func (x *User) Colleague() *jobs.Colleague {
	return &jobs.Colleague{
		UserId:        x.UserId,
		Identifier:    x.Identifier,
		Job:           x.Job,
		JobGrade:      x.JobGrade,
		Firstname:     x.Firstname,
		Lastname:      x.Lastname,
		Dateofbirth:   x.Dateofbirth,
		PhoneNumber:   x.PhoneNumber,
		JobLabel:      x.JobLabel,
		JobGradeLabel: x.JobGradeLabel,
	}
}
