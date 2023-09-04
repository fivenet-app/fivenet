package events

import (
	"context"

	"github.com/galexrt/fivenet/pkg/config"
	"github.com/nats-io/nats.go"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Subject string

type Topic string

type Type string

type Eventus struct {
	logger *zap.Logger

	NC *nats.Conn
	JS nats.JetStreamContext
}

type Params struct {
	fx.In

	LC     fx.Lifecycle
	Logger *zap.Logger
	Config *config.Config
}

func New(p Params) (*Eventus, error) {
	// Connect to NATS
	nc, err := nats.Connect(p.Config.NATS.URL, nats.Name("FiveNet"))
	if err != nil {
		return nil, err
	}

	// Half of `nats.go` `defaultAsyncPubAckInflight = 4000`
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(2000))
	if err != nil {
		return nil, err
	}

	events := &Eventus{
		logger: p.Logger.Named("eventus"),
		NC:     nc,
		JS:     js,
	}

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		return events.NC.Drain()
	}))

	return events, nil
}
