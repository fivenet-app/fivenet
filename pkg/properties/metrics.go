package properties

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	propertyExpiryCleanupRuns = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "fivenet",
		Subsystem: "properties",
		Name:      "expiry_cleanup_runs_total",
		Help:      "Completed dynamic-property expiry cleanup runs by scope and outcome.",
	}, []string{"scope", "outcome"})
	propertyExpiryCleanupCleared = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "fivenet",
		Subsystem: "properties",
		Name:      "expiry_cleanup_cleared_total",
		Help:      "Dynamic-property values cleared by expiry cleanup.",
	}, []string{"scope"})
	propertyExpiryCleanupFailures = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "fivenet",
		Subsystem: "properties",
		Name:      "expiry_cleanup_item_failures_total",
		Help:      "Individual dynamic-property expiry cleanup items that could not be cleared.",
	}, []string{"scope"})
	propertyExpiryCleanupDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "fivenet",
		Subsystem: "properties",
		Name:      "expiry_cleanup_duration_seconds",
		Help:      "Duration of dynamic-property expiry cleanup runs.",
		Buckets:   prometheus.DefBuckets,
	}, []string{"scope"})
)

// ObserveExpiryCleanup records aggregate, low-cardinality cleanup telemetry.
// Callers must use a stable scope name such as "user" or "vehicle"; never a
// property key or resource identifier.
func ObserveExpiryCleanup(scope string, duration time.Duration, cleared int, err error) {
	outcome := "success"
	if err != nil {
		outcome = "failure"
	}
	propertyExpiryCleanupRuns.WithLabelValues(scope, outcome).Inc()
	propertyExpiryCleanupDuration.WithLabelValues(scope).Observe(duration.Seconds())
	if cleared > 0 {
		propertyExpiryCleanupCleared.WithLabelValues(scope).Add(float64(cleared))
	}
}

// RecordExpiryCleanupItemFailure records a failed individual clear while the
// surrounding batch continues with other resources.
func RecordExpiryCleanupItemFailure(scope string) {
	propertyExpiryCleanupFailures.WithLabelValues(scope).Inc()
}
