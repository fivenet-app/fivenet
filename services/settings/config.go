package settings

import (
	"context"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/clientconfig"
	notificationsevents "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/notifications/events"
	pbsettings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/settings"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	"github.com/fivenet-app/fivenet/v2026/pkg/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils"
)

func (s *Server) GetAppConfig(
	ctx context.Context,
	req *pbsettings.GetAppConfigRequest,
) (*pbsettings.GetAppConfigResponse, error) {
	cfg, err := s.appCfg.Reload(ctx)
	if err != nil {
		return nil, err
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_VIEWED)

	return &pbsettings.GetAppConfigResponse{
		Config: cfg,
	}, nil
}

func (s *Server) UpdateAppConfig(
	ctx context.Context,
	req *pbsettings.UpdateAppConfigRequest,
) (*pbsettings.UpdateAppConfigResponse, error) {
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

	if err := s.store.UpdateAppConfig(ctx, req.GetConfig()); err != nil {
		return nil, err
	}

	// Update default perms
	cfgDefaultperms := req.GetConfig().GetPerms().GetDefault()
	defaultPerms := make([]string, len(req.GetConfig().GetPerms().GetDefault()))
	for i := range cfgDefaultperms {
		split := strings.Split(cfgDefaultperms[i].GetCategory(), ".")
		namespace := strings.Join(split[:len(split)-1], ".")
		svc := split[len(split)-1]

		defaultPerms[i] = perms.BuildGuard(
			perms.Namespace(namespace),
			perms.Service(svc),
			perms.Name(cfgDefaultperms[i].GetName()),
		)
	}
	if err := s.ps.SetDefaultRolePerms(ctx, defaultPerms); err != nil {
		return nil, err
	}

	// Update config state
	if err := s.appCfg.Update(ctx, req.GetConfig()); err != nil {
		return nil, err
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)

	config, err := s.appCfg.Reload(ctx)
	if err != nil {
		return nil, err
	}

	clientCfg := clientconfig.BuildClientConfig(
		s.cfg,
		clientconfig.BuildProviderList(s.cfg),
		s.appCfg.Get(),
	)

	s.notifi.SendSystemEvent(ctx, &notificationsevents.SystemEvent{
		Data: &notificationsevents.SystemEvent_ClientConfig{
			ClientConfig: clientCfg,
		},
	})

	return &pbsettings.UpdateAppConfigResponse{
		Config: config,
	}, nil
}
