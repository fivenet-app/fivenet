package userinfo

import (
	"github.com/galexrt/fivenet/query/fivenet/table"
)

var (
	tUsers           = table.Users.AS("userinfo")
	tFivenetAccounts = table.FivenetAccounts
)

type UserInfo struct {
	Enabled   bool
	AccountId uint64
	UserId    int32

	Group     string
	SuperUser bool

	Job          string
	JobGrade     int32
	OrigJob      string
	OrigJobGrade int32
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
		UserId:    u.UserId,

		Group:     u.Group,
		SuperUser: u.SuperUser,

		Job:          u.Job,
		JobGrade:     u.JobGrade,
		OrigJob:      u.OrigJob,
		OrigJobGrade: u.OrigJobGrade,
	}
}
