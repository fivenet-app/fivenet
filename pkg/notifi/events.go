package notifi

import (
	"fmt"
	"time"

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
		Discard:   nats.DiscardOld,
		MaxAge:    30 * time.Minute,
	}

	if _, err := n.events.JS.CreateOrUpdateStream(cfg); err != nil {
		return err
	}

	return nil
}
