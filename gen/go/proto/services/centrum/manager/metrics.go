package manager

import (
	"github.com/galexrt/fivenet/pkg/server/admin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	metricsDispatchLastID = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: admin.MetricsNamespace,
		Subsystem: "centrum",
		Name:      "dispatch_last_id",
		Help:      "Last dispatch ID.",
	}, []string{"job"})
)

// TODO what metrics are feasible to collect?
