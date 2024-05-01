package userinfo

import (
	"github.com/fivenet-app/fivenet/query/fivenet/table"
)

var (
	tUsers           = table.Users.AS("userinfo")
	tFivenetAccounts = table.FivenetAccounts
)

type UserInfo struct {
	Enabled   bool
	AccountId uint64
	License   string

	UserId   int32
	Job      string
	JobGrade int32

	Group      string
	CanBeSuper bool
	SuperUser  bool

	OverrideJob      *string
	OverrideJobGrade *int32
}

func (u *UserInfo) Equal(in *UserInfo) bool {
	if u == nil || in == nil {
		return false
	}

	return *u == *in
}

func (u *UserInfo) Clone() UserInfo {
	return UserInfo{
		Enabled:   u.Enabled,
		AccountId: u.AccountId,
		License:   u.License,

		Group:      u.Group,
		CanBeSuper: u.CanBeSuper,
		SuperUser:  u.SuperUser,

		UserId:   u.UserId,
		Job:      u.Job,
		JobGrade: u.JobGrade,

		OverrideJob:      u.OverrideJob,
		OverrideJobGrade: u.OverrideJobGrade,
	}
}
