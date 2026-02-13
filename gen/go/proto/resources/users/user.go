package users

import (
	jobscolleagues "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/colleagues"
)

func (x *User) UserShort() *UserShort {
	var profilePicture *string
	if x.ProfilePicture != nil {
		p := *x.ProfilePicture
		profilePicture = &p
	}
	var profilePictureFileId *int64
	if x.ProfilePictureFileId != nil {
		p := *x.ProfilePictureFileId
		profilePictureFileId = &p
	}

	return &UserShort{
		UserId:               x.GetUserId(),
		Job:                  x.GetJob(),
		JobGrade:             x.GetJobGrade(),
		Firstname:            x.GetFirstname(),
		Lastname:             x.GetLastname(),
		Dateofbirth:          x.GetDateofbirth(),
		PhoneNumber:          x.PhoneNumber,
		JobLabel:             x.JobLabel,
		JobGradeLabel:        x.JobGradeLabel,
		ProfilePicture:       profilePicture,
		ProfilePictureFileId: profilePictureFileId,
	}
}

func (x *User) Colleague() *jobscolleagues.Colleague {
	var profilePicture *string
	if x.ProfilePicture != nil {
		p := *x.ProfilePicture
		profilePicture = &p
	}
	var profilePictureFileId *int64
	if x.ProfilePictureFileId != nil {
		p := *x.ProfilePictureFileId
		profilePictureFileId = &p
	}

	return &jobscolleagues.Colleague{
		UserId:               x.GetUserId(),
		Job:                  x.GetJob(),
		JobGrade:             x.GetJobGrade(),
		Firstname:            x.GetFirstname(),
		Lastname:             x.GetLastname(),
		Dateofbirth:          x.GetDateofbirth(),
		PhoneNumber:          x.PhoneNumber,
		JobLabel:             x.JobLabel,
		JobGradeLabel:        x.JobGradeLabel,
		ProfilePicture:       profilePicture,
		ProfilePictureFileId: profilePictureFileId,
	}
}
