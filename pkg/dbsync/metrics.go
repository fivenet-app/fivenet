package dbsync

import (
	"sync"
	"time"

	"github.com/fivenet-app/fivenet/v2026/pkg/server/admin"
	"github.com/prometheus/client_golang/prometheus"
)

type syncMetrics struct {
	lastRun      *prometheus.GaugeVec
	runDuration  *prometheus.GaugeVec
	syncedRows   *prometheus.GaugeVec
	fetchedItems *prometheus.GaugeVec
	sentItems    *prometheus.GaugeVec
}

var (
	syncMetricsOnce sync.Once
	syncMetricsInst *syncMetrics
)

func getSyncMetrics() *syncMetrics {
	syncMetricsOnce.Do(func() {
		syncMetricsInst = &syncMetrics{
			lastRun: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: admin.MetricsNamespace,
				Subsystem: "dbsync",
				Name:      "last_run_unix",
				Help:      "UNIX timestamp of the last dbsync run per syncer and status.",
			}, []string{"syncer", "status"}),
			runDuration: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: admin.MetricsNamespace,
				Subsystem: "dbsync",
				Name:      "run_duration_seconds",
				Help:      "Duration of the latest dbsync run in seconds.",
			}, []string{"syncer"}),
			syncedRows: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: admin.MetricsNamespace,
				Subsystem: "dbsync",
				Name:      "synced_rows_count",
				Help:      "Number of rows synced during the latest dbsync run.",
			}, []string{"syncer"}),
			fetchedItems: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: admin.MetricsNamespace,
				Subsystem: "dbsync",
				Name:      "fetched_items_count",
				Help:      "Number of items fetched during the latest dbsync run.",
			}, []string{"syncer"}),
			sentItems: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: admin.MetricsNamespace,
				Subsystem: "dbsync",
				Name:      "sent_items_count",
				Help:      "Number of items sent during the latest dbsync run.",
			}, []string{"syncer"}),
		}

		prometheus.MustRegister(
			syncMetricsInst.lastRun,
			syncMetricsInst.runDuration,
			syncMetricsInst.syncedRows,
			syncMetricsInst.fetchedItems,
			syncMetricsInst.sentItems,
		)
	})

	return syncMetricsInst
}

func (s *Sync) recordSyncMetrics(
	syncer string,
	startedAt time.Time,
	fetched int64,
	sent int64,
	err error,
) {
	if syncer == "" {
		return
	}

	s.metrics.runDuration.WithLabelValues(syncer).Set(time.Since(startedAt).Seconds())

	status := "success"
	if err != nil {
		status = "failed"
	}
	s.metrics.lastRun.WithLabelValues(syncer, status).SetToCurrentTime()

	if err == nil {
		s.metrics.fetchedItems.WithLabelValues(syncer).Set(float64(fetched))
		s.metrics.sentItems.WithLabelValues(syncer).Set(float64(sent))
		s.metrics.syncedRows.WithLabelValues(syncer).Set(float64(sent))
	}
}
