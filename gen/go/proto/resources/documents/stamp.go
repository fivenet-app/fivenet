package documents

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
