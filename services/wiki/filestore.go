package wiki

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/file"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/wiki"
	"github.com/fivenet-app/fivenet/v2025/pkg/filestore"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2025/pkg/grpc/interceptors/audit"
	errorswiki "github.com/fivenet-app/fivenet/v2025/services/wiki/errors"
	logging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpc "google.golang.org/grpc"
)

const MaxFilesPerPage = 5

func (s *Server) UploadFile(
	srv grpc.ClientStreamingServer[file.UploadFileRequest, file.UploadFileResponse],
) error {
	ctx := srv.Context()

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	meta, err := s.fHandler.AwaitHandshake(srv)
	if err != nil {
		return errswrap.NewError(err, filestore.ErrInvalidUploadMeta)
	}
	meta.Namespace = "wiki"

	check, err := s.access.CanUserAccessTarget(
		ctx,
		meta.GetParentId(),
		userInfo,
		wiki.AccessLevel_ACCESS_LEVEL_EDIT,
	)
	if err != nil {
		return errswrap.NewError(err, errorswiki.ErrPageDenied)
	}
	if !check && !userInfo.GetSuperuser() {
		return errorswiki.ErrPageDenied
	}

	count, err := s.fHandler.CountFilesForParentID(ctx, meta.GetParentId())
	if err != nil {
		return errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	if count >= MaxFilesPerPage {
		return errorswiki.ErrMaxFilesReached
	}

	_, err = s.fHandler.UploadFromMeta(ctx, meta, meta.GetParentId(), srv)
	if err != nil {
		return err
	}

	logging.InjectFields(ctx, logging.Fields{
		"fivenet.file.namespace", meta.GetNamespace(),
		"fivenet.file.name", meta.GetOriginalName(),
	})

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)

	return nil
}
