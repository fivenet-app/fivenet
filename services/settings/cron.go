package settings

import (
	"cmp"
	"context"
	"slices"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/cron"
	pbsettings "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/settings"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
)

func (s *Server) ListCronjobs(ctx context.Context, req *pbsettings.ListCronjobsRequest) (*pbsettings.ListCronjobsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	defer s.aud.Log(&audit.AuditEntry{
		Service: pbsettings.SettingsService_ServiceDesc.ServiceName,
		Method:  "ListCronjobs",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_VIEWED,
	}, req)

	jobs := s.cronState.ListCronjobs(ctx)

	slices.SortFunc(jobs, func(a, b *cron.Cronjob) int {
		if a == nil {
			return 1
		} else if b == nil {
			return -1
		}

		return cmp.Compare(a.Name, b.Name)
	})

	return &pbsettings.ListCronjobsResponse{
		Jobs: jobs,
	}, nil
}
