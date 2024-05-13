package notifi

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/notifications"
	"github.com/fivenet-app/fivenet/pkg/events"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

var (
	tNots = table.FivenetNotifications
)

type INotifi interface {
	NotifyUser(ctx context.Context, not *notifications.Notification) error
}

type Notifi struct {
	logger *zap.Logger
	db     *sql.DB
	js     *events.JSWrapper
}

type Params struct {
	fx.In

	LC     fx.Lifecycle
	Logger *zap.Logger
	DB     *sql.DB
	JS     *events.JSWrapper
}

func New(p Params) INotifi {
	n := &Notifi{
		logger: p.Logger,
		db:     p.DB,
		js:     p.JS,
	}

	p.LC.Append(fx.StartHook(func(ctx context.Context) error {
		return n.registerEvents(ctx)
	}))

	return n
}

func (n *Notifi) NotifyUser(ctx context.Context, not *notifications.Notification) error {
	nId, err := n.insertNotification(ctx, not)
	if err != nil {
		n.logger.Error("failed to insert notification into database", zap.Error(err))
		return err
	}

	not.Id = uint64(nId)
	data, err := proto.Marshal(not)
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
