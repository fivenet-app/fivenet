package settings

import (
	"context"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/clientconfig"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/notifications"
	pbsettings "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/settings"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
)

var tConfig = table.FivenetConfig

func (s *Server) GetAppConfig(
	ctx context.Context,
	req *pbsettings.GetAppConfigRequest,
) (*pbsettings.GetAppConfigResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbsettings.ConfigService_ServiceDesc.ServiceName,
		Method:  "GetAppConfig",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	cfg, err := s.appCfg.Reload(ctx)
	if err != nil {
		return nil, err
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_VIEWED

	return &pbsettings.GetAppConfigResponse{
		Config: cfg,
	}, nil
}

func (s *Server) UpdateAppConfig(
	ctx context.Context,
	req *pbsettings.UpdateAppConfigRequest,
) (*pbsettings.UpdateAppConfigResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbsettings.ConfigService_ServiceDesc.ServiceName,
		Method:  "UpdateAppConfig",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	req.GetConfig().Default()
	if req.GetConfig().GetSystem().GetBannerMessage() != nil {
		var expiresAt time.Time
		if req.GetConfig().GetSystem().GetBannerMessage().GetExpiresAt() != nil {
			expiresAt = req.GetConfig().GetSystem().GetBannerMessage().GetExpiresAt().AsTime()
		}

		req.Config.System.BannerMessage.Id = utils.GetMD5HashFromString(
			req.GetConfig().GetSystem().GetBannerMessage().GetTitle() + "-" + expiresAt.String(),
		)
	}

	stmt := tConfig.
		INSERT(
			tConfig.Key,
			tConfig.AppConfig,
		).
		VALUES(
			1,
			req.GetConfig(),
		).
		ON_DUPLICATE_KEY_UPDATE(
			tConfig.AppConfig.SET(jet.RawString("VALUES(`app_config`)")),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	// Update default perms
	defaultPerms := make([]string, len(req.GetConfig().GetPerms().GetDefault()))
	for i := range req.GetConfig().GetPerms().GetDefault() {
		defaultPerms[i] = perms.BuildGuard(
			perms.Category(req.GetConfig().GetPerms().GetDefault()[i].GetCategory()),
			perms.Name(req.GetConfig().GetPerms().GetDefault()[i].GetName()),
		)
	}
	if err := s.ps.SetDefaultRolePerms(ctx, defaultPerms); err != nil {
		return nil, err
	}

	// Update config state
	if err := s.appCfg.Update(ctx, req.GetConfig()); err != nil {
		return nil, err
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED

	config, err := s.appCfg.Reload(ctx)
	if err != nil {
		return nil, err
	}

	clientCfg := clientconfig.BuildClientConfig(
		s.cfg,
		clientconfig.BuildProviderList(s.cfg),
		s.appCfg.Get(),
	)

	s.notifi.SendSystemEvent(ctx, &notifications.SystemEvent{
		Data: &notifications.SystemEvent_ClientConfig{
			ClientConfig: clientCfg,
		},
	})

	return &pbsettings.UpdateAppConfigResponse{
		Config: config,
	}, nil
}
