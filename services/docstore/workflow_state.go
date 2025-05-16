package docstore

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/cron"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/documents"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/notifications"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2025/pkg/access"
	"github.com/fivenet-app/fivenet/v2025/pkg/croner"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/v2025/pkg/notifi"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

var tWorkflow = table.FivenetDocumentsWorkflowState.AS("workflow_state")

type Workflow struct {
	logger *zap.Logger
	tracer trace.Tracer

	db    *sql.DB
	ui    userinfo.UserInfoRetriever
	notif notifi.INotifi

	access *access.Grouped[documents.DocumentJobAccess, *documents.DocumentJobAccess, documents.DocumentUserAccess, *documents.DocumentUserAccess, access.DummyQualificationAccess[documents.AccessLevel], *access.DummyQualificationAccess[documents.AccessLevel], documents.AccessLevel]
}

type WorkflowParams struct {
	fx.In

	LC fx.Lifecycle

	Logger *zap.Logger
	DB     *sql.DB
	TP     *tracesdk.TracerProvider
	Notif  notifi.INotifi
	Ui     userinfo.UserInfoRetriever

	Cron         croner.ICron
	CronHandlers *croner.Handlers
}

func NewWorkflow(p WorkflowParams) *Workflow {
	w := &Workflow{
		logger: p.Logger.Named("docstore.workflow"),
		tracer: p.TP.Tracer("docstore_workflow"),
		db:     p.DB,
		notif:  p.Notif,
		ui:     p.Ui,

		access: newAccess(p.DB),
	}

	p.LC.Append(fx.StartHook(func(ctx context.Context) error {
		if err := p.Cron.RegisterCronjob(ctx, &cron.Cronjob{
			Name:     "docstore.workflow_run",
			Schedule: "* * * * *", // Every minute
		}); err != nil {
			return err
		}

		if err := p.Cron.RegisterCronjob(ctx, &cron.Cronjob{
			Name:     "docstore.workflow_users_run",
			Schedule: "* * * * *", // Every minute
		}); err != nil {
			return err
		}

		return nil
	}))

	p.CronHandlers.Add("docstore.workflow_run", func(ctx context.Context, data *cron.CronjobData) error {
		ctx, span := w.tracer.Start(ctx, "docstore.workflow_run")
		defer span.End()

		dest := &documents.WorkflowCronData{
			LastDocId: 0,
		}
		if data.Data == nil {
			data.Data = &anypb.Any{}
		}

		if err := anypb.UnmarshalTo(data.Data, dest, proto.UnmarshalOptions{}); err != nil {
			w.logger.Warn("failed to unmarshal document workflow cron data", zap.Error(err))
		}

		if err := w.handleDocuments(ctx, dest); err != nil {
			return fmt.Errorf("error during docstore workflow handling. %w", err)
		}

		if err := data.Data.MarshalFrom(dest); err != nil {
			return fmt.Errorf("failed to marshal updated document workflow cron data. %w", err)
		}

		return nil
	})

	p.CronHandlers.Add("docstore.workflow_users_run", func(ctx context.Context, data *cron.CronjobData) error {
		ctx, span := w.tracer.Start(ctx, "docstore.workflow_users_run")
		defer span.End()

		dest := &documents.WorkflowCronData{
			LastDocId: 0,
		}
		if data.Data == nil {
			data.Data = &anypb.Any{}
		}

		if err := anypb.UnmarshalTo(data.Data, dest, proto.UnmarshalOptions{}); err != nil {
			w.logger.Error("failed to unmarshal document workflow cron data", zap.Error(err))
		}

		if err := w.handleDocumentsUsers(ctx, dest); err != nil {
			return fmt.Errorf("error during docstore workflow handling. %w", err)
		}

		if err := data.Data.MarshalFrom(dest); err != nil {
			return fmt.Errorf("failed to marshal updated document workflow cron data. %w", err)
		}

		return nil
	})

	return w
}

