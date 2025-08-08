package centrum

import (
	"context"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/jobs"
	pbcentrum "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/centrum"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
)

func (s *Server) ListDispatchTargetJobs(ctx context.Context, req *pbcentrum.ListDispatchTargetJobsRequest) (*pbcentrum.ListDispatchTargetJobsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	resp := &pbcentrum.ListDispatchTargetJobsResponse{
		Jobs: []*jobs.Job{},
	}

	if j := s.enricher.GetJobByName(userInfo.Job); j != nil {
		resp.Jobs = append(resp.Jobs, j)
	}

	return resp, nil
}
