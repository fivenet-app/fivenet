package cmd

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dlmiddlecote/sqlstats"
	"github.com/galexrt/arpanet/pkg/auth"
	"github.com/galexrt/arpanet/pkg/config"
	"github.com/galexrt/arpanet/pkg/perms"
	"github.com/galexrt/arpanet/pkg/routes"
	"github.com/galexrt/arpanet/proto"
	"github.com/galexrt/arpanet/query"
	"github.com/getsentry/sentry-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var db *sql.DB

var serverCmd = &cobra.Command{
	Use: "server",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Setup Sentry Integration
		if config.C.Sentry.DSN != "" {
			err := sentry.Init(sentry.ClientOptions{
				Dsn:         config.C.Sentry.DSN,
				Debug:       false,
				Environment: config.C.Sentry.Environment,
				Release:     version.Info(),
			})
			if err != nil {
				logger.Fatal("Sentry init failed", zap.Error(err))
			}
			defer sentry.Flush(5 * time.Second)
		}

		// Create sql collector for Prometheus metrics
		collector := sqlstats.NewStatsCollector(config.C.Database.DBName, db)
		prometheus.MustRegister(collector)

		// Setup permissions system
		p := perms.New(db)
		defer p.Stop()
		p.Register()

		// Create JWT Token TokenManager
		tm := auth.NewTokenManager(config.C.JWT.Secret)

		grpcServer, grpcLis := proto.NewGRPCServer(logger, db, tm, p)

		go func() {
			if err := grpcServer.Serve(grpcLis); err != nil {
				logger.Error("failed to serve grpc server", zap.Error(err))
			} else {
				logger.Info("grpc server started successfully")
			}
		}()
		logger.Info("grpc server listening", zap.String("address", config.C.GRPC.Listen))

		// Create HTTP Server for graceful shutdown handling
		srv := &http.Server{
			Addr:    config.C.HTTP.Listen,
			Handler: setupHTTPServer(),
		}
		// Initializing the server in a goroutine so that
		// it won't block the graceful shutdown handling below
		go func() {
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				logger.Error("http listen error", zap.Error(err))
			}
		}()
		logger.Info("http server listening", zap.String("address", config.C.HTTP.Listen))

		// Wait for interrupt signal to gracefully shutdown the server with
		// a timeout of 5 seconds.
		quit := make(chan os.Signal)
		// kill (no param) default send syscall.SIGTERM
		// kill -2 is syscall.SIGINT
		// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
		//lint:ignore SA1017 can be unbuffered because of signal channel usage
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		logger.Info("shutting down servers...")

		// The context is used to inform the servers, they have 5 seconds to finish
		// the requests they are currently handling
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		go func() {
			grpcServer.GracefulStop()
		}()

		if err := srv.Shutdown(ctx); err != nil {
			logger.Error("http server forced to shutdown", zap.Error(err))
		}

		grpcServer.Stop()

		logger.Info("http server exiting")
		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		db, err = query.SetupDB(logger)
		return err
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

func setupHTTPServer() *gin.Engine {
	// Gin HTTP Server
	gin.SetMode(config.C.Mode)
	e := gin.New()

	// Add Zap Logger to Gin
	e.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	e.Use(ginzap.RecoveryWithZap(logger, true))

	// Sessions
	sessStore := cookie.NewStore([]byte(config.C.HTTP.Sessions.CookieSecret))
	sessStore.Options(sessions.Options{
		Domain:   "localhost",
		Path:     "/",
		MaxAge:   int((10 * time.Hour).Seconds()),
		HttpOnly: true,
		Secure:   false,
	})
	e.Use(sessions.SessionsMany([]string{"arpanet_"}, sessStore))

	// Prometheus Metrics endpoint
	e.GET("/metrics", gin.WrapH(promhttp.Handler()))
	// Register app routes
	rs := routes.New(logger)
	rs.Register(e)
	// Register embed FS for assets and other static files
	if gin.Mode() == gin.DebugMode {
		e.StaticFS("/public", gin.Dir(".", false))
	} else {
		e.StaticFS("/public", http.FS(assets))
	}
	e.GET("favicon.ico", func(c *gin.Context) {
		file, _ := assets.ReadFile("assets/favicon.ico")
		c.Data(
			http.StatusOK,
			"image/x-icon",
			file,
		)
	})

	return e
}
