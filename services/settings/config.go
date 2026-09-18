package settings

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/clientconfig"
	notificationsevents "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/notifications/events"
	pbsettings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/settings"
	grpcauth "github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	"github.com/fivenet-app/fivenet/v2026/pkg/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils"
	errorssettings "github.com/fivenet-app/fivenet/v2026/services/settings/errors"
)

func (s *Server) GetAppConfig(
	ctx context.Context,
	req *pbsettings.GetAppConfigRequest,
) (*pbsettings.GetAppConfigResponse, error) {
	setAuditAccountMeta(ctx)

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
	setAuditAccountMeta(ctx)

	req.GetConfig().Default()
	if req.GetConfig().GetSystem().GetBannerMessage() != nil {
		var expiresAt time.Time
		if req.GetConfig().GetSystem().GetBannerMessage().GetExpiresAt() != nil {
			expiresAt = req.GetConfig().GetSystem().GetBannerMessage().GetExpiresAt().AsTime()
		}

		req.Config.System.BannerMessage.Id = utils.GetSHA256HashFromString(
			req.GetConfig().GetSystem().GetBannerMessage().GetTitle() + "-" + expiresAt.String(),
		)
	}

	// Update default perms
	cfgDefaultperms := req.GetConfig().GetPerms().GetDefault()
	defaultPerms := make([]string, len(req.GetConfig().GetPerms().GetDefault()))
	for i := range cfgDefaultperms {
		guard, err := perms.DefaultPermGuard(
			cfgDefaultperms[i].GetCategory(),
			cfgDefaultperms[i].GetName(),
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorssettings.ErrInvalidDefaultPerms)
		}

		separator := strings.LastIndexByte(cfgDefaultperms[i].GetCategory(), '.')
		permission, err := s.perms.GetPermission(
			ctx,
			perms.Namespace(cfgDefaultperms[i].GetCategory()[:separator]),
			perms.Service(cfgDefaultperms[i].GetCategory()[separator+1:]),
			perms.Name(cfgDefaultperms[i].GetName()),
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
		}
		if permission == nil {
			return nil, errswrap.NewError(
				fmt.Errorf("default permission not found: %s", guard),
				errorssettings.ErrInvalidDefaultPerms,
			)
		}

		defaultPerms[i] = guard
	}
	// Update config state
	if err := s.appCfg.Update(ctx, req.GetConfig()); err != nil {
		return nil, err
	}

	if err := s.perms.SetDefaultRolePerms(ctx, defaultPerms); err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)

	config, err := s.appCfg.Reload(ctx)
	if err != nil {
		return nil, err
	}

	clientCfg := clientconfig.BuildClientConfig(
		s.cfg,
		clientconfig.BuildProviderList(s.cfg),
		config,
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

func setAuditAccountMeta(ctx context.Context) {
	if userInfo, ok := grpcauth.GetUserInfoFromContext(ctx); ok && userInfo.GetAccountId() > 0 {
		grpc_audit.SetAccountID(ctx, userInfo.GetAccountId())
	}
}
