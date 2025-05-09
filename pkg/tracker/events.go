package tracker

import (
	"context"
	"fmt"
	"time"

	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/nats-io/nats.go/jetstream"
	"google.golang.org/protobuf/proto"
)

const (
	StreamName = "TRACKER"

	BaseSubject events.Subject = "tracker"

	UsersUpdate events.Type = "users_update"
)

func registerStreams(ctx context.Context, js *events.JSWrapper) error {
	cfg := jetstream.StreamConfig{
		Name:        StreamName,
		Description: "User Tracker Events",
		Retention:   jetstream.InterestPolicy,
		Subjects:    []string{fmt.Sprintf("%s.>", BaseSubject)},
		Discard:     jetstream.DiscardOld,
		MaxAge:      2 * time.Minute,
		Storage:     jetstream.MemoryStorage,
	}
	if _, err := js.CreateOrUpdateStream(ctx, cfg); err != nil {
		return err
	}

	return nil
}

func (m *Manager) sendUpdateEvent(ctx context.Context, tType events.Type, event proto.Message) error {
	data, err := proto.Marshal(event)
	if err != nil {
		return err
	}

	if _, err := m.js.Publish(ctx, fmt.Sprintf("%s.%s", BaseSubject, tType), data); err != nil {
		return err
	}

	return nil
}

func (s *Tracker) registerSubscriptions(ctxStartup context.Context, ctxCancel context.Context) error {
	consumer, err := s.js.CreateConsumer(ctxStartup, StreamName, jetstream.ConsumerConfig{
		DeliverPolicy: jetstream.DeliverNewPolicy,
		FilterSubject: fmt.Sprintf("%s.>", BaseSubject),
	})
	if err != nil {
		return err
	}

	if s.jsCons != nil {
		s.jsCons.Stop()
		s.jsCons = nil
	}

	s.jsCons, err = consumer.Consume(s.watchForChanges,
		s.js.ConsumeErrHandlerWithRestart(ctxCancel, s.logger,
			s.registerSubscriptions,
		))
	if err != nil {
		return err
	}

	return nil
}
