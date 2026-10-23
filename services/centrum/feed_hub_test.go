package centrum

import (
	"context"
	"slices"
	"sync"
	"testing"
	"time"

	centrum "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum"
	centrumdispatches "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatches"
	centrumsettings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/settings"
	centrumunits "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/units"
	pbcentrum "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/centrum"
	testnats "github.com/fivenet-app/fivenet/v2026/internal/tests/nats"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/broker"
	eventscentrum "github.com/fivenet-app/fivenet/v2026/services/centrum/events"
	centrummetrics "github.com/fivenet-app/fivenet/v2026/services/centrum/metrics"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

func newFeedHubTestServer(t *testing.T) (*Server, context.CancelFunc) {
	t.Helper()

	js := testnats.NewServer(t, testnats.ServerOptions{InProcess: true}).GetJS()

	ctx, cancel := context.WithCancel(t.Context())
	srv := &Server{
		logger:     zap.NewNop(),
		js:         js,
		wg:         sync.WaitGroup{},
		feedBroker: broker.NewWithResyncOnSlowSubscriber[*feedEvent](16),
		metrics:    centrummetrics.Get(),
	}

	_, err := eventscentrum.RegisterStream(ctx, js)
	require.NoError(t, err)
	_, err = js.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{Bucket: "centrum_settings"})
	require.NoError(t, err)

	srv.wg.Go(func() {
		srv.feedBroker.Start(ctx)
	})
	srv.wg.Go(func() {
		srv.runCentrumEventFeed(ctx)
	})
	srv.wg.Go(func() {
		srv.runKVFeed(ctx, feeds[0])
	})

	t.Cleanup(func() {
		cancel()
		srv.wg.Wait()
	})

	return srv, cancel
}

func TestKVFeedConsumerNamesAreUniquePerProjection(t *testing.T) {
	t.Parallel()

	names := make(map[string]string, len(feeds))
	for _, feed := range feeds {
		name := kvFeedConsumerName(feed)
		if previous, ok := names[name]; ok {
			t.Fatalf("feeds %q and %q share consumer name %q", previous, feed.KeyPrefix, name)
		}
		names[name] = feed.KeyPrefix
	}

	assert.Equal(t, "KV_centrum_settings", kvFeedConsumerName(feeds[0]))
	assert.Equal(t, "KV_centrum_units_spec", kvFeedConsumerName(feeds[3]))
	assert.Equal(t, "KV_centrum_dispatches_assignments", kvFeedConsumerName(feeds[8]))
}

func TestFeedHubPublishesResyncEvent(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	srv := &Server{
		feedBroker: broker.NewWithResyncOnSlowSubscriber[*feedEvent](10),
		metrics:    centrummetrics.Get(),
	}
	go srv.feedBroker.Start(ctx)

	feed := srv.feedBroker.Subscribe()
	srv.publishFeedResync("events")
	event := waitForFeedEvent(t, feed, func(event *feedEvent) bool { return event.Resync })
	assert.Equal(t, uint64(1), event.Sequence)
}

func TestFeedHubDoesNotAdvanceSequenceWithoutSubscribers(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	srv := &Server{
		feedBroker: broker.NewWithResyncOnSlowSubscriber[*feedEvent](10),
		metrics:    centrummetrics.Get(),
	}
	srv.publishFeedEvent("events", nil, &pbcentrum.StreamResponse{}, 0)
	assert.Equal(t, uint64(0), srv.feedSequence.Load())

	go srv.feedBroker.Start(ctx)
	feed := srv.feedBroker.Subscribe()
	srv.publishFeedEvent("events", nil, &pbcentrum.StreamResponse{}, 0)
	event := waitForFeedEvent(t, feed, func(event *feedEvent) bool { return !event.Resync })
	assert.Equal(t, uint64(1), event.Sequence)
}

