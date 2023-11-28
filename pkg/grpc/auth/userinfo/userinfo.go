package userinfo

import (
	"github.com/galexrt/fivenet/query/fivenet/table"
)

var (
	tUsers           = table.Users.AS("userinfo")
	tFivenetAccounts = table.FivenetAccounts
)

type UserInfo struct {
	Enabled bool
	AccId   uint64
	UserId  int32

	Job          string
	JobGrade     int32
	OrigJob      string
	OrigJobGrade int32

	Group     string
	SuperUser bool
}

func (u *UserInfo) Equal(in *UserInfo) bool {
	if u == nil || in == nil {
		return false
	}

	return *u == *in
}
