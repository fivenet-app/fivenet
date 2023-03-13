package cmd

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/galexrt/arpanet/pkg/auth"
	"github.com/galexrt/arpanet/pkg/config"
	"github.com/galexrt/arpanet/pkg/perms"
	"github.com/galexrt/arpanet/pkg/routes"
	"github.com/galexrt/arpanet/query"
	"github.com/getsentry/sentry-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	// GRPC + GRPC Middlewares
	grpc_auth "github.com/galexrt/arpanet/pkg/grpc/auth"
	grpc_permission "github.com/galexrt/arpanet/pkg/grpc/permission"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	// GRPC Services
	pbauth "github.com/galexrt/arpanet/proto/services/auth"
	pbcitizenstore "github.com/galexrt/arpanet/proto/services/citizenstore"
	pbcompletor "github.com/galexrt/arpanet/proto/services/completor"
	pbdispatcher "github.com/galexrt/arpanet/proto/services/dispatcher"
	pbdocstore "github.com/galexrt/arpanet/proto/services/docstore"
	pbjobs "github.com/galexrt/arpanet/proto/services/jobs"
	pblivmapper "github.com/galexrt/arpanet/proto/services/livemapper"
)

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

		// Create JWT Token TokenManager
		auth.Tokens = auth.NewTokenManager()

		// Setup permissions cache system
		perms.Setup()
		perms.Register()

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

		// Create GRPC Server
		lis, err := net.Listen("tcp", config.C.GRPC.Listen)
		if err != nil {
			logger.Fatal("failed to listen", zap.Error(err))
		}
		sentryRecoverFunc := func(p interface{}) (err error) {
			if e, ok := p.(error); ok {
				sentry.CaptureException(e)
			}
			return status.Errorf(codes.Internal, "%v", p)
		}
		grpcServer := grpc.NewServer(
			grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
				grpc_ctxtags.UnaryServerInterceptor(),
				grpc_prometheus.UnaryServerInterceptor,
				grpc_zap.UnaryServerInterceptor(logger),
				grpc_validator.UnaryServerInterceptor(),
				grpc_auth.UnaryServerInterceptor(auth.GRPCAuthFunc),
				grpc_permission.UnaryServerInterceptor(auth.GRPCPermissionUnaryFunc),
				grpc_recovery.UnaryServerInterceptor(
					grpc_recovery.WithRecoveryHandler(sentryRecoverFunc),
				),
			)),
			grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
				grpc_ctxtags.StreamServerInterceptor(),
				grpc_prometheus.StreamServerInterceptor,
				grpc_zap.StreamServerInterceptor(logger),
				grpc_validator.StreamServerInterceptor(),
				grpc_auth.StreamServerInterceptor(auth.GRPCAuthFunc),
				grpc_permission.StreamServerInterceptor(auth.GRPCPermissionStreamFunc),
				grpc_recovery.StreamServerInterceptor(
					grpc_recovery.WithRecoveryHandler(sentryRecoverFunc),
				),
			)),
		)
		// Only enable grpc server reflection when in debug mode
		if gin.Mode() == gin.DebugMode {
			reflection.Register(grpcServer)
		}

		// Attach our GRPC services
		pbauth.RegisterAuthServiceServer(grpcServer, pbauth.NewServer())
		pbcitizenstore.RegisterCitizenStoreServiceServer(grpcServer, pbcitizenstore.NewServer())
		pbcompletor.RegisterCompletorServiceServer(grpcServer, pbcompletor.NewServer())
		pbdispatcher.RegisterDispatcherServiceServer(grpcServer, pbdispatcher.NewServer())
		pbdocstore.RegisterDocStoreServiceServer(grpcServer, pbdocstore.NewServer())
		pbjobs.RegisterJobsServiceServer(grpcServer, pbjobs.NewServer())
		pblivmapper.RegisterLivemapperServiceServer(grpcServer, pblivmapper.NewServer(logger.Named("grpc_livemap")))

		go func() {
			if err := grpcServer.Serve(lis); err != nil {
				logger.Error("failed to serve grpc server", zap.Error(err))
			} else {
				logger.Info("grpc server started successfully")
			}
		}()
		logger.Info("grpc server listening", zap.String("address", config.C.GRPC.Listen))

		// Create HTTP Server for graceful shutdown handling
		srv := &http.Server{
			Addr:    config.C.HTTP.Listen,
			Handler: e,
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

		perms.Stop()

		logger.Info("http server exiting")
		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return query.SetupDB(logger)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
