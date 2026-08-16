package usershort

import jobscolleagues "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/colleagues"

func (x *UserShort) Colleague() *jobscolleagues.Colleague {
	var profilePicture *string
	if x.ProfilePicture != nil {
		p := x.GetProfilePicture()
		profilePicture = &p
	}
	var profilePictureFileId *int64
	if x.ProfilePictureFileId != nil {
		p := x.GetProfilePictureFileId()
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
