package rector

import (
	"context"

	pbrector "github.com/fivenet-app/fivenet/gen/go/proto/services/rector"
)

func (s *Server) ListCronjobs(ctx context.Context, req *pbrector.ListCronjobsRequest) (*pbrector.ListCronjobsResponse, error) {
	jobs := s.cronState.ListCronjobs(ctx)

	return &pbrector.ListCronjobsResponse{
		Jobs: jobs,
	}, nil
}
