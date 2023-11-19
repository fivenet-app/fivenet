package sentry

import (
	"context"
	"time"

	"github.com/galexrt/fivenet/pkg/config"
	"github.com/getsentry/sentry-go"
	"go.uber.org/fx"
)

var Module = fx.Module("sentry",
	fx.Provide(),
)

type Params struct {
	fx.In

	LC     fx.Lifecycle
	Config *config.Config
}

func New(p Params) (*sentry.Client, error) {
	if p.Config.Sentry.Environment != "dev" && p.Config.Sentry.ServerDSN != "" {
		client, err := sentry.NewClient(sentry.ClientOptions{
			Dsn: p.Config.Sentry.ServerDSN,
		})
		if err != nil {
			return nil, err
		}

		hub := sentry.CurrentHub()
		hub.BindClient(client)

		p.LC.Append(fx.StopHook(func(ctx context.Context) error {
			// Since sentry emits events in the background we need to make sure
			// they are sent before we shut down
			sentry.Flush(time.Second * 5)

			return nil
		}))

		return client, nil
	}

	return nil, nil
}
