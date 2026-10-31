package centrum

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"time"

	centrumdispatchers "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatchers"
	centrumdispatches "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatches"
	centrumsettings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/settings"
	centrumunits "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/units"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	pbcentrum "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/centrum"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/instance"
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
	KeyPrefix  string
	NoWildcard bool
	Unmarshal  func(ctx context.Context, s *Server, b []byte) (proto.Message, error)
	Job        func(ctx context.Context, s *Server, key string, value proto.Message) (string, error)
	Jobs       func(ctx context.Context, s *Server, key string, value proto.Message) ([]string, error)
	WrapPut    func(proto.Message) *pbcentrum.StreamResponse
	WrapDelete func(key string) *pbcentrum.StreamResponse
}

func kvFeedConsumerName(feed feedCfg) string {
	name := "KV_" + feed.StreamName
	if feed.KeyPrefix != "" {
		name += "_" + feed.KeyPrefix
	}
	if instanceID := instance.ID(); instanceID != "" {
		name += "_" + instanceID
	}
	return name
}

func centrumEventConsumerName() string {
	name := "CENTRUM_events"
	if instanceID := instance.ID(); instanceID != "" {
		name += "_" + instanceID
	}
	return name
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
		StreamName: "centrum_units", Bucket: "centrum_units", KeyPrefix: "job",
		Unmarshal: func(ctx context.Context, s *Server, b []byte) (proto.Message, error) {
			var mapping common.IDMapping
			if err := proto.Unmarshal(b, &mapping); err != nil {
				return nil, fmt.Errorf("failed to unmarshal unit id mapping. %w", err)
			}
			return s.units.Get(ctx, mapping.GetId())
		},
		// The job index remains the authoritative delete feed. Component puts
		// are consumed from their own projections below.
		WrapPut: func(proto.Message) *pbcentrum.StreamResponse { return nil },
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
		StreamName: "centrum_units",
		Bucket:     "centrum_units",
		KeyPrefix:  "spec",
		NoWildcard: true,
		Unmarshal: func(ctx context.Context, _ *Server, b []byte) (proto.Message, error) {
			var value centrumunits.UnitSpec
			return &value, proto.Unmarshal(b, &value)
		},
		Job: func(_ context.Context, _ *Server, _ string, value proto.Message) (string, error) {
			return value.(*centrumunits.UnitSpec).GetJob(), nil
		},
		WrapPut: func(value proto.Message) *pbcentrum.StreamResponse {
			return &pbcentrum.StreamResponse{
				Change: &pbcentrum.StreamResponse_UnitSpecUpdated{
					UnitSpecUpdated: value.(*centrumunits.UnitSpec),
				},
			}
		},
	},
	{
		StreamName: "centrum_units",
		Bucket:     "centrum_units",
		KeyPrefix:  "members",
		NoWildcard: true,
		Unmarshal: func(ctx context.Context, _ *Server, b []byte) (proto.Message, error) {
			var value centrumunits.UnitAssignments
			return &value, proto.Unmarshal(b, &value)
		},
		Job: func(_ context.Context, _ *Server, _ string, value proto.Message) (string, error) {
			return value.(*centrumunits.UnitAssignments).GetJob(), nil
		},
		WrapPut: func(value proto.Message) *pbcentrum.StreamResponse {
			return &pbcentrum.StreamResponse{
				Change: &pbcentrum.StreamResponse_UnitMembersUpdated{
					UnitMembersUpdated: value.(*centrumunits.UnitAssignments),
				},
			}
		},
	},
	{
		StreamName: "centrum_units",
		Bucket:     "centrum_units",
		KeyPrefix:  "status",
		NoWildcard: true,
		Unmarshal: func(ctx context.Context, _ *Server, b []byte) (proto.Message, error) {
			var value centrumunits.UnitStatus
			return &value, proto.Unmarshal(b, &value)
		},
		Job: func(ctx context.Context, s *Server, _ string, value proto.Message) (string, error) {
			unit, err := s.units.Get(ctx, value.(*centrumunits.UnitStatus).GetUnitId())
			if err != nil || unit == nil {
				return "", err
			}
			return unit.GetJob(), nil
		},
		WrapPut: func(value proto.Message) *pbcentrum.StreamResponse {
			return &pbcentrum.StreamResponse{
				Change: &pbcentrum.StreamResponse_UnitStatus{
					UnitStatus: value.(*centrumunits.UnitStatus),
				},
			}
		},
	},
	{
		StreamName: "centrum_dispatches", Bucket: "centrum_dispatches", KeyPrefix: "job",
		Unmarshal: func(ctx context.Context, s *Server, b []byte) (proto.Message, error) {
			var mapping common.IDMapping
			if err := proto.Unmarshal(b, &mapping); err != nil {
				return nil, fmt.Errorf("failed to unmarshal dispatch id mapping. %w", err)
			}
			return s.dispatches.Get(ctx, mapping.GetId())
		},
		WrapPut: func(proto.Message) *pbcentrum.StreamResponse { return nil },
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
	{
		StreamName: "centrum_dispatches",
		Bucket:     "centrum_dispatches",
		KeyPrefix:  "spec",
		NoWildcard: true,
		Unmarshal: func(ctx context.Context, _ *Server, b []byte) (proto.Message, error) {
			var value centrumdispatches.DispatchSpec
			return &value, proto.Unmarshal(b, &value)
		},
		Jobs: func(_ context.Context, _ *Server, _ string, value proto.Message) ([]string, error) {
			return nonEmptyJobs(
				value.(*centrumdispatches.DispatchSpec).GetJobs().GetJobStrings(),
			), nil
		},
		WrapPut: func(value proto.Message) *pbcentrum.StreamResponse {
			return &pbcentrum.StreamResponse{
				Change: &pbcentrum.StreamResponse_DispatchSpecUpdated{
					DispatchSpecUpdated: value.(*centrumdispatches.DispatchSpec),
				},
			}
		},
	},
	{
		StreamName: "centrum_dispatches",
		Bucket:     "centrum_dispatches",
		KeyPrefix:  "assignments",
		NoWildcard: true,
		Unmarshal: func(ctx context.Context, _ *Server, b []byte) (proto.Message, error) {
			var value centrumdispatches.DispatchAssignments
			return &value, proto.Unmarshal(b, &value)
		},
		Jobs: func(ctx context.Context, s *Server, _ string, value proto.Message) ([]string, error) {
			return dispatchJobs(
				ctx,
				s,
				value.(*centrumdispatches.DispatchAssignments).GetDispatchId(),
			)
		},
		WrapPut: func(value proto.Message) *pbcentrum.StreamResponse {
			return &pbcentrum.StreamResponse{
				Change: &pbcentrum.StreamResponse_DispatchAssignmentsUpdated{
					DispatchAssignmentsUpdated: value.(*centrumdispatches.DispatchAssignments),
				},
			}
		},
	},
	{
		StreamName: "centrum_dispatches",
		Bucket:     "centrum_dispatches",
		KeyPrefix:  "status",
		NoWildcard: true,
		Unmarshal: func(ctx context.Context, _ *Server, b []byte) (proto.Message, error) {
			var value centrumdispatches.DispatchStatus
			return &value, proto.Unmarshal(b, &value)
		},
		Jobs: func(ctx context.Context, s *Server, _ string, value proto.Message) ([]string, error) {
			return dispatchJobs(ctx, s, value.(*centrumdispatches.DispatchStatus).GetDispatchId())
		},
		WrapPut: func(value proto.Message) *pbcentrum.StreamResponse {
			return &pbcentrum.StreamResponse{
				Change: &pbcentrum.StreamResponse_DispatchStatus{
					DispatchStatus: value.(*centrumdispatches.DispatchStatus),
				},
			}
		},
	},
}

