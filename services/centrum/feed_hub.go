package centrum

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	centrumdispatchers "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatchers"
	centrumdispatches "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatches"
	centrumsettings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/settings"
	centrumunits "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/units"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	pbcentrum "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/centrum"
	eventscentrum "github.com/fivenet-app/fivenet/v2026/services/centrum/events"
	centrumutils "github.com/fivenet-app/fivenet/v2026/services/centrum/utils"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// feedCfg describes one source projection and its wire representation.
type feedCfg struct {
	StreamName string
	Bucket     string
	NoWildcard bool
	Unmarshal  func(ctx context.Context, s *Server, b []byte) (proto.Message, error)
	WrapPut    func(proto.Message) *pbcentrum.StreamResponse
	WrapDelete func(key string) *pbcentrum.StreamResponse
}

var feeds = []feedCfg{
	{
		StreamName: "centrum_settings", Bucket: "centrum_settings", NoWildcard: true,
		Unmarshal: func(ctx context.Context, _ *Server, b []byte) (proto.Message, error) {
			var value centrumsettings.Settings
			return &value, proto.Unmarshal(b, &value)
		},
		WrapPut: func(value proto.Message) *pbcentrum.StreamResponse {
			return &pbcentrum.StreamResponse{
				Change: &pbcentrum.StreamResponse_Settings{
					Settings: value.(*centrumsettings.Settings),
				},
			}
		},
		WrapDelete: func(key string) *pbcentrum.StreamResponse {
			return &pbcentrum.StreamResponse{
				Change: &pbcentrum.StreamResponse_SettingsDeleted{SettingsDeleted: key},
			}
		},
	},
	{
		StreamName: "centrum_dispatchers", Bucket: "centrum_dispatchers", NoWildcard: true,
		Unmarshal: func(ctx context.Context, _ *Server, b []byte) (proto.Message, error) {
			var value centrumdispatchers.Dispatchers
			return &value, proto.Unmarshal(b, &value)
		},
		WrapPut: func(value proto.Message) *pbcentrum.StreamResponse {
			return &pbcentrum.StreamResponse{
				Change: &pbcentrum.StreamResponse_Dispatchers{
					Dispatchers: value.(*centrumdispatchers.Dispatchers),
				},
			}
		},
		WrapDelete: func(key string) *pbcentrum.StreamResponse {
			return &pbcentrum.StreamResponse{
				Change: &pbcentrum.StreamResponse_Dispatchers{
					Dispatchers: &centrumdispatchers.Dispatchers{Job: key},
				},
			}
		},
	},
	{
		StreamName: "centrum_units", Bucket: "centrum_units.job",
		Unmarshal: func(ctx context.Context, s *Server, b []byte) (proto.Message, error) {
			var mapping common.IDMapping
			if err := proto.Unmarshal(b, &mapping); err != nil {
				return nil, fmt.Errorf("failed to unmarshal unit id mapping. %w", err)
			}
			return s.units.Get(ctx, mapping.GetId())
		},
		WrapPut: func(value proto.Message) *pbcentrum.StreamResponse {
			return &pbcentrum.StreamResponse{
				Change: &pbcentrum.StreamResponse_UnitUpdated{
					UnitUpdated: value.(*centrumunits.Unit),
				},
			}
		},
		WrapDelete: func(key string) *pbcentrum.StreamResponse {
			id, err := centrumutils.ExtractID(key)
			if err != nil {
				return nil
			}
			return &pbcentrum.StreamResponse{
				Change: &pbcentrum.StreamResponse_UnitDeleted{UnitDeleted: id},
			}
		},
	},
	{
		StreamName: "centrum_dispatches", Bucket: "centrum_dispatches.job",
		Unmarshal: func(ctx context.Context, s *Server, b []byte) (proto.Message, error) {
			var mapping common.IDMapping
			if err := proto.Unmarshal(b, &mapping); err != nil {
				return nil, fmt.Errorf("failed to unmarshal dispatch id mapping. %w", err)
			}
			return s.dispatches.Get(ctx, mapping.GetId())
		},
		WrapPut: func(value proto.Message) *pbcentrum.StreamResponse {
			return &pbcentrum.StreamResponse{
				Change: &pbcentrum.StreamResponse_DispatchUpdated{
					DispatchUpdated: value.(*centrumdispatches.Dispatch),
				},
			}
		},
		WrapDelete: func(key string) *pbcentrum.StreamResponse {
			id, err := centrumutils.ExtractID(key)
			if err != nil {
				return nil
			}
			return &pbcentrum.StreamResponse{
				Change: &pbcentrum.StreamResponse_DispatchDeleted{DispatchDeleted: id},
			}
		},
	},
}

type projectionFeedEventKey struct {
	feed    string
	id      int64
	deleted bool
}

type pendingProjectionFeedEvent struct {
	jobs       map[string]struct{}
	response   *pbcentrum.StreamResponse
	kvRevision uint64
}

