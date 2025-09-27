package documents

func (x *ApprovalJobAccess) GetJobGrade() int32 {
	return x.GetMinimumGrade()
}

func (x *ApprovalJobAccess) SetJobGrade(grade int32) {
	x.MinimumGrade = grade
}

func (x *ApprovalJobAccess) SetJobGradeLabel(label string) {
	x.JobGradeLabel = &label
}

func (x *ApprovalJobAccess) SetMinimumGrade(grade int32) {
	x.MinimumGrade = grade
}

func (x *ApprovalJobAccess) SetAccess(access ApprovalAccessLevel) {
	x.Access = access
}
