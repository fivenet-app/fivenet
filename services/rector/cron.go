package rector

import (
	"cmp"
	"context"
	"slices"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/cron"
	pbrector "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/rector"
)

func (s *Server) ListCronjobs(ctx context.Context, req *pbrector.ListCronjobsRequest) (*pbrector.ListCronjobsResponse, error) {
	jobs := s.cronState.ListCronjobs(ctx)

	slices.SortFunc(jobs, func(a, b *cron.Cronjob) int {
		if a == nil {
			return 1
		} else if b == nil {
			return -1
		}

		return cmp.Compare(a.Name, b.Name)
	})

	return &pbrector.ListCronjobsResponse{
		Jobs: jobs,
	}, nil
}
