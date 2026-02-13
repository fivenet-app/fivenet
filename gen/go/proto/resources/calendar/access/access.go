package calendaraccess

func (x *CalendarJobAccess) GetJobGrade() int32 {
	return x.GetMinimumGrade()
}

func (x *CalendarJobAccess) SetJobGrade(grade int32) {
	x.MinimumGrade = grade
}
