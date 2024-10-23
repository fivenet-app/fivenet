package centrummanager

import (
	"context"
	"fmt"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/centrum"
	eventscentrum "github.com/fivenet-app/fivenet/gen/go/proto/services/centrum/events"
	"github.com/fivenet-app/fivenet/pkg/events"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (s *Manager) registerStream(ctx context.Context) (jetstream.StreamConfig, error) {
	streamCfg, err := eventscentrum.RegisterStream(ctx, s.js)
	if err != nil {
		return streamCfg, fmt.Errorf("failed to register stream: %w", err)
	}

	return streamCfg, nil
}

func (s *Manager) registerSubscriptions(ctxStartup context.Context, ctxCancel context.Context) error {
	streamCfg, err := s.registerStream(ctxStartup)
	if err != nil {
		return err
	}

	consumer, err := s.js.CreateConsumer(ctxStartup, streamCfg.Name, jetstream.ConsumerConfig{
		DeliverPolicy: jetstream.DeliverLastPerSubjectPolicy,
		FilterSubject: fmt.Sprintf("%s.*.%s.>", eventscentrum.BaseSubject, eventscentrum.TopicGeneral),
	})
	if err != nil {
		return err
	}

	if s.jsCons != nil {
		s.jsCons.Stop()
		s.jsCons = nil
	}

	s.jsCons, err = consumer.Consume(s.watchTopicGeneralFunc(ctxCancel),
		s.js.ConsumeErrHandlerWithRestart(ctxCancel, s.logger, s.registerSubscriptions))
	if err != nil {
		s.logger.Error("failed to subscribe to centrum general topic", zap.Error(err))
		return err
	}

	return nil
}

func (s *Manager) watchTopicGeneralFunc(ctx context.Context) jetstream.MessageHandler {
	return func(msg jetstream.Msg) {
		remoteCtx, _ := events.GetJetstreamMsgContext(msg)
		_, span := s.tracer.Start(trace.ContextWithRemoteSpanContext(ctx, remoteCtx), msg.Subject())
		defer span.End()

		if err := msg.Ack(); err != nil {
			s.logger.Error("failed to ack message", zap.Error(err))
		}

		job, _, tType := eventscentrum.SplitSubject(msg.Subject())
		if job == "" || tType == "" {
			if err := msg.TermWithReason("invalid centrum general subject"); err != nil {
				s.logger.Error("invalid centrum general subject", zap.String("subject", msg.Subject()), zap.Error(err))
			}
			return
		}

		meta, err := msg.Metadata()
		if err != nil {
			s.logger.Error("failed to read message metadata in centrum general topic subscription", zap.Error(err))
			return
		}
		s.logger.Debug("received general message", zap.Uint64("stream_sequence_id", meta.Sequence.Stream),
			zap.String("job", job), zap.String("type", string(tType)))

		switch tType {
		case eventscentrum.TypeGeneralSettings:
			var dest centrum.Settings
			if err := proto.Unmarshal(msg.Data(), &dest); err != nil {
				s.logger.Error("failed to unmarshal settings message", zap.Error(err))
				return
			}

			if err := s.State.UpdateSettings(ctx, job, &dest); err != nil {
				s.logger.Error("failed to update settings", zap.Error(err))
				return
			}
		}

		s.logger.Debug("handled general message", zap.String("job", job), zap.String("type", string(tType)))
	}
}
