package server

import (
	"context"
	"database/sql"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/server/api"
	"github.com/galexrt/fivenet/pkg/server/oauth2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/static"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var HTTPServerModule = fx.Module("httpserver",
	fx.Provide(
		NewHTTP,
	),
	fx.Decorate(wrapLogger),
)

type ServerParams struct {
	fx.In

	LC fx.Lifecycle

	Logger   *zap.Logger
	Config   *config.Config
	DB       *sql.DB
	TP       *tracesdk.TracerProvider
	TokenMgr *auth.TokenMgr
}

type ServerResult struct {
	fx.Out

	Server *http.Server
}

// NewHTTP builds an HTTP server that will begin serving requests
// when the Fx application starts.
func NewHTTP(p ServerParams) (ServerResult, error) {
	// Create HTTP Server for graceful shutdown handling
	srv := &http.Server{
		Addr:    p.Config.HTTP.Listen,
		Handler: setupHTTPServer(p),
	}

	p.LC.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			p.Logger.Info("http server listening", zap.String("address", p.Config.HTTP.Listen))
			go srv.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})

	return ServerResult{
		Server: srv,
	}, nil
}

func setupHTTPServer(p ServerParams) *gin.Engine {
	// Gin HTTP Server
	gin.SetMode(p.Config.Mode)
	e := gin.New()

	// Add Zap Logger to Gin
	e.Use(ginzap.Ginzap(p.Logger, time.RFC3339, true))
	e.Use(ginzap.RecoveryWithZap(p.Logger, true))

	// Sessions
	sessStore := cookie.NewStore([]byte(p.Config.HTTP.Sessions.CookieSecret))
	sessStore.Options(sessions.Options{
		Domain:   p.Config.HTTP.Sessions.Domain,
		Path:     "/",
		MaxAge:   int((24 * time.Hour).Seconds()),
		HttpOnly: true,
		Secure:   true,
	})
	e.Use(sessions.SessionsMany([]string{"fivenet_oauth2_state"}, sessStore))

	// GZIP
	e.Use(gzip.Gzip(gzip.DefaultCompression))

	// Prometheus Metrics endpoint
	e.GET("/metrics", gin.WrapH(promhttp.InstrumentMetricHandler(
		prometheus.DefaultRegisterer, promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{
			// Opt into OpenMetrics e.g. to support exemplars
			EnableOpenMetrics: true,
		}),
	)))

	// Tracing
	e.Use(otelgin.Middleware("fivenet", otelgin.WithTracerProvider(p.TP)))

	// Register HTTP API routes
	rs := api.New(p.Logger, p.Config)
	rs.Register(e)

	if len(p.Config.OAuth2.Providers) > 0 {
		oauth := oauth2.New(p.Logger.Named("oauth"), p.DB, p.TokenMgr, p.Config.OAuth2.Providers)
		oauth.Register(e)
	}

	// Setup nuxt generated files serving
	fs := static.LocalFile(".output/public/", false)
	fileserver := http.FileServer(fs)
	fileserver = http.StripPrefix("/", fileserver)

	e.NoRoute(func(c *gin.Context) {
		requestPath := c.Request.URL.Path
		if strings.HasPrefix(requestPath, "/api") || requestPath == "/" {
			return
		}

		if strings.HasSuffix(requestPath, "/") || !strings.Contains(requestPath, ".") {
			c.Request.URL.Path = "/"
			fileserver.ServeHTTP(c.Writer, c.Request)
			c.Abort()
			return
		}

		fileserver.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	})
	// Register output dir for assets and other static files
	e.Use(static.Serve("/", fs))

	return e
}

func wrapLogger(log *zap.Logger) *zap.Logger {
	return log.Named("http_server")
}
