package events

import (
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

func NewEventus(logger *zap.Logger, url string) (*Eventus, error) {
	// Connect to NATS
	nc, err := nats.Connect(url, nats.Name("FiveNet"),
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
