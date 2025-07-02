package wiki

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/file"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/wiki"
	pbwiki "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/wiki"
	"github.com/fivenet-app/fivenet/v2025/pkg/filestore"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	errorswiki "github.com/fivenet-app/fivenet/v2025/services/wiki/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	grpc "google.golang.org/grpc"
)

const MaxFilesPerPage = 5

func (s *Server) UploadFile(srv grpc.ClientStreamingServer[file.UploadFileRequest, file.UploadFileResponse]) error {
	ctx := srv.Context()

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbwiki.WikiService_ServiceDesc.ServiceName,
		Method:  "UploadFile",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}

	meta, err := s.fHandler.AwaitHandshake(srv)
	defer s.aud.Log(auditEntry, meta)
	if err != nil {
		return errswrap.NewError(err, filestore.ErrInvalidUploadMeta)
	}
	meta.Namespace = "wiki"

	check, err := s.access.CanUserAccessTarget(ctx, meta.ParentId, userInfo, wiki.AccessLevel_ACCESS_LEVEL_EDIT)
	if err != nil {
		return errswrap.NewError(err, errorswiki.ErrPageDenied)
	}
	if !check && !userInfo.Superuser {
		return errorswiki.ErrPageDenied
	}

	count, err := s.fHandler.CountFilesForParentID(ctx, meta.ParentId)
	if err != nil {
		return errswrap.NewError(err, errorswiki.ErrFailedQuery)
	}

	if count >= MaxFilesPerPage {
		return errorswiki.ErrMaxFilesReached
	}

	_, err = s.fHandler.UploadFromMeta(ctx, meta, meta.ParentId, srv)
	if err != nil {
		return err
	}

	trace.SpanFromContext(ctx).SetAttributes(attribute.String("fivenet.file.namespace", meta.Namespace))
	trace.SpanFromContext(ctx).SetAttributes(attribute.String("fivenet.file.name", meta.OriginalName))

	auditEntry.State = audit.EventType_EVENT_TYPE_CREATED

	return nil
}
