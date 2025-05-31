package documents

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/documents"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/file"
	pbdocuments "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/documents"
	"github.com/fivenet-app/fivenet/v2025/pkg/filestore"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	errorsdocuments "github.com/fivenet-app/fivenet/v2025/services/documents/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	grpc "google.golang.org/grpc"
)

func (s *Server) UploadFile(srv grpc.ClientStreamingServer[file.UploadPacket, file.UploadResponse]) error {
	ctx := srv.Context()

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbdocuments.DocumentsService_ServiceDesc.ServiceName,
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
	meta.Namespace = "documents"

	check, err := s.access.CanUserAccessTarget(ctx, meta.ParentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_EDIT)
	if err != nil {
		return errswrap.NewError(err, errorsdocuments.ErrNotFoundOrNoPerms)
	}
	if !check && !userInfo.Superuser {
		return errorsdocuments.ErrDocViewDenied
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
