package notifi

import (
	"context"
	"database/sql"

	"github.com/galexrt/fivenet/gen/go/proto/resources/notifications"
	"github.com/galexrt/fivenet/query/fivenet/table"
	"go.uber.org/zap"
)

var (
	tNots = table.FivenetNotifications
)

type Types string

const (
	SuccessType Types = "success"
	InfoType    Types = "info"
	WarningType Types = "warning"
	ErrorType   Types = "error"
)

type INotifi interface {
	NotifyUser(ctx context.Context, not *notifications.Notification)
}

type Notifi struct {
	logger *zap.Logger
	db     *sql.DB
}

func New(logger *zap.Logger, db *sql.DB) INotifi {
	return &Notifi{
		logger: logger,
		db:     db,
	}
}

func (n *Notifi) NotifyUser(ctx context.Context, not *notifications.Notification) {
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

	if _, err := stmt.ExecContext(ctx, n.db); err != nil {
		n.logger.Error("failed to insert notification into database", zap.Error(err))
		return
	}
}
