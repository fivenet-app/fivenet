package qualifications

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/file"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/qualifications"
	"github.com/fivenet-app/fivenet/v2025/pkg/filestore"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2025/pkg/grpc/interceptors/audit"
	errorsdocuments "github.com/fivenet-app/fivenet/v2025/services/documents/errors"
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
	if meta.GetNamespace() != "qualifications" &&
		meta.GetNamespace() != "qualifications-exam-questions" {
		return errswrap.NewError(err, filestore.ErrInvalidUploadMeta)
	}

	check, err := s.access.CanUserAccessTarget(
		ctx,
		meta.GetParentId(),
		userInfo,
		qualifications.AccessLevel_ACCESS_LEVEL_EDIT,
	)
	if err != nil {
		return errswrap.NewError(err, errorsdocuments.ErrNotFoundOrNoPerms)
	}
	if !check && !userInfo.GetSuperuser() {
		return errorsdocuments.ErrDocViewDenied
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
