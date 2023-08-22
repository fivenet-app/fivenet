package centrum

import (
	"context"

	dispatch "github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/query/fivenet/model"
	jet "github.com/go-jet/jet/v2/mysql"
	"google.golang.org/protobuf/proto"
)

func (s *Server) GetSettings(ctx context.Context, req *GetSettingsRequest) (*dispatch.Settings, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	settings := s.getSettings(userInfo.Job)

	return settings, nil
}

func (s *Server) UpdateSettings(ctx context.Context, req *dispatch.Settings) (*dispatch.Settings, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: CentrumService_ServiceDesc.ServiceName,
		Method:  "UpdateSettings",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EVENT_TYPE_ERRORED),
	}
	defer s.auditer.Log(auditEntry, req)

	stmt := tCentrumSettings.
		INSERT(
			tCentrumSettings.Job,
			tCentrumSettings.Enabled,
			tCentrumSettings.Mode,
			tCentrumSettings.FallbackMode,
		).
		VALUES(
			userInfo.Job,
			req.Enabled,
			req.Mode,
			req.FallbackMode,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tCentrumSettings.Job.SET(jet.String(userInfo.Job)),
			tCentrumSettings.Enabled.SET(jet.Bool(req.Enabled)),
			tCentrumSettings.Mode.SET(jet.Int32(int32(req.Mode))),
			tCentrumSettings.FallbackMode.SET(jet.Int32(int32(req.FallbackMode))),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	// Load settings from database so they are updated in the "cache"
	if err := s.loadSettings(ctx, userInfo.Job); err != nil {
		return nil, err
	}

	settings := s.getSettings(userInfo.Job)

	data, err := proto.Marshal(settings)
	if err != nil {
		return nil, err
	}
	s.broadcastToAllUnits(TopicGeneral, TypeGeneralSettings, userInfo.Job, data)

	auditEntry.State = int16(rector.EVENT_TYPE_UPDATED)

	return settings, nil
}
