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
		// TODO

	case *sync.AddActivity_UserProps:
		// TODO

	case *sync.AddActivity_JobsUserActivity:
		// TODO

	case *sync.AddActivity_JobsTimeclock:
		// TODO
	}

	return resp, nil
}
