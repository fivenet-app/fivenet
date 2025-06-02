package centrum

import (
	"context"
	"fmt"
	"time"

	centrum "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	pbcentrum "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/centrum"
	eventscentrum "github.com/fivenet-app/fivenet/v2025/services/centrum/events"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func (s *Server) registerSubscriptions(ctxStartup context.Context, ctxCancel context.Context) error {
	streamCfg, err := eventscentrum.RegisterStream(ctxStartup, s.js)
	if err != nil {
		return fmt.Errorf("failed to register events: %w", err)
	}

	consumer, err := s.js.CreateConsumer(ctxStartup, streamCfg.Name, jetstream.ConsumerConfig{
		DeliverPolicy: jetstream.DeliverLastPerSubjectPolicy,
		FilterSubject: fmt.Sprintf("%s.>", eventscentrum.BaseSubject),
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

func (s *Server) watchForChanges(msg jetstream.Msg) {
	startTime := time.Now()

	if err := msg.Ack(); err != nil {
		s.logger.Error("failed to ack message", zap.Error(err))
	}

	job, topic, tType := eventscentrum.SplitSubject(msg.Subject())
	if job == "" || topic == "" || tType == "" {
		if err := msg.TermWithReason("invalid centrum subject"); err != nil {
			s.logger.Error("invalid centrum subject", zap.String("subject", msg.Subject()), zap.Error(err))
		}
		return
	}

	broker, ok := s.brokers.GetJobBroker(job)
	if !ok {
		s.logger.Debug("no broker found for job", zap.String("job", job))
		return
	}

	resp := &pbcentrum.StreamResponse{}
	switch topic {
	case eventscentrum.TopicGeneral:
		switch tType {
		case eventscentrum.TypeGeneralDisponents:
			dest := &centrum.Disponents{}
			if err := proto.Unmarshal(msg.Data(), dest); err != nil {
				s.logger.Error("failed to unmarshal nats change response", zap.Error(err))
				return
			}

			resp.Change = &pbcentrum.StreamResponse_Disponents{
				Disponents: dest,
			}

		case eventscentrum.TypeGeneralSettings:
			dest := &centrum.Settings{}
			if err := proto.Unmarshal(msg.Data(), dest); err != nil {
				s.logger.Error("failed to unmarshal nats change response", zap.Error(err))
				return
			}

			resp.Change = &pbcentrum.StreamResponse_Settings{
				Settings: dest,
			}

			if dest.Enabled {
				s.brokers.GetOrCreateJobBroker(dest.Job)
			} else {
				s.brokers.RemoveJobBroker(dest.Job)
			}
		}

	case eventscentrum.TopicUnit:
		switch tType {
		case eventscentrum.TypeUnitCreated:
			dest := &centrum.Unit{}
			if err := proto.Unmarshal(msg.Data(), dest); err != nil {
				s.logger.Error("failed to unmarshal nats change response", zap.Error(err))
				return
			}

			resp.Change = &pbcentrum.StreamResponse_UnitCreated{
				UnitCreated: dest,
			}

		case eventscentrum.TypeUnitDeleted:
			dest := &centrum.Unit{}
			if err := proto.Unmarshal(msg.Data(), dest); err != nil {
				s.logger.Error("failed to unmarshal nats change response", zap.Error(err))
				return
			}

			resp.Change = &pbcentrum.StreamResponse_UnitDeleted{
				UnitDeleted: dest,
			}

		case eventscentrum.TypeUnitUpdated:
			dest := &centrum.Unit{}
			if err := proto.Unmarshal(msg.Data(), dest); err != nil {
				s.logger.Error("failed to unmarshal nats change response", zap.Error(err))
				return
			}

			resp.Change = &pbcentrum.StreamResponse_UnitUpdated{
				UnitUpdated: dest,
			}

		case eventscentrum.TypeUnitStatus:
			dest := &centrum.UnitStatus{}
			if err := proto.Unmarshal(msg.Data(), dest); err != nil {
				s.logger.Error("failed to unmarshal nats change response", zap.Error(err))
				return
			}

			resp.Change = &pbcentrum.StreamResponse_UnitStatus{
				UnitStatus: dest,
			}
		}

	case eventscentrum.TopicDispatch:
		switch tType {
		case eventscentrum.TypeDispatchCreated:
			dest := &centrum.Dispatch{}
			if err := proto.Unmarshal(msg.Data(), dest); err != nil {
				s.logger.Error("failed to unmarshal nats change response", zap.Error(err))
				return
			}

			resp.Change = &pbcentrum.StreamResponse_DispatchCreated{
				DispatchCreated: dest,
			}

		case eventscentrum.TypeDispatchDeleted:
			dest := &centrum.Dispatch{}
			if err := proto.Unmarshal(msg.Data(), dest); err != nil {
				s.logger.Error("failed to unmarshal nats change response", zap.Error(err))
				return
			}

			resp.Change = &pbcentrum.StreamResponse_DispatchDeleted{
				DispatchDeleted: dest,
			}

		case eventscentrum.TypeDispatchUpdated:
			dest := &centrum.Dispatch{}
			if err := proto.Unmarshal(msg.Data(), dest); err != nil {
				s.logger.Error("failed to unmarshal nats change response", zap.Error(err))
				return
			}

			resp.Change = &pbcentrum.StreamResponse_DispatchUpdated{
				DispatchUpdated: dest,
			}

		case eventscentrum.TypeDispatchStatus:
			dest := &centrum.DispatchStatus{}
			if err := proto.Unmarshal(msg.Data(), dest); err != nil {
				s.logger.Error("failed to unmarshal nats change response", zap.Error(err))
				return
			}

			resp.Change = &pbcentrum.StreamResponse_DispatchStatus{
				DispatchStatus: dest,
			}
		}

	default:
		s.logger.Error("unknown centrum nats change response", zap.String("subject", msg.Subject()))
	}

	broker.Publish(resp)

	meta, err := msg.Metadata()
	if err != nil {
		s.logger.Error("sent centrum message broker, but failed to get msg metadata ", zap.Uint64("stream_sequence_id", meta.Sequence.Stream),
			zap.String("job", job), zap.String("topic", string(topic)), zap.String("type", string(tType)), zap.Error(err))
		return
	}
	s.logger.Debug("sent centrum message broker", zap.Uint64("stream_sequence_id", meta.Sequence.Stream),
		zap.String("job", job), zap.String("topic", string(topic)), zap.String("type", string(tType)),
		zap.Duration("duration", time.Since(startTime)))
}
