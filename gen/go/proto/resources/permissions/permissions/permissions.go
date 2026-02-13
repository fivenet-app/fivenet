package permissionspermissions

func (x *Role) GetJobGrade() int32 {
	return x.GetGrade()
}

func (x *Role) SetJobGrade(grade int32) {
	x.Grade = grade
}
