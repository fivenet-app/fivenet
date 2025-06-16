// Package filestore provides file storage utilities and related metrics.
package filestore

import (
	"github.com/fivenet-app/fivenet/v2025/pkg/server/admin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// metricDeletedFiles tracks the number of files deleted by the filestore housekeeper.
// This metric is registered as a Prometheus gauge and is labeled under the
// 'filestore_housekeeper' subsystem.
var metricDeletedFiles = promauto.NewGauge(prometheus.GaugeOpts{
	Namespace: admin.MetricsNamespace,
	Subsystem: "filestore_housekeeper",
	Name:      "deleted_files_count",
	Help:      "Number of deleted files.",
})
