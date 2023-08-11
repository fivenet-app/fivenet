package notifi

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/galexrt/fivenet/gen/go/proto/resources/notifications"
	"github.com/galexrt/fivenet/pkg/events"
	"github.com/galexrt/fivenet/query/fivenet/table"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
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
	events *events.Eventus
}

type Params struct {
	fx.In

	LC     fx.Lifecycle
	Logger *zap.Logger
	DB     *sql.DB
	Events *events.Eventus
}

func New(p Params) INotifi {
	n := &Notifi{
		logger: p.Logger,
		db:     p.DB,
		events: p.Events,
	}

	p.LC.Append(fx.StartHook(func(_ context.Context) error {
		return n.registerEvents()
	}))

	return n
}

func (n *Notifi) NotifyUser(ctx context.Context, not *notifications.Notification) {
	nId, err := n.insertNotification(ctx, not)
	if err != nil {
		n.logger.Error("failed to insert notification into database", zap.Error(err))
		return
	}

	not.Id = uint64(nId)
	data, err := proto.Marshal(not)
	if err != nil {
		n.logger.Error("failed to proto marshal notification", zap.Error(err))
		return
	}

	n.events.JS.Publish(fmt.Sprintf("%s.%s.%d", BaseSubject, UserNotification, not.UserId), data)
}

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
