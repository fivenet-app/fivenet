package server

import (
	"context"
	"database/sql"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	grpcws "github.com/fivenet-app/fivenet/pkg/grpc/grpcws"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/static"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
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

	GRPCSrv *grpc.Server
}

func NewEngine(p EngineParams) *gin.Engine {
	// Gin HTTP Server
	gin.SetMode(p.Config.Mode)
	e := gin.New()

	// Enable forwarded by client ip headers when one or more trusted proxies specified
	e.ForwardedByClientIP = len(p.Config.HTTP.TrustedProxies) > 0
	e.SetTrustedProxies(p.Config.HTTP.TrustedProxies)

	// Add Zap logger and panic recovery to Gin
	e.Use(ginzap.GinzapWithConfig(p.Logger, &ginzap.Config{
		UTC:        true,
		TimeFormat: time.RFC3339,
		Context: ginzap.Fn(func(c *gin.Context) []zapcore.Field {
			fields := []zapcore.Field{}
			// Log request ID
			if requestID := c.Writer.Header().Get("X-Request-Id"); requestID != "" {
				fields = append(fields, zap.String("request_id", requestID))
			}

			// Log trace and span ID
			if trace.SpanFromContext(c.Request.Context()).SpanContext().IsValid() {
				fields = append(fields, zap.String("traceId", trace.SpanFromContext(c.Request.Context()).SpanContext().TraceID().String()))
				fields = append(fields, zap.String("spanId", trace.SpanFromContext(c.Request.Context()).SpanContext().SpanID().String()))
			}

			return fields
		}),
	}))
	e.Use(ginzap.RecoveryWithZap(p.Logger, true))

	// CORS
	e.Use(cors.New(cors.Config{
		AllowOrigins:           p.Config.HTTP.Origins,
		AllowMethods:           []string{"GET", "POST", "HEAD", "OPTIONS"},
		AllowHeaders:           []string{"Origin", "Content-Length", "Content-Type", "Cookie"},
		AllowBrowserExtensions: true,
		ExposeHeaders:          []string{"Content-Length", "Content-Type", "Accept-Encoding"},
		AllowCredentials:       true,
		MaxAge:                 1 * time.Hour,
	}))

	// Sessions
	sessStore := cookie.NewStore([]byte(p.Config.HTTP.Sessions.CookieSecret))
	sessStore.Options(sessions.Options{
		Domain:   p.Config.HTTP.Sessions.Domain,
		Path:     "/",
		MaxAge:   int((6 * time.Hour).Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})
	e.Use(sessions.SessionsMany([]string{"fivenet_oauth2_state"}, sessStore))

	// Tracing
	e.Use(otelgin.Middleware("fivenet", otelgin.WithTracerProvider(p.TP), otelgin.WithGinFilter(func(c *gin.Context) bool {
		// Skip `/images/*` requests
		return !strings.HasPrefix(c.FullPath(), "/images/")
	})))
	e.Use(InjectToHeaders(p.TP))

	for _, service := range p.Services {
		if service == nil {
			continue
		}

		service.RegisterHTTP(e)
	}

	// Setup Nuxt generated files serving
	frontendFS := static.LocalFile(".output/public/", true)
	fileServer := http.FileServer(frontendFS)
	// Register output dir for assets and other static files
	e.Use(static.Serve("/", frontendFS))

	// GRPC web + websocket handling
	wrapperGrpc := grpcws.WrapServer(
		p.GRPCSrv,
		grpcws.WithAllowedRequestHeaders([]string{"Origin", "Content-Length", "Content-Type", "Cookie", "Keep-Alive"}), // Allow cookie header
		grpcws.WithWebsocketChannelMaxStreamCount(10000),
		grpcws.WithWebsocketOriginFunc(func(req *http.Request) bool {
			return true
		}),
		grpcws.WithCorsForRegisteredEndpointsOnly(false),
		grpcws.WithAllowNonRootResource(true),
		grpcws.WithWebsocketPingInterval(40*time.Second),
	)
	ginWrappedGrpc := gin.WrapH(wrapperGrpc)
	e.Any("/api/grpc", ginWrappedGrpc)
	e.Any("/api/grpc/*path", ginWrappedGrpc)

	// Setup not found handler
	notFoundPage := []byte("404 page not found")
	notFoundPageFile, err := frontendFS.Open("404.html")
	if err != nil {
		p.Logger.Error("failed to open 404.html file, falling back to 404 text page", zap.Error(err))
	} else {
		notFoundPage, err = io.ReadAll(notFoundPageFile)
		if err != nil {
			p.Logger.Error("failed to read 404.html file contents", zap.Error(err))
		}
	}

	// 404 handling
	e.NoRoute(func(c *gin.Context) {
		requestPath := c.Request.URL.Path
		if requestPath == "/" {
			return
		}

		// If the target is a directory (e.g., `/livemap`), load root `index.html``
		if strings.HasSuffix(requestPath, "/") || !strings.Contains(requestPath, ".") {
			c.Request.URL.Path = "/"
			fileServer.ServeHTTP(c.Writer, c.Request)
			c.Abort()
			return
		}

		// Check if file exists
		if frontendFS.Exists("/", c.Request.URL.Path) {
			fileServer.ServeHTTP(c.Writer, c.Request)
		} else {
			c.Data(http.StatusNotFound, "text/html; charset=utf-8", notFoundPage)
		}
		c.Abort()
	})

	return e
}
