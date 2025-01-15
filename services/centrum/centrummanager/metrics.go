package centrummanager

import (
	"github.com/fivenet-app/fivenet/pkg/server/admin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var metricDispatchLastID = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: admin.MetricsNamespace,
	Subsystem: "centrum",
	Name:      "dispatch_last_id",
	Help:      "Last dispatch ID.",
}, []string{"job_name"})
