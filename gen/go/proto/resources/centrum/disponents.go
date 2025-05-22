package centrum

import jobs "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/jobs"

func (x *Disponents) SetJobLabel(jobLabel string) {
	x.JobLabel = &jobLabel
}

func (x *Disponents) Merge(in *Disponents) *Disponents {
	x.Job = in.Job

	if in.JobLabel == nil {
		in.JobLabel = x.JobLabel
	} else if x.JobLabel != nil {
		*in.JobLabel = *x.JobLabel
	}

	if len(in.Disponents) == 0 {
		x.Disponents = []*jobs.Colleague{}
	} else {
		x.Disponents = in.Disponents
	}

	return x
}
