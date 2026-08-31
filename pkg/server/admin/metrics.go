package admin

import (
	"context"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/fivenet-app/fivenet/v2026/pkg/config"
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

	MetricsJobNameLabel = "job_name"
)

// Module provides the metrics server as an fx module.
var Module = fx.Module("metricsserver",
	fx.Provide(
		NewReadiness,
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

// Readiness tracks whether the main HTTP server is ready to receive traffic.
type Readiness struct {
	ready atomic.Bool
}

// NewReadiness creates the shared application readiness state.
func NewReadiness() *Readiness {
	return &Readiness{}
}

// Ready reports whether the application is ready to receive traffic.
func (r *Readiness) Ready() bool {
	return r.ready.Load()
}

// SetReady updates the application readiness state.
func (r *Readiness) SetReady(ready bool) {
	r.ready.Store(ready)
}

// Params contains dependencies for constructing the metrics server.
type Params struct {
	fx.In

	// LC is the fx lifecycle for managing server start/stop hooks.
	LC fx.Lifecycle

	// Logger is the zap logger instance for logging.
	Logger *zap.Logger
	// Config is the application configuration.
	Config *config.Config
	// Readiness is the shared application readiness state.
	Readiness *Readiness
}

// Result is the output struct for the metrics server constructor.
type Result struct {
	fx.Out

	// Server is the constructed admin HTTP server.
	Server AdminServer
}

// NewServer creates and configures the metrics (admin) HTTP server with Prometheus metrics, readiness, and pprof endpoints.
func NewServer(p Params) (Result, error) {
	if p.Config.HTTP.AdminListen == "" {
		p.Logger.Info("admin server disabled (adminListen is empty)")
		return Result{}, nil
	}

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

	// Liveness probe endpoint
	e.GET("/liveness", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Readiness probe endpoint
	e.GET("/readiness", func(c *gin.Context) {
		if !p.Readiness.Ready() {
			c.String(http.StatusServiceUnavailable, "NOT READY")
			return
		}

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
	p.LC.Append(fx.StartHook(func(ctx context.Context) error {
		//nolint:noctx // net.Listen is shutdown via the server's Shutdown method
		ln, err := net.Listen("tcp", srv.Addr)
		if err != nil {
			return err
		}
		p.Logger.Info("metrics server listening", zap.String("address", srv.Addr))
		go srv.Serve(ln)

		return nil
	}))
	p.LC.Append(fx.StopHook(func(ctx context.Context) error {
		return srv.Shutdown(ctx)
	}))

	return Result{
		Server: srv,
	}, nil
}
