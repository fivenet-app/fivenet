package livemapper

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/pkg/events"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

const (
	BaseSubject events.Subject = "livemap"

	MarkerUpdate events.Type = "marker_update"
)

func (s *Server) registerSubscriptions(ctxStartup context.Context, ctxCancel context.Context) error {
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

	consumer, err := s.js.CreateConsumer(ctxStartup, cfg.Name, jetstream.ConsumerConfig{
		DeliverPolicy: jetstream.DeliverNewPolicy,
		FilterSubject: fmt.Sprintf("%s.>", BaseSubject),
	})
	if err != nil {
		return err
	}

	if s.jsCons != nil {
		s.jsCons.Stop()
	}

	s.jsCons, err = consumer.Consume(s.watchForEventsFunc(ctxCancel),
		s.js.ConsumeErrHandlerWithRestart(ctxCancel, s.logger, s.registerSubscriptions))
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) sendUpdateEvent(ctx context.Context, tType events.Type, event proto.Message) error {
	data, err := proto.Marshal(event)
	if err != nil {
		return err
	}

	if _, err := s.js.Publish(ctx, fmt.Sprintf("%s.%s", BaseSubject, tType), data); err != nil {
		return err
	}

	return nil
}

func (s *Server) watchForEventsFunc(ctx context.Context) jetstream.MessageHandler {
	return func(msg jetstream.Msg) {
		remoteCtx, err := events.GetJetstreamMsgContext(msg)
		if err != nil {
			s.logger.Error("failed to get js msg context", zap.Error(err))
		}
		_, span := s.tracer.Start(trace.ContextWithRemoteSpanContext(ctx, remoteCtx), msg.Subject())
		defer span.End()

		if err := msg.Ack(); err != nil {
			s.logger.Error("failed to ack message", zap.Error(err))
		}

		split := strings.Split(msg.Subject(), ".")
		if len(split) < 2 {
			return
		}

		tType := events.Type(split[1])
		if tType == MarkerUpdate {
			if err := s.refreshData(ctx); err != nil {
				s.logger.Error("failed to refresh livemap markers cache", zap.Error(err))
				return
			}

			// Send marker update when data has been refreshed and we have at least one subscriber
			if s.broker.SubCount() <= 0 {
				return
			}
			s.broker.Publish(&brokerEvent{Send: MarkerUpdate})
		}
	}
}
