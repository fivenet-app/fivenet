package tracker

import (
	"context"
	"fmt"
	"time"

	"github.com/galexrt/fivenet/pkg/events"
	natsutils "github.com/galexrt/fivenet/pkg/nats"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

const (
	BaseSubject events.Subject = "tracker"

	UsersUpdate events.Type = "users_update"
)

func (m *Manager) registerEvents(ctx context.Context) error {
	cfg := &nats.StreamConfig{
		Name:        "TRACKER",
		Description: natsutils.Description,
		Retention:   nats.InterestPolicy,
		Subjects:    []string{fmt.Sprintf("%s.>", BaseSubject)},
		Discard:     nats.DiscardOld,
		MaxAge:      2 * time.Minute,
	}
	if _, err := natsutils.CreateOrUpdateStream(ctx, m.js, cfg); err != nil {
		return err
	}

	return nil
}

func (m *Manager) sendUpdateEvent(tType events.Type, event proto.Message) error {
	data, err := proto.Marshal(event)
	if err != nil {
		return err
	}

	if _, err := m.js.Publish(fmt.Sprintf("%s.%s", BaseSubject, tType), data); err != nil {
		return err
	}

	return nil
}
