package documentsstore

import (
	"context"
	"time"

	documentsworkflow "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/workflow"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) UpsertWorkflowState(ctx context.Context, tx qrm.DB, documentID int64, workflow *documentsworkflow.Workflow) error {
	if workflow == nil || (!workflow.GetAutoClose() && !workflow.GetReminder()) {
		return nil
	}

	now := time.Now()
	autoCloseTime := mysql.TimestampExp(mysql.NULL)
	if workflow.GetAutoClose() && workflow.GetAutoCloseSettings() != nil && workflow.GetAutoCloseSettings().GetDuration() != nil {
		autoCloseTime = mysql.TimestampT(now.Add(workflow.GetAutoCloseSettings().GetDuration().AsDuration()))
	}

	nextReminderTime := mysql.TimestampExp(mysql.NULL)
	if workflow.GetReminder() && workflow.GetReminderSettings() != nil && len(workflow.GetReminderSettings().GetReminders()) > 0 {
		reminder := workflow.GetReminderSettings().GetReminders()[0]
		nextReminderTime = mysql.TimestampT(now.Add(reminder.GetDuration().AsDuration()))
	}

	tWorkflowState := table.FivenetDocumentsWorkflowState
	stmt := tWorkflowState.
		INSERT(
			tWorkflowState.DocumentID,
			tWorkflowState.AutoCloseTime,
			tWorkflowState.NextReminderTime,
			tWorkflowState.NextReminderCount,
		).
		VALUES(
			documentID,
			autoCloseTime,
			nextReminderTime,
			mysql.NULL,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tWorkflowState.AutoCloseTime.SET(autoCloseTime),
			tWorkflowState.NextReminderTime.SET(nextReminderTime),
			tWorkflowState.NextReminderCount.SET(mysql.IntExp(mysql.NULL)),
		)

	_, err := stmt.ExecContext(ctx, tx)
	return err
}

func (s *Store) UpsertWorkflowUserState(ctx context.Context, tx qrm.DB, state *documentsworkflow.WorkflowUserState) error {
	reminderTime := mysql.TimestampExp(mysql.NULL)
	if state.GetManualReminderTime() != nil {
		reminderTime = mysql.TimestampT(state.GetManualReminderTime().AsTime())
	}

	reminderMessage := mysql.StringExp(mysql.NULL)
	if state.ManualReminderMessage != nil && state.GetManualReminderMessage() != "" {
		reminderMessage = mysql.String(state.GetManualReminderMessage())
	}

	tUserWorkflow := table.FivenetDocumentsWorkflowUsers
	stmt := tUserWorkflow.
		INSERT(
			tUserWorkflow.DocumentID,
			tUserWorkflow.UserID,
			tUserWorkflow.ManualReminderTime,
			tUserWorkflow.ManualReminderMessage,
			tUserWorkflow.ReminderCount,
			tUserWorkflow.MaxReminderCount,
		).
		VALUES(
			state.GetDocumentId(),
			state.GetUserId(),
			state.GetManualReminderTime(),
			state.ManualReminderMessage,
			state.GetReminderCount(),
			state.GetMaxReminderCount(),
		).
		ON_DUPLICATE_KEY_UPDATE(
			tUserWorkflow.ManualReminderTime.SET(reminderTime),
			tUserWorkflow.ManualReminderMessage.SET(reminderMessage),
			tUserWorkflow.ReminderCount.SET(mysql.Int32(state.GetReminderCount())),
			tUserWorkflow.MaxReminderCount.SET(mysql.Int32(state.GetMaxReminderCount())),
		)

	_, err := stmt.ExecContext(ctx, tx)
	return err
}

func (s *Store) DeleteWorkflowUserState(ctx context.Context, tx qrm.DB, documentID int64, userID int32) error {
	tUserWorkflow := table.FivenetDocumentsWorkflowUsers
	stmt := tUserWorkflow.
		DELETE().
		WHERE(mysql.AND(
			tUserWorkflow.DocumentID.EQ(mysql.Int64(documentID)),
			tUserWorkflow.UserID.EQ(mysql.Int32(userID)),
		)).
		LIMIT(1)

	_, err := stmt.ExecContext(ctx, tx)
	return err
}
