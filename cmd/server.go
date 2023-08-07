package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/events"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/notifi"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/pkg/server/audit"
	"github.com/galexrt/fivenet/query"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/common/version"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var db *sql.DB

var serverCmd = &cobra.Command{
	Use: "server",
	RunE: func(cmd *cobra.Command, args []string) error {
		start := time.Now()
		// Setup Sentry Integration
		if config.C.Sentry.ServerDSN != "" && config.C.Mode != gin.DebugMode {
			err := sentry.Init(sentry.ClientOptions{
				Dsn:         config.C.Sentry.ServerDSN,
				Debug:       false,
				Environment: config.C.Sentry.Environment,
				Release:     version.Info(),
			})
			if err != nil {
				logger.Fatal("Sentry init failed", zap.Error(err))
			}
			defer sentry.Flush(5 * time.Second)
		}

		// Central context for cancelling any "background" services
		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()

		// Setup SQL Prometheus metrics collector
		prometheus.MustRegister(collectors.NewDBStatsCollector(db, config.C.Database.DBName))

		// Set up OTLP tracing
		tp, err := tracerProvider()
		if err != nil {
			logger.Fatal("failed to setup tracing provider", zap.Error(err))
		}
		defer func(ctx context.Context) {
			<-ctx.Done()
			ctx, cancel := context.WithTimeout(ctx, time.Second*5)
			defer cancel()
			if err := tp.Shutdown(ctx); err != nil {
				logger.Error("failed to cleanly shut down tracing", zap.Error(err))
			}
		}(ctx)

		// Create JWT Token TokenMgr
		tm := auth.NewTokenMgr(config.C.JWT.Secret)

		// Setup Event bus
		eventus, err := events.New(logger.Named("eventus"), config.C.NATS.URL)
		if err != nil {
			logger.Fatal("failed to setup event bus", zap.Error(err))
		}
		defer eventus.Stop()

		// Setup permissions system
		p, err := perms.New(ctx, logger.Named("perms"), db, tp, eventus)
		if err != nil {
			logger.Fatal("failed to setup permission system", zap.Error(err))
		}
		defer p.Stop()

		cfgDefaultPerms := config.C.Game.DefaultPermissions
		defaultPerms := make([]string, len(config.C.Game.DefaultPermissions))
		for i := 0; i < len(config.C.Game.DefaultPermissions); i++ {
			defaultPerms[i] = perms.BuildGuard(perms.Category(cfgDefaultPerms[i].Category), perms.Name(cfgDefaultPerms[i].Name))
		}

		if err := p.Register(defaultPerms); err != nil {
			return fmt.Errorf("failed to register permissions. %w", err)
		}

		// Audit Storer
		aud := audit.New(logger.Named("audit"), tp, db)
		aud.Start()

		// Notifier
		notif := notifi.New(logger.Named("notifi"), db, ctx, eventus, config.C.NATS.WorkerCount)

		// Wrap the server parts to try to isolate the actual "run servers" logic
		server := &server{
			tp:     tp,
			db:     db,
			events: eventus,
			p:      p,
			tm:     tm,
			audit:  aud,
			notif:  notif,
		}

		logger.Info("server start preparations took", zap.Duration("duration", time.Since(start)))
		return server.runServers(ctx)
	},
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		start := time.Now()
		db, err = query.SetupDB(logger)
		logger.Info("database setup took", zap.Duration("duration", time.Since(start)))
		return err
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
