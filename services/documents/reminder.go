package documents

import (
	"context"
	"time"

	documentsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/access"
	documentsworkflow "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/workflow"
	pbdocuments "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/documents"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorsdocuments "github.com/fivenet-app/fivenet/v2026/services/documents/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

func (s *Server) createOrUpdateWorkflowState(
	ctx context.Context,
	tx qrm.DB,
	documentId int64,
	workflow *documentsworkflow.Workflow,
) error {
	if workflow == nil || (!workflow.GetAutoClose() && !workflow.GetReminder()) {
		return nil
	}

	now := time.Now()

	autoCloseTime := mysql.TimestampExp(mysql.NULL)
	if workflow.GetAutoClose() && workflow.GetAutoCloseSettings() != nil &&
		workflow.GetAutoCloseSettings().GetDuration() != nil {
		autoCloseTime = mysql.TimestampT(
			now.Add(workflow.GetAutoCloseSettings().GetDuration().AsDuration()),
		)
	}

	nextReminderTime := mysql.TimestampExp(mysql.NULL)
	if workflow.GetReminder() && workflow.GetReminderSettings() != nil &&
		len(workflow.GetReminderSettings().GetReminders()) > 0 {
		reminder := workflow.GetReminderSettings().GetReminders()[0]

		nextReminderTime = mysql.TimestampT(now.Add(reminder.GetDuration().AsDuration()))
	}

	tWorkflow := table.FivenetDocumentsWorkflowState
	stmt := tWorkflow.
		INSERT(
			tWorkflow.DocumentID,
			tWorkflow.AutoCloseTime,
			tWorkflow.NextReminderTime,
			tWorkflow.NextReminderCount,
		).
		VALUES(
			documentId,
			autoCloseTime,
			nextReminderTime,
			mysql.NULL,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tWorkflow.AutoCloseTime.SET(autoCloseTime),
			tWorkflow.NextReminderTime.SET(nextReminderTime),
			tWorkflow.NextReminderCount.SET(mysql.IntExp(mysql.NULL)),
		)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	return nil
}

func (s *Server) SetDocumentReminder(
	ctx context.Context,
	req *pbdocuments.SetDocumentReminderRequest,
) (*pbdocuments.SetDocumentReminderResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.id", req.GetDocumentId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.access.CanUserAccessTarget(
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
		if err := deleteWorkflowUserState(ctx, s.db, req.GetDocumentId(), userInfo.GetUserId()); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	} else {
		if err := updateWorkflowUserState(ctx, s.db, &documentsworkflow.WorkflowUserState{
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
