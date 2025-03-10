package centrum

import jobs "github.com/fivenet-app/fivenet/gen/go/proto/resources/jobs"

func (x *Disponents) Merge(in *Disponents) *Disponents {
	if len(in.Disponents) == 0 {
		x.Disponents = []*jobs.Colleague{}
	} else {
		x.Disponents = in.Disponents
	}

	return x
}
