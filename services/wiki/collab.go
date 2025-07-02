package wiki

import (
	pbcollab "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/collab"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/wiki"
	pbdocuments "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/documents"
	"github.com/fivenet-app/fivenet/v2025/pkg/collab"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/grpcws"
	errorswiki "github.com/fivenet-app/fivenet/v2025/services/wiki/errors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/metadata"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (s *Server) JoinRoom(srv pbdocuments.CollabService_JoinRoomServer) error {
	ctx := srv.Context()

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Prepare Client Id from user and connection info
	meta := metadata.ExtractIncoming(ctx)
	connId := meta.Get(grpcws.ConnectionIdHeader)
	clientId := collab.MakeClientID(userInfo.UserId, connId)

	pageId, err := s.collabServer.HandleFirstMsg(ctx, clientId, srv)
	if err != nil {
		return err
	}

	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.page.id", int64(pageId)))

	check, err := s.access.CanUserAccessTarget(ctx, pageId, userInfo, wiki.AccessLevel_ACCESS_LEVEL_ACCESS)
	if err != nil {
		return errswrap.NewError(err, errorswiki.ErrPageDenied)
	}
	if !check && !userInfo.Superuser {
		return errorswiki.ErrPageDenied
	}

	return s.collabServer.HandleClient(ctx, pageId, userInfo.UserId, clientId, pbcollab.ClientRole_CLIENT_ROLE_WRITER, srv)
}
