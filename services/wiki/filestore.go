package wiki

import (
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/file"
	wikiaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/wiki/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/filestore"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	errorswiki "github.com/fivenet-app/fivenet/v2026/services/wiki/errors"
	logging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpc "google.golang.org/grpc"
)

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
		int32(wikiaccess.AccessLevel_ACCESS_LEVEL_EDIT),
	)
	if err != nil {
		return errswrap.NewError(err, errorswiki.ErrPageDenied)
	}
	if !check && !userInfo.GetJobAdmin() {
		return errorswiki.ErrPageDenied
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
