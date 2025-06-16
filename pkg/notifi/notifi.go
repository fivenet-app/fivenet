package notifi

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/notifications"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// tNots is a reference to the notifications table in the database.
var tNots = table.FivenetNotifications

// INotifi defines the interface for sending notifications to users.
type INotifi interface {
	// NotifyUser inserts a notification for a user and publishes it asynchronously.
	NotifyUser(ctx context.Context, not *notifications.Notification) error
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
		return n.registerEvents(ctx)
	}))

	return n
}

// NotifyUser inserts a notification for a user and publishes it asynchronously.
func (n *Notifi) NotifyUser(ctx context.Context, not *notifications.Notification) error {
	nId, err := n.insertNotification(ctx, not)
	if err != nil {
		n.logger.Error("failed to insert notification into database", zap.Error(err))
		return err
	}

	not.Id = uint64(nId)
	data, err := proto.Marshal(&notifications.UserEvent{
		Data: &notifications.UserEvent_Notification{
			Notification: not,
		},
	})
	if err != nil {
		n.logger.Error("failed to proto marshal notification", zap.Error(err))
		return err
	}

	if _, err := n.js.PublishAsync(ctx, fmt.Sprintf("%s.%s.%d", BaseSubject, UserTopic, not.UserId), data); err != nil {
		n.logger.Error("failed to publish notification message", zap.Error(err))
		return err
	}

	return nil
}

// insertNotification inserts a notification into the database and returns the new notification ID.
func (n *Notifi) insertNotification(ctx context.Context, not *notifications.Notification) (int64, error) {
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
			not.UserId,
			not.Title,
			not.Type,
			not.Content,
			not.Category,
			not.Data,
		)

	res, err := stmt.ExecContext(ctx, n.db)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}
