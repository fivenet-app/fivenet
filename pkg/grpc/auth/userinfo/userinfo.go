package userinfo

// UserInfo holds information about a user and their account.
type UserInfo struct {
	Enabled   bool
	AccountId uint64
	License   string // License is a string, not a pointer
	LastChar  *int32

	UserId   int32
	Job      string
	JobGrade int32

	Group          string
	CanBeSuperuser bool
	Superuser      bool

	OverrideJob      *string
	OverrideJobGrade *int32
}

// Equal returns true if all fields of u and in are equal.
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
	if !equalInt32Ptr(u.LastChar, in.LastChar) {
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
	if u.CanBeSuperuser != in.CanBeSuperuser {
		return false
	}
	if u.Superuser != in.Superuser {
		return false
	}
	if !equalStringPtr(u.OverrideJob, in.OverrideJob) {
		return false
	}
	if !equalInt32Ptr(u.OverrideJobGrade, in.OverrideJobGrade) {
		return false
	}
	return true
}

// equalInt32Ptr compares two *int32 pointers for equality.
func equalInt32Ptr(a, b *int32) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

// equalStringPtr compares two *string pointers for equality.
func equalStringPtr(a, b *string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

// Clone returns a deep copy of the UserInfo struct.
func (u *UserInfo) Clone() UserInfo {
	clone := UserInfo{
		Enabled:        u.Enabled,
		AccountId:      u.AccountId,
		License:        u.License,
		UserId:         u.UserId,
		Job:            u.Job,
		JobGrade:       u.JobGrade,
		Group:          u.Group,
		CanBeSuperuser: u.CanBeSuperuser,
		Superuser:      u.Superuser,
	}
	if u.LastChar != nil {
		val := *u.LastChar
		clone.LastChar = &val
	}
	if u.OverrideJob != nil {
		val := *u.OverrideJob
		clone.OverrideJob = &val
	}
	if u.OverrideJobGrade != nil {
		val := *u.OverrideJobGrade
		clone.OverrideJobGrade = &val
	}
	return clone
}
