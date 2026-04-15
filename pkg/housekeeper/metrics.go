// Package housekeeper provides metrics for tracking housekeeping job results.
package housekeeper

import (
	"github.com/fivenet-app/fivenet/v2026/pkg/server/admin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// metricSoftDeleteRowsAffected tracks the number of rows affected by the soft delete operation.
	metricSoftDeleteRowsAffected = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: admin.MetricsNamespace,
		Subsystem: "housekeeper",
		Name:      "soft_delete_rows_affected_count",
		Help:      "Number of rows affected by the soft delete operation.",
	})

	// metricHardDeleteRowsAffected tracks the number of rows affected by the hard delete operation, labeled by table name.
	metricHardDeleteRowsAffected = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: admin.MetricsNamespace,
		Subsystem: "housekeeper",
		Name:      "hard_delete_rows_affected_count",
		Help:      "Number of rows affected by the hard delete operation.",
	}, []string{"table_name"})
)
