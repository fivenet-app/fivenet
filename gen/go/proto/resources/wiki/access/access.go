package wikiaccess

func (x *PageJobAccess) GetJobGrade() int32 {
	return x.GetMinimumGrade()
}

func (x *PageJobAccess) SetJobGrade(grade int32) {
	x.MinimumGrade = grade
}

func (x *PageAccess) IsEmpty() bool {
	return len(x.GetJobs()) == 0 && len(x.GetUsers()) == 0
}
