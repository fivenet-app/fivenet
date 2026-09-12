package metrics

import (
	"sync"

	"github.com/fivenet-app/fivenet/v2026/pkg/server/admin"
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	activeStreams       prometheus.Gauge
	feedMessages        *prometheus.CounterVec
	feedResyncs         *prometheus.CounterVec
	feedDecodeFailures  *prometheus.CounterVec
	housekeeperEvents   *prometheus.CounterVec
	housekeeperDuration *prometheus.HistogramVec
	housekeeperWork     *prometheus.GaugeVec
}

var (
	instance *Metrics
	once     sync.Once
)

const metricsSubsystem = "centrum"

func Get() *Metrics {
	once.Do(func() {
		instance = &Metrics{
			activeStreams: prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace: admin.MetricsNamespace,
				Subsystem: metricsSubsystem,
				Name:      "active_streams",
				Help:      "Number of active Centrum streams served by this process.",
			}),
			feedMessages: prometheus.NewCounterVec(prometheus.CounterOpts{
				Namespace: admin.MetricsNamespace,
				Subsystem: metricsSubsystem,
				Name:      "feed_messages_total",
				Help:      "Centrum feed messages by source and outcome.",
			}, []string{"feed", "outcome"}),
			feedResyncs: prometheus.NewCounterVec(prometheus.CounterOpts{
				Namespace: admin.MetricsNamespace,
				Subsystem: metricsSubsystem,
				Name:      "feed_resyncs_total",
				Help:      "Centrum stream resnapshots by delivery or authorization reason.",
			}, []string{"reason"}),
			feedDecodeFailures: prometheus.NewCounterVec(prometheus.CounterOpts{
				Namespace: admin.MetricsNamespace,
				Subsystem: metricsSubsystem,
				Name:      "feed_decode_failures_total",
				Help:      "Centrum feed payloads that could not be decoded.",
			}, []string{"feed"}),
			housekeeperEvents: prometheus.NewCounterVec(prometheus.CounterOpts{
				Namespace: admin.MetricsNamespace,
				Subsystem: metricsSubsystem,
				Name:      "housekeeper_events_total",
				Help:      "Centrum housekeeper watcher events and corrections by outcome.",
			}, []string{"watcher", "outcome"}),
			housekeeperDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
				Namespace: admin.MetricsNamespace,
				Subsystem: metricsSubsystem,
				Name:      "housekeeper_duration_seconds",
				Help:      "Duration of Centrum housekeeper cleanup and audit tasks.",
				Buckets:   prometheus.ExponentialBuckets(0.01, 2, 14),
			}, []string{"task"}),
			housekeeperWork: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: admin.MetricsNamespace,
				Subsystem: metricsSubsystem,
				Name:      "housekeeper_work_items",
				Help:      "Work performed by the latest Centrum housekeeper run.",
			}, []string{"task", "operation"}),
		}

		prometheus.MustRegister(
			instance.activeStreams,
			instance.feedMessages,
			instance.feedResyncs,
			instance.feedDecodeFailures,
			instance.housekeeperEvents,
			instance.housekeeperDuration,
			instance.housekeeperWork,
		)
	})

	return instance
}

func (m *Metrics) IncActiveStreams() { m.activeStreams.Inc() }

func (m *Metrics) DecActiveStreams() { m.activeStreams.Dec() }

func (m *Metrics) IncFeedMessage(feed string, outcome string) {
	m.feedMessages.WithLabelValues(feed, outcome).Inc()
}

func (m *Metrics) IncFeedResync(reason string) { m.feedResyncs.WithLabelValues(reason).Inc() }

func (m *Metrics) IncFeedDecodeFailure(feed string) {
	m.feedDecodeFailures.WithLabelValues(feed).Inc()
}

func (m *Metrics) IncHousekeeperEvent(watcher string, outcome string) {
	m.housekeeperEvents.WithLabelValues(watcher, outcome).Inc()
}

func (m *Metrics) AddHousekeeperEvents(watcher string, outcome string, count int) {
	if count > 0 {
		m.housekeeperEvents.WithLabelValues(watcher, outcome).Add(float64(count))
	}
}

func (m *Metrics) ObserveHousekeeperDuration(task string, seconds float64) {
	m.housekeeperDuration.WithLabelValues(task).Observe(seconds)
}

func (m *Metrics) SetHousekeeperWork(task string, operation string, count int) {
	m.housekeeperWork.WithLabelValues(task, operation).Set(float64(count))
}