func nonEmptyJobs(jobs []string) []string {
	return slices.DeleteFunc(slices.Clone(jobs), func(job string) bool {
		return strings.TrimSpace(job) == ""
	})
}

func dispatchJobs(ctx context.Context, s *Server, dispatchID int64) ([]string, error) {
	dispatch, err := s.dispatches.Get(ctx, dispatchID)
	if err != nil || dispatch == nil {
		return nil, err
	}
	return nonEmptyJobs(dispatch.GetJobs().GetJobStrings()), nil
}

type projectionFeedEventKey struct {
	feed      string
	id        int64
	component string
	deleted   bool
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
		return projectionFeedEventKey{
			feed:      "centrum_dispatches",
			id:        dispatch.GetId(),
			component: "full",
		}, true
	}
	if spec := response.GetDispatchSpecUpdated(); spec != nil && spec.GetId() > 0 {
		return projectionFeedEventKey{
			feed:      "centrum_dispatches",
			id:        spec.GetId(),
			component: "spec",
		}, true
	}
	if assignments := response.GetDispatchAssignmentsUpdated(); assignments != nil &&
		assignments.GetDispatchId() > 0 {
		return projectionFeedEventKey{
			feed:      "centrum_dispatches",
			id:        assignments.GetDispatchId(),
			component: "assignments",
		}, true
	}
	if status := response.GetDispatchStatus(); status != nil && status.GetDispatchId() > 0 {
		return projectionFeedEventKey{
			feed:      "centrum_dispatches",
			id:        status.GetDispatchId(),
			component: "status",
		}, true
	}
	if id := response.GetDispatchDeleted(); id > 0 {
		return projectionFeedEventKey{feed: "centrum_dispatches", id: id, deleted: true}, true
	}

	return projectionFeedEventKey{}, false
}

