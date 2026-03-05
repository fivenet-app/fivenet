package dbsync

import (
	"time"

	"github.com/fivenet-app/fivenet/v2026/pkg/server/admin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	metricLastRun = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: admin.MetricsNamespace,
		Subsystem: "dbsync",
		Name:      "last_run_unix",
		Help:      "UNIX timestamp of the last dbsync run per syncer and status.",
	}, []string{"syncer", "status"})

	metricRunDuration = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: admin.MetricsNamespace,
		Subsystem: "dbsync",
		Name:      "run_duration_seconds",
		Help:      "Duration of the latest dbsync run in seconds.",
	}, []string{"syncer"})

	metricSyncedRows = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: admin.MetricsNamespace,
		Subsystem: "dbsync",
		Name:      "synced_rows_count",
		Help:      "Number of rows synced during the latest dbsync run.",
	}, []string{"syncer"})

	metricFetchedItems = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: admin.MetricsNamespace,
		Subsystem: "dbsync",
		Name:      "fetched_items_count",
		Help:      "Number of items fetched during the latest dbsync run.",
	}, []string{"syncer"})

	metricSentItems = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: admin.MetricsNamespace,
		Subsystem: "dbsync",
		Name:      "sent_items_count",
		Help:      "Number of items sent during the latest dbsync run.",
	}, []string{"syncer"})
)

func recordSyncMetrics(syncer string, startedAt time.Time, fetched int64, sent int64, err error) {
	if syncer == "" {
		return
	}

	metricRunDuration.WithLabelValues(syncer).Set(time.Since(startedAt).Seconds())

	status := "success"
	if err != nil {
		status = "failed"
	}
	metricLastRun.WithLabelValues(syncer, status).SetToCurrentTime()

	if err == nil {
		metricFetchedItems.WithLabelValues(syncer).Set(float64(fetched))
		metricSentItems.WithLabelValues(syncer).Set(float64(sent))
		metricSyncedRows.WithLabelValues(syncer).Set(float64(sent))
	}
}
