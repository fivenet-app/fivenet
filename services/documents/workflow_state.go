package documents

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/cron"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents"
	documentsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/access"
	documentsactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/activity"
	documentsworkflow "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/workflow"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/notifications"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/croner"
	"github.com/fivenet-app/fivenet/v2026/pkg/notifi"
	"github.com/fivenet-app/fivenet/v2026/pkg/userinfo"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var tDWorkflow = table.FivenetDocumentsWorkflowState.AS("workflow_state")

var WorkflowModule = fx.Module(
	"documents.workflow",
	fx.Provide(
		NewWorkflow,
	),
)

type Workflow struct {
	logger *zap.Logger
	tracer trace.Tracer

	db    *sql.DB
	ui    userinfo.UserInfoRetriever
	notif notifi.INotifi

	access *access.Grouped[documentsaccess.DocumentJobAccess, *documentsaccess.DocumentJobAccess, documentsaccess.DocumentUserAccess, *documentsaccess.DocumentUserAccess, access.DummyQualificationAccess[documentsaccess.AccessLevel], *access.DummyQualificationAccess[documentsaccess.AccessLevel], documentsaccess.AccessLevel]
}

type WorkflowParams struct {
	fx.In

	LC fx.Lifecycle

	Logger *zap.Logger
	DB     *sql.DB
	TP     *tracesdk.TracerProvider
	Notif  notifi.INotifi
	Ui     userinfo.UserInfoRetriever
}

type WorkflowResult struct {
	fx.Out

	Workflow     *Workflow
	CronRegister croner.CronRegister `group:"cronjobregister"`
}

func NewWorkflow(p WorkflowParams) WorkflowResult {
	w := &Workflow{
		logger: p.Logger.Named("documents.workflow"),
		tracer: p.TP.Tracer("documents.workflow"),
		db:     p.DB,
		notif:  p.Notif,
		ui:     p.Ui,

		access: newAccess(p.DB),
	}

	return WorkflowResult{
		Workflow:     w,
		CronRegister: w,
	}
}

func (w *Workflow) RegisterCronjobs(ctx context.Context, registry croner.IRegistry) error {
	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     "documents.workflow_run",
		Schedule: "* * * * *", // Every minute
	}); err != nil {
		return err
	}

	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     "documents.workflow_users_run",
		Schedule: "* * * * *", // Every minute
	}); err != nil {
		return err
	}

	return nil
}

func (w *Workflow) RegisterCronjobHandlers(h *croner.Handlers) error {
	h.Add("documents.workflow_run", func(ctx context.Context, data *cron.CronjobData) error {
		ctx, span := w.tracer.Start(ctx, "documents.workflow_run")
		defer span.End()

		dest := &documentsworkflow.WorkflowCronData{
			LastDocId: 0,
		}
		if err := data.Unmarshal(dest); err != nil {
			w.logger.Warn("failed to unmarshal document workflow cron data", zap.Error(err))
		}

		if err := w.handleDocuments(ctx, dest); err != nil {
			return fmt.Errorf("error during documents workflow handling. %w", err)
		}

		// Marshal the updated cron data
		if err := data.MarshalFrom(dest); err != nil {
			return fmt.Errorf("failed to marshal updated document workflow cron data. %w", err)
		}

		return nil
	})

	h.Add("documents.workflow_users_run", func(ctx context.Context, data *cron.CronjobData) error {
		ctx, span := w.tracer.Start(ctx, "documents.workflow_users_run")
		defer span.End()

		dest := &documentsworkflow.WorkflowCronData{
			LastDocId: 0,
		}
		if err := data.Unmarshal(dest); err != nil {
			w.logger.Error("failed to unmarshal document workflow cron data", zap.Error(err))
		}

		if err := w.handleDocumentsUsers(ctx, dest); err != nil {
			return fmt.Errorf("error during documents workflow handling. %w", err)
		}

		// Marshal the updated cron data
		if err := data.MarshalFrom(dest); err != nil {
			return fmt.Errorf("failed to marshal updated document workflow cron data. %w", err)
		}

		return nil
	})

	return nil
}

