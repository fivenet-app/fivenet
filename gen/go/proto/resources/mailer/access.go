package mailer

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

func (x *UserAccess) SetUserId(id int32) {
	x.UserId = id
}

func (x *UserAccess) SetAccess(access AccessLevel) {
	x.Access = access
}

func (x *QualificationAccess) SetQualificationId(id uint64) {
	x.QualificationId = id
}

func (x *QualificationAccess) SetAccess(access AccessLevel) {
	x.Access = access
}
