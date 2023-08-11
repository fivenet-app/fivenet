package notifi

import (
	"errors"
	"fmt"

	"github.com/galexrt/fivenet/pkg/events"
	"github.com/nats-io/nats.go"
)

const (
	BaseSubject events.Subject = "notifi"

	UserNotification events.Type = "user"
)

func (n *Notifi) registerEvents() error {
	cfg := &nats.StreamConfig{
		Name:      "NOTIFI",
		Retention: nats.InterestPolicy,
		Subjects:  []string{fmt.Sprintf("%s.>", BaseSubject)},
	}

	if _, err := n.events.JS.AddStream(cfg); err != nil {
		if !errors.Is(nats.ErrStreamNameAlreadyInUse, err) {
			return err
		}

		if _, err := n.events.JS.UpdateStream(cfg); err != nil {
			return err
		}
	}

	return nil
}
