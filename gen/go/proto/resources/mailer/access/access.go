package maileraccess

func (x *Access) IsEmpty() bool {
	return len(x.GetJobs()) == 0 && len(x.GetUsers()) == 0 && len(x.GetQualifications()) == 0
}
