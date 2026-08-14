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
	if err != nil {
		t.Fatalf("gather metrics: %v", err)
	}

	for _, family := range families {
		if family.GetName() != familyName {
			continue
		}
		for _, metric := range family.GetMetric() {
			if !labelsMatch(metric.GetLabel(), labels) {
				continue
			}
			if metric.GetHistogram() == nil {
				t.Fatalf("metric family %s label set found but histogram missing", familyName)
			}
			return metric.GetHistogram().GetSampleCount()
		}
	}

	t.Fatalf("metric family %s with labels %v not found", familyName, labels)
	return 0
}

func metricFamilyCounterValue(t *testing.T, familyName string, labels map[string]string) float64 {
	t.Helper()

	families, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		t.Fatalf("gather metrics: %v", err)
	}

	for _, family := range families {
		if family.GetName() != familyName {
			continue
		}
		for _, metric := range family.GetMetric() {
			if !labelsMatch(metric.GetLabel(), labels) {
				continue
			}
			if metric.GetCounter() == nil {
				t.Fatalf("metric family %s label set found but counter missing", familyName)
			}
			return metric.GetCounter().GetValue()
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

func cronFamilyName(subsystem, metric string) string {
	return prometheus.BuildFQName("fivenet", subsystem, metric)
}

func TestSchedulerRunCronjobRecordsHandoffLatency(t *testing.T) {
	js := &fakeCronJS{}
	s := &Scheduler{
		logger:    zap.NewNop(),
		publisher: js,
		metrics:   getSchedulerMetrics(),
	}

	job := &cron.Cronjob{Name: "croner.scheduler.test"}
	if _, err := s.runCronjob(t.Context(), job); err != nil {
		t.Fatalf("runCronjob: %v", err)
	}

	if got := len(js.published); got != 1 {
		t.Fatalf("published messages = %d, want 1", got)
	}
	wantSubject := string(CronScheduleSubject) + "." + string(CronScheduleTopic)
	if got := js.published[0].subject; got != wantSubject {
		t.Fatalf("published subject = %q, want %q", got, wantSubject)
	}

	family := cronFamilyName(schedulerMetricSubsystem, schedulerMetricHandoffLatency)
	if got := metricFamilyHistogramCount(
		t,
		family,
		map[string]string{"job": job.GetName()},
	); got != 1 {
		t.Fatalf("handoff latency observations = %d, want 1", got)
	}
}

func TestExecutorWatchForEventsRecordsSuccessAndDurations(t *testing.T) {
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
	if err != nil {
		t.Fatalf("marshal event: %v", err)
	}

	msg := &fakeMsg{
		data:    payload,
		subject: wantSubject,
	}
	exec.watchForEvents(msg)

	if !msg.acked {
		t.Fatal("expected message ack")
	}
	if got := len(js.published); got != 1 {
		t.Fatalf("published messages = %d, want 1", got)
	}

	completed, ok := js.published[0].msg.(*cron.CronjobCompletedEvent)
	if !ok {
		t.Fatalf(
			"published message type = %T, want *cron.CronjobCompletedEvent",
			js.published[0].msg,
		)
	}
	if !completed.GetSuccess() {
		t.Fatal("expected successful completion event")
	}

	startFamily := cronFamilyName(executorMetricSubsystem, executorMetricStartLatency)
	if got := metricFamilyHistogramCount(
		t,
		startFamily,
		map[string]string{"job": jobName},
	); got != 1 {
		t.Fatalf("start latency observations = %d, want 1", got)
	}

	durationFamily := cronFamilyName(executorMetricSubsystem, executorMetricHandlerDuration)
	if got := metricFamilyHistogramCount(
		t,
		durationFamily,
		map[string]string{"job": jobName},
	); got != 1 {
		t.Fatalf("handler duration observations = %d, want 1", got)
	}

	runFamily := cronFamilyName(executorMetricSubsystem, executorMetricRunsTotal)
	if got := metricFamilyCounterValue(
		t,
		runFamily,
		map[string]string{"job": jobName, "status": "success"},
	); got != 1 {
		t.Fatalf("success counter = %v, want 1", got)
	}
}

func TestExecutorWatchForEventsRecordsFailureOnPanic(t *testing.T) {
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
	if err != nil {
		t.Fatalf("marshal event: %v", err)
	}

	msg := &fakeMsg{
		data:    payload,
		subject: wantSubject,
	}
	exec.watchForEvents(msg)

	if !msg.acked {
		t.Fatal("expected message ack")
	}
	if got := len(js.published); got != 1 {
		t.Fatalf("published messages = %d, want 1", got)
	}

	completed, ok := js.published[0].msg.(*cron.CronjobCompletedEvent)
	if !ok {
		t.Fatalf(
			"published message type = %T, want *cron.CronjobCompletedEvent",
			js.published[0].msg,
		)
	}
	if completed.GetSuccess() {
		t.Fatal("expected failed completion event")
	}
	if completed.GetErrorMessage() == "" {
		t.Fatal("expected error message on failed completion event")
	}

	runFamily := cronFamilyName(executorMetricSubsystem, executorMetricRunsTotal)
	if got := metricFamilyCounterValue(
		t,
		runFamily,
		map[string]string{"job": jobName, "status": "failure"},
	); got != 1 {
		t.Fatalf("failure counter = %v, want 1", got)
	}
}
