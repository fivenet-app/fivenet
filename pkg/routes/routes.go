package routes

import (
	"net/http"

	"github.com/galexrt/fivenet/pkg/api"
	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/oauth2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Routes struct {
	logger *zap.Logger

	clientCfg *api.ClientConfig
}

func New(logger *zap.Logger) *Routes {
	providers := make([]*api.ProviderConfig, len(config.C.OAuth2.Providers))

	for k, p := range config.C.OAuth2.Providers {
		providers[k] = &api.ProviderConfig{
			Name:  p.Name,
			Label: p.Label,
		}
	}

	return &Routes{
		logger: logger,

		clientCfg: &api.ClientConfig{
			Version:   "TODO",
			SentryDSN: config.C.Sentry.ClientDSN,
			Login: api.LoginConfig{
				Providers: providers,
			},
		},
	}
}

func (r *Routes) Register(e *gin.Engine, oa2 *oauth2.OAuth2) {
	e.GET("/readiness", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
	// API Base
	g := e.Group("/api")
	{
		g.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, api.PingResponse)
		})
		g.POST("/config", func(c *gin.Context) {
			c.JSON(http.StatusOK, r.clientCfg)
		})
		g.GET("/clear-site-data", func(c *gin.Context) {
			c.Header("Clear-Site-Data", "\"cache\", \"cookies\", \"storage\"")
			c.String(http.StatusOK, "Your local site data should be cleared now, please go back to the FiveNet homepage yourself.")
		})
	}

	// OAuth2
	oauth := g.Group("/oauth2")
	{
		oauth.GET("/login/:provider", oa2.Login)
		oauth.POST("/login/:provider", oa2.Login)
		oauth.GET("/callback/:provider", oa2.Callback)
		oauth.POST("/callback/:provider", oa2.Callback)
	}
}
