package centrum

import (
	"context"
	"sync"
	"testing"
	"time"

	centrumsettings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/settings"
	centrumunits "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/units"
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

	_, js, cleanup, err := testnats.NewInProcessNATSServer()
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, cleanup())
	})

	ctx, cancel := context.WithCancel(t.Context())
	srv := &Server{
		logger:     zap.NewNop(),
		js:         js,
		wg:         sync.WaitGroup{},
		feedBroker: broker.NewWithResyncOnSlowSubscriber[*feedEvent](16),
		metrics:    centrummetrics.Get(),
	}

	_, err = eventscentrum.RegisterStream(ctx, js)
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

func waitForFeedHubConsumer(t *testing.T, srv *Server, streamName string) {
	t.Helper()

	require.Eventually(t, func() bool {
		stream, err := srv.js.Stream(t.Context(), streamName)
		if err != nil {
			return false
		}
		_, ok := <-stream.ConsumerNames(t.Context()).Name()
		return ok
	}, 2*time.Second, 10*time.Millisecond)
}

func waitForFeedEvent(t *testing.T, feed <-chan *feedEvent, match func(*feedEvent) bool) *feedEvent {
	t.Helper()

	timeout := time.NewTimer(2 * time.Second)
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

	status := &centrumunits.UnitStatus{Id: 42, UnitId: 7, Status: centrumunits.StatusUnit_STATUS_UNIT_BUSY}
	data, err := proto.Marshal(status)
	require.NoError(t, err)
	_, err = srv.js.Publish(t.Context(), eventscentrum.BuildSubject(eventscentrum.TopicUnit, eventscentrum.TypeUnitStatus, "ambulance"), data)
	require.NoError(t, err)

	event := waitForFeedEvent(t, feed, func(event *feedEvent) bool {
		return event.Job == "ambulance" && event.Response.GetUnitStatus().GetId() == status.GetId()
	})
	assert.Equal(t, status.GetUnitId(), event.Response.GetUnitStatus().GetUnitId())

	settingsKV, err := srv.js.KeyValue(t.Context(), "centrum_settings")
	require.NoError(t, err)
	settingsData, err := proto.Marshal(&centrumsettings.Settings{Job: "police", Enabled: true})
	require.NoError(t, err)
	_, err = settingsKV.Put(t.Context(), "police", settingsData)
	require.NoError(t, err)

	event = waitForFeedEvent(t, feed, func(event *feedEvent) bool {
		return event.Job == "police" && event.Response.GetSettings().GetJob() == "police"
	})
	assert.True(t, event.Response.GetSettings().GetEnabled())
}
