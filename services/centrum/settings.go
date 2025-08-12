package centrum

import (
	"context"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	pbcentrum "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/centrum"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	errorscentrum "github.com/fivenet-app/fivenet/v2025/services/centrum/errors"
)

func (s *Server) GetSettings(
	ctx context.Context,
	req *pbcentrum.GetSettingsRequest,
) (*pbcentrum.GetSettingsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	settings, err := s.settings.Get(ctx, userInfo.GetJob())
	if err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	return &pbcentrum.GetSettingsResponse{
		Settings: settings,
	}, nil
}

func (s *Server) UpdateSettings(
	ctx context.Context,
	req *pbcentrum.UpdateSettingsRequest,
) (*pbcentrum.UpdateSettingsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbcentrum.CentrumService_ServiceDesc.ServiceName,
		Method:  "UpdateSettings",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	settings, err := s.settings.Update(ctx, userInfo.GetJob(), req.GetSettings())
	if err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED

	return &pbcentrum.UpdateSettingsResponse{
		Settings: settings,
	}, nil
}
