package unitsaccess

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

func (x *UnitJobAccess) GetJobGrade() int32 {
	return x.GetMinimumGrade()
}

func (x *UnitJobAccess) SetJobGrade(grade int32) {
	x.MinimumGrade = grade
}

// pkg/access compatibility

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
