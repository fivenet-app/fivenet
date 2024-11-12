package mailer

// pkg/access compatibility

func (x *ThreadJobAccess) GetId() uint64 {
	return 0
}

func (x *ThreadJobAccess) GetTargetId() uint64 {
	return 0
}

func (x *ThreadJobAccess) GetJob() string {
	return ""
}

func (x *ThreadJobAccess) GetMinimumGrade() int32 {
	return 0
}

func (x *ThreadJobAccess) SetMinimumGrade(grade int32) {}

func (x *ThreadJobAccess) GetAccess() AccessLevel {
	return AccessLevel_ACCESS_LEVEL_UNSPECIFIED
}

func (x *ThreadJobAccess) SetAccess(access AccessLevel) {}

func (x *ThreadUserAccess) SetUserId(id int32) {
	x.UserId = id
}

func (x *ThreadUserAccess) SetAccess(access AccessLevel) {
	x.Access = access
}
