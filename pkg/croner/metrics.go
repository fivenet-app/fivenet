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
	executorMetricRunsTotal       = "runs_total"
)

type schedulerMetrics struct {
	handoffLatency *prometheus.HistogramVec
}

type executorMetrics struct {
	startLatency    *prometheus.HistogramVec
	handlerDuration *prometheus.HistogramVec
	runsTotal       *prometheus.CounterVec
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
			}, []string{"job"}),
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
			}, []string{"job"}),
			handlerDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
				Namespace: admin.MetricsNamespace,
				Subsystem: executorMetricSubsystem,
				Name:      executorMetricHandlerDuration,
				Help:      "Duration of cron handler execution.",
				Buckets:   prometheus.ExponentialBuckets(0.005, 2, 17),
			}, []string{"job"}),
			runsTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
				Namespace: admin.MetricsNamespace,
				Subsystem: executorMetricSubsystem,
				Name:      executorMetricRunsTotal,
				Help:      "Total number of cron handler runs by outcome.",
			}, []string{"job", "status"}),
		}

		prometheus.MustRegister(
			executorMetricsInst.startLatency,
			executorMetricsInst.handlerDuration,
			executorMetricsInst.runsTotal,
		)
	})

	return executorMetricsInst
}
