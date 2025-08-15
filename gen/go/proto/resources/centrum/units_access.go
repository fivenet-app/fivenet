package centrum

func (x *UnitAccess) IsEmpty() bool {
	return len(x.GetJobs()) == 0 && len(x.GetQualifications()) == 0
}

func (x *UnitAccess) ClearQualificationResults() {
	for _, quali := range x.GetQualifications() {
		if quali.GetQualification() != nil && quali.GetQualification().GetResult() != nil {
			quali.Qualification.Result = nil
		}
	}
}

func (x *UnitJobAccess) SetJobLabel(label string) {
	x.JobLabel = &label
}

func (x *UnitJobAccess) GetJobGrade() int32 {
	return x.GetMinimumGrade()
}

func (x *UnitJobAccess) SetJobGrade(grade int32) {
	x.MinimumGrade = grade
}

func (x *UnitJobAccess) SetJobGradeLabel(jobLabel string) {
	x.JobGradeLabel = &jobLabel
}

// pkg/access compatibility

func (x *UnitJobAccess) SetJob(job string) {
	x.Job = job
}

func (x *UnitJobAccess) SetMinimumGrade(grade int32) {
	x.MinimumGrade = grade
}

func (x *UnitJobAccess) SetAccess(access UnitAccessLevel) {
	x.Access = access
}

func (x *UnitUserAccess) GetId() int64 {
	return 0
}

func (x *UnitUserAccess) GetTargetId() int64 {
	return 0
}

func (x *UnitUserAccess) GetUserId() int32 {
	return 0
}

func (x *UnitUserAccess) GetAccess() UnitAccessLevel {
	return UnitAccessLevel_UNIT_ACCESS_LEVEL_JOIN
}

func (x *UnitUserAccess) SetUserId(id int32) {}

func (x *UnitUserAccess) SetAccess(access UnitAccessLevel) {}

func (x *UnitQualificationAccess) SetQualificationId(id int64) {
	x.QualificationId = id
}

func (x *UnitQualificationAccess) SetAccess(access UnitAccessLevel) {
	x.Access = access
}
