package centrum

import (
	"context"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/rector"
	pbcentrum "github.com/fivenet-app/fivenet/gen/go/proto/services/centrum"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/query/fivenet/model"
	errorscentrum "github.com/fivenet-app/fivenet/services/centrum/errors"
)

func (s *Server) GetSettings(ctx context.Context, req *pbcentrum.GetSettingsRequest) (*pbcentrum.GetSettingsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	settings := s.state.GetSettings(ctx, userInfo.Job)

	return &pbcentrum.GetSettingsResponse{
		Settings: settings,
	}, nil
}

func (s *Server) UpdateSettings(ctx context.Context, req *pbcentrum.UpdateSettingsRequest) (*pbcentrum.UpdateSettingsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbcentrum.CentrumService_ServiceDesc.ServiceName,
		Method:  "UpdateSettings",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	settings, err := s.state.UpdateSettingsInDB(ctx, userInfo.Job, req.Settings)
	if err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return &pbcentrum.UpdateSettingsResponse{
		Settings: settings,
	}, nil
}
