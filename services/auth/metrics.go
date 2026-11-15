package auth

import (
	"sync"

	"github.com/fivenet-app/fivenet/v2026/pkg/server/admin"
	"github.com/prometheus/client_golang/prometheus"
)

type authMetrics struct {
	loginAttempts *prometheus.CounterVec
}

var (
	authMetricsOnce sync.Once
	authMetricsInst *authMetrics
)

func getAuthMetrics() *authMetrics {
	authMetricsOnce.Do(func() {
		authMetricsInst = newAuthMetrics(prometheus.DefaultRegisterer)
	})

	return authMetricsInst
}

func newAuthMetrics(registerer prometheus.Registerer) *authMetrics {
	metrics := &authMetrics{
		loginAttempts: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: admin.MetricsNamespace,
			Subsystem: "auth",
			Name:      "login_attempts_total",
			Help:      "Login attempts by outcome and failure reason.",
		}, []string{"outcome", "reason"}),
	}

	if registerer != nil {
		registerer.MustRegister(metrics.loginAttempts)
	}

	return metrics
}
