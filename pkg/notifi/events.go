package notifi

import (
	"context"
	"fmt"
	"time"

	"github.com/galexrt/fivenet/pkg/events"
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
		Description: "User Notifications events",
		Retention:   jetstream.InterestPolicy,
		Subjects:    []string{fmt.Sprintf("%s.>", BaseSubject)},
		Discard:     jetstream.DiscardOld,
		MaxAge:      30 * time.Minute,
		Replicas:    2,
	}
	if _, err := n.js.CreateOrUpdateStream(ctx, cfg); err != nil {
		return err
	}

	return nil
}
