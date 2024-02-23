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
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/static"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var HTTPServerModule = fx.Module("httpserver",
	fx.Provide(
		New,
	),
	fx.Decorate(wrapLogger),
)

type HTTPServer *http.Server

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger *zap.Logger
	Config *config.Config
	Engine *gin.Engine
}

type Result struct {
	fx.Out

	Server HTTPServer
}

func wrapLogger(log *zap.Logger) *zap.Logger {
	return log.Named("http_server")
}

// New builds an HTTP server that will begin serving requests
// when the Fx application starts.
func New(p Params) (Result, error) {
	// Create HTTP Server for graceful shutdown handling
	srv := &http.Server{
		Addr:    p.Config.HTTP.Listen,
		Handler: p.Engine,
	}

	p.LC.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			p.Logger.Info("http server listening", zap.String("address", srv.Addr))
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

var HTTPEngineModule = fx.Module("httpengine",
	fx.Provide(
		NewEngine,
	),
	fx.Decorate(wrapLogger),
)

type EngineParams struct {
	fx.In

	Logger   *zap.Logger
	Config   *config.Config
	DB       *sql.DB
	TP       *tracesdk.TracerProvider
	TokenMgr *auth.TokenMgr

	Services []Service `group:"httpservices"`
}

func NewEngine(p EngineParams) *gin.Engine {
	// Gin HTTP Server
	gin.SetMode(p.Config.Mode)
	e := gin.New()

	// Add Zap Logger to Gin
	e.Use(ginzap.GinzapWithConfig(p.Logger, &ginzap.Config{
		UTC:        true,
		TimeFormat: time.RFC3339,
		Context: ginzap.Fn(func(c *gin.Context) []zapcore.Field {
			fields := []zapcore.Field{}
			// log request ID
			if requestID := c.Writer.Header().Get("X-Request-Id"); requestID != "" {
				fields = append(fields, zap.String("request_id", requestID))
			}

			// log trace and span ID
			if trace.SpanFromContext(c.Request.Context()).SpanContext().IsValid() {
				fields = append(fields, zap.String("traceId", trace.SpanFromContext(c.Request.Context()).SpanContext().TraceID().String()))
				fields = append(fields, zap.String("spanId", trace.SpanFromContext(c.Request.Context()).SpanContext().SpanID().String()))
			}

			return fields
		}),
	}))
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

	// Tracing
	e.Use(otelgin.Middleware("fivenet", otelgin.WithTracerProvider(p.TP)))
	e.Use(InjectToHeaders(p.TP))

	for _, service := range p.Services {
		if service == nil {
			continue
		}

		service.RegisterHTTP(e)
	}

	// Setup nuxt generated files serving
	fs := static.LocalFile(".output/public/", false)
	fileServer := http.StripPrefix("/", http.FileServer(fs))

	e.NoRoute(func(c *gin.Context) {
		requestPath := c.Request.URL.Path
		if strings.HasPrefix(requestPath, "/api") || requestPath == "/" {
			return
		}

		if strings.HasSuffix(requestPath, "/") || !strings.Contains(requestPath, ".") {
			c.Request.URL.Path = "/"
			fileServer.ServeHTTP(c.Writer, c.Request)
			c.Abort()
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	})

	// Register output dir for assets and other static files
	e.Use(static.Serve("/", fs))

	return e
}
