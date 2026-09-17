package notifi

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/notifications"
	notificationsclientview "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/notifications/clientview"
	notificationsevents "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/notifications/events"
	pbtimestamp "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/pkg/events"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	notificationsstore "github.com/fivenet-app/fivenet/v2026/stores/notifications"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// INotifi defines the interface for sending notifications to users.
type INotifi interface {
	// NotifyUser inserts a notification for a user and publishes it asynchronously.
	NotifyUser(ctx context.Context, not *notifications.Notification) error
	// PrepareUserNotification stores the inbox record on db and returns a
	// publisher which must only be called after its surrounding transaction commits.
	PrepareUserNotification(
		ctx context.Context,
		db qrm.DB,
		not *notifications.Notification,
	) (func(context.Context) error, error)
	// SendObjectEvent publishes an object event notification to the event system.
	SendObjectEvent(ctx context.Context, event *notificationsclientview.ObjectEvent) error
	// SendUserEvent publishes a user event notification to the event system.
	SendUserEvent(ctx context.Context, userId int32, event *notificationsevents.UserEvent) error
	// SendSystemEvent publishes a system-wide event notification to the event system.
	SendSystemEvent(ctx context.Context, event *notificationsevents.SystemEvent) error
}

// UserNotificationParams contains the standard metadata shared by durable and
// transient user notifications.
type UserNotificationParams struct {
	UserID      int32
	Type        notifications.NotificationType
	Category    notifications.NotificationCategory
	Kind        notifications.NotificationKind
	ActorUserID *int32
	EntityType  *string
	EntityID    *int64
	Title       *common.I18NItem
	Content     *common.I18NItem
	Data        *notifications.Data
}

func NewUserNotification(p UserNotificationParams) *notifications.Notification {
	return &notifications.Notification{
		UserId:      p.UserID,
		Type:        p.Type,
		Category:    p.Category,
		Kind:        p.Kind,
		ActorUserId: p.ActorUserID,
		EntityType:  p.EntityType,
		EntityId:    p.EntityID,
		Title:       p.Title,
		Content:     p.Content,
		Data:        p.Data,
	}
}

// PublishAfterCommit runs callbacks prepared inside a successfully committed
// transaction. Publishing is best-effort: a failure cannot roll back the
// already-persisted mutation, so it is logged rather than returned to callers.
func PublishAfterCommit(
	ctx context.Context,
	logger *zap.Logger,
	operation string,
	callbacks ...func(context.Context) error,
) {
	for _, publish := range callbacks {
		if publish == nil {
			continue
		}
		if err := publish(ctx); err != nil {
			logger.Warn(
				"failed to publish notification",
				zap.String("notification_operation", operation),
				zap.Error(err),
			)
		}
	}
}

func noopPublish(context.Context) error {
	return nil
}

// Notifi implements the INotifi interface for managing user notifications.
type Notifi struct {
	// logger is used for logging errors and information.
	logger *zap.Logger
	// db is the database connection used for storing notifications.
	db *sql.DB
	// js is the event system wrapper for publishing notifications.
	js          *events.JSWrapper
	preferences notificationsstore.IStore
}

// Params contains dependencies for constructing a Notifi instance.
type Params struct {
	fx.In

	// LC is the application lifecycle for registering hooks.
	LC fx.Lifecycle
	// Logger is the logger instance for logging.
	Logger *zap.Logger
	// DB is the database connection.
	DB *sql.DB
	// JS is the event system wrapper.
	JS                *events.JSWrapper
	NotificationStore notificationsstore.IStore
}

// New creates a new Notifi instance and registers event hooks.
func New(p Params) INotifi {
	n := &Notifi{
		logger:      p.Logger,
		db:          p.DB,
		js:          p.JS,
		preferences: p.NotificationStore,
	}

	// Register event hooks on application start.
	p.LC.Append(fx.StartHook(func(ctx context.Context) error {
		return n.registerStream(ctx)
	}))

	return n
}

// NotifyUser inserts a notification for a user and publishes it asynchronously.
func (n *Notifi) NotifyUser(ctx context.Context, not *notifications.Notification) error {
	publish, err := n.PrepareUserNotification(ctx, n.db, not)
	if err != nil {
		return err
	}
	return publish(ctx)
}

