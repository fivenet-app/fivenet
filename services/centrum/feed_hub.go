package centrum

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	centrumdispatches "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatches"
	centrumunits "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/units"
	pbcentrum "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/centrum"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/protoutils"
	eventscentrum "github.com/fivenet-app/fivenet/v2026/services/centrum/events"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type feedEvent struct {
	Sequence uint64
	Job      string
	Response *pbcentrum.StreamResponse
	// Resync tells connected clients that a source-feed interruption may have
	// skipped updates. They must reload a snapshot before accepting more events.
	Resync bool
}

func (s *Server) startFeedHub(ctx context.Context) {
	s.wg.Go(func() {
		s.runCentrumEventFeed(ctx)
	})

	for _, feed := range feeds {
		s.wg.Go(func() {
			s.runKVFeed(ctx, feed)
		})
	}
}

func (s *Server) publishFeedEvent(
	feed string,
	job string,
	response *pbcentrum.StreamResponse,
	kvRevision uint64,
) {
	if response == nil {
		return
	}
	response.KvRevision = kvRevision

	// Feed workers run concurrently. Keep sequence allocation and broker enqueue
	// ordered so clients never mistake concurrent publications for a gap.
	s.feedMu.Lock()
	defer s.feedMu.Unlock()
	if s.feedBroker.SubCount() == 0 {
		s.metrics.IncFeedMessage(feed, "skipped_no_subscribers")
		return
	}

	s.feedBroker.Publish(&feedEvent{
		Sequence: s.feedSequence.Add(1),
		Job:      job,
		Response: response,
	})
	s.metrics.IncFeedMessage(feed, "published")
}

func (s *Server) publishFeedResync(feed string) {
	s.feedMu.Lock()
	defer s.feedMu.Unlock()
	if s.feedBroker.SubCount() == 0 {
		return
	}

	s.feedBroker.Publish(&feedEvent{
		Sequence: s.feedSequence.Add(1),
		Resync:   true,
	})
	s.metrics.IncFeedResync(feed + "_worker_restart")
}

func (s *Server) runCentrumEventFeed(ctx context.Context) {
	for {
		err := s.consumeCentrumEvents(ctx)
		if ctx.Err() == nil && !protoutils.IsContextCanceled(err) {
			if err != nil {
				s.logger.Error("centrum event feed stopped", zap.Error(err))
			}
			s.publishFeedResync("events")
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(2 * time.Second):
		}
	}
}

func (s *Server) consumeCentrumEvents(ctx context.Context) error {
	consumer, err := s.js.CreateConsumer(ctx, "CENTRUM", jetstream.ConsumerConfig{
		FilterSubject: string(eventscentrum.BaseSubject) + ".>",
		DeliverPolicy: jetstream.DeliverNewPolicy,
		AckPolicy:     jetstream.AckNonePolicy,
	})
	if err != nil {
		return fmt.Errorf("failed to create centrum event consumer. %w", err)
	}

	msgs, err := consumer.Messages(
		jetstream.PullMaxMessages(feedFetch),
		jetstream.WithMessagesErrOnMissingHeartbeat(false),
	)
	if err != nil {
		return err
	}
	defer msgs.Stop()

	for {
		msg, err := msgs.Next(jetstream.NextContext(ctx))
		if err != nil {
			if errors.Is(err, jetstream.ErrMsgIteratorClosed) {
				return nil
			}
			return err
		}

		job, topic, eventType := eventscentrum.SplitSubject(msg.Subject())
		if s.feedBroker.SubCount() == 0 {
			s.metrics.IncFeedMessage("events", "skipped_no_subscribers")
			continue
		}
		var response *pbcentrum.StreamResponse
		switch topic {
		case eventscentrum.TopicDispatch:
			if eventType != eventscentrum.TypeDispatchStatus {
				continue
			}
			status := &centrumdispatches.DispatchStatus{}
			if err := proto.Unmarshal(msg.Data(), status); err != nil {
				s.metrics.IncFeedDecodeFailure("dispatch_status")
				s.logger.Warn(
					"failed to unmarshal dispatch status feed event",
					zap.Error(err),
					zap.String("subject", msg.Subject()),
				)
				continue
			}
			response = &pbcentrum.StreamResponse{
				Change: &pbcentrum.StreamResponse_DispatchStatus{DispatchStatus: status},
			}

		case eventscentrum.TopicUnit:
			if eventType != eventscentrum.TypeUnitStatus {
				continue
			}
			status := &centrumunits.UnitStatus{}
			if err := proto.Unmarshal(msg.Data(), status); err != nil {
				s.metrics.IncFeedDecodeFailure("unit_status")
				s.logger.Warn(
					"failed to unmarshal unit status feed event",
					zap.Error(err),
					zap.String("subject", msg.Subject()),
				)
				continue
			}
			response = &pbcentrum.StreamResponse{
				Change: &pbcentrum.StreamResponse_UnitStatus{UnitStatus: status},
			}
		}

		s.publishFeedEvent("events", job, response, 0)
	}
}

