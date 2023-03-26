package routes

import (
	"net/http"

	"github.com/galexrt/arpanet/pkg/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Routes struct {
	logger *zap.Logger

	clientCfg *ClientConfig
}

func New(logger *zap.Logger) *Routes {
	return &Routes{
		logger: logger,

		clientCfg: &ClientConfig{
			SentryDSN:   config.C.Sentry.ClientDSN,
			APIProtoURL: config.C.GRPC.ClientURL,
		},
	}
}

func (r *Routes) Register(e *gin.Engine) {
	g := e.Group("/api")
	{
		g.GET("/ping", r.PingGET)
		g.POST("/config", r.ConfigPOST)
	}
}

func (r *Routes) PingGET(c *gin.Context) {
	c.JSON(http.StatusOK, "Pong")
}

type ClientConfig struct {
	SentryDSN   string `json:"sentryDSN"`
	APIProtoURL string `json:"apiProtoURL"`
}

func (r *Routes) ConfigPOST(c *gin.Context) {
	c.JSON(http.StatusOK, r.clientCfg)
}
