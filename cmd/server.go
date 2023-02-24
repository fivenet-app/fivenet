package cmd

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	apiv1 "github.com/galexrt/arpanet/api/v1"
	"github.com/galexrt/arpanet/pkg/auth"
	"github.com/galexrt/arpanet/pkg/config"
	gormsessions "github.com/galexrt/arpanet/pkg/gormsessions"
	"github.com/galexrt/arpanet/query"
	"github.com/gin-contrib/sessions"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var serverCmd = &cobra.Command{
	Use: "server",
	Run: func(cmd *cobra.Command, args []string) {
		// Gin HTTP Server
		gin.SetMode(config.C.Mode)
		r := gin.New()
		// Add Zap Logger to Gin
		r.Use(ginzap.Ginzap(logger, time.RFC3339, true))
		r.Use(ginzap.RecoveryWithZap(logger, true))
		// Sessions
		store := gormsessions.NewStore(query.DB, true, []byte("secret"))
		store.Options(sessions.Options{
			Domain:   "172.16.1.111",
			Path:     "/",
			MaxAge:   int((10 * time.Minute).Seconds()),
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
		})
		r.Use(sessions.SessionsMany([]string{auth.SessionName}, store))

		// Prometheus Metrics endpoint
		r.GET("/metrics", gin.WrapH(promhttp.Handler()))
		// Register API routes
		apiv1.Register(r)

		// Register embed FS for assets and other static files
		if gin.Mode() == gin.DebugMode {
			r.StaticFS("/public", gin.Dir(".", false))
		} else {
			r.StaticFS("/public", http.FS(assets))
		}
		r.GET("favicon.ico", func(c *gin.Context) {
			file, _ := assets.ReadFile("assets/favicon.ico")
			c.Data(
				http.StatusOK,
				"image/x-icon",
				file,
			)
		})
		r.Static("/livemap", "./livemap/dist/")

		// Create HTTP Server for graceful shutdown handling
		srv := &http.Server{
			Addr:    config.C.HTTP.Listen,
			Handler: r,
		}
		// Initializing the server in a goroutine so that
		// it won't block the graceful shutdown handling below
		go func() {
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				logger.Error("listen error", zap.Error(err))
			}
		}()
		logger.Info("server listening", zap.String("address", config.C.HTTP.Listen))

		// Wait for interrupt signal to gracefully shutdown the server with
		// a timeout of 5 seconds.
		quit := make(chan os.Signal)
		// kill (no param) default send syscall.SIGTERM
		// kill -2 is syscall.SIGINT
		// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
		//lint:ignore SA1017 can be unbuffered because of signal channel usage
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		logger.Info("shutting down server...")

		// The context is used to inform the server it has 5 seconds to finish
		// the request it is currently handling
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logger.Error("server forced to shutdown", zap.Error(err))
		}

		logger.Info("server exiting")
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return query.SetupDB(logger)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
