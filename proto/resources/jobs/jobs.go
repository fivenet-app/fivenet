package jobs

const DefaultLivemapMarkerColor = "5C7AFF"

func (x *Job) GetJob() string {
	return x.Name
}

func (x *Job) SetJobLabel(label string) {
	x.Label = label
}

func (x *JobProps) Default(job string) {
	if x.Job == "" {
		x.Job = job
	}
	if x.Theme == "" {
		x.Theme = "default"
	}
	if x.LivemapMarkerColor == "" {
		x.LivemapMarkerColor = DefaultLivemapMarkerColor
	}
}
