package userinfo

import (
	"github.com/fivenet-app/fivenet/v2025/pkg/utils"
	"google.golang.org/protobuf/proto"
)

func (x *UserInfoChanged) GetJob() string {
	return x.GetNewJob()
}

func (x *UserInfoChanged) SetJob(job string) {
	x.NewJob = &job
}

func (x *UserInfoChanged) SetJobLabel(label string) {
	x.NewJobLabel = &label
}

func (x *UserInfoChanged) GetJobGrade() int32 {
	return x.GetNewJobGrade()
}

func (x *UserInfoChanged) SetJobGrade(grade int32) {
	x.NewJobGrade = &grade
}

func (x *UserInfoChanged) SetJobGradeLabel(label string) {
	x.NewJobGradeLabel = &label
}

// Equal returns true if all fields of u and in are equal.
func (u *UserInfo) Equal(in *UserInfo) bool {
	if u == nil || in == nil {
		return false
	}

	if u.GetEnabled() != in.GetEnabled() {
		return false
	}
	if u.GetAccountId() != in.GetAccountId() {
		return false
	}
	if u.GetLicense() != in.GetLicense() {
		return false
	}
	if !utils.EqualInt32Ptr(u.LastChar, in.LastChar) {
		return false
	}
	if u.GetUserId() != in.GetUserId() {
		return false
	}
	if u.GetJob() != in.GetJob() {
		return false
	}
	if u.GetJobGrade() != in.GetJobGrade() {
		return false
	}
	if u.GetGroup() != in.GetGroup() {
		return false
	}
	if u.GetCanBeSuperuser() != in.GetCanBeSuperuser() {
		return false
	}
	if u.GetSuperuser() != in.GetSuperuser() {
		return false
	}
	if !utils.EqualStringPtr(u.OverrideJob, in.OverrideJob) {
		return false
	}
	if !utils.EqualInt32Ptr(u.OverrideJobGrade, in.OverrideJobGrade) {
		return false
	}
	return true
}

// Clone returns a deep copy of the UserInfo struct.
func (u *UserInfo) Clone() *UserInfo {
	//nolint:forcetypeassert // Value type is guaranteed to be *UserInfo
	return proto.Clone(u).(*UserInfo)
}