// PrepareUserNotification persists an inbox notification in db. The returned
// callback deliberately defers NATS publication until the caller has committed
// the transaction which caused the notification.
func (n *Notifi) PrepareUserNotification(
	ctx context.Context,
	db qrm.DB,
	not *notifications.Notification,
) (func(context.Context) error, error) {
	delivery, err := n.preferences.ResolveDelivery(
		ctx,
		not.GetUserId(),
		not.GetCategory(),
		not.GetKind(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve notification delivery preferences. %w", err)
	}
	not.Delivery = delivery
	if not.GetCreatedAt() == nil {
		not.CreatedAt = pbtimestamp.Now()
	}

	if delivery.GetInboxEnabled() {
		nId, err := n.insertNotification(ctx, db, not)
		if err != nil {
			return nil, fmt.Errorf("failed to insert notification into database. %w", err)
		}
		not.Id = nId
	}

	if !delivery.GetInboxEnabled() && !delivery.GetToastEnabled() && !delivery.GetSoundEnabled() {
		// Keep the post-commit callback contract total: callers can always invoke
		// the returned function, even when this notification has no delivery path.
		return noopPublish, nil
	}

	return func(ctx context.Context) error {
		var unreadCount *int64
		if delivery.GetInboxEnabled() {
			count, err := n.preferences.CountUnread(ctx, not.GetUserId())
			if err != nil {
				n.logger.Warn(
					"failed to count unread notifications for event",
					zap.Int32("user_id", not.GetUserId()),
					zap.Error(err),
				)
			} else {
				unreadCount = &count
			}
		}

		return n.publishUserNotification(ctx, not, unreadCount)
	}, nil
}

func (n *Notifi) publishUserNotification(
	ctx context.Context,
	not *notifications.Notification,
	unreadCount *int64,
) error {
	data, err := proto.Marshal(&notificationsevents.UserEvent{
		UnreadCount: unreadCount,
		Data: &notificationsevents.UserEvent_Notification{
			Notification: not,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to proto marshal notification. %w", err)
	}

	if _, err := n.js.PublishAsync(
		ctx,
		fmt.Sprintf("%s.%s.%d", BaseSubject, UserTopic, not.GetUserId()),
		data,
	); err != nil {
		return fmt.Errorf("failed to publish notification message. %w", err)
	}

	return nil
}

// insertNotification inserts a notification into the database and returns the new notification ID.
func (n *Notifi) insertNotification(
	ctx context.Context,
	db qrm.DB,
	not *notifications.Notification,
) (int64, error) {
	tNots := table.FivenetNotifications

	stmt := tNots.
		INSERT(
			tNots.UserID,
			tNots.ActorUserID,
			tNots.EntityType,
			tNots.EntityID,
			tNots.Title,
			tNots.Type,
			tNots.Content,
			tNots.Category,
			tNots.Kind,
			tNots.Data,
			tNots.Starred,
		).
		VALUES(
			not.GetUserId(),
			dbutils.Int32P(not.GetActorUserId()),
			dbutils.StringPP(not.EntityType),
			dbutils.Int64P(not.GetEntityId()),
			not.GetTitle(),
			not.GetType(),
			not.GetContent(),
			not.GetCategory(),
			not.GetKind(),
			not.GetData(),
			not.GetStarred(),
		)

	res, err := stmt.ExecContext(ctx, db)
	if err != nil {
		return 0, fmt.Errorf("failed to insert notification into database. %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID of notification. %w", err)
	}

	return id, nil
}

func (n *Notifi) SendObjectEvent(
	ctx context.Context,
	event *notificationsclientview.ObjectEvent,
) error {
	if event == nil {
		return errors.New("object event is required")
	}

	if event.Id == nil {
		return errors.New("object event ID is required")
	}

	spec, ok := event.GetType().Spec()
	if !ok {
		return fmt.Errorf("unsupported object event type: %s", event.GetType().String())
	}

	if _, err := n.js.PublishAsyncProto(
		ctx,
		fmt.Sprintf(
			"%s.%s.%s.%d",
			BaseSubject,
			ObjectTopic,
			spec.NatsKey,
			event.GetId(),
		),
		event,
	); err != nil {
		return fmt.Errorf(
			"failed to publish object %s event message. %w",
			event.GetType().String(),
			err,
		)
	}

	return nil
}

func (n *Notifi) SendUserEvent(
	ctx context.Context,
	userId int32,
	event *notificationsevents.UserEvent,
) error {
	if event == nil || event.Data == nil {
		return errors.New("user event data is required")
	}

	if _, err := n.js.PublishAsyncProto(
		ctx,
		fmt.Sprintf("%s.%s.%d", BaseSubject, UserTopic, userId),
		event,
	); err != nil {
		return fmt.Errorf("failed to publish user %d event message. %w", userId, err)
	}

	return nil
}

func (n *Notifi) SendSystemEvent(
	ctx context.Context,
	event *notificationsevents.SystemEvent,
) error {
	if event == nil || event.Data == nil {
		return errors.New("system event data is required")
	}

	if _, err := n.js.PublishAsyncProto(
		ctx,
		fmt.Sprintf("%s.%s", BaseSubject, SystemTopic),
		event,
	); err != nil {
		return fmt.Errorf("failed to publish system event message. %w", err)
	}

	return nil
}
