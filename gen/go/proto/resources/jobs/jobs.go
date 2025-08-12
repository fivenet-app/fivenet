package jobs

func (x *Job) Merge(in *Job) *Job {
	if in != nil {
		x.Name = in.GetName()
		x.Label = in.GetLabel()
		x.Grades = in.GetGrades()
	}

	return x
}

func (x *Job) GetJob() string {
	return x.GetName()
}

func (x *Job) SetJobLabel(label string) {
	x.Label = label
}
