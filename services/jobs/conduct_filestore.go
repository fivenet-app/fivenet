package jobs

import (
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/file"
	"github.com/fivenet-app/fivenet/v2026/pkg/filestore"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorsjobs "github.com/fivenet-app/fivenet/v2026/services/jobs/errors"
	logging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpc "google.golang.org/grpc"
)

var tConductFiles = table.FivenetJobConductFiles

func (s *Server) UploadFile(
	srv grpc.ClientStreamingServer[file.UploadFileRequest, file.UploadFileResponse],
) error {
	ctx := srv.Context()

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	meta, err := s.fHandler.AwaitHandshake(srv)
	if err != nil {
		return errswrap.NewError(err, filestore.ErrInvalidUploadMeta)
	}
	if meta.GetNamespace() != "jobs-conduct" {
		return errswrap.NewError(err, filestore.ErrInvalidUploadMeta)
	}

	entry, err := s.getConductEntry(ctx, meta.GetParentId(), false)
	if err != nil {
		return errswrap.NewError(err, errorsjobs.ErrNotFoundOrNoPerms)
	}

	if entry.GetDeletedAt() != nil && !userInfo.GetSuperuser() {
		return errorsjobs.ErrNotFoundOrNoPerms
	}

	if entry.GetJob() != userInfo.GetJob() && !userInfo.GetSuperuser() {
		return errorsjobs.ErrNotFoundOrNoPerms
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
