package users

import "google.golang.org/protobuf/proto"

func (x *Job) Merge(in *Job) *Job {
	if in != nil {
		// Nil grades list to ensure it is updated
		x.Grades = nil

		proto.Merge(x, in)
	}

	return x
}

func (x *Job) GetJob() string {
	return x.Name
}

func (x *Job) SetJobLabel(label string) {
	x.Label = label
}
