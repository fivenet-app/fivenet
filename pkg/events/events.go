package events

import (
	"context"

	"github.com/galexrt/fivenet/pkg/config"
	"github.com/nats-io/nats.go"
	"go.uber.org/fx"
)

type Subject string

type Topic string

type Type string

type Params struct {
	fx.In

	LC     fx.Lifecycle
	Config *config.Config
}

type Result struct {
	fx.Out

	JS nats.JetStreamContext
}

func New(p Params) (res Result, err error) {
	// Connect to NATS
	nc, err := nats.Connect(p.Config.NATS.URL, nats.Name("FiveNet"))
	if err != nil {
		return res, err
	}

	// Default `defaultAsyncPubAckInflight` is `4000` (`nats.go`)
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		return res, err
	}
	res.JS = js

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		return nc.Drain()
	}))

	return res, nil
}
