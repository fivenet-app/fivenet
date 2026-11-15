package sync

import (
	"sync"

	"github.com/fivenet-app/fivenet/v2026/pkg/server/admin"
	"github.com/prometheus/client_golang/prometheus"
)

type syncMetrics struct {
	apiTokenAuthAttempts *prometheus.CounterVec
}

var (
	syncMetricsOnce sync.Once
	syncMetricsInst *syncMetrics
)

func getSyncMetrics() *syncMetrics {
	syncMetricsOnce.Do(func() {
		syncMetricsInst = newSyncMetrics(prometheus.DefaultRegisterer)
	})

	return syncMetricsInst
}

func newSyncMetrics(registerer prometheus.Registerer) *syncMetrics {
	metrics := &syncMetrics{
		apiTokenAuthAttempts: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: admin.MetricsNamespace,
			Subsystem: "sync",
			Name:      "api_token_auth_attempts_total",
			Help:      "Sync API token authentication attempts by outcome.",
		}, []string{"outcome"}),
	}

	if registerer != nil {
		registerer.MustRegister(metrics.apiTokenAuthAttempts)
	}

	return metrics
}
