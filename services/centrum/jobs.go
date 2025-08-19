package centrum

import (
	"context"
	"slices"

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
		Jobs: s.settings.GetPublicJobs(),
	}

	// Add user's job to the job list if dispatch center is enabled (and it's not already in the list)
	if slices.IndexFunc(resp.Jobs, func(j *jobs.Job) bool {
		return j.GetName() == userInfo.GetJob()
	}) == -1 {
		settings, err := s.settings.Get(ctx, userInfo.GetJob())
		if err != nil {
			return nil, err
		}

		if settings.Enabled {
			if j := s.enricher.GetJobByName(userInfo.GetJob()); j != nil {
				resp.Jobs = append(resp.Jobs, j)
			}
		}
	}

	return resp, nil
}
