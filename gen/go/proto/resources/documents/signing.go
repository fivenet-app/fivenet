package documents

func (x *SignatureJobAccess) GetJobGrade() int32 {
	return x.GetMinimumGrade()
}

func (x *SignatureJobAccess) SetJobGrade(grade int32) {
	x.MinimumGrade = grade
}

func (x *SignatureJobAccess) SetJobGradeLabel(label string) {
	x.JobGradeLabel = &label
}

func (x *SignatureJobAccess) SetMinimumGrade(grade int32) {
	x.MinimumGrade = grade
}

func (x *SignatureJobAccess) SetAccess(access SignatureAccessLevel) {
	x.Access = access
}

func (x *SignatureUserAccess) SetUserId(userId int32) {
	x.UserId = userId
}

func (x *SignatureUserAccess) SetAccess(access SignatureAccessLevel) {
	x.Access = access
}

func (x *StampJobAccess) GetJobGrade() int32 {
	return x.GetMinimumGrade()
}

func (x *StampJobAccess) SetJobGrade(grade int32) {
	x.MinimumGrade = grade
}

func (x *StampJobAccess) SetJobGradeLabel(label string) {
	x.JobGradeLabel = &label
}

func (x *StampJobAccess) SetMinimumGrade(grade int32) {
	x.MinimumGrade = grade
}

func (x *StampJobAccess) SetAccess(access StampAccessLevel) {
	x.Access = access
}
