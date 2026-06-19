package resourcesaccess

// GetJobGrade keeps shared job access entries compatible with common job-info enrichment.
func (x *JobAccess) GetJobGrade() int32 {
	return x.GetMinimumGrade()
}

func (x *JobAccess) SetJobGrade(grade int32) {
	x.MinimumGrade = grade
}

func (x *Access) IsEmpty() bool {
	return len(x.GetJobs()) == 0 && len(x.GetUsers()) == 0 && len(x.GetQualifications()) == 0
}
