package sync

import "context"

func (s *Server) AddActivity(ctx context.Context, req *AddActivityRequest) (*AddActivityResponse, error) {
	resp := &AddActivityResponse{}

	switch req.Activity.(type) {
	case *AddActivityRequest_UserActivity:
		// TODO

	case *AddActivityRequest_UserProps:
		// TODO

	case *AddActivityRequest_JobsUserActivity:
		// TODO

	case *AddActivityRequest_JobsTimeclock:
		// TODO
	}

	return resp, nil
}
