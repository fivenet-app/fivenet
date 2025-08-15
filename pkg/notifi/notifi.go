package notifi

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/notifications"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// INotifi defines the interface for sending notifications to users.
type INotifi interface {
	// NotifyUser inserts a notification for a user and publishes it asynchronously.
	NotifyUser(ctx context.Context, not *notifications.Notification) error
	// SendObjectEvent publishes an object event notification to the event system.
	SendObjectEvent(ctx context.Context, event *notifications.ObjectEvent) error
	// SendUserEvent publishes a user event notification to the event system.
	SendUserEvent(ctx context.Context, userId int32, event *notifications.UserEvent) error
	// SendSystemEvent publishes a system-wide event notification to the event system.
	SendSystemEvent(ctx context.Context, event *notifications.SystemEvent) error
}

// Notifi implements the INotifi interface for managing user notifications.
type Notifi struct {
	// logger is used for logging errors and information.
	logger *zap.Logger
	// db is the database connection used for storing notifications.
	db *sql.DB
	// js is the event system wrapper for publishing notifications.
	js *events.JSWrapper
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
	JS *events.JSWrapper
}

// New creates a new Notifi instance and registers event hooks.
func New(p Params) INotifi {
	n := &Notifi{
		logger: p.Logger,
		db:     p.DB,
		js:     p.JS,
	}

	// Register event hooks on application start.
	p.LC.Append(fx.StartHook(func(ctx context.Context) error {
		return n.registerStream(ctx)
	}))

	return n
}

// NotifyUser inserts a notification for a user and publishes it asynchronously.
func (n *Notifi) NotifyUser(ctx context.Context, not *notifications.Notification) error {
	nId, err := n.insertNotification(ctx, not)
	if err != nil {
		return fmt.Errorf("failed to insert notification into database. %w", err)
	}

	not.Id = nId
	data, err := proto.Marshal(&notifications.UserEvent{
		Data: &notifications.UserEvent_Notification{
			Notification: not,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to proto marshal notification. %w", err)
	}

	if _, err := n.js.PublishAsync(ctx, fmt.Sprintf("%s.%s.%d", BaseSubject, UserTopic, not.GetUserId()), data); err != nil {
		return fmt.Errorf("failed to publish notification message. %w", err)
	}

	return nil
}

// insertNotification inserts a notification into the database and returns the new notification ID.
func (n *Notifi) insertNotification(
	ctx context.Context,
	not *notifications.Notification,
) (int64, error) {
	tNots := table.FivenetNotifications

	stmt := tNots.
		INSERT(
			tNots.UserID,
			tNots.Title,
			tNots.Type,
			tNots.Content,
			tNots.Category,
			tNots.Data,
		).
		VALUES(
			not.GetUserId(),
			not.GetTitle(),
			not.GetType(),
			not.GetContent(),
			not.GetCategory(),
			not.GetData(),
		)

	res, err := stmt.ExecContext(ctx, n.db)
	if err != nil {
		return 0, fmt.Errorf("failed to insert notification into database. %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID of notification. %w", err)
	}

	return id, nil
}

func (n *Notifi) SendObjectEvent(ctx context.Context, event *notifications.ObjectEvent) error {
	if event.Id == nil {
		return errors.New("object event ID is required")
	}

	if _, err := n.js.PublishAsyncProto(ctx, fmt.Sprintf("%s.%s.%s.%d", BaseSubject, ObjectTopic, event.GetType().ToNatsKey(), event.GetId()), event); err != nil {
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
	event *notifications.UserEvent,
) error {
	if event == nil || event.Data == nil {
		return errors.New("user event data is required")
	}

	if _, err := n.js.PublishAsyncProto(ctx, fmt.Sprintf("%s.%s.%d", BaseSubject, UserTopic, userId), event); err != nil {
		return fmt.Errorf("failed to publish user %d event message. %w", userId, err)
	}

	return nil
}

func (n *Notifi) SendSystemEvent(ctx context.Context, event *notifications.SystemEvent) error {
	if event == nil || event.Data == nil {
		return errors.New("system event data is required")
	}

	if _, err := n.js.PublishAsyncProto(ctx, fmt.Sprintf("%s.%s", BaseSubject, SystemTopic), event); err != nil {
		return fmt.Errorf("failed to publish system event message. %w", err)
	}

	return nil
}
