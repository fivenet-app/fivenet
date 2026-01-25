package centrumdispatchers

import (
	jobscolleagues "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/colleagues"
)

func (x *Dispatchers) SetJobLabel(jobLabel string) {
	x.JobLabel = &jobLabel
}

func (x *Dispatchers) Merge(in *Dispatchers) *Dispatchers {
	if len(in.GetDispatchers()) == 0 {
		x.Dispatchers = []*jobscolleagues.Colleague{}
	} else {
		x.Dispatchers = in.GetDispatchers()
	}

	return x
}

func (x *Dispatchers) IsEmpty() bool {
	return x == nil || len(x.GetDispatchers()) == 0
}
