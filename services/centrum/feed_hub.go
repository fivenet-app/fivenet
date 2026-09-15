package centrum

import (
	"context"
	"errors"
	"fmt"
	"slices"
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
	Jobs     []string
	Response *pbcentrum.StreamResponse
	// Resync tells connected clients that a source-feed interruption may have
	// skipped updates. They must reload a snapshot before accepting more events.
	Resync bool
}

const feedCoalesceWindow = 5 * time.Millisecond

type dispatchFeedEventKey struct {
	id      int64
	deleted bool
}

type pendingDispatchFeedEvent struct {
	jobs       map[string]struct{}
	response   *pbcentrum.StreamResponse
	kvRevision uint64
}

type unitFeedEventKey struct {
	id      int64
	deleted bool
}

type pendingUnitFeedEvent struct {
	jobs       map[string]struct{}
	response   *pbcentrum.StreamResponse
	kvRevision uint64
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
	jobs []string,
	response *pbcentrum.StreamResponse,
	kvRevision uint64,
) {
	if response == nil {
		return
	}

	// Feed workers run concurrently. Keep sequence allocation and broker enqueue
	// ordered so clients never mistake concurrent publications for a gap.
	s.feedMu.Lock()
	defer s.feedMu.Unlock()
	s.publishFeedEventLocked(feed, jobs, response, kvRevision)
}

func (s *Server) publishFeedEventLocked(
	feed string,
	jobs []string,
	response *pbcentrum.StreamResponse,
	kvRevision uint64,
) {
	response.KvRevision = kvRevision
	if s.feedBroker.SubCount() == 0 {
		s.metrics.IncFeedMessage(feed, "skipped_no_subscribers")
		return
	}

	s.feedBroker.Publish(&feedEvent{
		Sequence: s.feedSequence.Add(1),
		Jobs:     jobs,
		Response: response,
	})
	s.metrics.IncFeedMessage(feed, "published")
}

func dispatchFeedKey(response *pbcentrum.StreamResponse) (dispatchFeedEventKey, bool) {
	if dispatch := response.GetDispatchUpdated(); dispatch != nil && dispatch.GetId() > 0 {
		return dispatchFeedEventKey{id: dispatch.GetId()}, true
	}
	if id := response.GetDispatchDeleted(); id > 0 {
		return dispatchFeedEventKey{id: id, deleted: true}, true
	}

	return dispatchFeedEventKey{}, false
}

func (s *Server) queueDispatchFeedEvent(job string, response *pbcentrum.StreamResponse, kvRevision uint64) {
	if response == nil {
		return
	}

	key, ok := dispatchFeedKey(response)
	if !ok {
		s.publishFeedEvent("centrum_dispatches", []string{job}, response, kvRevision)
		return
	}

	s.feedMu.Lock()
	defer s.feedMu.Unlock()
	if s.feedBroker.SubCount() == 0 {
		s.metrics.IncFeedMessage("centrum_dispatches", "skipped_no_subscribers")
		return
	}

	// Do not let an update and a deletion for the same dispatch overtake each other.
	if key.deleted {
		s.flushDispatchFeedEventLocked(dispatchFeedEventKey{id: key.id})
	}

	if s.pendingDispatchFeedEvents == nil {
		s.pendingDispatchFeedEvents = make(map[dispatchFeedEventKey]*pendingDispatchFeedEvent)
	}
	pending := s.pendingDispatchFeedEvents[key]
	if pending == nil {
		pending = &pendingDispatchFeedEvent{jobs: make(map[string]struct{})}
		s.pendingDispatchFeedEvents[key] = pending
		time.AfterFunc(feedCoalesceWindow, func() {
			s.flushDispatchFeedEvent(key)
		})
	}

	pending.jobs[job] = struct{}{}
	pending.response = response
	pending.kvRevision = max(pending.kvRevision, kvRevision)
}

func (s *Server) flushDispatchFeedEvent(key dispatchFeedEventKey) {
	s.feedMu.Lock()
	defer s.feedMu.Unlock()
	s.flushDispatchFeedEventLocked(key)
}

