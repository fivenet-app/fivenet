package events

import (
	"github.com/galexrt/fivenet/pkg/config"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type IEvent interface {
	GetName() string
}

type Eventus struct {
	logger *zap.Logger

	NC *nats.Conn
	JS nats.JetStreamContext
}

func NewEventus(logger *zap.Logger) (*Eventus, error) {
	// Connect to NATS
	nc, err := nats.Connect(config.C.NATS.URL, nats.Name("FiveNet"),
		nats.NoEcho())
	if err != nil {
		return nil, err
	}

	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		return nil, err
	}

	return &Eventus{
		logger: logger,
		NC:     nc,
		JS:     js,
	}, nil
}

func (e *Eventus) Stop() error {
	return e.NC.Drain()
}