func TestFeedHubCoalescesDispatchMappingEvents(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	srv := &Server{
		feedBroker: broker.NewWithResyncOnSlowSubscriber[*feedEvent](10),
		metrics:    centrummetrics.Get(),
	}
	go srv.feedBroker.Start(ctx)

	feed := srv.feedBroker.Subscribe()
	defer srv.feedBroker.Unsubscribe(feed)

	for _, job := range []string{"ambulance", "fire", "police"} {
		srv.queueDispatchFeedEventContext(ctx, job, &pbcentrum.StreamResponse{
			Change: &pbcentrum.StreamResponse_DispatchUpdated{
				DispatchUpdated: &centrumdispatches.Dispatch{Id: 42},
			},
		}, 7)
	}

	event := waitForFeedEvent(t, feed, func(event *feedEvent) bool {
		return event.Response.GetDispatchSpecUpdated().GetId() == 42
	})
	assert.Equal(t, []string{"ambulance", "fire", "police"}, event.Jobs)
	assert.Equal(t, uint64(7), event.Response.GetKvRevision())
}

func TestFeedHubDoesNotReuseSpecRevisionForComposedDispatchComponents(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	srv := &Server{
		feedBroker: broker.NewWithResyncOnSlowSubscriber[*feedEvent](10),
		metrics:    centrummetrics.Get(),
	}
	go srv.feedBroker.Start(ctx)

	feed := srv.feedBroker.Subscribe()
	defer srv.feedBroker.Unsubscribe(feed)

	srv.queueDispatchFeedEventContext(ctx, "police", &pbcentrum.StreamResponse{
		Change: &pbcentrum.StreamResponse_DispatchUpdated{
			DispatchUpdated: &centrumdispatches.Dispatch{
				Id:     42,
				Status: &centrumdispatches.DispatchStatus{DispatchId: 42},
			},
		},
	}, 101)

	seen := map[string]uint64{}
	deadline := time.NewTimer(time.Second)
	defer deadline.Stop()
	for len(seen) < 3 {
		select {
		case event := <-feed:
			response := event.Response
			switch {
			case response.GetDispatchSpecUpdated() != nil:
				seen["spec"] = response.GetKvRevision()
			case response.GetDispatchAssignmentsUpdated() != nil:
				seen["assignments"] = response.GetKvRevision()
			case response.GetDispatchStatus() != nil:
				seen["status"] = response.GetKvRevision()
			}
		case <-deadline.C:
			require.FailNow(t, "timed out waiting for composed dispatch components")
		}
	}

	require.Equal(t, uint64(101), seen["spec"])
	require.Zero(t, seen["assignments"])
	require.Zero(t, seen["status"])
}

func TestFeedHubPublishesDispatchComponentToEveryJob(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	srv := &Server{
		feedBroker: broker.NewWithResyncOnSlowSubscriber[*feedEvent](10),
		metrics:    centrummetrics.Get(),
	}
	go srv.feedBroker.Start(ctx)

	feed := srv.feedBroker.Subscribe()
	defer srv.feedBroker.Unsubscribe(feed)

	spec := &centrumdispatches.DispatchSpec{
		Id:   42,
		Jobs: &centrum.JobList{Jobs: []*centrum.JobListEntry{{Name: "ambulance"}, {Name: "fire"}}},
	}
	config := slices.IndexFunc(feeds, func(feed feedCfg) bool {
		return feed.StreamName == "centrum_dispatches" && feed.KeyPrefix == "spec"
	})
	require.GreaterOrEqual(t, config, 0)
	jobs, err := feeds[config].Jobs(ctx, srv, "42", spec)
	require.NoError(t, err)

	response := &pbcentrum.StreamResponse{
		Change: &pbcentrum.StreamResponse_DispatchSpecUpdated{DispatchSpecUpdated: spec},
	}
	for _, job := range jobs {
		srv.queueDispatchFeedEventContext(ctx, job, response, 7)
	}

	event := waitForFeedEvent(t, feed, func(event *feedEvent) bool {
		return event.Response.GetDispatchSpecUpdated().GetId() == 42
	})
	assert.Equal(t, []string{"ambulance", "fire"}, event.Jobs)
}

