package rector

import (
	"cmp"
	"context"
	"slices"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/cron"
	rector "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/rector"
	pbrector "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/rector"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/model"
)

func (s *Server) ListCronjobs(ctx context.Context, req *pbrector.ListCronjobsRequest) (*pbrector.ListCronjobsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	defer s.aud.Log(&model.FivenetAuditLog{
		Service: pbrector.RectorService_ServiceDesc.ServiceName,
		Method:  "ListCronjobs",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_VIEWED),
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

	return &pbrector.ListCronjobsResponse{
		Jobs: jobs,
	}, nil
}
