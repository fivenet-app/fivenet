package settings

import (
	"context"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	pbsettings "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/settings"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/pkg/version"
	errorssettings "github.com/fivenet-app/fivenet/v2025/services/settings/errors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

func (s *Server) GetStatus(
	ctx context.Context,
	req *pbsettings.GetStatusRequest,
) (*pbsettings.GetStatusResponse, error) {
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
		Version: &pbsettings.VersionStatus{
			Current:    version.Version,
			NewVersion: nil,
		},
	}

	// Get version info from the update checker if it's enabled
	if s.updateChecker != nil {
		current, latestVersion, url, releaseDate := s.updateChecker.GetNewVersionInfo()
		if latestVersion != current {
			resp.Version.NewVersion = &pbsettings.NewVersionInfo{
				Version: latestVersion,
				Url:     url,
			}
			if !releaseDate.IsZero() {
				resp.Version.NewVersion.ReleaseDate = timestamp.New(releaseDate)
			}
		}
	}

	return resp, nil
}

func (s *Server) GetAllPermissions(
	ctx context.Context,
	req *pbsettings.GetAllPermissionsRequest,
) (*pbsettings.GetAllPermissionsResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.settings.job", req.GetJob()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbsettings.SettingsService_ServiceDesc.ServiceName,
		Method:  "GetAllPermissions",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	job := s.enricher.GetJobByName(req.GetJob())
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

func (s *Server) GetJobLimits(
	ctx context.Context,
	req *pbsettings.GetJobLimitsRequest,
) (*pbsettings.GetJobLimitsResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.settings.job", req.GetJob()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbsettings.SettingsService_ServiceDesc.ServiceName,
		Method:  "GetJobLimits",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	job := s.enricher.GetJobByName(req.GetJob())
	if job == nil {
		return nil, errorssettings.ErrInvalidRequest
	}

	resp := &pbsettings.GetJobLimitsResponse{}

	perms, err := s.ps.GetJobPermissions(ctx, job.GetName())
	if err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}
	resp.Permissions = perms

	attrs, _ := s.ps.GetJobAttributes(ctx, job.GetName())
	resp.Attributes = attrs

	resp.Job = job.GetName()
	resp.JobLabel = &job.Label

	return resp, nil
}

func (s *Server) UpdateJobLimits(
	ctx context.Context,
	req *pbsettings.UpdateJobLimitsRequest,
) (*pbsettings.UpdateJobLimitsResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.settings.job", req.GetJob()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbsettings.SettingsService_ServiceDesc.ServiceName,
		Method:  "UpdateJobLimits",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	job := s.enricher.GetJobByName(req.GetJob())
	if job == nil {
		return nil, errorssettings.ErrInvalidRequest
	}

	if err := s.ps.UpdateJobPermissions(ctx, job.GetName(), req.GetPerms().GetToUpdate()...); err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	if err := s.ps.UpdateJobAttributes(ctx, job.GetName(), req.GetAttrs().GetToUpdate()...); err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	if err := s.ps.UpdateJobPermissions(ctx, job.GetName(), req.GetPerms().GetToRemove()...); err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	if err := s.ps.ApplyJobPermissions(ctx, job.GetName()); err != nil {
		return nil, errswrap.NewError(err, errorssettings.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED

	return &pbsettings.UpdateJobLimitsResponse{}, nil
}
