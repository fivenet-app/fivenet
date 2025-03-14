package rector

import (
	"context"
	"fmt"
	"time"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/notifications"
	rector "github.com/fivenet-app/fivenet/gen/go/proto/resources/rector"
	pbrector "github.com/fivenet-app/fivenet/gen/go/proto/services/rector"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/notifi"
	"github.com/fivenet-app/fivenet/pkg/perms"
	"github.com/fivenet-app/fivenet/pkg/utils"
	"github.com/fivenet-app/fivenet/query/fivenet/model"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
)

var tConfig = table.FivenetConfig

func (s *Server) GetAppConfig(ctx context.Context, req *pbrector.GetAppConfigRequest) (*pbrector.GetAppConfigResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbrector.RectorConfigService_ServiceDesc.ServiceName,
		Method:  "GetAppConfig",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	config, err := s.appCfg.Reload(ctx)
	if err != nil {
		return nil, err
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_VIEWED)

	return &pbrector.GetAppConfigResponse{
		Config: config,
	}, nil
}

func (s *Server) UpdateAppConfig(ctx context.Context, req *pbrector.UpdateAppConfigRequest) (*pbrector.UpdateAppConfigResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbrector.RectorConfigService_ServiceDesc.ServiceName,
		Method:  "UpdateAppConfig",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	currentConfig := s.appCfg.Get()

	req.Config.Default()
	if req.Config.System.BannerMessage != nil {
		var expiresAt time.Time
		if req.Config.System.BannerMessage.ExpiresAt != nil {
			expiresAt = req.Config.System.BannerMessage.ExpiresAt.AsTime()
		}

		req.Config.System.BannerMessage.Id = utils.GetMD5HashFromString(req.Config.System.BannerMessage.Title + "-" + expiresAt.String())
	}

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

	// Update default perms
	defaultPerms := make([]string, len(req.Config.Perms.Default))
	for i := range req.Config.Perms.Default {
		defaultPerms[i] = perms.BuildGuard(perms.Category(req.Config.Perms.Default[i].Category), perms.Name(req.Config.Perms.Default[i].Name))
	}
	if err := s.ps.SetDefaultRolePerms(ctx, defaultPerms); err != nil {
		return nil, err
	}

	// If banner is disabled and was previously enabled, send "nil" banner message to remove it from clients
	if !req.Config.System.BannerMessageEnabled && currentConfig.System.BannerMessageEnabled != req.Config.System.BannerMessageEnabled {
		s.js.PublishProto(ctx, fmt.Sprintf("%s.%s", notifi.BaseSubject, notifi.SystemTopic), &notifications.SystemEvent{
			Data: &notifications.SystemEvent_BannerMessage{
				BannerMessage: &notifications.BannerMessageWrapper{
					BannerMessage: nil,
				},
			},
		})
	} else if req.Config.System.BannerMessageEnabled {
		// Check if an updated banner message event is needed by md5 hashing the title and using that as the ID
		if currentConfig.System.BannerMessage == nil || req.Config.System.BannerMessage != nil && (currentConfig.System.BannerMessage.Id != req.Config.System.BannerMessage.Id ||
			(req.Config.System.BannerMessage.ExpiresAt != nil &&
				(currentConfig.System.BannerMessage.ExpiresAt == nil ||
					req.Config.System.BannerMessage.ExpiresAt.AsTime().Compare(currentConfig.System.BannerMessage.ExpiresAt.AsTime()) != 0))) {
			s.js.PublishProto(ctx, fmt.Sprintf("%s.%s", notifi.BaseSubject, notifi.SystemTopic), &notifications.SystemEvent{
				Data: &notifications.SystemEvent_BannerMessage{
					BannerMessage: &notifications.BannerMessageWrapper{
						BannerMessage: req.Config.System.BannerMessage,
					},
				},
			})
		}
	}

	// Update config state
	if err := s.appCfg.Update(ctx, req.Config); err != nil {
		return nil, err
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	config, err := s.appCfg.Reload(ctx)
	if err != nil {
		return nil, err
	}

	return &pbrector.UpdateAppConfigResponse{
		Config: config,
	}, nil
}
