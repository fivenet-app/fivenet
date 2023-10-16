package centrum

import (
	"context"

	dispatch "github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	errorscentrum "github.com/galexrt/fivenet/gen/go/proto/services/centrum/errors"
	eventscentrum "github.com/galexrt/fivenet/gen/go/proto/services/centrum/events"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"google.golang.org/protobuf/proto"
)

var tCentrumSettings = table.FivenetCentrumSettings

func (s *Server) GetSettings(ctx context.Context, req *GetSettingsRequest) (*dispatch.Settings, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	settings := s.state.GetSettings(userInfo.Job)

	return settings, nil
}

func (s *Server) UpdateSettings(ctx context.Context, req *dispatch.Settings) (*dispatch.Settings, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: CentrumService_ServiceDesc.ServiceName,
		Method:  "UpdateSettings",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
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
		return nil, errorscentrum.ErrFailedQuery
	}

	// Load settings from database so they are updated in the "cache"
	if err := s.state.LoadSettings(ctx, userInfo.Job); err != nil {
		return nil, errorscentrum.ErrFailedQuery
	}

	settings := s.state.GetSettings(userInfo.Job)

	data, err := proto.Marshal(settings)
	if err != nil {
		return nil, errorscentrum.ErrFailedQuery
	}
	s.events.JS.PublishAsync(eventscentrum.BuildSubject(eventscentrum.TopicGeneral, eventscentrum.TypeGeneralSettings, userInfo.Job, 0), data)

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return settings, nil
}