type feedEvent struct {
	Sequence uint64
	// Jobs is a de-duplicated list of jobs authorized to receive Response.
	Jobs     []string
	Response *pbcentrum.StreamResponse
	// Resync tells connected clients that a source-feed interruption may have
	// skipped updates. They must reload a snapshot before accepting more events.
	Resync bool
}

const feedCoalesceWindow = 5 * time.Millisecond

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

func dispatchFeedKey(response *pbcentrum.StreamResponse) (projectionFeedEventKey, bool) {
	if dispatch := response.GetDispatchUpdated(); dispatch != nil && dispatch.GetId() > 0 {
		return projectionFeedEventKey{feed: "centrum_dispatches", id: dispatch.GetId()}, true
	}
	if id := response.GetDispatchDeleted(); id > 0 {
		return projectionFeedEventKey{feed: "centrum_dispatches", id: id, deleted: true}, true
	}

	return projectionFeedEventKey{}, false
}

func (s *Server) queueDispatchFeedEvent(
	job string,
	response *pbcentrum.StreamResponse,
	kvRevision uint64,
) {
	key, ok := dispatchFeedKey(response)
	if !ok {
		s.publishFeedEvent("centrum_dispatches", []string{job}, response, kvRevision)
		return
	}
	s.queueProjectionFeedEvent(job, response, kvRevision, key)
}

func unitFeedKey(response *pbcentrum.StreamResponse) (projectionFeedEventKey, bool) {
	if unit := response.GetUnitUpdated(); unit != nil && unit.GetId() > 0 {
		return projectionFeedEventKey{feed: "centrum_units", id: unit.GetId()}, true
	}
	if id := response.GetUnitDeleted(); id > 0 {
		return projectionFeedEventKey{feed: "centrum_units", id: id, deleted: true}, true
	}

	return projectionFeedEventKey{}, false
}

func (s *Server) queueUnitFeedEvent(
	job string,
	response *pbcentrum.StreamResponse,
	kvRevision uint64,
) {
	key, ok := unitFeedKey(response)
	if !ok {
		s.publishFeedEvent("centrum_units", []string{job}, response, kvRevision)
		return
	}
	s.queueProjectionFeedEvent(job, response, kvRevision, key)
}

// queueProjectionFeedEvent coalesces job-mapping updates while preserving an
// update-before-delete order for the same live projection.
func (s *Server) queueProjectionFeedEvent(
	job string,
	response *pbcentrum.StreamResponse,
	kvRevision uint64,
	key projectionFeedEventKey,
) {
	if response == nil {
		return
	}
	s.feedMu.Lock()
	defer s.feedMu.Unlock()
	if s.feedBroker.SubCount() == 0 {
		s.metrics.IncFeedMessage(key.feed, "skipped_no_subscribers")
		return
	}

	if key.deleted {
		s.flushProjectionFeedEventLocked(projectionFeedEventKey{feed: key.feed, id: key.id})
	}

	if s.pendingProjectionFeedEvents == nil {
		s.pendingProjectionFeedEvents = make(map[projectionFeedEventKey]*pendingProjectionFeedEvent)
	}
	pending := s.pendingProjectionFeedEvents[key]
	if pending == nil {
		pending = &pendingProjectionFeedEvent{jobs: make(map[string]struct{})}
		s.pendingProjectionFeedEvents[key] = pending
		time.AfterFunc(feedCoalesceWindow, func() {
			s.flushProjectionFeedEvent(key)
		})
	}

	pending.jobs[job] = struct{}{}
	pending.response = response
	pending.kvRevision = max(pending.kvRevision, kvRevision)
}

func (s *Server) flushProjectionFeedEvent(key projectionFeedEventKey) {
	s.feedMu.Lock()
	defer s.feedMu.Unlock()
	s.flushProjectionFeedEventLocked(key)
}

func (s *Server) flushProjectionFeedEventLocked(key projectionFeedEventKey) {
	pending := s.pendingProjectionFeedEvents[key]
	if pending == nil {
		return
	}
	delete(s.pendingProjectionFeedEvents, key)

	jobs := make([]string, 0, len(pending.jobs))
	for job := range pending.jobs {
		jobs = append(jobs, job)
	}
	slices.Sort(jobs)
	s.publishFeedEventLocked(key.feed, jobs, pending.response, pending.kvRevision)
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
			switch feed.StreamName {
			case "centrum_dispatches":
				s.queueDispatchFeedEvent(job, response, kvRevision)
			case "centrum_units":
				s.queueUnitFeedEvent(job, response, kvRevision)
			default:
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
		switch feed.StreamName {
		case "centrum_dispatches":
			s.queueDispatchFeedEvent(job, response, kvRevision)
		case "centrum_units":
			s.queueUnitFeedEvent(job, response, kvRevision)
		default:
			s.publishFeedEvent(feed.StreamName, []string{job}, response, kvRevision)
		}
	}
}
