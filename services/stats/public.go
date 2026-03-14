package stats

import (
	"context"

	pbstats "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/stats"
)

func (s *Server) GetPublicStats(
	ctx context.Context,
	req *pbstats.GetPublicStatsRequest,
) (*pbstats.GetPublicStatsResponse, error) {
	stats := s.worker.GetStats()
	if stats == nil {
		return &pbstats.GetPublicStatsResponse{}, nil
	}

	return &pbstats.GetPublicStatsResponse{
		Stats: *stats,
	}, nil
}
