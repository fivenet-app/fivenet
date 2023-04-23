package jobs

const DefaultLivemapMarkerColor = "5C7AFF"

func (x *Job) GetJob() string {
	return x.Name
}

func (x *Job) SetJobLabel(label string) {
	x.Label = label
}
