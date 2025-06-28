package livemap

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/livemap"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/instance"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/protoutils"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

const (
	BaseSubject events.Subject = "livemap"

	MarkerTopic events.Topic = "marker"

	MarkerUpdate events.Type = "update"
	MarkerDelete events.Type = "delete"
)

func (s *Server) registerStreamAndConsumer(ctxStartup context.Context, ctxCancel context.Context) error {
	cfg := jetstream.StreamConfig{
		Name:        "LIVEMAP",
		Description: "Livemapper Service events",
		Retention:   jetstream.InterestPolicy,
		Subjects:    []string{fmt.Sprintf("%s.>", BaseSubject)},
		Discard:     jetstream.DiscardOld,
		MaxAge:      2 * time.Minute,
		Storage:     jetstream.MemoryStorage,
	}
	if _, err := s.js.CreateOrUpdateStream(ctxStartup, cfg); err != nil {
		return err
	}

	consumer, err := s.js.CreateOrUpdateConsumer(ctxStartup, cfg.Name, jetstream.ConsumerConfig{
		Durable:           instance.ID() + "_livemap",
		FilterSubject:     fmt.Sprintf("%s.>", BaseSubject),
		DeliverPolicy:     jetstream.DeliverNewPolicy,
		InactiveThreshold: 1 * time.Minute, // Close consumer if inactive for 1 minute
	})
	if err != nil {
		return err
	}

	if s.jsCons != nil {
		s.jsCons.Stop()
	}

	s.jsCons, err = consumer.Consume(s.watchForEventsFunc,
		s.js.ConsumeErrHandlerWithRestart(ctxCancel, s.logger, s.registerStreamAndConsumer))
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) sendUpdateEvent(ctx context.Context, topic events.Topic, tType events.Type, job string, msg proto.Message) error {
	if _, err := s.js.PublishProto(ctx, fmt.Sprintf("%s.%s.%s.%s", BaseSubject, topic, tType, job), msg); err != nil {
		return err
	}

	return nil
}

func (s *Server) watchForEventsFunc(msg jetstream.Msg) {
	if err := msg.Ack(); err != nil {
		s.logger.Error("failed to ack message", zap.Error(err))
	}

	split := strings.Split(msg.Subject(), ".")
	if len(split) < 4 {
		return
	}

	topic := events.Topic(split[1])
	tType := events.Type(split[2])
	job := events.Type(split[3])
	_ = job
	switch topic {
	case MarkerTopic:
		switch tType {
		case MarkerUpdate:
			// Send marker update when there is at least one subscriber
			if s.broker.SubCount() <= 0 {
				return
			}

			marker := &livemap.MarkerMarker{}
			if err := protoutils.UnmarshalPartialPJSON(msg.Data(), marker); err != nil {
				s.logger.Error("failed to unmarshal livemap marker update data", zap.Error(err))
				return
			}

			s.broker.Publish(&brokerEvent{
				MarkerUpdate: marker,
			})

		case MarkerDelete:
			// Send marker deletion when there is at least one subscriber
			if s.broker.SubCount() <= 0 {
				return
			}

			marker := &livemap.MarkerMarker{}
			if err := protoutils.UnmarshalPartialPJSON(msg.Data(), marker); err != nil {
				s.logger.Error("failed to unmarshal livemap marker update data", zap.Error(err))
				return
			}

			s.broker.Publish(&brokerEvent{
				MarkerDelete: &marker.Id,
			})
		}
	}
}
