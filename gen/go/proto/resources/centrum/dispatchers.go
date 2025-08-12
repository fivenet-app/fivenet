package centrum

import jobs "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/jobs"

func (x *Dispatchers) SetJobLabel(jobLabel string) {
	x.JobLabel = &jobLabel
}

func (x *Dispatchers) Merge(in *Dispatchers) *Dispatchers {
	if len(in.GetDispatchers()) == 0 {
		x.Dispatchers = []*jobs.Colleague{}
	} else {
		x.Dispatchers = in.GetDispatchers()
	}

	return x
}

func (x *Dispatchers) IsEmpty() bool {
	return x == nil || len(x.GetDispatchers()) == 0
}
