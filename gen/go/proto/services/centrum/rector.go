package centrum

import (
	"context"

	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	errorscentrum "github.com/galexrt/fivenet/gen/go/proto/services/centrum/errors"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/grpc/errswrap"
	"github.com/galexrt/fivenet/query/fivenet/model"
)

func (s *Server) GetSettings(ctx context.Context, req *GetSettingsRequest) (*GetSettingsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	settings := s.state.GetSettings(userInfo.Job)

	return &GetSettingsResponse{
		Settings: settings,
	}, nil
}

func (s *Server) UpdateSettings(ctx context.Context, req *UpdateSettingsRequest) (*UpdateSettingsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: CentrumService_ServiceDesc.ServiceName,
		Method:  "UpdateSettings",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.auditer.Log(auditEntry, req)

	settings, err := s.state.UpdateSettingsInDB(ctx, userInfo.Job, req.Settings)
	if err != nil {
		return nil, errswrap.NewError(errorscentrum.ErrFailedQuery, err)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return &UpdateSettingsResponse{
		Settings: settings,
	}, nil
}
