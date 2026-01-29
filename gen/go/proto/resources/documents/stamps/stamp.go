package documentsstamps

func (x *StampJobAccess) GetJobGrade() int32 {
	return x.GetMinimumGrade()
}

func (x *StampJobAccess) SetJobGrade(grade int32) {
	x.MinimumGrade = grade
}
