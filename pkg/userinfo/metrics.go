package userinfo

import (
	"sync"

	"github.com/fivenet-app/fivenet/v2026/pkg/server/admin"
	"github.com/prometheus/client_golang/prometheus"
)

type changesMetrics struct {
	publishFailures prometheus.Counter
}

var (
	changesMetricsInstance *changesMetrics
	changesMetricsOnce     sync.Once
)

func getChangesMetrics() *changesMetrics {
	changesMetricsOnce.Do(func() {
		changesMetricsInstance = &changesMetrics{
			publishFailures: prometheus.NewCounter(prometheus.CounterOpts{
				Namespace: admin.MetricsNamespace,
				Subsystem: "userinfo",
				Name:      "change_publish_failures_total",
				Help:      "Failed JetStream publishes of canonical user-info change events.",
			}),
		}

		prometheus.MustRegister(changesMetricsInstance.publishFailures)
	})

	return changesMetricsInstance
}