func (s *Server) queueDispatchFeedEventContext(
	ctx context.Context,
	job string,
	response *pbcentrum.StreamResponse,
	kvRevision uint64,
) {
	if dispatch := response.GetDispatchUpdated(); dispatch != nil {
		s.queueDispatchComponentFeedEvent(
			job,
			&pbcentrum.StreamResponse{Change: &pbcentrum.StreamResponse_DispatchSpecUpdated{
				DispatchSpecUpdated: dispatchSpecFromDispatch(dispatch),
			}},
			kvRevision,
		)
		s.queueDispatchComponentFeedEvent(
			job,
			&pbcentrum.StreamResponse{Change: &pbcentrum.StreamResponse_DispatchAssignmentsUpdated{
				DispatchAssignmentsUpdated: &centrumdispatches.DispatchAssignments{
					DispatchId: dispatch.GetId(),
					Job:        job,
					Units:      dispatch.GetUnits(),
				},
			}},
			componentRevision(ctx, s.dispatches, dispatch.GetId(), "assignments"),
		)
		if dispatch.GetStatus() != nil {
			s.queueDispatchComponentFeedEvent(
				job,
				&pbcentrum.StreamResponse{
					Change: &pbcentrum.StreamResponse_DispatchStatus{
						DispatchStatus: dispatch.GetStatus(),
					},
				},
				componentRevision(ctx, s.dispatches, dispatch.GetId(), "status"),
			)
		}
		return
	}
	if spec := response.GetDispatchSpecUpdated(); spec != nil && s.dispatches != nil {
		if dispatch, err := s.dispatches.Get(ctx, spec.GetId()); err == nil && dispatch != nil {
			for _, authorizedJob := range nonEmptyJobs(spec.GetJobs().GetJobStrings()) {
				s.queueDispatchComponentFeedEvent(authorizedJob, &pbcentrum.StreamResponse{
					Change: &pbcentrum.StreamResponse_DispatchAssignmentsUpdated{
						DispatchAssignmentsUpdated: &centrumdispatches.DispatchAssignments{
							DispatchId: dispatch.GetId(),
							Job:        authorizedJob,
							Units:      dispatch.GetUnits(),
						},
					},
				}, componentRevision(ctx, s.dispatches, dispatch.GetId(), "assignments"))
				if dispatch.GetStatus() != nil {
					s.queueDispatchComponentFeedEvent(authorizedJob, &pbcentrum.StreamResponse{
						Change: &pbcentrum.StreamResponse_DispatchStatus{
							DispatchStatus: dispatch.GetStatus(),
						},
					}, componentRevision(ctx, s.dispatches, dispatch.GetId(), "status"))
				}
			}
		}
	}
	key, ok := dispatchFeedKey(response)
	if !ok {
		s.publishFeedEvent("centrum_dispatches", []string{job}, response, kvRevision)
		return
	}
	s.queueProjectionFeedEvent(job, response, kvRevision, key)
}

type dispatchComponentRevisionReader interface {
	ComponentRevision(context.Context, int64, string) (uint64, error)
}