type workflowState struct {
	DocumentId int64                            `alias:"document_id"`
	State      *documentsworkflow.WorkflowState `alias:"workflow_state"`
	Document   *documents.DocumentShort         `alias:"document_short"`
}

func (w *Workflow) handleDocuments(
	ctx context.Context,
	data *documentsworkflow.WorkflowCronData,
) error {
	nowTs := mysql.TimestampT(time.Now())

	tDTemplates := table.FivenetDocumentsTemplates.AS("template")

	dest := []*workflowState{}
	for {
		stmt := tDWorkflow.
			SELECT(
				tDWorkflow.DocumentID.AS("document_id"),
				tDWorkflow.DocumentID,
				tDWorkflow.NextReminderTime,
				tDWorkflow.NextReminderCount,
				tDWorkflow.AutoCloseTime,
				tDWorkflow.ReminderCount,
				tDTemplates.Workflow.AS("workflow_state.workflow"),
				tDocumentShort.Title,
				tDocumentShort.CreatorID,
				tDocumentShort.CreatorJob,
			).
			FROM(
				tDWorkflow.
					INNER_JOIN(tDocumentShort,
						mysql.AND(
							tDocumentShort.ID.EQ(tDWorkflow.DocumentID),
							tDocumentShort.DeletedAt.IS_NULL(),
						),
					).
					LEFT_JOIN(tDTemplates,
						mysql.AND(
							tDTemplates.ID.EQ(tDocumentShort.TemplateID),
							tDTemplates.DeletedAt.IS_NULL(),
						),
					),
			).
			WHERE(mysql.AND(
				tDWorkflow.DocumentID.GT(mysql.Int64(data.GetLastDocId())),
				mysql.AND( // Only auto close and auto remind docs that aren't closed and have an owner
					tDocumentShort.Closed.IS_FALSE(),
					mysql.OR(
						tDWorkflow.NextReminderTime.LT_EQ(nowTs),
						tDWorkflow.AutoCloseTime.LT_EQ(nowTs),
					),
				),
			)).
			LIMIT(250)

		if err := stmt.QueryContext(ctx, w.db, &dest); err != nil {
			if !errors.Is(err, qrm.ErrNoRows) {
				return err
			}
		}

		if data.GetLastDocId() == 0 && len(dest) == 0 {
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

	var wg sync.WaitGroup

	// Run at max 3 handlers at once
	workChannel := make(chan *workflowState, 3)

	wg.Add(1)
	go func() {
		defer wg.Done()

		for state := range workChannel {
			wg.Add(1)
			go func() {
				defer wg.Done()

				if state == nil || state.State == nil {
					return
				}

				if err := w.handleWorkflowState(ctx, state); err != nil {
					w.logger.Error("error during workflow state handling",
						zap.Int64("document_id", state.DocumentId), zap.Error(err))
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

func (w *Workflow) handleWorkflowState(
	ctx context.Context,
	st *workflowState,
) error {
	state := st.State
	doc := st.Document

	if state.GetAutoCloseTime() != nil && time.Since(state.GetAutoCloseTime().AsTime()) > 0 {
		if state.GetWorkflow() != nil && state.GetWorkflow().GetAutoCloseSettings() != nil &&
			state.GetWorkflow().GetAutoClose() &&
			state.GetWorkflow().GetAutoCloseSettings().GetMessage() != "" {
			// Auto close document and null "next reminder time"
			if err := w.autoCloseDocument(ctx, st, state.GetWorkflow().GetAutoCloseSettings().GetMessage()); err != nil {
				return fmt.Errorf("failed to auto close document. %w", err)
			}
		}

		// Delete document workflow state, auto reminders are not sent for a closed document
		return w.deleteWorkflowState(ctx, state)
	} else if state.GetNextReminderTime() != nil && time.Since(state.GetNextReminderTime().AsTime()) > 0 {
		if doc != nil && doc.GetCreatorId() != 0 {
			var reminderMessage string
			if reminder := w.getAutoReminder(state); reminder != nil && reminder.GetMessage() != "" {
				reminderMessage = reminder.GetMessage()
			}

			// Send notification when the document has a creator that is still part of the document's job
			if err := w.sendDocumentReminder(ctx, state.GetDocumentId(), doc.GetCreatorId(), doc, reminderMessage, false); err != nil {
				return fmt.Errorf("failed to send document reminder. %w", err)
			}

			w.updateAutoReminderTime(state)
			state.ReminderCount++
		} else {
			state.NextReminderTime = nil
			state.NextReminderCount = nil
			state.ReminderCount = 0
		}
	}

	// Make sure to delete the document workflow state under one of the following conditions:
	// * document doesn't exist anymore
	// * if we don't have a doc creator anymore
	// * reached the max reminder count
	if doc == nil ||
		doc.GetCreatorId() == 0 ||
		state.GetReminderCount() >= documentsworkflow.MaxReminderCount {
		return w.deleteWorkflowState(ctx, state)
	}

	if err := w.updateWorkflowState(ctx, state); err != nil {
		return fmt.Errorf("failed to update workflow state. %w", err)
	}

	return nil
}

func (w *Workflow) getAutoReminder(
	state *documentsworkflow.WorkflowState,
) *documentsworkflow.Reminder {
	count := int32(0)
	if state.GetNextReminderCount() > 0 {
		count = state.GetNextReminderCount()
	}

	if state.GetWorkflow() == nil || state.GetWorkflow().GetReminderSettings() == nil ||
		len(state.GetWorkflow().GetReminderSettings().GetReminders()) <= int(count) {
		return nil
	}

	return state.GetWorkflow().GetReminderSettings().GetReminders()[count]
}

func (w *Workflow) updateAutoReminderTime(state *documentsworkflow.WorkflowState) {
	if state.GetWorkflow() == nil || state.GetWorkflow().GetReminderSettings() == nil ||
		!state.GetWorkflow().GetReminder() ||
		len(state.GetWorkflow().GetReminderSettings().GetReminders()) == 0 {
		state.NextReminderTime = nil
		state.NextReminderCount = nil
		return
	}

	if state.NextReminderCount == nil {
		zero := int32(0)
		state.NextReminderCount = &zero
	} else {
		//nolint:protogetter // The value is updated via the pointer
		*state.NextReminderCount++
	}

	if len(
		state.GetWorkflow().GetReminderSettings().GetReminders(),
	) <= int(
		state.GetNextReminderCount(),
	) {
		*state.NextReminderCount = 0
	}

	// No reminders? How did we end up here? Unset reminder time
	if len(state.GetWorkflow().GetReminderSettings().GetReminders()) == 0 {
		state.NextReminderTime = nil
	}

	reminder := state.GetWorkflow().GetReminderSettings().GetReminders()[state.GetNextReminderCount()]

	// Now + reminder duration = next reminder time
	state.NextReminderTime = timestamp.New(time.Now().Add(reminder.GetDuration().AsDuration()))
}

func (w *Workflow) updateWorkflowState(
	ctx context.Context,
	state *documentsworkflow.WorkflowState,
) error {
	tWorkflow := table.FivenetDocumentsWorkflowState

	stmt := tWorkflow.
		UPDATE(
			tWorkflow.NextReminderTime,
			tWorkflow.NextReminderCount,
			tWorkflow.AutoCloseTime,
			tWorkflow.ReminderCount,
		).
		SET(
			state.GetNextReminderTime(),
			state.GetNextReminderCount(),
			state.GetAutoCloseTime(),
			state.GetReminderCount(),
		).
		WHERE(tWorkflow.DocumentID.EQ(mysql.Int64(state.GetDocumentId()))).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, w.db); err != nil {
		return err
	}

	return nil
}

func (w *Workflow) deleteWorkflowState(
	ctx context.Context,
	state *documentsworkflow.WorkflowState,
) error {
	tWorkflow := table.FivenetDocumentsWorkflowState

	stmt := tWorkflow.
		DELETE().
		WHERE(tWorkflow.DocumentID.EQ(mysql.Int64(state.GetDocumentId()))).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, w.db); err != nil {
		return err
	}

	return nil
}

func (w *Workflow) autoCloseDocument(
	ctx context.Context,
	st *workflowState,
	message string,
) error {
	state := st.State
	doc := st.Document

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
			tDocument.Closed.SET(mysql.Bool(true)),
		).
		WHERE(mysql.AND(
			tDocument.ID.EQ(mysql.Int64(state.DocumentId)),
		))

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	if _, err := addDocumentActivity(ctx, tx, &documentsactivity.DocActivity{
		DocumentId:   state.DocumentId,
		ActivityType: documentsactivity.DocActivityType_DOC_ACTIVITY_TYPE_STATUS_CLOSED,
		Reason:       &message,
		CreatorId:    doc.CreatorId,
		CreatorJob:   doc.GetCreatorJob(),
	}); err != nil {
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	if doc == nil || doc.GetCreatorId() == 0 {
		return nil
	}

	// Make sure user has access to document
	userInfo, err := w.ui.GetUserInfoWithoutAccountId(ctx, doc.GetCreatorId())
	if err != nil {
		return err
	}

	// Don't send "auto reminders" if job doesn't match document
	if doc.GetCreatorJob() != userInfo.GetJob() {
		return nil
	}

	check, err := w.access.CanUserAccessTarget(
		ctx,
		state.GetDocumentId(),
		userInfo,
		documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW,
	)
	if err != nil {
		return err
	}
	if !check {
		return nil
	}

	not := &notifications.Notification{
		UserId: userInfo.GetUserId(),
		Title: &common.I18NItem{
			Key:        "notifications.documents.document_auto_closed.title",
			Parameters: map[string]string{"id": strconv.FormatInt(state.GetDocumentId(), 10)},
		},
		Content: &common.I18NItem{
			Key: "notifications.documents.document_auto_closed.content",
			Parameters: map[string]string{
				"id":      strconv.FormatInt(state.GetDocumentId(), 10),
				"title":   doc.GetTitle(),
				"message": message,
			},
		},
		Type:     notifications.NotificationType_NOTIFICATION_TYPE_INFO,
		Category: notifications.NotificationCategory_NOTIFICATION_CATEGORY_DOCUMENT,
		Data: &notifications.Data{
			Link: &notifications.Link{
				To: fmt.Sprintf("/documents/%d", state.GetDocumentId()),
			},
		},
	}

	if err := w.notif.NotifyUser(ctx, not); err != nil {
		return err
	}

	return nil
}

func (w *Workflow) sendDocumentReminder(
	ctx context.Context,
	documentId int64,
	userId int32,
	doc *documents.DocumentShort,
	message string,
	singleReminder bool,
) error {
	// Make sure user has access to document
	userInfo, err := w.ui.GetUserInfoWithoutAccountId(ctx, userId)
	if err != nil {
		return err
	}

	// Don't send "auto reminders" if job doesn't match document
	if !singleReminder && doc.GetCreatorJob() != userInfo.GetJob() {
		return nil
	}

	check, err := w.access.CanUserAccessTarget(
		ctx,
		documentId,
		userInfo,
		documentsaccess.AccessLevel_ACCESS_LEVEL_VIEW,
	)
	if err != nil {
		return err
	}
	if !check {
		return nil
	}

	not := &notifications.Notification{
		UserId: userId,
		Title: &common.I18NItem{
			Key:        "notifications.documents.document_reminder.title",
			Parameters: map[string]string{"id": strconv.FormatInt(documentId, 10)},
		},
		Content: &common.I18NItem{
			Key: "notifications.documents.document_reminder.content",
			Parameters: map[string]string{
				"id":    strconv.FormatInt(documentId, 10),
				"title": doc.GetTitle(),
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
	if message != "" {
		not.Title.Key = "notifications.documents.document_reminder_with_message.title"

		not.Content.Key = "notifications.documents.document_reminder_with_message.content"
		not.Content.Parameters["message"] = message
	}

	if err := w.notif.NotifyUser(ctx, not); err != nil {
		return err
	}

	return nil
}
