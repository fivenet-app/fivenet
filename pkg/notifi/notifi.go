package notifi

import (
	"context"
	"database/sql"

	"github.com/galexrt/fivenet/gen/go/proto/resources/notifications"
	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/events"
	"github.com/galexrt/fivenet/query/fivenet/table"
	"github.com/nats-io/nats.go"
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
	Add(not *notifications.Notification)
}

type Notifi struct {
	logger *zap.Logger
	db     *sql.DB
	ctx    context.Context

	events    *events.Eventus
	streamCfg *nats.StreamConfig
	subs      []*nats.Subscription
}

func New(logger *zap.Logger, db *sql.DB, ctx context.Context, events *events.Eventus) *Notifi {
	return &Notifi{
		logger: logger,
		db:     db,
		ctx:    ctx,
		events: events,
		streamCfg: &nats.StreamConfig{
			Name:     "NOTIFI",
			Subjects: []string{"notifi.>"},
		},
		subs: make([]*nats.Subscription, config.C.NATS.WorkerCount),
	}
}

func (n *Notifi) Add(not *notifications.Notification) {
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

	if _, err := stmt.ExecContext(n.ctx, n.db); err != nil {
		n.logger.Error("failed to insert notification into database", zap.Error(err))
		return
	}
}
