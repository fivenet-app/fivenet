package croner

import (
	"sync"

	"github.com/fivenet-app/fivenet/v2026/pkg/server/admin"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	schedulerMetricSubsystem = "cron_scheduler"
	executorMetricSubsystem  = "cron_executor"

	schedulerMetricHandoffLatency = "handoff_latency_seconds"

	executorMetricStartLatency    = "start_latency_seconds"
	executorMetricHandlerDuration = "handler_duration_seconds"
	executorMetricLastRunSuccess  = "last_run_success"
)

type schedulerMetrics struct {
	handoffLatency *prometheus.HistogramVec
}

type executorMetrics struct {
	startLatency    *prometheus.HistogramVec
	handlerDuration *prometheus.HistogramVec
	lastRunSuccess  *prometheus.GaugeVec
}

var (
	schedulerMetricsOnce sync.Once
	schedulerMetricsInst *schedulerMetrics

	executorMetricsOnce sync.Once
	executorMetricsInst *executorMetrics
)

func getSchedulerMetrics() *schedulerMetrics {
	schedulerMetricsOnce.Do(func() {
		schedulerMetricsInst = &schedulerMetrics{
			handoffLatency: prometheus.NewHistogramVec(prometheus.HistogramOpts{
				Namespace: admin.MetricsNamespace,
				Subsystem: schedulerMetricSubsystem,
				Name:      schedulerMetricHandoffLatency,
				Help:      "Latency of publishing a cron job to the scheduler stream.",
				Buckets:   prometheus.ExponentialBuckets(0.005, 2, 17),
			}, []string{admin.MetricsJobNameLabel}),
		}

		prometheus.MustRegister(schedulerMetricsInst.handoffLatency)
	})

	return schedulerMetricsInst
}

func getExecutorMetrics() *executorMetrics {
	executorMetricsOnce.Do(func() {
		executorMetricsInst = &executorMetrics{
			startLatency: prometheus.NewHistogramVec(prometheus.HistogramOpts{
				Namespace: admin.MetricsNamespace,
				Subsystem: executorMetricSubsystem,
				Name:      executorMetricStartLatency,
				Help:      "Latency between cron job start and handler execution start.",
				Buckets:   prometheus.ExponentialBuckets(0.005, 2, 17),
			}, []string{admin.MetricsJobNameLabel}),
			handlerDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
				Namespace: admin.MetricsNamespace,
				Subsystem: executorMetricSubsystem,
				Name:      executorMetricHandlerDuration,
				Help:      "Duration of cron handler execution.",
				Buckets:   prometheus.ExponentialBuckets(0.005, 2, 17),
			}, []string{admin.MetricsJobNameLabel}),
			lastRunSuccess: prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: admin.MetricsNamespace,
				Subsystem: executorMetricSubsystem,
				Name:      executorMetricLastRunSuccess,
				Help:      "Whether the last cron handler run succeeded (1) or failed (0).",
			}, []string{admin.MetricsJobNameLabel}),
		}

		prometheus.MustRegister(
			executorMetricsInst.startLatency,
			executorMetricsInst.handlerDuration,
			executorMetricsInst.lastRunSuccess,
		)
	})

	return executorMetricsInst
}
