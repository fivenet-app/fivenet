package sync

import (
	"context"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/sync"
)

func (s *Server) AddActivity(ctx context.Context, req *AddActivityRequest) (*AddActivityResponse, error) {
	resp := &AddActivityResponse{}

	switch d := req.Activity.Activity.(type) {
	case *sync.AddActivity_UserActivity:
		_ = d
		if d.UserActivity.UserProps != nil {
			// Set user props
		}

		// TODO handle traffic points

	case *sync.AddActivity_JobsUserActivity:
		if d.JobsUserActivity.JobsUserProps != nil {
			// Set user props
		}

		// TODO

	case *sync.AddActivity_JobsTimeclock:
		// TODO

	}

	return resp, nil
}
