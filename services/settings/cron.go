package settings

import (
	"cmp"
	"context"
	"slices"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/cron"
	pbsettings "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/settings"
)

func (s *Server) ListCronjobs(
	ctx context.Context,
	req *pbsettings.ListCronjobsRequest,
) (*pbsettings.ListCronjobsResponse, error) {
	jobs := s.cronState.ListCronjobs(ctx)

	slices.SortFunc(jobs, func(a, b *cron.Cronjob) int {
		if a == nil {
			return 1
		} else if b == nil {
			return -1
		}

		return cmp.Compare(a.GetName(), b.GetName())
	})

	return &pbsettings.ListCronjobsResponse{
		Jobs: jobs,
	}, nil
}
