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

var Module = fx.Module("metricsserver",
	fx.Provide(
		NewServer,
	),
	fx.Decorate(wrapLogger),
)

func wrapLogger(log *zap.Logger) *zap.Logger {
	return log.Named("metrics_server")
}

type AdminServer *http.Server

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger *zap.Logger
	Config *config.Config
}

type Result struct {
	fx.Out

	Server AdminServer
}

func NewServer(p Params) (Result, error) {
	// Gin HTTP Server
	gin.SetMode(gin.ReleaseMode)
	e := gin.New()

	// Add Zap Logger to Gin
	e.Use(ginzap.Ginzap(p.Logger, time.RFC3339, true))
	e.Use(ginzap.RecoveryWithZap(p.Logger, true))

	// Prometheus Metrics endpoint
	e.GET("/metrics", gin.WrapH(promhttp.InstrumentMetricHandler(
		prometheus.DefaultRegisterer, promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{
			// Opt into OpenMetrics e.g. to support exemplars
			EnableOpenMetrics: true,
		}),
	)))

	e.GET("/readiness", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	pprof.Register(e)

	// Create HTTP Server for graceful shutdown handling
	srv := &http.Server{
		Addr:    p.Config.HTTP.AdminListen,
		Handler: e,
	}

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
