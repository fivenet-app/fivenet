package filestore

import (
	"github.com/fivenet-app/fivenet/v2025/pkg/server/admin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var metricDeletedFiles = promauto.NewGauge(prometheus.GaugeOpts{
	Namespace: admin.MetricsNamespace,
	Subsystem: "filestore_housekeeper",
	Name:      "deleted_files_count",
	Help:      "Number of deleted files.",
})
