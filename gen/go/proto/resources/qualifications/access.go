package qualifications

// pkg/access compatibility

func (x *QualificationJobAccess) SetMinimumGrade(grade int32) {
	x.MinimumGrade = grade
}

func (x *QualificationJobAccess) SetAccess(access AccessLevel) {
	x.Access = access
}

func (x *QualificationUserAccess) GetId() int64 {
	return 0
}

func (x *QualificationUserAccess) GetTargetId() int64 {
	return 0
}

func (x *QualificationUserAccess) GetAccess() AccessLevel {
	return AccessLevel_ACCESS_LEVEL_UNSPECIFIED
}

func (x *QualificationUserAccess) SetAccess(access AccessLevel) {}

func (x *QualificationUserAccess) GetUserId() int32 {
	return 0
}

func (x *QualificationUserAccess) SetUserId(userId int32) {}
