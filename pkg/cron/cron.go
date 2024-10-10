package cron

import (
	"context"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common/cron"
	"github.com/fivenet-app/fivenet/pkg/events"
	"github.com/fivenet-app/fivenet/pkg/nats/store"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger *zap.Logger
	JS     *events.JSWrapper
}

type Cron struct {
	js *events.JSWrapper
	cs *store.Store[cron.Cronjob, *cron.Cronjob]
}

func New(p Params) (*Cron, error) {
	c := &Cron{
		js: p.JS,
	}

	p.LC.Append(fx.StartHook(func(ctx context.Context) error {
		if err := c.registerSubscriptions(ctx); err != nil {
			return err
		}

		cron, err := store.New[cron.Cronjob, *cron.Cronjob](ctx, p.Logger, p.JS, "cron", nil)
		if err != nil {
			return err
		}
		c.cs = cron

		if err := cron.Start(ctx); err != nil {
			return err
		}

		return nil
	}))

	return c, nil
}