func (s *Server) flushDispatchFeedEventLocked(key dispatchFeedEventKey) {
	pending := s.pendingDispatchFeedEvents[key]
	if pending == nil {
		return
	}
	delete(s.pendingDispatchFeedEvents, key)

	jobs := make([]string, 0, len(pending.jobs))
	for job := range pending.jobs {
		jobs = append(jobs, job)
	}
	slices.Sort(jobs)
	s.publishFeedEventLocked("centrum_dispatches", jobs, pending.response, pending.kvRevision)
}

func unitFeedKey(response *pbcentrum.StreamResponse) (unitFeedEventKey, bool) {
	if unit := response.GetUnitUpdated(); unit != nil && unit.GetId() > 0 {
		return unitFeedEventKey{id: unit.GetId()}, true
	}
	if id := response.GetUnitDeleted(); id > 0 {
		return unitFeedEventKey{id: id, deleted: true}, true
	}

	return unitFeedEventKey{}, false
}

func (s *Server) queueUnitFeedEvent(job string, response *pbcentrum.StreamResponse, kvRevision uint64) {
	if response == nil {
		return
	}

	key, ok := unitFeedKey(response)
	if !ok {
		s.publishFeedEvent("centrum_units", []string{job}, response, kvRevision)
		return
	}

	s.feedMu.Lock()
	defer s.feedMu.Unlock()
	if s.feedBroker.SubCount() == 0 {
		s.metrics.IncFeedMessage("centrum_units", "skipped_no_subscribers")
		return
	}

	if key.deleted {
		s.flushUnitFeedEventLocked(unitFeedEventKey{id: key.id})
	}

	if s.pendingUnitFeedEvents == nil {
		s.pendingUnitFeedEvents = make(map[unitFeedEventKey]*pendingUnitFeedEvent)
	}
	pending := s.pendingUnitFeedEvents[key]
	if pending == nil {
		pending = &pendingUnitFeedEvent{jobs: make(map[string]struct{})}
		s.pendingUnitFeedEvents[key] = pending
		time.AfterFunc(feedCoalesceWindow, func() {
			s.flushUnitFeedEvent(key)
		})
	}

	pending.jobs[job] = struct{}{}
	pending.response = response
	pending.kvRevision = max(pending.kvRevision, kvRevision)
}

func (s *Server) flushUnitFeedEvent(key unitFeedEventKey) {
	s.feedMu.Lock()
	defer s.feedMu.Unlock()
	s.flushUnitFeedEventLocked(key)
}

func (s *Server) flushUnitFeedEventLocked(key unitFeedEventKey) {
	pending := s.pendingUnitFeedEvents[key]
	if pending == nil {
		return
	}
	delete(s.pendingUnitFeedEvents, key)

	jobs := make([]string, 0, len(pending.jobs))
	for job := range pending.jobs {
		jobs = append(jobs, job)
	}
	slices.Sort(jobs)
	s.publishFeedEventLocked("centrum_units", jobs, pending.response, pending.kvRevision)
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

		s.publishFeedEvent("events", []string{job}, response, 0)
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
			response := feed.WrapDelete(key)
			if feed.StreamName == "centrum_dispatches" {
				s.queueDispatchFeedEvent(job, response, kvRevision)
			} else if feed.StreamName == "centrum_units" {
				s.queueUnitFeedEvent(job, response, kvRevision)
			} else {
				s.publishFeedEvent(feed.StreamName, []string{job}, response, kvRevision)
			}
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

		response := feed.WrapPut(obj)
		if feed.StreamName == "centrum_dispatches" {
			s.queueDispatchFeedEvent(job, response, kvRevision)
		} else if feed.StreamName == "centrum_units" {
			s.queueUnitFeedEvent(job, response, kvRevision)
		} else {
			s.publishFeedEvent(feed.StreamName, []string{job}, response, kvRevision)
		}
	}
}
