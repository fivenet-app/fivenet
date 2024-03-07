package notifi

import (
	"context"
	"fmt"
	"time"

	"github.com/galexrt/fivenet/pkg/events"
	natsutils "github.com/galexrt/fivenet/pkg/nats"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	StreamName = "NOTIFI"

	BaseSubject events.Subject = "notifi"

	UserNotification events.Type = "user"
)

func (n *Notifi) registerEvents(ctx context.Context) error {
	cfg := jetstream.StreamConfig{
		Name:        StreamName,
		Description: natsutils.Description,
		Retention:   jetstream.InterestPolicy,
		Subjects:    []string{fmt.Sprintf("%s.>", BaseSubject)},
		Discard:     jetstream.DiscardOld,
		MaxAge:      30 * time.Minute,
	}
	if _, err := natsutils.CreateOrUpdateStream(ctx, n.js, cfg); err != nil {
		return err
	}

	return nil
}
