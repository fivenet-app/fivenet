package documents

import (
	pbcollab "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/collab"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/documents"
	pbdocuments "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/documents"
	permsdocuments "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/documents/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/access"
	"github.com/fivenet-app/fivenet/v2025/pkg/collab"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/grpcws"
	errorsdocuments "github.com/fivenet-app/fivenet/v2025/services/documents/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/metadata"
)

func (s *Server) JoinRoom(srv pbdocuments.CollabService_JoinRoomServer) error {
	ctx := srv.Context()

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Prepare Client Id from user and connection info
	meta := metadata.ExtractIncoming(ctx)
	connId := meta.Get(grpcws.ConnectionIdHeader)
	clientId := collab.MakeClientID(userInfo.GetUserId(), connId)

	docId, err := s.collabServer.HandleFirstMsg(ctx, clientId, srv)
	if err != nil {
		return err
	}

	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.id", docId})

	check, err := s.access.CanUserAccessTarget(
		ctx,
		docId,
		userInfo,
		documents.AccessLevel_ACCESS_LEVEL_ACCESS,
	)
	if err != nil {
		return errswrap.NewError(err, errorsdocuments.ErrNotFoundOrNoPerms)
	}
	if !check && !userInfo.GetSuperuser() {
		return errorsdocuments.ErrNotFoundOrNoPerms
	}

	doc, err := s.getDocument(ctx,
		tDocument.ID.EQ(jet.Int64(docId)),
		userInfo, true)
	if err != nil {
		return errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}

	// Field Permission Check for same job handling
	fields, err := s.ps.AttrStringList(
		userInfo,
		permsdocuments.DocumentsServicePerm,
		permsdocuments.DocumentsServiceUpdateDocumentPerm,
		permsdocuments.DocumentsServiceUpdateDocumentAccessPermField,
	)
	if err != nil {
		return errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !access.CheckIfHasOwnJobAccess(fields, userInfo, doc.GetCreatorJob(), doc.GetCreator()) {
		return errorsdocuments.ErrNotFoundOrNoPerms
	}

	return s.collabServer.HandleClient(
		ctx,
		docId,
		userInfo.GetUserId(),
		clientId,
		pbcollab.ClientRole_CLIENT_ROLE_WRITER,
		srv,
	)
}
