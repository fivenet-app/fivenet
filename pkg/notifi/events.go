package notifi

import (
	"context"
	"fmt"
	"time"

	"github.com/galexrt/fivenet/pkg/events"
	natsutils "github.com/galexrt/fivenet/pkg/nats"
	"github.com/nats-io/nats.go"
)

const (
	BaseSubject events.Subject = "notifi"

	UserNotification events.Type = "user"
)

func (n *Notifi) registerEvents(ctx context.Context) error {
	cfg := &nats.StreamConfig{
		Name:        "NOTIFI",
		Description: natsutils.Description,
		Retention:   nats.InterestPolicy,
		Subjects:    []string{fmt.Sprintf("%s.>", BaseSubject)},
		Discard:     nats.DiscardOld,
		MaxAge:      30 * time.Minute,
	}
	if _, err := natsutils.CreateOrUpdateStream(n.js, cfg); err != nil {
		return err
	}

	return nil
}