func componentRevision(
	ctx context.Context,
	dispatches dispatchComponentRevisionReader,
	id int64,
	component string,
) uint64 {
	if dispatches == nil ||
		(reflect.ValueOf(dispatches).Kind() == reflect.Pointer && reflect.ValueOf(dispatches).IsNil()) {
		return 0
	}
	revision, err := dispatches.ComponentRevision(ctx, id, component)
	if err != nil {
		return 0
	}
	return revision
}

func dispatchSpecFromDispatch(
	dispatch *centrumdispatches.Dispatch,
) *centrumdispatches.DispatchSpec {
	return &centrumdispatches.DispatchSpec{
		Id:          dispatch.GetId(),
		CreatedAt:   dispatch.CreatedAt,
		UpdatedAt:   dispatch.UpdatedAt,
		Jobs:        dispatch.GetJobs(),
		Message:     dispatch.GetMessage(),
		Description: dispatch.Description,
		Attributes:  dispatch.GetAttributes(),
		X:           dispatch.GetX(),
		Y:           dispatch.GetY(),
		Postal:      dispatch.Postal,
		Anon:        dispatch.GetAnon(),
		CreatorId:   dispatch.CreatorId,
		Creator:     dispatch.Creator,
		References:  dispatch.GetReferences(),
	}
}

func (s *Server) queueDispatchComponentFeedEvent(
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
		return projectionFeedEventKey{
			feed:      "centrum_units",
			id:        unit.GetId(),
			component: "full",
		}, true
	}
	if spec := response.GetUnitSpecUpdated(); spec != nil && spec.GetId() > 0 {
		return projectionFeedEventKey{
			feed:      "centrum_units",
			id:        spec.GetId(),
			component: "spec",
		}, true
	}
	if members := response.GetUnitMembersUpdated(); members != nil && members.GetUnitId() > 0 {
		return projectionFeedEventKey{
			feed:      "centrum_units",
			id:        members.GetUnitId(),
			component: "members",
		}, true
	}
	if status := response.GetUnitStatus(); status != nil && status.GetUnitId() > 0 {
		return projectionFeedEventKey{
			feed:      "centrum_units",
			id:        status.GetUnitId(),
			component: "status",
		}, true
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
	// The job mapping is refreshed for every component write. Fan the composed
	// value into independent wire updates so clients can reconcile only the
	// component that changed.
	if unit := response.GetUnitUpdated(); unit != nil {
		s.queueUnitComponentFeedEvent(job, &pbcentrum.StreamResponse{
			Change: &pbcentrum.StreamResponse_UnitSpecUpdated{
				UnitSpecUpdated: unitSpecFromUnit(unit),
			},
		}, kvRevision)
		s.queueUnitComponentFeedEvent(job, &pbcentrum.StreamResponse{
			Change: &pbcentrum.StreamResponse_UnitMembersUpdated{
				UnitMembersUpdated: &centrumunits.UnitAssignments{
					UnitId: unit.GetId(), Job: unit.GetJob(), Users: unit.GetUsers(),
				},
			},
		}, kvRevision)
		if unit.GetStatus() != nil {
			s.queueUnitComponentFeedEvent(job, &pbcentrum.StreamResponse{
				Change: &pbcentrum.StreamResponse_UnitStatus{UnitStatus: unit.GetStatus()},
			}, kvRevision)
		}
		return
	}

	key, ok := unitFeedKey(response)
	if !ok {
		s.publishFeedEvent("centrum_units", []string{job}, response, kvRevision)
		return
	}
	s.queueProjectionFeedEvent(job, response, kvRevision, key)
}

func unitSpecFromUnit(unit *centrumunits.Unit) *centrumunits.UnitSpec {
	return &centrumunits.UnitSpec{
		Id:          unit.GetId(),
		CreatedAt:   unit.GetCreatedAt(),
		UpdatedAt:   unit.GetUpdatedAt(),
		Job:         unit.GetJob(),
		JobLabel:    unit.JobLabel,
		SortOrder:   unit.GetSortOrder(),
		Name:        unit.GetName(),
		Initials:    unit.GetInitials(),
		Color:       unit.GetColor(),
		Icon:        unit.Icon,
		Description: unit.Description,
		Attributes:  unit.GetAttributes(),
		HomePostal:  unit.HomePostal,
		Access:      unit.GetAccess(),
	}
}

