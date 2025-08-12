package wiki

func (x *PageJobAccess) SetJob(job string) {
	x.Job = job
}

func (x *PageJobAccess) SetJobLabel(label string) {
	x.JobLabel = &label
}

func (x *PageJobAccess) GetJobGrade() int32 {
	return x.GetMinimumGrade()
}

func (x *PageJobAccess) SetJobGrade(grade int32) {
	x.MinimumGrade = grade
}

func (x *PageJobAccess) SetJobGradeLabel(label string) {
	x.JobGradeLabel = &label
}

// pkg/access compatibility

func (x *PageJobAccess) SetMinimumGrade(grade int32) {
	x.MinimumGrade = grade
}

func (x *PageJobAccess) SetAccess(access AccessLevel) {
	x.Access = access
}

func (x *PageUserAccess) SetAccess(access AccessLevel) {}

func (x *PageUserAccess) SetUserId(userId int32) {
	x.UserId = userId
}
