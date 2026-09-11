package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	activeStreams       prometheus.Gauge
	feedMessages        *prometheus.CounterVec
	feedResyncs         *prometheus.CounterVec
	feedDecodeFailures  *prometheus.CounterVec
	housekeeperDuration *prometheus.HistogramVec
	housekeeperWork     *prometheus.GaugeVec
}

var (
	instance *Metrics
	once     sync.Once
)

func Get() *Metrics {
	once.Do(func() {
		instance = &Metrics{
			activeStreams: prometheus.NewGauge(prometheus.GaugeOpts{
				Namespace: "fivenet",
				Subsystem: "centrum",
				Name:      "active_streams",
				Help:      "Number of active Centrum streams served by this process.",
			}),
			feedMessages: prometheus.NewCounterVec(prometheus.CounterOpts{
				Namespace: "fivenet",
				Subsystem: "centrum",
				Name:      "feed_messages_total",
				Help:      "Centrum feed messages by source and outcome.",
			}, []string{"feed", "outcome"}),
			feedResyncs: prometheus.NewCounterVec(prometheus.CounterOpts{
				Namespace: "fivenet",
				Subsystem: "centrum",
				Name:      "feed_resyncs_total",
				Help:      "Centrum stream resyncs caused by feed delivery gaps.",
			}, []string{"reason"}),
			feedDecodeFailures: prometheus.NewCounterVec(prometheus.CounterOpts{
				Namespace: "fivenet",
				Subsystem: "centrum",
				Name:      "feed_decode_failures_total",
				Help:      "Centrum feed payloads that could not be decoded.",
			}, []string{"feed"}),
			housekeeperDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
				Namespace: "fivenet",
				Subsystem: "centrum",
				Name:      "housekeeper_duration_seconds",
				Help:      "Duration of Centrum housekeeper cleanup and audit tasks.",
				Buckets:   prometheus.ExponentialBuckets(0.01, 2, 14),
			}, []string{"task"}),
			housekeeperWork: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: "fivenet",
				Subsystem: "centrum",
				Name:      "housekeeper_work_items",
				Help:      "Work performed by the latest Centrum housekeeper run.",
			}, []string{"task", "operation"}),
		}

		prometheus.MustRegister(
			instance.activeStreams,
			instance.feedMessages,
			instance.feedResyncs,
			instance.feedDecodeFailures,
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

func (m *Metrics) ObserveHousekeeperDuration(task string, seconds float64) {
	m.housekeeperDuration.WithLabelValues(task).Observe(seconds)
}

func (m *Metrics) SetHousekeeperWork(task string, operation string, count int) {
	m.housekeeperWork.WithLabelValues(task, operation).Set(float64(count))
}