func (s *Server) queueUnitComponentFeedEvent(
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
		// A delete supersedes every component update queued for the resource.
		for _, component := range []string{"", "full", "spec", "members", "assignments", "status"} {
			s.flushProjectionFeedEventLocked(
				projectionFeedEventKey{feed: key.feed, id: key.id, component: component},
			)
		}
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
	consumer, err := s.js.CreateOrUpdateConsumer(ctx, "CENTRUM", jetstream.ConsumerConfig{
		Name:              centrumEventConsumerName(),
		FilterSubject:     string(eventscentrum.BaseSubject) + ".>",
		DeliverPolicy:     jetstream.DeliverNewPolicy,
		AckPolicy:         jetstream.AckNonePolicy,
		InactiveThreshold: time.Minute,
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

func (s *Server) consumeKVFeed(ctx context.Context, feed feedCfg) (bool, error) {
	filterSubject := "$KV." + feed.Bucket + ".>"
	if feed.KeyPrefix != "" {
		filterSubject = "$KV." + feed.Bucket + "." + feed.KeyPrefix + ".>"
	}
	consumer, err := s.js.CreateOrUpdateConsumer(
		ctx,
		"KV_"+feed.StreamName,
		jetstream.ConsumerConfig{
			Name:          kvFeedConsumerName(feed),
			FilterSubject: filterSubject,
			DeliverPolicy: jetstream.DeliverNewPolicy,
			AckPolicy:     jetstream.AckNonePolicy,
			MaxWaiting:    8,
			// Keep the named consumer stable across feed-worker restarts, but let
			// abandoned consumers from terminated instances be removed by NATS.
			InactiveThreshold: time.Minute,
		},
	)
	if err != nil {
		return false, fmt.Errorf("failed to create KV feed consumer. %w", err)
	}

	msgs, err := consumer.Messages(
		jetstream.PullMaxMessages(feedFetch),
		jetstream.WithMessagesErrOnMissingHeartbeat(false),
	)
	if err != nil {
		return true, err
	}
	defer msgs.Stop()

	for {
		msg, err := msgs.Next(jetstream.NextContext(ctx))
		if err != nil {
			if errors.Is(err, jetstream.ErrMsgIteratorClosed) {
				return true, nil
			}
			return true, err
		}

		key := strings.TrimPrefix(msg.Subject(), "$KV."+feed.Bucket+".")
		if feed.KeyPrefix != "" {
			key = strings.TrimPrefix(key, feed.KeyPrefix+".")
		}
		metadata, err := msg.Metadata()
		if err != nil {
			return true, fmt.Errorf("failed to read KV feed message metadata. %w", err)
		}
		kvRevision := metadata.Sequence.Stream
		job := key
		if !feed.NoWildcard {
			job, _, _ = strings.Cut(key, ".")
		}

		if op := msg.Headers().Get("KV-Operation"); op == "DEL" || op == "PURGE" {
			if feed.WrapDelete == nil {
				continue
			}
			response := feed.WrapDelete(key)
			switch feed.StreamName {
			case "centrum_dispatches":
				s.queueDispatchFeedEventContext(ctx, job, response, kvRevision)
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
		if response == nil {
			continue
		}
		jobs := []string{job}
		if feed.Job != nil {
			job, err = feed.Job(ctx, s, key, obj)
			if err != nil || job == "" {
				if err != nil {
					s.logger.Warn(
						"failed to determine KV feed authorization job",
						zap.Error(err),
						zap.String("bucket", feed.Bucket),
						zap.String("key", key),
					)
				}
				continue
			}
			jobs = []string{job}
		}
		if feed.Jobs != nil {
			jobs, err = feed.Jobs(ctx, s, key, obj)
			if err != nil || len(jobs) == 0 {
				if err != nil {
					s.logger.Warn(
						"failed to determine KV feed authorization jobs",
						zap.Error(err),
						zap.String("bucket", feed.Bucket),
						zap.String("key", key),
					)
				}
				continue
			}
		}
		switch feed.StreamName {
		case "centrum_dispatches":
			for _, job := range jobs {
				s.queueDispatchFeedEventContext(ctx, job, response, kvRevision)
			}
		case "centrum_units":
			for _, job := range jobs {
				s.queueUnitFeedEvent(job, response, kvRevision)
			}
		default:
			s.publishFeedEvent(feed.StreamName, jobs, response, kvRevision)
		}
	}
}
