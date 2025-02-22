package internet

func (x *PageJobAccess) SetJob(job string) {
	x.Job = job
}

func (x *PageJobAccess) SetJobLabel(label string) {
	x.JobLabel = &label
}

// pkg/access compatibility

func (x *PageJobAccess) GetJobGrade() int32 {
	return x.MinimumGrade
}

func (x *PageJobAccess) SetJobGrade(grade int32) {
	x.MinimumGrade = grade
}

func (x *PageJobAccess) SetJobGradeLabel(label string) {
	x.JobGradeLabel = &label
}

func (x *PageJobAccess) SetMinimumGrade(grade int32) {
	x.MinimumGrade = grade
}

func (x *PageJobAccess) SetAccess(access AccessLevel) {
	x.Access = access
}

func (x *PageUserAccess) SetUserId(id int32) {
	x.UserId = id
}

func (x *PageUserAccess) SetAccess(access AccessLevel) {
	x.Access = access
}
