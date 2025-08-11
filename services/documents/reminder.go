package documents

import (
	"context"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/documents"
	pbdocuments "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/documents"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsdocuments "github.com/fivenet-app/fivenet/v2025/services/documents/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

func (s *Server) createOrUpdateWorkflowState(ctx context.Context, tx qrm.DB, documentId uint64, workflow *documents.Workflow) error {
	if workflow == nil || (!workflow.AutoClose && !workflow.Reminder) {
		return nil
	}

	now := time.Now()

	autoCloseTime := jet.TimestampExp(jet.NULL)
	if workflow.AutoClose && workflow.AutoCloseSettings != nil && workflow.AutoCloseSettings.Duration != nil {
		autoCloseTime = jet.TimestampT(now.Add(workflow.AutoCloseSettings.Duration.AsDuration()))
	}

	nextReminderTime := jet.TimestampExp(jet.NULL)
	if workflow.Reminder && workflow.ReminderSettings != nil && len(workflow.ReminderSettings.Reminders) > 0 {
		reminder := workflow.ReminderSettings.Reminders[0]

		nextReminderTime = jet.TimestampT(now.Add(reminder.Duration.AsDuration()))
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
			jet.NULL,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tWorkflow.AutoCloseTime.SET(autoCloseTime),
			tWorkflow.NextReminderTime.SET(nextReminderTime),
			tWorkflow.NextReminderCount.SET(jet.IntExp(jet.NULL)),
		)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	return nil
}

func (s *Server) SetDocumentReminder(ctx context.Context, req *pbdocuments.SetDocumentReminderRequest) (*pbdocuments.SetDocumentReminderResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.documents.id", req.DocumentId})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbdocuments.DocumentsService_ServiceDesc.ServiceName,
		Method:  "SetDocumentReminder",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(ctx, req.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_VIEW)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
	}
	if !check {
		return nil, errorsdocuments.ErrDocViewDenied
	}

	if req.ReminderTime == nil {
		if err := deleteWorkflowUserState(ctx, s.db, &documents.WorkflowUserState{
			DocumentId: req.DocumentId,
			UserId:     userInfo.UserId,
		}); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	} else {
		if err := updateWorkflowUserState(ctx, s.db, &documents.WorkflowUserState{
			DocumentId:            req.DocumentId,
			UserId:                userInfo.UserId,
			ManualReminderTime:    req.ReminderTime,
			ManualReminderMessage: req.Message,
		}); err != nil {
			return nil, errswrap.NewError(err, errorsdocuments.ErrFailedQuery)
		}
	}

	return &pbdocuments.SetDocumentReminderResponse{}, nil
}
