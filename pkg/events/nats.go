package events

import (
	"github.com/galexrt/fivenet/pkg/config"
	"github.com/nats-io/nats.go"
)

func Connect() error {
	// Connect to NATS
	nc, err := nats.Connect(config.C.NATS.URL, nats.Name("FiveNet"))
	if err != nil {
		return err
	}

	// Create JetStream Context
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		return err
	}

	_ = js

	return nil
}
