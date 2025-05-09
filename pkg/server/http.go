package server

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	grpcws "github.com/fivenet-app/fivenet/v2025/pkg/grpc/grpcws"
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
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

var HTTPServerModule = fx.Module("httpserver",
	fx.Provide(
		New,
	),
	fx.Decorate(wrapLogger),
)

var allowedHeaders = []string{
	"Origin", "Content-Length", "Content-Type", "Cookie", "Keep-Alive",
	// For GRPC-Web User agent
	"U-A",
}

type HTTPServer *http.Server

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger  *zap.Logger
	Config  *config.Config
	Engine  *gin.Engine
	GRPCSrv *grpc.Server
}

type Result struct {
	fx.Out

	Server HTTPServer
}

func wrapLogger(log *zap.Logger) *zap.Logger {
	return log.Named("http")
}

type handler struct {
	gin  http.Handler
	grpc *grpc.Server
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.ProtoMajor == 2 && strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc") {
		h.grpc.ServeHTTP(w, r)
		return
	}

	h.gin.ServeHTTP(w, r)
}

// New builds an HTTP server that will begin serving requests
// when the Fx application starts.
func New(p Params) (Result, error) {
	// Create HTTP Server for graceful shutdown handling and h2c wrapped handler
	srv := &http.Server{
		Addr: p.Config.HTTP.Listen,
		Handler: h2c.NewHandler(
			&handler{
				gin:  p.Engine,
				grpc: p.GRPCSrv,
			},
			&http2.Server{},
		),
		ErrorLog: zap.NewStdLog(p.Logger),
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

func NewEngine(p EngineParams) (*gin.Engine, error) {
	// Gin HTTP Server
	gin.SetMode(p.Config.Mode)
	e := gin.New()

	// Enable forwarded by client ip headers when one or more trusted proxies specified
	e.ForwardedByClientIP = len(p.Config.HTTP.TrustedProxies) > 0
	if err := e.SetTrustedProxies(p.Config.HTTP.TrustedProxies); err != nil {
		return nil, fmt.Errorf("failed to set trusted proxies list. %w", err)
	}

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
		AllowHeaders:           allowedHeaders,
		ExposeHeaders:          []string{"Content-Length", "Content-Type", "Accept-Encoding"},
		AllowBrowserExtensions: true,
		AllowWebSockets:        true,
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
	e.Use(otelgin.Middleware("fivenet",
		otelgin.WithTracerProvider(p.TP),
		otelgin.WithGinFilter(func(c *gin.Context) bool {
			// Skip `/images/*`, image proxy and GRPC-Web + GRPC-websocket requests
			fullPath := c.FullPath()
			return !strings.HasPrefix(fullPath, "/images/") &&
				!strings.HasPrefix(fullPath, "/api/grpc") &&
				!strings.HasPrefix(fullPath, "/api/image_proxy") &&
				!strings.HasPrefix(fullPath, "/api/filestore")
		}),
	))
	e.Use(InjectToHeaders(p.TP))

	for _, service := range p.Services {
		if service == nil {
			continue
		}

		service.RegisterHTTP(e)
	}

	// Setup assets and other static files serving
	frontendFS := static.LocalFile(".output/public/", true)
	fileServer := http.FileServer(frontendFS)

	// GRPC-web and websocket handling
	wrapperGrpc := grpcws.WrapServer(p.GRPCSrv,
		grpcws.WithAllowedRequestHeaders(allowedHeaders),
	)
	e.GET("/api/grpcws", func(ctx *gin.Context) {
		resp, req := ctx.Writer, ctx.Request
		if grpcws.IsGrpcWebSocketChannelRequest(req) {
			wrapperGrpc.HandleGrpcWebsocketChannelRequest(resp, req)
		} else {
			ctx.AbortWithStatus(http.StatusBadRequest)
		}
	})
	e.POST("/api/grpc/*path", gin.WrapH(http.StripPrefix("/api/grpc", wrapperGrpc)))

	// Setup not found handler
	notFoundPage := []byte("404 page not found")
	notFoundPageFile, err := frontendFS.Open("404.html")
	if err != nil {
		p.Logger.Warn("failed to open 404.html file, falling back to 404 text page")
	} else {
		notFoundPage, err = io.ReadAll(notFoundPageFile)
		if err != nil {
			p.Logger.Error("failed to read 404.html file contents", zap.Error(err))
		}
	}

	// 404 handling
	e.NoRoute(func(c *gin.Context) {
		requestPath := c.Request.URL.Path

		// If the target is a directory (e.g., `/livemap`), load root `index.html``
		if strings.HasSuffix(requestPath, "/") || !strings.Contains(requestPath, ".") {
			c.Request.URL.Path = "/"
			fileServer.ServeHTTP(c.Writer, c.Request)
			c.Abort()
			return
		}

		// Check if index file exists
		if frontendFS.Exists("/", requestPath) {
			fileServer.ServeHTTP(c.Writer, c.Request)
			c.Abort()
			return
		}

		c.Data(http.StatusNotFound, "text/html; charset=utf-8", notFoundPage)
	})

	return e, nil
}
