package jobs

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/file"
	jobsgroups "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/groups"
	pbjobs "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/jobs"
	"github.com/fivenet-app/fivenet/v2026/pkg/filestore"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	errorsjobs "github.com/fivenet-app/fivenet/v2026/services/jobs/errors"
	jobsstore "github.com/fivenet-app/fivenet/v2026/stores/jobs"
	logging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.uber.org/zap"
	grpc "google.golang.org/grpc"
)

const groupLogoUploadNamespace = "jobgrouplogos"

func (s *Server) UploadGroupLogo(
	srv grpc.ClientStreamingServer[file.UploadFileRequest, file.UploadFileResponse],
) error {
	ctx := srv.Context()

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	meta, err := s.groupLogoFileHandler.AwaitHandshake(srv)
	if err != nil {
		return errswrap.NewError(err, filestore.ErrInvalidUploadMeta)
	}
	if meta.GetNamespace() != groupLogoUploadNamespace {
		return errswrap.NewError(err, filestore.ErrInvalidUploadMeta)
	}

	group, err := s.store.GetGroup(ctx, s.db, jobsstore.GroupQuery{
		Job:             userInfo.GetJob(),
		IncludeArchived: false,
	}, meta.GetParentId())
	if err != nil {
		return errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if group == nil {
		return errorsjobs.ErrNotFoundOrNoPerms
	}

	var oldLogoFileID int64
	if group.GetLogoFileId() > 0 {
		oldLogoFileID = group.GetLogoFileId()
	} else if group.GetLogoFile() != nil && group.GetLogoFile().GetId() > 0 {
		oldLogoFileID = group.GetLogoFile().GetId()
	}

	name := filepath.Base(meta.GetOriginalName())
	ext := filepath.Ext(name)
	key := fmt.Sprintf(
		"%s/%s/%d%s",
		groupLogoUploadNamespace,
		userInfo.GetJob(),
		group.GetId(),
		ext,
	)

	resp, err := s.groupLogoFileHandler.UploadFile(
		ctx,
		group.GetId(),
		key,
		meta.GetSize(),
		meta.GetContentType(),
		srv,
	)
	if err != nil {
		return err
	}
	if oldLogoFileID > 0 && resp.GetId() != oldLogoFileID {
		if err := s.groupLogoFileHandler.Delete(ctx, group.GetId(), oldLogoFileID); err != nil {
			s.logger.Warn(
				"failed to cleanup old job group logo",
				zap.Int64("group_id", group.GetId()),
				zap.Int64("file_id", oldLogoFileID),
				zap.Error(err),
			)
		}
	}

	logging.InjectFields(ctx, logging.Fields{
		"fivenet.jobs.groups.id", group.GetId(),
		"fivenet.file.namespace", meta.GetNamespace(),
		"fivenet.file.name", meta.GetOriginalName(),
	})
	if err := s.addGroupActivity(
		ctx,
		s.db,
		userInfo.GetJob(),
		group.GetId(),
		jobsgroups.GroupActivityType_GROUP_ACTIVITY_TYPE_LOGO_UPDATED,
		userInfo.GetUserId(),
		0,
		0,
		nil,
		nil,
	); err != nil {
		s.logger.Warn(
			"failed to create job group logo activity",
			zap.Int64("group_id", group.GetId()),
			zap.Error(err),
		)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)

	return nil
}

func (s *Server) DeleteGroupLogo(
	ctx context.Context,
	req *pbjobs.DeleteGroupLogoRequest,
) (*pbjobs.DeleteGroupLogoResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	group, err := s.store.GetGroup(ctx, s.db, jobsstore.GroupQuery{
		Job:             userInfo.GetJob(),
		IncludeArchived: false,
	}, req.GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if group == nil {
		return nil, errorsjobs.ErrNotFoundOrNoPerms
	}

	logoFileID := group.GetLogoFileId()
	if logoFileID <= 0 && group.GetLogoFile() != nil {
		logoFileID = group.GetLogoFile().GetId()
	}
	if logoFileID > 0 {
		if err := s.groupLogoFileHandler.Delete(ctx, group.GetId(), logoFileID); err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
		if err := s.addGroupActivity(
			ctx,
			s.db,
			userInfo.GetJob(),
			group.GetId(),
			jobsgroups.GroupActivityType_GROUP_ACTIVITY_TYPE_LOGO_UPDATED,
			userInfo.GetUserId(),
			0,
			0,
			nil,
			nil,
		); err != nil {
			s.logger.Warn(
				"failed to create job group logo activity",
				zap.Int64("group_id", group.GetId()),
				zap.Error(err),
			)
		}
	}

	updated, err := s.store.GetGroup(ctx, s.db, jobsstore.GroupQuery{
		Job:             userInfo.GetJob(),
		IncludeArchived: false,
	}, req.GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if updated == nil {
		return nil, errorsjobs.ErrNotFoundOrNoPerms
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)

	if err := s.hydrateGroupColleagues(ctx, userInfo, updated); err != nil {
		return nil, err
	}

	return &pbjobs.DeleteGroupLogoResponse{Group: updated}, nil
}
