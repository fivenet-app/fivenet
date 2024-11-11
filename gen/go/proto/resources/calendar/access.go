package calendar

func (x *CalendarJobAccess) SetJob(job string) {
	x.Job = job
}

func (x *CalendarJobAccess) SetJobLabel(label string) {
	x.JobLabel = &label
}

func (x *CalendarJobAccess) GetJobGrade() int32 {
	return x.MinimumGrade
}

func (x *CalendarJobAccess) SetJobGrade(grade int32) {
	x.MinimumGrade = grade
}

func (x *CalendarJobAccess) SetJobGradeLabel(label string) {
	x.JobGradeLabel = &label
}

// pkg/access compatibility

func (x *CalendarJobAccess) SetMinimumGrade(grade int32) {
	x.MinimumGrade = grade
}

func (x *CalendarJobAccess) SetAccess(access AccessLevel) {
	x.Access = access
}

func (x *CalendarUserAccess) SetUserId(id int32) {
	x.UserId = id
}

func (x *CalendarUserAccess) SetAccess(access AccessLevel) {
	x.Access = access
}
