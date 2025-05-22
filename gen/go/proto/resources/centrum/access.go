package centrum

func (x *Access) IsEmpty() bool {
	return len(x.Jobs) == 0
}

func (x *JobAccess) SetJobLabel(label string) {
	x.JobLabel = &label
}

func (x *JobAccess) GetJobGrade() int32 {
	return x.MinimumGrade
}

func (x *JobAccess) SetJobGrade(grade int32) {
	x.MinimumGrade = grade
}

func (x *JobAccess) SetJobGradeLabel(jobLabel string) {
	x.JobGradeLabel = &jobLabel
}

// pkg/access compatibility

func (x *JobAccess) SetJob(job string) {
	x.Job = job
}

func (x *JobAccess) SetMinimumGrade(grade int32) {
	x.MinimumGrade = grade
}

func (x *JobAccess) SetAccess(access AccessLevel) {
	x.Access = access
}

func (x *UserAccess) GetId() uint64 {
	return 0
}

func (x *UserAccess) GetTargetId() uint64 {
	return 0
}

func (x *UserAccess) GetUserId() int32 {
	return 0
}

func (x *UserAccess) GetAccess() AccessLevel {
	return AccessLevel_ACCESS_LEVEL_BLOCKED
}

func (x *UserAccess) SetUserId(id int32) {}

func (x *UserAccess) SetAccess(access AccessLevel) {}

func (x *QualificationAccess) SetQualificationId(id uint64) {}

func (x *QualificationAccess) SetAccess(access AccessLevel) {}
