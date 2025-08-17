package centrum

import (
	"context"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/jobs"
	pbcentrum "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/centrum"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
)

func (s *Server) ListDispatchTargetJobs(
	ctx context.Context,
	req *pbcentrum.ListDispatchTargetJobsRequest,
) (*pbcentrum.ListDispatchTargetJobsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	resp := &pbcentrum.ListDispatchTargetJobsResponse{
		Jobs: []*jobs.Job{},
	}

	// TODO retrieve public jobs if the dispatch center is enabled

	// Add user's job to the job list if dispatch center is enabled
	settings, err := s.settings.Get(ctx, userInfo.GetJob())
	if err != nil {
		return nil, err
	}
	if settings.Enabled {
		if j := s.enricher.GetJobByName(userInfo.GetJob()); j != nil {
			resp.Jobs = append(resp.Jobs, j)
		}
	}

	return resp, nil
}
