package dbsyncconfig

import (
	"time"

	"github.com/fivenet-app/fivenet/v2026/pkg/server/admin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var metricCursorLastCheckUnix = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: admin.MetricsNamespace,
	Subsystem: "dbsync",
	Name:      "cursor_last_check_unix",
	Help:      "UNIX timestamp currently stored in the dbsync cursor last_check state.",
}, []string{"table"})

func setCursorMetrics(table string, lastCheck *time.Time) {
	if table == "" {
		return
	}

	if lastCheck != nil {
		metricCursorLastCheckUnix.WithLabelValues(table).Set(float64(lastCheck.Unix()))
	} else {
		metricCursorLastCheckUnix.WithLabelValues(table).Set(0)
	}
}
