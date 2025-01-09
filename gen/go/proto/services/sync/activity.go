package sync

import (
	"context"
)

func (s *Server) AddActivity(ctx context.Context, req *AddActivityRequest) (*AddActivityResponse, error) {
	resp := &AddActivityResponse{}

	switch d := req.Activity.(type) {
	case *AddActivityRequest_UserOauth2:
		_ = d

		// TODO handle each activity type
	}

	return resp, nil
}
