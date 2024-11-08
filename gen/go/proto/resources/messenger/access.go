package messenger

// pkg/access compatibility

func (x *ThreadJobAccess) GetTargetId() uint64 {
	return x.ThreadId
}

func (x *ThreadJobAccess) SetMinimumGrade(grade int32) {
	x.MinimumGrade = grade
}

func (x *ThreadJobAccess) SetAccess(access AccessLevel) {
	x.Access = access
}

func (x *ThreadUserAccess) GetTargetId() uint64 {
	return x.ThreadId
}

func (x *ThreadUserAccess) SetUserId(id int32) {
	x.UserId = id
}

func (x *ThreadUserAccess) SetAccess(access AccessLevel) {
	x.Access = access
}
