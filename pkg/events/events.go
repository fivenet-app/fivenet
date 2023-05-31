package events

import (
	"fmt"

	"github.com/galexrt/fivenet/pkg/config"
	"github.com/nats-io/nats.go"
)

type IEvent interface {
	GetName() string
}

type Eventus struct {
	nc  *nats.Conn
	js  nats.JetStreamContext
	cfg *nats.StreamConfig

	subs []*nats.Subscription
}

func NewEventus() (*Eventus, error) {
	// Connect to NATS
	nc, err := nats.Connect(config.C.NATS.URL, nats.Name("FiveNet"))
	if err != nil {
		return nil, err
	}

	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		return nil, err
	}

	cfg := &nats.StreamConfig{
		Name:      "EVENTS",
		Retention: nats.WorkQueuePolicy,
		Subjects:  []string{"events.>"},
	}

	if _, err := js.AddStream(cfg); err != nil {
		return nil, err
	}

	return &Eventus{
		nc:   nc,
		js:   js,
		cfg:  cfg,
		subs: make([]*nats.Subscription, config.C.NATS.WorkerCount),
	}, nil
}

func (e *Eventus) Start() error {
	e.js.AddConsumer(e.cfg.Name, &nats.ConsumerConfig{
		Durable:        "event-processor",
		DeliverSubject: "my-subject",
		DeliverGroup:   "event-processor",
		AckPolicy:      nats.AckExplicitPolicy,
	})
	defer e.js.DeleteConsumer(e.cfg.Name, "event-processor")

	for i := 0; i < len(e.subs); i++ {
		sub, err := e.nc.QueueSubscribe("my-subject", "event-processor", func(msg *nats.Msg) {
			fmt.Printf("sub1: received message %q\n", msg.Subject)
			msg.Ack()
		})
		if err != nil {
			return err
		}
		e.subs[i] = sub
	}

	return nil
}

func (e *Eventus) Stop() error {
	return e.nc.Drain()
}
