package centrumaccess

func (x *CentrumJobAccess) SetJobLabel(label string) {
	if x == nil {
		return
	}
	x.JobLabel = &label
}
