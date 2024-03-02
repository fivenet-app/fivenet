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
		UPDATE(
			tConfig.AppConfig,
		).
		SET(
			req.Config,
		).
		WHERE(
			tConfig.Key.EQ(jet.Uint64(1)),
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
