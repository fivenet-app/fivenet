package docstore

import (
	"context"
	"errors"
	sync "sync"
	"time"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/documents"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/zap"
)

var tUserWorkflow = table.FivenetDocumentsWorkflowUsers.AS("workflow_user_state")

func (w *Workflow) handleDocumentsUsers(ctx context.Context, data *documents.WorkflowCronData) error {
	nowTs := jet.TimestampT(time.Now())

	stmt := tUserWorkflow.
		SELECT(
			tUserWorkflow.DocumentID,
			tUserWorkflow.UserID,
			tUserWorkflow.ManualReminderTime,
			tUserWorkflow.ManualReminderMessage,
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
		WHERE(jet.AND(
			tUserWorkflow.DocumentID.GT(jet.Uint64(data.LastDocId)),
			tUserWorkflow.ManualReminderTime.LT_EQ(nowTs),
		)).
		LIMIT(100)

	dest := []*documents.WorkflowUserState{}
	if err := stmt.QueryContext(ctx, w.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
	}

	wg := sync.WaitGroup{}

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
					w.logger.Error("error during workflow user state handling",
						zap.Uint64("document_id", state.DocumentId), zap.Int32("user_id", state.UserId), zap.Error(err))
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

func (w *Workflow) handleWorkflowUserState(ctx context.Context, state *documents.WorkflowUserState) error {
	if state.ManualReminderTime != nil && time.Since(state.ManualReminderTime.AsTime()) > 0 {
		// Send reminder and null reminder time
		if err := w.sendDocumentReminder(ctx, state.DocumentId, state.UserId, state.Document, state.ManualReminderMessage, true); err != nil {
			return err
		}

		if err := deleteWorkflowUserState(ctx, w.db, state); err != nil {
			return err
		}
	}

	return nil
}

func updateWorkflowUserState(ctx context.Context, tx qrm.DB, state *documents.WorkflowUserState) error {
	reminderTime := jet.TimestampExp(jet.NULL)
	if state.ManualReminderTime != nil {
		reminderTime = jet.TimestampT(state.ManualReminderTime.AsTime())
	}

	reminderMessage := jet.StringExp(jet.NULL)
	if state.ManualReminderMessage != nil {
		reminderMessage = jet.String(*state.ManualReminderMessage)
	}

	stmt := tUserWorkflow.
		INSERT(
			tUserWorkflow.DocumentID,
			tUserWorkflow.UserID,
			tUserWorkflow.ManualReminderTime,
			tUserWorkflow.ManualReminderMessage,
		).
		VALUES(
			state.DocumentId,
			state.UserId,
			state.ManualReminderTime,
			state.ManualReminderMessage,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tUserWorkflow.ManualReminderTime.SET(reminderTime),
			tUserWorkflow.ManualReminderMessage.SET(reminderMessage),
		)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	return nil
}

func deleteWorkflowUserState(ctx context.Context, tx qrm.DB, state *documents.WorkflowUserState) error {
	stmt := tUserWorkflow.
		DELETE().
		WHERE(jet.AND(
			tUserWorkflow.DocumentID.EQ(jet.Uint64(state.DocumentId)),
			tUserWorkflow.UserID.EQ(jet.Int32(state.UserId)),
		)).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	return nil
}