func (w *Workflow) handleDocuments(ctx context.Context, data *documents.WorkflowCronData) error {
	nowTs := jet.TimestampT(time.Now())

	tDTemplates := tDTemplates.AS("template")

	dest := []*documents.WorkflowState{}
	for {
		stmt := tWorkflow.
			SELECT(
				tWorkflow.DocumentID,
				tWorkflow.NextReminderTime,
				tWorkflow.NextReminderCount,
				tWorkflow.AutoCloseTime,
				tDTemplates.Workflow.AS("workflowstate.workflow"),
				tDocumentShort.Title,
				tDocumentShort.CreatorID,
				tDocumentShort.CreatorJob,
			).
			FROM(
				tWorkflow.
					INNER_JOIN(tDocumentShort,
						tDocumentShort.ID.EQ(tWorkflow.DocumentID).
							AND(tDocumentShort.DeletedAt.IS_NULL()),
					).
					LEFT_JOIN(tDTemplates,
						tDTemplates.ID.EQ(tDocumentShort.TemplateID).
							AND(tDTemplates.DeletedAt.IS_NULL()),
					),
			).
			WHERE(jet.AND(
				tWorkflow.DocumentID.GT(jet.Uint64(data.LastDocId)),
				jet.AND( // Only auto close and auto remind docs that aren't closed and have an owner
					tDocumentShort.Closed.IS_FALSE(),
					jet.OR(
						tWorkflow.NextReminderTime.LT_EQ(nowTs),
						tWorkflow.AutoCloseTime.LT_EQ(nowTs),
					),
				),
			)).
			LIMIT(250)

		if err := stmt.QueryContext(ctx, w.db, &dest); err != nil {
			if !errors.Is(err, qrm.ErrNoRows) {
				return err
			}
		}

		if data.LastDocId == 0 && len(dest) == 0 {
			// No entries match condition
			break
		} else {
			// Less than 250 entries? Reset last id to 0
			if len(dest) < 250 {
				data.LastDocId = 0
				break
				// No entries, reset last id to 0 and try again
			} else if len(dest) == 0 {
				data.LastDocId = 0
				continue
			}

			break
		}
	}

	wg := sync.WaitGroup{}

	// Run at max 3 handlers at once
	workChannel := make(chan *documents.WorkflowState, 3)

	wg.Add(1)
	go func() {
		defer wg.Done()

		for state := range workChannel {
			wg.Add(1)
			go func() {
				defer wg.Done()

				if err := w.handleWorkflowState(ctx, state); err != nil {
					w.logger.Error("error during workflow state handling",
						zap.Uint64("document_id", state.DocumentId), zap.Error(err))
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

func (w *Workflow) handleWorkflowState(ctx context.Context, state *documents.WorkflowState) error {
	if state.AutoCloseTime != nil && time.Since(state.AutoCloseTime.AsTime()) > 0 {
		if state.Workflow != nil && state.Workflow.AutoCloseSettings != nil && state.Workflow.AutoClose && state.Workflow.AutoCloseSettings.Message != "" {
			// Auto close document and null "next reminder time"
			if err := w.autoCloseDocument(ctx, state, state.Workflow.AutoCloseSettings.Message); err != nil {
				return err
			}
		}

		// Delete document workflow state, auto reminders are not sent for a closed document
		return w.deleteWorkflowState(ctx, state)
	} else if state.NextReminderTime != nil && time.Since(state.NextReminderTime.AsTime()) > 0 {
		if state.Document != nil && state.Document.CreatorId != nil {
			var reminderMessage *string
			if reminder := w.getAutoReminder(state); reminder != nil && reminder.Message != "" {
				reminderMessage = &reminder.Message
			}

			// Send notification when the document has a creator that is still part of the document's job
			if err := w.sendDocumentReminder(ctx, state.DocumentId, *state.Document.CreatorId, state.Document, reminderMessage, false); err != nil {
				return err
			}

			w.updateAutoReminderTime(state)
		} else {
			state.NextReminderTime = nil
			state.NextReminderCount = nil
		}
	}

	// Make sure to delete the document workflow state as we don't have a creator anymore
	if state.Document == nil || state.Document.CreatorId == nil {
		return w.deleteWorkflowState(ctx, state)
	}

	if err := w.updateWorkflowState(ctx, state); err != nil {
		return err
	}

	return nil
}

func (w *Workflow) getAutoReminder(state *documents.WorkflowState) *documents.Reminder {
	count := int32(0)
	if state.NextReminderCount != nil {
		count = *state.NextReminderCount
	}

	if state.Workflow == nil || state.Workflow.ReminderSettings == nil || len(state.Workflow.ReminderSettings.Reminders) < int(count) {
		return nil
	}

	return state.Workflow.ReminderSettings.Reminders[count]
}

func (w *Workflow) updateAutoReminderTime(state *documents.WorkflowState) {
	if state.Workflow == nil || state.Workflow.ReminderSettings == nil || !state.Workflow.Reminder || len(state.Workflow.ReminderSettings.Reminders) == 0 {
		state.NextReminderTime = nil
		state.NextReminderCount = nil
		return
	}

	if state.NextReminderCount == nil {
		zero := int32(0)
		state.NextReminderCount = &zero
	} else {
		*state.NextReminderCount++
	}

	if len(state.Workflow.ReminderSettings.Reminders) >= int(*state.NextReminderCount) {
		reminder := state.Workflow.ReminderSettings.Reminders[*state.NextReminderCount]

		// Now + reminder duration = next reminder time
		state.NextReminderTime = timestamp.New(time.Now().Add(reminder.Duration.AsDuration()))
	}
}

func (w *Workflow) updateWorkflowState(ctx context.Context, state *documents.WorkflowState) error {
	tWorkflow := table.FivenetDocumentsWorkflowState

	stmt := tWorkflow.
		UPDATE(
			tWorkflow.NextReminderTime,
			tWorkflow.NextReminderCount,
			tWorkflow.AutoCloseTime,
		).
		SET(
			state.NextReminderTime,
			state.NextReminderCount,
			state.AutoCloseTime,
		).
		WHERE(jet.AND(
			tWorkflow.DocumentID.EQ(jet.Uint64(state.DocumentId)),
		))

	if _, err := stmt.ExecContext(ctx, w.db); err != nil {
		return err
	}

	return nil
}

func (w *Workflow) deleteWorkflowState(ctx context.Context, state *documents.WorkflowState) error {
	tWorkflow := table.FivenetDocumentsWorkflowState

	stmt := tWorkflow.
		DELETE().
		WHERE(jet.AND(
			tWorkflow.DocumentID.EQ(jet.Uint64(state.DocumentId)),
		)).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, w.db); err != nil {
		return err
	}

	return nil
}

func (w *Workflow) autoCloseDocument(ctx context.Context, state *documents.WorkflowState, message string) error {
	// Begin transaction
	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	// Close document and add activity
	stmt := tDocument.
		UPDATE().
		SET(
			tDocument.Closed.SET(jet.Bool(true)),
		).
		WHERE(jet.AND(
			tDocument.ID.EQ(jet.Uint64(state.DocumentId)),
		))

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	if _, err := addDocumentActivity(ctx, tx, &documents.DocActivity{
		DocumentId:   state.DocumentId,
		ActivityType: documents.DocActivityType_DOC_ACTIVITY_TYPE_STATUS_CLOSED,
		Reason:       &message,
		CreatorId:    state.Document.CreatorId,
		CreatorJob:   state.Document.CreatorJob,
	}); err != nil {
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	if state.Document == nil || state.Document.CreatorId == nil {
		return nil
	}

	// Make sure user has access to document
	userInfo, err := w.ui.GetUserInfoWithoutAccountId(ctx, *state.Document.CreatorId)
	if err != nil {
		return err
	}

	// Don't send "auto reminders" if job doesn't match document
	if state.Document.CreatorJob != userInfo.Job {
		return nil
	}

	check, err := w.access.CanUserAccessTarget(ctx, state.DocumentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_VIEW)
	if err != nil {
		return err
	}
	if !check {
		return nil
	}

	not := &notifications.Notification{
		UserId: userInfo.UserId,
		Title: &common.TranslateItem{
			Key:        "notifications.docstore.document_auto_closed.title",
			Parameters: map[string]string{"id": strconv.FormatUint(state.DocumentId, 10)},
		},
		Content: &common.TranslateItem{
			Key: "notifications.docstore.document_auto_closed.content",
			Parameters: map[string]string{
				"id":      strconv.FormatUint(state.DocumentId, 10),
				"title":   state.Document.Title,
				"message": message,
			},
		},
		Type:     notifications.NotificationType_NOTIFICATION_TYPE_INFO,
		Category: notifications.NotificationCategory_NOTIFICATION_CATEGORY_DOCUMENT,
		Data: &notifications.Data{
			Link: &notifications.Link{
				To: fmt.Sprintf("/documents/%d", state.DocumentId),
			},
		},
	}

	if err := w.notif.NotifyUser(ctx, not); err != nil {
		return err
	}

	return nil
}

func (w *Workflow) sendDocumentReminder(ctx context.Context, documentId uint64, userId int32, document *documents.DocumentShort, message *string, singleReminder bool) error {
	// Make sure user has access to document
	userInfo, err := w.ui.GetUserInfoWithoutAccountId(ctx, userId)
	if err != nil {
		return err
	}

	// Don't send "auto reminders" if job doesn't match document
	if !singleReminder && document.CreatorJob != userInfo.Job {
		return nil
	}

	check, err := w.access.CanUserAccessTarget(ctx, documentId, userInfo, documents.AccessLevel_ACCESS_LEVEL_VIEW)
	if err != nil {
		return err
	}
	if !check {
		return nil
	}

	not := &notifications.Notification{
		UserId: userId,
		Title: &common.TranslateItem{
			Key:        "notifications.docstore.document_reminder.title",
			Parameters: map[string]string{"id": strconv.FormatUint(documentId, 10)},
		},
		Content: &common.TranslateItem{
			Key: "notifications.docstore.document_reminder.content",
			Parameters: map[string]string{
				"id":    strconv.FormatUint(documentId, 10),
				"title": document.Title,
			},
		},
		Type:     notifications.NotificationType_NOTIFICATION_TYPE_INFO,
		Category: notifications.NotificationCategory_NOTIFICATION_CATEGORY_DOCUMENT,
		Data: &notifications.Data{
			Link: &notifications.Link{
				To: fmt.Sprintf("/documents/%d", documentId),
			},
		},
	}
	if message != nil {
		not.Title.Key = "notifications.docstore.document_reminder_with_message.title"

		not.Content.Key = "notifications.docstore.document_reminder_with_message.content"
		not.Content.Parameters["message"] = *message
	}

	if err := w.notif.NotifyUser(ctx, not); err != nil {
		return err
	}

	return nil
}
