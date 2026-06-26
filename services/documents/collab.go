package documents

import (
	"context"

	pbcollab "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/collab"
	documentsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/access"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pbdocuments "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/documents"
	"github.com/fivenet-app/fivenet/v2026/pkg/collab"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/grpcws"
	errorsdocuments "github.com/fivenet-app/fivenet/v2026/services/documents/errors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/metadata"
)

func (s *Server) canUserJoinDocumentCollab(
	ctx context.Context,
	docID int64,
	userInfo *userinfo.UserInfo,
) (bool, error) {
	for _, level := range []documentsaccess.AccessLevel{
		documentsaccess.AccessLevel_ACCESS_LEVEL_ACCESS,
		documentsaccess.AccessLevel_ACCESS_LEVEL_STATUS,
		documentsaccess.AccessLevel_ACCESS_LEVEL_EDIT,
	} {
		check, err := s.canUserAccessDocument(ctx, docID, userInfo, level)
		if err != nil {
			return false, err
		}
		if check {
			return true, nil
		}
	}

	return false, nil
}

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

	logging.InjectFields(ctx, logging.Fields{documentIDLogFieldKey, docId})

	check, err := s.canUserJoinDocumentCollab(ctx, docId, userInfo)
	if err != nil {
		return errswrap.NewError(err, errorsdocuments.ErrNotFoundOrNoPerms)
	}
	if !check && !userInfo.GetSuperuser() {
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
