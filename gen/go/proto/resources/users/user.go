package users

import (
	jobscolleagues "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/colleagues"
)

func (x *User) UserShort() *UserShort {
	return &UserShort{
		UserId:        x.GetUserId(),
		Identifier:    x.Identifier,
		Job:           x.GetJob(),
		JobGrade:      x.GetJobGrade(),
		Firstname:     x.GetFirstname(),
		Lastname:      x.GetLastname(),
		Dateofbirth:   x.GetDateofbirth(),
		PhoneNumber:   x.PhoneNumber,
		JobLabel:      x.JobLabel,
		JobGradeLabel: x.JobGradeLabel,
	}
}

func (x *User) Colleague() *jobscolleagues.Colleague {
	return &jobscolleagues.Colleague{
		UserId:        x.GetUserId(),
		Identifier:    x.Identifier,
		Job:           x.GetJob(),
		JobGrade:      x.GetJobGrade(),
		Firstname:     x.GetFirstname(),
		Lastname:      x.GetLastname(),
		Dateofbirth:   x.GetDateofbirth(),
		PhoneNumber:   x.PhoneNumber,
		JobLabel:      x.JobLabel,
		JobGradeLabel: x.JobGradeLabel,
	}
}
