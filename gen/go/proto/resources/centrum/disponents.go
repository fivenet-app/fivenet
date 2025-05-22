package centrum

import jobs "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/jobs"

func (x *Disponents) SetJobLabel(jobLabel string) {
	x.JobLabel = &jobLabel
}

func (x *Disponents) Merge(in *Disponents) *Disponents {
	x.Job = in.Job

	if in.JobLabel != nil {
		if x.JobLabel == nil {
			x.JobLabel = in.JobLabel
		} else {
			*x.JobLabel = *in.JobLabel
		}
	}

	if len(in.Disponents) == 0 {
		x.Disponents = []*jobs.Colleague{}
	} else {
		x.Disponents = in.Disponents
	}

	return x
}