func TestFeedHubCoalescesUnitMappingEvents(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	srv := &Server{
		feedBroker: broker.NewWithResyncOnSlowSubscriber[*feedEvent](10),
		metrics:    centrummetrics.Get(),
	}
	go srv.feedBroker.Start(ctx)

	feed := srv.feedBroker.Subscribe()
	defer srv.feedBroker.Unsubscribe(feed)

	for _, job := range []string{"ambulance", "fire", "police"} {
		srv.queueUnitFeedEvent(job, &pbcentrum.StreamResponse{
			Change: &pbcentrum.StreamResponse_UnitUpdated{
				UnitUpdated: &centrumunits.Unit{Id: 42},
			},
		}, 7)
	}

	event := waitForFeedEvent(t, feed, func(event *feedEvent) bool {
		return event.Response.GetUnitSpecUpdated().GetId() == 42
	})
	assert.Equal(t, []string{"ambulance", "fire", "police"}, event.Jobs)
	assert.Equal(t, uint64(7), event.Response.GetKvRevision())

	// The composed update is fanned out into spec, members and status events.
}

func TestFeedHubPreservesDispatchUpdateBeforeDeletion(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	srv := &Server{
		feedBroker: broker.NewWithResyncOnSlowSubscriber[*feedEvent](10),
		metrics:    centrummetrics.Get(),
	}
	go srv.feedBroker.Start(ctx)

	feed := srv.feedBroker.Subscribe()
	defer srv.feedBroker.Unsubscribe(feed)

	srv.queueDispatchFeedEventContext(ctx, "ambulance", &pbcentrum.StreamResponse{
		Change: &pbcentrum.StreamResponse_DispatchUpdated{
			DispatchUpdated: &centrumdispatches.Dispatch{Id: 42},
		},
	}, 7)
	srv.queueDispatchFeedEventContext(ctx, "ambulance", &pbcentrum.StreamResponse{
		Change: &pbcentrum.StreamResponse_DispatchDeleted{DispatchDeleted: 42},
	}, 8)

	updated := waitForFeedEvent(t, feed, func(event *feedEvent) bool {
		return event.Response.GetDispatchSpecUpdated().GetId() == 42
	})
	deleted := waitForFeedEvent(t, feed, func(event *feedEvent) bool {
		return event.Response.GetDispatchDeleted() == 42
	})
	assert.Less(t, updated.Sequence, deleted.Sequence)
}

func waitForFeedHubConsumer(t *testing.T, srv *Server, streamName string) {
	t.Helper()

	require.Eventually(t, func() bool {
		stream, err := srv.js.Stream(t.Context(), streamName)
		if err != nil {
			return false
		}
		_, ok := <-stream.ConsumerNames(t.Context()).Name()
		return ok
	}, 4*time.Second, 10*time.Millisecond)
}

func waitForFeedEvent(
	t *testing.T,
	feed <-chan *feedEvent,
	match func(*feedEvent) bool,
) *feedEvent {
	t.Helper()
	return waitForFeedEventWithin(t, feed, match, 2*time.Second)
}

func waitForFeedEventWithin(
	t *testing.T,
	feed <-chan *feedEvent,
	match func(*feedEvent) bool,
	timeoutDuration time.Duration,
) *feedEvent {
	t.Helper()

	timeout := time.NewTimer(timeoutDuration)
	defer timeout.Stop()
	for {
		select {
		case event, ok := <-feed:
			require.True(t, ok, "feed closed before the expected event arrived")
			if match(event) {
				return event
			}
		case <-timeout.C:
			require.FailNow(t, "timed out waiting for feed event")
		}
	}
}

