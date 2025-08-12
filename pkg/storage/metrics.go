package storage

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/events"
	"github.com/fivenet-app/fivenet/v2025/pkg/nats/leaderelection"
	"github.com/fivenet-app/fivenet/v2025/pkg/server/admin"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/instance"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var metricSpaceUsage = prometheus.NewGauge(prometheus.GaugeOpts{
	Namespace: admin.MetricsNamespace,
	Subsystem: "storage",
	Name:      "space_usage_total_bytes",
	Help:      "Total space used by files in the storage.",
})

var MetricsCollectorModule = fx.Module("storage.metrics_collector",
	fx.Provide(NewMetricsCollector),
)

type MetricsCollector struct {
	logger  *zap.Logger
	le      *leaderelection.LeaderElector
	storage IStorage
}

type MetricsCollectorParams struct {
	fx.In

	LC fx.Lifecycle

	Logger  *zap.Logger
	JS      *events.JSWrapper
	Cfg     *config.Config
	Storage IStorage
}

func NewMetricsCollector(p MetricsCollectorParams) *MetricsCollector {
	if !p.Cfg.Storage.MetricsEnabled {
		p.Logger.Info("Metrics collection is disabled in configuration")
		return nil
	}

	prometheus.DefaultRegisterer.MustRegister(metricSpaceUsage)

	ctxCancel, cancel := context.WithCancel(context.Background())

	mc := &MetricsCollector{
		logger:  p.Logger.Named("storage.metrics"),
		storage: p.Storage,
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		nodeName := instance.ID() + "_storage_metrics_collector"

		var err error
		mc.le, err = leaderelection.New(
			ctxCancel, mc.logger, p.JS,
			"leader_election",           // Bucket
			"storage_metrics_collector", // Key
			32*time.Second,              // TTL for the lock
			16*time.Second,              // Heartbeat interval
			func(ctx context.Context) {
				mc.logger.Info("housekeeper started", zap.String("node_name", nodeName))

				mc.start(ctx, p.Cfg.Storage.MetricsInterval)
			},
			nil, // No on stopped function
		)
		if err != nil {
			return fmt.Errorf("failed to create leader elector. %w", err)
		}
		mc.le.Start()

		return nil
	}))

	p.LC.Append(fx.StopHook(func(ctxShutdown context.Context) error {
		cancel()

		return nil
	}))

	return mc
}

func (mc *MetricsCollector) start(ctx context.Context, interval time.Duration) {
	mc.logger.Info("Starting metrics collector")

	// Wait for a random delay before collecting the first metrics
	//nolint:gosec // G404 - The random delay is not security sensitive, it's just to avoid all instances collecting metrics at the same time.
	delay := time.Duration(1+rand.Intn(15)) * time.Second
	time.Sleep(delay)

	for {
		// Collect space usage regularly
		if err := mc.CollectMetrics(ctx); err != nil {
			mc.logger.Error("Failed to collect metrics", zap.Error(err))
		}

		select {
		case <-ctx.Done():
			return

		case <-time.After(interval):
		}
	}
}

func (mc *MetricsCollector) CollectMetrics(ctx context.Context) error {
	mc.logger.Info("Starting storage metrics collection")

	usage, err := mc.storage.GetSpaceUsage(ctx)
	if err != nil {
		return fmt.Errorf("failed to get space usage: %w", err)
	}

	mc.logger.Info("Collected storage metrics", zap.Int64("space_usage", usage))

	metricSpaceUsage.Set(float64(usage))

	return nil
}
