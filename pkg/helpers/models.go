package helpers

import (
	"github.com/galexrt/arpanet/model"
	"github.com/galexrt/arpanet/proto/common"
)

func ConvertModelUserToCommonCharacter(user *model.User) *common.Character {
	licenses := make([]*common.License, len(user.UserLicenses))
	for i := 0; i < len(user.UserLicenses); i++ {
		licenses[i] = &common.License{
			Name: string(user.UserLicenses[i].Type),
		}
	}

	return &common.Character{
		Identifier:  user.Identifier,
		Firstname:   user.Firstname,
		Lastname:    user.Lastname,
		Dateofbirth: user.Dateofbirth,
		Sex:         string(user.Sex),
		Height:      user.Height,
		Job:         user.Job,
		JobGrade:    int32(user.JobGrade),
		Visum:       int64(user.Visum),
		Playtime:    int64(user.Playtime),
		Licenses:    licenses,
	}
}
