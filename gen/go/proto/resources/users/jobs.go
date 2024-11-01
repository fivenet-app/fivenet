package users

func (x *Job) Merge(in *Job) *Job {
	if in != nil {
		x = in
	}

	return x
}

func (x *Job) GetJob() string {
	return x.Name
}

func (x *Job) SetJobLabel(label string) {
	x.Label = label
}
