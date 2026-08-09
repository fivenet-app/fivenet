package dbsyncconfig

import (
	"sync"
	"time"

	"github.com/fivenet-app/fivenet/v2026/pkg/server/admin"
	"github.com/prometheus/client_golang/prometheus"
)

type cursorMetrics struct {
	lastCheckUnix *prometheus.GaugeVec
}

var (
	cursorMetricsOnce sync.Once
	cursorMetricsInst *cursorMetrics
)

func getCursorMetrics() *cursorMetrics {
	cursorMetricsOnce.Do(func() {
		cursorMetricsInst = &cursorMetrics{
			lastCheckUnix: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: admin.MetricsNamespace,
				Subsystem: "dbsync",
				Name:      "cursor_last_check_unix",
				Help:      "UNIX timestamp currently stored in the dbsync cursor last_check state.",
			}, []string{"table"}),
		}

		prometheus.MustRegister(cursorMetricsInst.lastCheckUnix)
	})

	return cursorMetricsInst
}

func (m *cursorMetrics) setCursorMetrics(table string, lastCheck *time.Time) {
	if table == "" {
		return
	}

	if lastCheck != nil {
		m.lastCheckUnix.WithLabelValues(table).Set(float64(lastCheck.Unix()))
	} else {
		m.lastCheckUnix.WithLabelValues(table).Set(0)
	}
}
