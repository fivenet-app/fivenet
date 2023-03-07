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
	"github.com/galexrt/arpanet/pkg/session"
	pbauth "github.com/galexrt/arpanet/proto/auth"
	pbdispatches "github.com/galexrt/arpanet/proto/dispatches"
	pbdocuments "github.com/galexrt/arpanet/proto/documents"
	pbjob "github.com/galexrt/arpanet/proto/job"
	pblivemap "github.com/galexrt/arpanet/proto/livemap"
	pbusers "github.com/galexrt/arpanet/proto/users"
	"github.com/galexrt/arpanet/query"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var serverCmd = &cobra.Command{
	Use: "server",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create JWT Token TokenManager
		session.Tokens = session.NewTokenManager()

		// Setup and register Permissions
		perms.Setup()
		perms.Register()

		// Gin HTTP Server
		gin.SetMode(config.C.Mode)
		e := gin.New()

		// Add Zap Logger to Gin
		e.Use(ginzap.Ginzap(logger, time.RFC3339, true))
		e.Use(ginzap.RecoveryWithZap(logger, true))

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
		grpcServer := grpc.NewServer(
			grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
				grpc_ctxtags.UnaryServerInterceptor(),
				grpc_prometheus.UnaryServerInterceptor,
				grpc_zap.UnaryServerInterceptor(logger),
				grpc_validator.UnaryServerInterceptor(),
				grpc_auth.UnaryServerInterceptor(auth.GRPCAuthFunc),
				grpc_recovery.UnaryServerInterceptor(),
			)),
			grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
				grpc_ctxtags.StreamServerInterceptor(),
				grpc_prometheus.StreamServerInterceptor,
				grpc_zap.StreamServerInterceptor(logger),
				grpc_validator.StreamServerInterceptor(),
				grpc_auth.StreamServerInterceptor(auth.GRPCAuthFunc),
				grpc_recovery.StreamServerInterceptor(),
			)),
		)
		// Only enable grpc server reflection when in debug mode
		if gin.Mode() == gin.DebugMode {
			reflection.Register(grpcServer)
		}

		// Attach our GRPC services
		pbauth.RegisterAccountServiceServer(grpcServer, pbauth.NewServer())
		pbdispatches.RegisterDispatchesServiceServer(grpcServer, pbdispatches.NewServer())
		pbdocuments.RegisterDocumentsServiceServer(grpcServer, pbdocuments.NewServer())
		pbjob.RegisterJobServiceServer(grpcServer, pbjob.NewServer())
		pblivemap.RegisterLivemapServiceServer(grpcServer, pblivemap.NewServer())
		pbusers.RegisterUsersServiceServer(grpcServer, pbusers.NewServer())

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
