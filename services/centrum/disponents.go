package centrum

import (
	"context"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/rector"
	pbcentrum "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/centrum"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/model"
)

func (s *Server) TakeControl(ctx context.Context, req *pbcentrum.TakeControlRequest) (*pbcentrum.TakeControlResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbcentrum.CentrumService_ServiceDesc.ServiceName,
		Method:  "TakeControl",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	if err := s.state.DisponentSignOn(ctx, userInfo.Job, userInfo.UserId, req.Signon); err != nil {
		return nil, err
	}

	if req.Signon {
		auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)
	} else {
		auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)
	}

	return &pbcentrum.TakeControlResponse{}, nil
}
