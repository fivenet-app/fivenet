package notifi

import (
	"errors"

	"github.com/nats-io/nats.go"
)

const BaseSubject = "notifi"

const UserNotification = "user"

func (n *Notifi) registerEvents() error {
	cfg := &nats.StreamConfig{
		Name:      "NOTIFI",
		Retention: nats.InterestPolicy,
		Subjects:  []string{BaseSubject + ".>"},
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
