package documents

import (
	"context"
	"errors"
	sync "sync"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/documents"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/zap"
)

var tUserWorkflow = table.FivenetDocumentsWorkflowUsers.AS("workflow_user_state")

func (w *Workflow) handleDocumentsUsers(
	ctx context.Context,
	data *documents.WorkflowCronData,
) error {
	tDTemplates := table.FivenetDocumentsTemplates.AS("template_short")

	stmt := tUserWorkflow.
		SELECT(
			tUserWorkflow.DocumentID,
			tUserWorkflow.UserID,
			tUserWorkflow.ManualReminderTime,
			tUserWorkflow.ManualReminderMessage,
			tUserWorkflow.ReminderCount,
			tUserWorkflow.MaxReminderCount,
			tDTemplates.Workflow,
			tDocumentShort.Title,
			tDocumentShort.CreatorID,
			tDocumentShort.CreatorJob,
		).
		FROM(
			tUserWorkflow.
				INNER_JOIN(tDocumentShort,
					tDocumentShort.ID.EQ(tUserWorkflow.DocumentID).
						AND(tDocumentShort.DeletedAt.IS_NULL()),
				).
				LEFT_JOIN(tDTemplates,
					tDTemplates.ID.EQ(tDocumentShort.TemplateID).
						AND(tDTemplates.DeletedAt.IS_NULL()),
				),
		).
		WHERE(mysql.AND(
			tUserWorkflow.DocumentID.GT(mysql.Int64(data.GetLastDocId())),
			tUserWorkflow.ManualReminderTime.LT_EQ(mysql.TimestampT(time.Now())),
		)).
		LIMIT(100)

	var dest []*documents.WorkflowUserState
	if err := stmt.QueryContext(ctx, w.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
	}

	var wg sync.WaitGroup

	// Run at max 3 handlers at once
	workChannel := make(chan *documents.WorkflowUserState, 3)

	wg.Add(1)
	go func() {
		defer wg.Done()

		for state := range workChannel {
			wg.Add(1)
			go func() {
				defer wg.Done()

				if err := w.handleWorkflowUserState(ctx, state); err != nil {
					w.logger.Error(
						"error during workflow user state handling",
						zap.Int64(
							"document_id",
							state.GetDocumentId(),
						),
						zap.Int32("user_id", state.GetUserId()),
						zap.Error(err),
					)
				}
			}()
		}
	}()

	for _, ws := range dest {
		workChannel <- ws
	}

	close(workChannel)

	wg.Wait()

	return nil
}

func (w *Workflow) handleWorkflowUserState(
	ctx context.Context,
	state *documents.WorkflowUserState,
) error {
	if state.GetManualReminderTime() != nil &&
		time.Since(state.GetManualReminderTime().AsTime()) > 0 {
		// Send reminder and null reminder time
		if err := w.sendDocumentReminder(ctx, state.GetDocumentId(), state.GetUserId(), state.GetDocument(), state.GetManualReminderMessage(), true); err != nil {
			return err
		}

		if err := deleteWorkflowUserState(ctx, w.db, state); err != nil {
			return err
		}
	}

	return nil
}

func updateWorkflowUserState(
	ctx context.Context,
	tx qrm.DB,
	state *documents.WorkflowUserState,
) error {
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

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	return nil
}

func deleteWorkflowUserState(
	ctx context.Context,
	tx qrm.DB,
	state *documents.WorkflowUserState,
) error {
	tUserWorkflow := table.FivenetDocumentsWorkflowUsers

	stmt := tUserWorkflow.
		DELETE().
		WHERE(mysql.AND(
			tUserWorkflow.DocumentID.EQ(mysql.Int64(state.GetDocumentId())),
			tUserWorkflow.UserID.EQ(mysql.Int32(state.GetUserId())),
		)).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	return nil
}
