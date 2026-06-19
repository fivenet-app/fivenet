package documents

import (
	"context"

	documentsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/access"
	documentsworkflow "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/workflow"
	pbdocuments "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/documents"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	errorsdocuments "github.com/fivenet-app/fivenet/v2026/services/documents/errors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

func (s *Server) SetDocumentReminder(
	ctx context.Context,
	req *pbdocuments.SetDocumentReminderRequest,
) (*pbdocuments.SetDocumentReminderResponse, error) {
	logging.InjectFields(ctx, logging.Fields{documentIDLogFieldKey, req.GetDocumentId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.canUserAccessDocument(
		ctx,
		req.GetDocumentId(),
		userInfo,
		documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check {
		return nil, errorsdocuments.ErrDocViewDenied
	}

	if req.GetReminderTime() == nil {
		if err := s.store.DeleteWorkflowUserState(
			ctx,
			s.db,
			req.GetDocumentId(),
			userInfo.GetUserId(),
		); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	} else {
		if err := s.store.UpsertWorkflowUserState(ctx, s.db, &documentsworkflow.WorkflowUserState{
			DocumentId:            req.GetDocumentId(),
			UserId:                userInfo.GetUserId(),
			ManualReminderTime:    req.GetReminderTime(),
			ManualReminderMessage: req.Message,
			ReminderCount:         0,
			MaxReminderCount:      req.MaxReminderCount,
		}); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	return &pbdocuments.SetDocumentReminderResponse{}, nil
}
