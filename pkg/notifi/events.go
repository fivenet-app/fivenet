package notifi

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/galexrt/fivenet/pkg/events"
	"github.com/nats-io/nats.go"
)

const (
	BaseSubject events.Subject = "notifi"

	UserNotification events.Type = "user"
)

func (n *Notifi) registerEvents(ctx context.Context) error {
	cfg := &nats.StreamConfig{
		Name:      "NOTIFI",
		Retention: nats.InterestPolicy,
		Subjects:  []string{fmt.Sprintf("%s.>", BaseSubject)},
		Discard:   nats.DiscardOld,
		MaxAge:    30 * time.Minute,
	}

	if _, err := n.js.UpdateStream(cfg); err != nil {
		if !errors.Is(nats.ErrStreamNotFound, err) {
			return err
		}

		if _, err := n.js.AddStream(cfg); err != nil {
			return err
		}
	}

	return nil
}
