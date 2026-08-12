package croner

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/cron"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/pkg/events"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type fakeCronJS struct {
	mu sync.Mutex

	published []publishedMsg
}

type publishedMsg struct {
	subject string
	msg     proto.Message
}

func (f *fakeCronJS) CreateOrUpdateConsumer(
	context.Context,
	string,
	jetstream.ConsumerConfig,
) (jetstream.Consumer, error) {
	return nil, errors.New("not implemented")
}

func (f *fakeCronJS) ConsumeErrHandlerWithRestart(
	context.Context,
	*zap.Logger,
	events.ConsumeErrRestartFn,
) jetstream.PullConsumeOpt {
	return nil
}

func (f *fakeCronJS) PublishProto(
	_ context.Context,
	subject string,
	msg proto.Message,
	_ ...jetstream.PublishOpt,
) (*jetstream.PubAck, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.published = append(f.published, publishedMsg{
		subject: subject,
		msg:     msg,
	})

	return &jetstream.PubAck{Stream: "cron", Sequence: uint64(len(f.published))}, nil
}

type fakeMsg struct {
	data    []byte
	subject string
	acked   bool
}

func (m *fakeMsg) Metadata() (*jetstream.MsgMetadata, error) { return nil, nil }
func (m *fakeMsg) Data() []byte                              { return m.data }
func (m *fakeMsg) Headers() nats.Header                      { return nil }
func (m *fakeMsg) Subject() string                           { return m.subject }
func (m *fakeMsg) Reply() string                             { return "" }
func (m *fakeMsg) Ack() error                                { m.acked = true; return nil }
func (m *fakeMsg) DoubleAck(context.Context) error           { return nil }
func (m *fakeMsg) Nak() error                                { return nil }
func (m *fakeMsg) NakWithDelay(time.Duration) error          { return nil }
func (m *fakeMsg) InProgress() error                         { return nil }
func (m *fakeMsg) Term() error                               { return nil }
func (m *fakeMsg) TermWithReason(string) error               { return nil }

func metricFamilyHistogramCount(t *testing.T, familyName string, labels map[string]string) uint64 {
	t.Helper()

	families, err := prometheus.DefaultGatherer.Gather()
	require.NoError(t, err)

	for _, family := range families {
		if family.GetName() != familyName {
			continue
		}
		for _, metric := range family.GetMetric() {
			if !labelsMatch(metric.GetLabel(), labels) {
				continue
			}
			require.NotNil(
				t,
				metric.GetHistogram(),
				"metric family %s label set found but histogram missing",
				familyName,
			)
			return metric.GetHistogram().GetSampleCount()
		}
	}

	t.Fatalf("metric family %s with labels %v not found", familyName, labels)
	return 0
}

func metricFamilyGaugeValue(t *testing.T, familyName string, labels map[string]string) float64 {
	t.Helper()

	families, err := prometheus.DefaultGatherer.Gather()
	require.NoError(t, err)

	for _, family := range families {
		if family.GetName() != familyName {
			continue
		}
		for _, metric := range family.GetMetric() {
			if !labelsMatch(metric.GetLabel(), labels) {
				continue
			}
			require.NotNil(
				t,
				metric.GetGauge(),
				"metric family %s label set found but gauge missing",
				familyName,
			)
			return metric.GetGauge().GetValue()
		}
	}

	t.Fatalf("metric family %s with labels %v not found", familyName, labels)
	return 0
}

func labelsMatch(actual []*io_prometheus_client.LabelPair, expected map[string]string) bool {
	if len(actual) != len(expected) {
		return false
	}
	for _, label := range actual {
		if expected[label.GetName()] != label.GetValue() {
			return false
		}
	}
	return true
}

func cronFamilyName(t *testing.T, subsystem, metric string) string {
	t.Helper()
	return prometheus.BuildFQName("fivenet", subsystem, metric)
}

func TestSchedulerRunCronjobRecordsHandoffLatency(t *testing.T) {
	t.Parallel()

	js := &fakeCronJS{}
	s := &Scheduler{
		logger:    zap.NewNop(),
		publisher: js,
		metrics:   getSchedulerMetrics(),
	}

	job := &cron.Cronjob{Name: "croner.scheduler.test"}
	_, err := s.runCronjob(t.Context(), job)
	require.NoError(t, err)

	require.Len(t, js.published, 1)
	wantSubject := string(CronScheduleSubject) + "." + string(CronScheduleTopic)
	require.Equal(t, wantSubject, js.published[0].subject)

	family := cronFamilyName(t, schedulerMetricSubsystem, schedulerMetricHandoffLatency)
	latency := metricFamilyHistogramCount(
		t,
		family,
		map[string]string{"job": job.GetName()},
	)
	require.Equal(t, uint64(1), latency)
}