func TestFeedHubForwardsJetStreamEventsAndKVUpdates(t *testing.T) {
	t.Parallel()

	srv, _ := newFeedHubTestServer(t)
	waitForFeedHubConsumer(t, srv, "CENTRUM")
	waitForFeedHubConsumer(t, srv, "KV_centrum_settings")

	feed := srv.feedBroker.Subscribe()
	defer srv.feedBroker.Unsubscribe(feed)

	status := &centrumunits.UnitStatus{
		Id:     42,
		UnitId: 7,
		Status: centrumunits.StatusUnit_STATUS_UNIT_BUSY,
	}
	data, err := proto.Marshal(status)
	require.NoError(t, err)
	_, err = srv.js.Publish(
		t.Context(),
		eventscentrum.BuildSubject(
			eventscentrum.TopicUnit,
			eventscentrum.TypeUnitStatus,
			"ambulance",
		),
		data,
	)
	require.NoError(t, err)

	event := waitForFeedEvent(t, feed, func(event *feedEvent) bool {
		return assert.Equal(t, []string{"ambulance"}, event.Jobs) &&
			event.Response.GetUnitStatus().GetId() == status.GetId()
	})
	assert.Equal(t, status.GetUnitId(), event.Response.GetUnitStatus().GetUnitId())
	assert.Zero(t, event.Response.GetKvRevision())

	settingsKV, err := srv.js.KeyValue(t.Context(), "centrum_settings")
	require.NoError(t, err)
	settingsData, err := proto.Marshal(&centrumsettings.Settings{Job: "police", Enabled: true})
	require.NoError(t, err)
	_, err = settingsKV.Put(t.Context(), "police", settingsData)
	require.NoError(t, err)

	event = waitForFeedEvent(t, feed, func(event *feedEvent) bool {
		return assert.Equal(t, []string{"police"}, event.Jobs) &&
			event.Response.GetSettings().GetJob() == "police"
	})
	assert.True(t, event.Response.GetSettings().GetEnabled())
	assert.NotZero(t, event.Response.GetKvRevision())
}

func TestFeedHubResyncsAndRecoversAfterKVConsumerDeletion(t *testing.T) {
	t.Parallel()

	srv, _ := newFeedHubTestServer(t)
	waitForFeedHubConsumer(t, srv, "KV_centrum_settings")

	feed := srv.feedBroker.Subscribe()
	defer srv.feedBroker.Unsubscribe(feed)

	stream, err := srv.js.Stream(t.Context(), "KV_centrum_settings")
	require.NoError(t, err)
	deleted := false
	for name := range stream.ConsumerNames(t.Context()).Name() {
		require.NoError(t, stream.DeleteConsumer(t.Context(), name))
		deleted = true
		break
	}
	require.True(t, deleted, "expected the settings feed consumer to exist")

	settingsKV, err := srv.js.KeyValue(t.Context(), "centrum_settings")
	require.NoError(t, err)
	data, err := proto.Marshal(&centrumsettings.Settings{Job: "fire", Enabled: true})
	require.NoError(t, err)
	_, err = settingsKV.Put(t.Context(), "fire", data)
	require.NoError(t, err)

	// A client receiving this event rebuilds its handshake/latest-state snapshot
	// before accepting events from the replacement worker.
	resync := waitForFeedEventWithin(
		t,
		feed,
		func(event *feedEvent) bool { return event.Resync },
		7*time.Second,
	)
	require.NotZero(t, resync.Sequence)
	waitForFeedHubConsumer(t, srv, "KV_centrum_settings")

	data, err = proto.Marshal(&centrumsettings.Settings{Job: "ambulance", Enabled: true})
	require.NoError(t, err)
	_, err = settingsKV.Put(t.Context(), "ambulance", data)
	require.NoError(t, err)

	event := waitForFeedEvent(t, feed, func(event *feedEvent) bool {
		return event.Response.GetSettings().GetJob() == "ambulance"
	})
	assert.True(t, event.Response.GetSettings().GetEnabled())
}
