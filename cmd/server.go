package cmd

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/galexrt/fivenet/gen/go/proto"
	"github.com/galexrt/fivenet/pkg/audit"
	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/oauth2"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/pkg/routes"
	"github.com/galexrt/fivenet/query"
	"github.com/getsentry/sentry-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/static"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
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

		// Create JWT Token TokenMgr
		tm := auth.NewTokenMgr(config.C.JWT.Secret)

		// Setup permissions system
		p := perms.New(ctx, db)

		cfgDefaultPerms := config.C.Game.DefaultPermissions
		defaultPerms := make([]string, len(config.C.Game.DefaultPermissions))
		for i := 0; i < len(config.C.Game.DefaultPermissions); i++ {
			defaultPerms[i] = perms.BuildGuard(perms.Category(cfgDefaultPerms[i].Category), perms.Name(cfgDefaultPerms[i].Name))
		}

		if err := p.Register(defaultPerms); err != nil {
			return fmt.Errorf("failed to register permissions. %w", err)
		}

		// Audit Storer
		aud := audit.New(logger.Named("audit"), db)
		aud.Start()

		// Wrap the server parts to try to isolate the actual "run servers" logic
		server := &server{
			db:    db,
			tm:    tm,
			p:     p,
			audit: aud,
		}

		return server.runServers(ctx)
	},
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		db, err = query.SetupDB(logger)
		return err
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

type server struct {
	db    *sql.DB
	tm    *auth.TokenMgr
	p     perms.Permissions
	audit *audit.AuditStorer
}

func (s *server) runServers(bctx context.Context) error {
	grpcServer, grpcLis := proto.NewGRPCServer(bctx, logger.Named("grpc/server"), s.db, s.tm, s.p, s.audit)

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
		Handler: s.setupHTTPServer(),
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

	s.audit.Stop()

	logger.Info("http server exiting")
	return nil
}

func (s *server) setupHTTPServer() *gin.Engine {
	// Gin HTTP Server
	gin.SetMode(config.C.Mode)
	e := gin.New()

	// Add Zap Logger to Gin
	e.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	e.Use(ginzap.RecoveryWithZap(logger, true))

	// Sessions
	sessStore := cookie.NewStore([]byte(config.C.HTTP.Sessions.CookieSecret))
	sessStore.Options(sessions.Options{
		Domain:   config.C.HTTP.Sessions.Domain,
		Path:     "/",
		MaxAge:   int((24 * time.Hour).Seconds()),
		HttpOnly: true,
		Secure:   true,
	})
	e.Use(sessions.SessionsMany([]string{"fivenet_oauth2_state", "fivenet_token"}, sessStore))

	// GZIP
	e.Use(gzip.Gzip(gzip.DefaultCompression))

	// Prometheus Metrics endpoint
	e.GET("/metrics", gin.WrapH(promhttp.Handler()))

	oauth := oauth2.New(logger.Named("oauth"), s.db, s.tm, config.C.OAuth2.Providers)

	// Register app routes
	rs := routes.New(logger)
	rs.Register(e, oauth)

	fs := static.LocalFile(".output/public/", false)
	fileserver := http.FileServer(fs)
	fileserver = http.StripPrefix("/", fileserver)

	e.NoRoute(func(c *gin.Context) {
		requestPath := c.Request.URL.Path
		if strings.HasPrefix(requestPath, "/api") || requestPath == "/" {
			return
		}

		pathComponents := strings.Split(requestPath[1:], "/")
		newPath := "/"
		if requestPath != "/" && len(pathComponents) > 0 {
			newPath += path.Join(pathComponents[:len(pathComponents)-1]...) + "/"
		}

		if fs.Exists("/", newPath) {
			c.Request.URL.Path = newPath
			fileserver.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	})
	// Register output dir for assets and other static files
	e.Use(static.Serve("/", fs))

	return e
}
