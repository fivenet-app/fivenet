package api

import (
	"net/http"
	"runtime/debug"

	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/server/oauth2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Routes struct {
	logger *zap.Logger

	clientCfg *ClientConfig
}

func New(logger *zap.Logger, cfg *config.Config) *Routes {
	providers := make([]*ProviderConfig, len(cfg.OAuth2.Providers))

	for k, p := range cfg.OAuth2.Providers {
		providers[k] = &ProviderConfig{
			Name:  p.Name,
			Label: p.Label,
		}
	}

	var version string
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		version = "UNKNOWN"
	} else {
		version = buildInfo.Main.Version
	}

	clientCfg := &ClientConfig{
		Version:   version,
		SentryDSN: cfg.Sentry.ClientDSN,
		Login: LoginConfig{
			SignupEnabled: cfg.Game.SignupEnabled,
			Providers:     providers,
		},
		Discord: &Discord{},
	}

	if cfg.Discord.Enabled {
		clientCfg.Discord.BotInviteURL = cfg.Discord.Bot.InviteURL
	}

	return &Routes{
		logger:    logger,
		clientCfg: clientCfg,
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
			c.JSON(http.StatusOK, PingResponse)
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
