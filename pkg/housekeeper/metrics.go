// Package housekeeper provides metrics for tracking housekeeping job results.
package housekeeper

import (
	"github.com/fivenet-app/fivenet/v2025/pkg/server/admin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// metricSoftDeleteAffectedRows tracks the number of rows affected by the soft delete operation.
	metricSoftDeleteAffectedRows = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: admin.MetricsNamespace,
		Subsystem: "housekeeper",
		Name:      "soft_delete_affected_rows_count",
		Help:      "Number of rows affected by the soft delete operation.",
	})

	// metricHardDeleteAffectedRows tracks the number of rows affected by the hard delete operation, labeled by table name.
	metricHardDeleteAffectedRows = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: admin.MetricsNamespace,
		Subsystem: "housekeeper",
		Name:      "hard_delete_affected_rows_count",
		Help:      "Number of rows affected by the hard delete operation.",
	}, []string{"table_name"})
)
