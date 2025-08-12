package admin

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/gin-contrib/pprof"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	MetricsNamespace = "fivenet"
)

// Module provides the metrics server as an fx module.
var Module = fx.Module("metricsserver",
	fx.Provide(
		NewServer,
	),
	fx.Decorate(wrapLogger),
)

// wrapLogger returns a logger named "server.metrics" for metrics server logging.
func wrapLogger(log *zap.Logger) *zap.Logger {
	return log.Named("server.metrics")
}

// AdminServer is a type alias for *http.Server, representing the admin HTTP server.
type AdminServer *http.Server

// Params contains dependencies for constructing the metrics server.
type Params struct {
	fx.In

	// LC is the fx lifecycle for managing server start/stop hooks.
	LC fx.Lifecycle

	// Logger is the zap logger instance for logging.
	Logger *zap.Logger
	// Config is the application configuration.
	Config *config.Config
}

// Result is the output struct for the metrics server constructor.
type Result struct {
	fx.Out

	// Server is the constructed admin HTTP server.
	Server AdminServer
}

// NewServer creates and configures the metrics (admin) HTTP server with Prometheus metrics, readiness, and pprof endpoints.
func NewServer(p Params) (Result, error) {
	// Gin HTTP Server
	gin.SetMode(gin.ReleaseMode)
	e := gin.New()

	// Add Zap Logger to Gin
	e.Use(ginzap.Ginzap(p.Logger, time.RFC3339, true))
	e.Use(ginzap.RecoveryWithZap(p.Logger, true))

	// Prometheus Metrics endpoint
	e.GET("/metrics", gin.WrapH(promhttp.InstrumentMetricHandler(
		prometheus.DefaultRegisterer,
		promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{
			// Opt into OpenMetrics e.g. to support exemplars
			EnableOpenMetrics: true,
		}),
	)))

	// Readiness probe endpoint
	e.GET("/readiness", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Register pprof endpoints for profiling
	pprof.Register(e)

	// Create HTTP Server for graceful shutdown handling
	srv := &http.Server{
		ReadHeaderTimeout: 5 * time.Second,
		Addr:              p.Config.HTTP.AdminListen,
		Handler:           e,
	}

	// Register lifecycle hooks for server start and stop
	p.LC.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			p.Logger.Info("metrics server listening", zap.String("address", srv.Addr))
			go srv.Serve(ln)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})

	return Result{
		Server: srv,
	}, nil
}
