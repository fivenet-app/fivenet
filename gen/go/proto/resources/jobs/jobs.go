package jobs

func (x *Job) Merge(in *Job) *Job {
	if in != nil {
		x.Name = in.Name
		x.Label = in.Label
		x.Grades = in.Grades
	}

	return x
}

func (x *Job) GetJob() string {
	return x.Name
}

func (x *Job) SetJobLabel(label string) {
	x.Label = label
}