func TestExecutorWatchForEventsRecordsSuccessAndDurations(t *testing.T) {
	t.Parallel()

	jobName := "croner.executor.success"
	wantSubject := string(CronScheduleSubject) + "." + string(CronScheduleTopic)
	js := &fakeCronJS{}
	exec := &Executor{
		logger:    zap.NewNop(),
		ctx:       t.Context(),
		publisher: js,
		handlers: &Handlers{
			handlers: map[string]CronjobHandlerFn{
				jobName: func(ctx context.Context, data *cron.CronjobData) error {
					time.Sleep(5 * time.Millisecond)
					return nil
				},
			},
		},
		metrics: getExecutorMetrics(),
	}

	startedAt := time.Now().Add(-20 * time.Millisecond)
	payload, err := protojson.Marshal(&cron.CronjobSchedulerEvent{
		Cronjob: &cron.Cronjob{
			Name:        jobName,
			StartedTime: timestamp.New(startedAt),
			Data:        &cron.CronjobData{Data: nil},
		},
	})
	require.NoError(t, err)

	msg := &fakeMsg{
		data:    payload,
		subject: wantSubject,
	}
	exec.watchForEvents(msg)

	require.True(t, msg.acked, "expected message ack")
	require.Len(t, js.published, 1)

	completed, ok := js.published[0].msg.(*cron.CronjobCompletedEvent)
	require.True(
		t,
		ok,
		"published message type = %T, want *cron.CronjobCompletedEvent",
		js.published[0].msg,
	)
	require.True(t, completed.GetSuccess(), "expected successful completion event")

	startFamily := cronFamilyName(t, executorMetricSubsystem, executorMetricStartLatency)
	histoCount := metricFamilyHistogramCount(
		t,
		startFamily,
		map[string]string{"job": jobName},
	)
	require.Equal(t, uint64(1), histoCount)

	durationFamily := cronFamilyName(t, executorMetricSubsystem, executorMetricHandlerDuration)
	observations := metricFamilyHistogramCount(
		t,
		durationFamily,
		map[string]string{"job": jobName},
	)
	require.Equal(t, uint64(1), observations)

	runFamily := cronFamilyName(t, executorMetricSubsystem, executorMetricLastRunSuccess)
	success := metricFamilyGaugeValue(
		t,
		runFamily,
		map[string]string{"job": jobName},
	)
	require.InDelta(t, float64(1), success, 0.000001)
}

func TestExecutorWatchForEventsRecordsFailureOnPanic(t *testing.T) {
	t.Parallel()

	jobName := "croner.executor.failure"
	wantSubject := string(CronScheduleSubject) + "." + string(CronScheduleTopic)
	js := &fakeCronJS{}
	exec := &Executor{
		logger:    zap.NewNop(),
		ctx:       t.Context(),
		publisher: js,
		handlers: &Handlers{
			handlers: map[string]CronjobHandlerFn{
				jobName: func(ctx context.Context, data *cron.CronjobData) error {
					panic("boom")
				},
			},
		},
		metrics: getExecutorMetrics(),
	}

	payload, err := protojson.Marshal(&cron.CronjobSchedulerEvent{
		Cronjob: &cron.Cronjob{
			Name:        jobName,
			StartedTime: timestamp.New(time.Now().Add(-10 * time.Millisecond)),
		},
	})
	require.NoError(t, err)

	msg := &fakeMsg{
		data:    payload,
		subject: wantSubject,
	}
	exec.watchForEvents(msg)

	require.True(t, msg.acked, "expected message ack")
	require.Len(t, js.published, 1)

	completed, ok := js.published[0].msg.(*cron.CronjobCompletedEvent)
	require.True(
		t,
		ok,
		"published message type = %T, want *cron.CronjobCompletedEvent",
		js.published[0].msg,
	)
	require.False(t, completed.GetSuccess(), "expected failed completion event")
	require.NotEmpty(
		t,
		completed.GetErrorMessage(),
		"expected error message on failed completion event",
	)

	runFamily := cronFamilyName(t, executorMetricSubsystem, executorMetricLastRunSuccess)
	failure := metricFamilyGaugeValue(
		t,
		runFamily,
		map[string]string{"job": jobName},
	)
	require.InDelta(t, float64(0), failure, 0.000001)
}
