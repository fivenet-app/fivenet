package settings

import (
	"context"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	pbsettings "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/settings"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	errorssettings "github.com/fivenet-app/fivenet/v2025/services/settings/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (s *Server) GetStatus(ctx context.Context, req *pbsettings.GetStatusRequest) (*pbsettings.GetStatusResponse, error) {
	dbCharset, dbCollation := s.dbReq.GetDBCharsetAndCollation()
	migrationVersion, migrationDirty := s.dbReq.GetMigrationState()

	resp := &pbsettings.GetStatusResponse{
		Database: &pbsettings.Database{
			Version:          s.dbReq.GetVersion(),
			Connected:        true,
			DbCharset:        dbCharset,
			DbCollation:      dbCollation,
			MigrationVersion: uint64(migrationVersion),
			MigrationDirty:   migrationDirty,
			TablesOk:         len(s.dbReq.GetTables()) == 0,
		},
		Nats: &pbsettings.Nats{
			Version:   s.natsReq.GetVersion(),
			Connected: !s.js.Conn().IsClosed(),
		},
		Dbsync: s.syncServer.GetSyncTimes(),
	}

	return resp, nil
}

func (s *Server) GetAllPermissions(ctx context.Context, req *pbsettings.GetAllPermissionsRequest) (*pbsettings.GetAllPermissionsResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.String("fivenet.settings.job", req.Job))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbsettings.SettingsService_ServiceDesc.ServiceName,
		Method:  "GetAllPermissions",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	job := s.enricher.GetJobByName(req.Job)
	if job == nil {
		return nil, errorssettings.ErrInvalidRequest
	}

	perms, err := s.ps.GetAllPermissions(ctx)
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	attrs, err := s.ps.GetAllAttributes(ctx)
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	resp := &pbsettings.GetAllPermissionsResponse{}
	resp.Permissions = perms
	resp.Attributes = attrs

	return resp, nil
}

func (s *Server) GetJobLimits(ctx context.Context, req *pbsettings.GetJobLimitsRequest) (*pbsettings.GetJobLimitsResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.String("fivenet.settings.job", req.Job))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbsettings.SettingsService_ServiceDesc.ServiceName,
		Method:  "GetJobLimits",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	job := s.enricher.GetJobByName(req.Job)
	if job == nil {
		return nil, errorssettings.ErrInvalidRequest
	}

	resp := &pbsettings.GetJobLimitsResponse{}

	perms, err := s.ps.GetJobPermissions(ctx, job.Name)
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}
	resp.Permissions = perms

	attrs, _ := s.ps.GetJobAttributes(ctx, job.Name)
	resp.Attributes = attrs

	resp.Job = job.Name
	resp.JobLabel = &job.Label

	return resp, nil
}

func (s *Server) UpdateJobLimits(ctx context.Context, req *pbsettings.UpdateJobLimitsRequest) (*pbsettings.UpdateJobLimitsResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.String("fivenet.settings.job", req.Job))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbsettings.SettingsService_ServiceDesc.ServiceName,
		Method:  "UpdateJobLimits",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	job := s.enricher.GetJobByName(req.Job)
	if job == nil {
		return nil, errorssettings.ErrInvalidRequest
	}

	if err := s.ps.UpdateJobPermissions(ctx, job.Name, req.Perms.ToUpdate...); err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	if err := s.ps.UpdateJobAttributes(ctx, job.Name, req.Attrs.ToUpdate...); err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	if err := s.ps.UpdateJobPermissions(ctx, job.Name, req.Perms.ToRemove...); err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	if err := s.ps.ApplyJobPermissions(ctx, job.Name); err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED

	return &pbsettings.UpdateJobLimitsResponse{}, nil
}
