package rector

import (
	"context"
	"errors"

	rector "github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"google.golang.org/protobuf/types/known/durationpb"
)

var (
	tConfig = table.FivenetConfig
)

func (s *Server) getAppConfig(ctx context.Context) (*rector.AppConfig, error) {
	stmt := tConfig.
		SELECT(
			tConfig.AppConfig,
		).
		FROM(tConfig).
		LIMIT(1)

	dest := &rector.AppConfig{}
	if err := stmt.QueryContext(ctx, s.db, dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	dest.Default()

	return dest, nil
}

func (s *Server) GetAppConfig(ctx context.Context, req *GetAppConfigRequest) (*GetAppConfigResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: RectorConfigService_ServiceDesc.ServiceName,
		Method:  "GetAppConfig",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	config, err := s.getAppConfig(ctx)
	if err != nil {
		return nil, err
	}

	config.Auth.SignupEnabled = s.cfg.Game.Auth.SignupEnabled

	config.Website.Links.Imprint = s.cfg.HTTP.Links.Imprint
	config.Website.Links.PrivacyPolicy = s.cfg.HTTP.Links.PrivacyPolicy

	config.JobInfo.HiddenJobs = s.cfg.Game.HiddenJobs
	config.JobInfo.PublicJobs = s.cfg.Game.PublicJobs

	config.UserTracker.RefreshTime = durationpb.New(s.cfg.Game.Livemap.RefreshTime)
	config.UserTracker.DbRefreshTime = durationpb.New(s.cfg.Game.Livemap.DBRefreshTime)
	config.UserTracker.LivemapJobs = s.cfg.Game.Livemap.Jobs
	config.UserTracker.TimeclockJobs = []string{}

	config.Discord.Enabled = s.cfg.Discord.Enabled

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_VIEWED)

	return &GetAppConfigResponse{
		Config: config,
	}, nil
}

func (s *Server) UpdateAppConfig(ctx context.Context, req *UpdateAppConfigRequest) (*UpdateAppConfigResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: RectorConfigService_ServiceDesc.ServiceName,
		Method:  "UpdateAppConfig",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	stmt := tConfig.
		INSERT(
			tConfig.Key,
			tConfig.AppConfig,
		).
		VALUES(
			1,
			req.Config,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tConfig.AppConfig.SET(jet.RawString("VALUES(`app_config`)")),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	config, err := s.getAppConfig(ctx)
	if err != nil {
		return nil, err
	}

	return &UpdateAppConfigResponse{
		Config: config,
	}, nil
}