func (s *Server) runKVFeed(ctx context.Context, feed feedCfg) {
	for {
		err := s.consumeKVFeed(ctx, feed)
		if ctx.Err() == nil && !protoutils.IsContextCanceled(err) {
			if err != nil {
				s.logger.Error(
					"centrum KV feed stopped",
					zap.Error(err),
					zap.String("bucket", feed.Bucket),
				)
			}
			s.publishFeedResync(feed.StreamName)
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(2 * time.Second):
		}
	}
}

func (s *Server) consumeKVFeed(ctx context.Context, feed feedCfg) error {
	consumer, err := s.js.CreateConsumer(ctx, "KV_"+feed.StreamName, jetstream.ConsumerConfig{
		FilterSubject: "$KV." + feed.Bucket + ".>",
		DeliverPolicy: jetstream.DeliverNewPolicy,
		AckPolicy:     jetstream.AckNonePolicy,
		MaxWaiting:    8,
	})
	if err != nil {
		return fmt.Errorf("failed to create KV feed consumer. %w", err)
	}

	msgs, err := consumer.Messages(
		jetstream.PullMaxMessages(feedFetch),
		jetstream.WithMessagesErrOnMissingHeartbeat(false),
	)
	if err != nil {
		return err
	}
	defer msgs.Stop()

	for {
		msg, err := msgs.Next(jetstream.NextContext(ctx))
		if err != nil {
			if errors.Is(err, jetstream.ErrMsgIteratorClosed) {
				return nil
			}
			return err
		}

		key := strings.TrimPrefix(msg.Subject(), "$KV."+feed.Bucket+".")
		metadata, err := msg.Metadata()
		if err != nil {
			return fmt.Errorf("failed to read KV feed message metadata. %w", err)
		}
		kvRevision := metadata.Sequence.Stream
		job := key
		if !feed.NoWildcard {
			job, _, _ = strings.Cut(key, ".")
		}

		if op := msg.Headers().Get("KV-Operation"); op == "DEL" || op == "PURGE" {
			s.publishFeedEvent(feed.StreamName, job, feed.WrapDelete(key), kvRevision)
			continue
		}
		if s.feedBroker.SubCount() == 0 {
			s.metrics.IncFeedMessage(feed.StreamName, "skipped_no_subscribers")
			continue
		}

		obj, err := feed.Unmarshal(ctx, s, msg.Data())
		if err != nil {
			s.metrics.IncFeedDecodeFailure(feed.StreamName)
			s.logger.Warn(
				"failed to unmarshal KV feed message",
				zap.Error(err),
				zap.String("bucket", feed.Bucket),
				zap.String("subject", msg.Subject()),
			)
			continue
		}

		s.publishFeedEvent(feed.StreamName, job, feed.WrapPut(obj), kvRevision)
	}
}
