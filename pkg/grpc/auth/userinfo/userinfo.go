package userinfo

type UserInfo struct {
	Enabled   bool
	AccountId uint64
	License   string
	LastChar  *int32

	UserId   int32
	Job      string
	JobGrade int32

	Group      string
	CanBeSuper bool
	Superuser  bool

	OverrideJob      *string
	OverrideJobGrade *int32
}

func (u *UserInfo) Equal(in *UserInfo) bool {
	if u == nil || in == nil {
		return false
	}

	if u.Enabled != in.Enabled {
		return false
	}
	if u.AccountId != in.AccountId {
		return false
	}
	if u.License != in.License {
		return false
	}
	if (u.LastChar == nil && in.LastChar != nil) ||
		(u.LastChar != nil && in.LastChar == nil) ||
		(u.LastChar != nil && in.LastChar != nil && *u.LastChar != *in.LastChar) {
		return false
	}

	if u.UserId != in.UserId {
		return false
	}
	if u.Job != in.Job {
		return false
	}
	if u.JobGrade != in.JobGrade {
		return false
	}

	if u.Group != in.Group {
		return false
	}
	if u.CanBeSuper != in.CanBeSuper {
		return false
	}
	if u.Superuser != in.Superuser {
		return false
	}

	if (u.OverrideJob == nil && in.OverrideJob != nil) ||
		(u.OverrideJob != nil && in.OverrideJob == nil) ||
		(u.OverrideJob != nil && in.OverrideJob != nil && *u.OverrideJob != *in.OverrideJob) {
		return false
	}
	if (u.OverrideJobGrade == nil && in.OverrideJobGrade != nil) ||
		(u.OverrideJobGrade != nil && in.OverrideJobGrade == nil) ||
		(u.OverrideJobGrade != nil && in.OverrideJobGrade != nil && *u.OverrideJobGrade != *in.OverrideJobGrade) {
		return false
	}

	return true
}

func (u *UserInfo) Clone() UserInfo {
	return UserInfo{
		Enabled:   u.Enabled,
		AccountId: u.AccountId,
		License:   u.License,

		LastChar:   u.LastChar,
		Group:      u.Group,
		CanBeSuper: u.CanBeSuper,
		Superuser:  u.Superuser,

		UserId:   u.UserId,
		Job:      u.Job,
		JobGrade: u.JobGrade,

		OverrideJob:      u.OverrideJob,
		OverrideJobGrade: u.OverrideJobGrade,
	}
}
