package api

import (
	"net/http"

	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/version"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Routes struct {
	logger *zap.Logger

	clientCfg *ClientConfig
}

func New(logger *zap.Logger, cfg *config.Config) *Routes {
	providers := make([]*ProviderConfig, len(cfg.OAuth2.Providers))

	for i, p := range cfg.OAuth2.Providers {
		providers[i] = &ProviderConfig{
			Name:  p.Name,
			Label: p.Label,
		}
	}

	clientCfg := &ClientConfig{
		Version:   version.Version,
		SentryDSN: cfg.Sentry.ClientDSN,
		SentryEnv: cfg.Sentry.Environment,
		Login: LoginConfig{
			SignupEnabled: cfg.Game.Auth.SignupEnabled,
			Providers:     providers,
		},
		Discord: &Discord{},
		Links:   Links{},
	}

	if cfg.Discord.Bot.Enabled {
		clientCfg.Discord.BotInviteURL = cfg.Discord.Bot.InviteURL
	}

	if cfg.HTTP.Links.Imprint != nil {
		clientCfg.Links.Imprint = cfg.HTTP.Links.Imprint
	}
	if cfg.HTTP.Links.PrivacyPolicy != nil {
		clientCfg.Links.PrivacyPolicy = cfg.HTTP.Links.PrivacyPolicy
	}

	return &Routes{
		logger:    logger,
		clientCfg: clientCfg,
	}
}

func (r *Routes) Register(e *gin.Engine) {
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

		ver := &Version{
			Version: r.clientCfg.Version,
		}
		g.GET("/version", func(c *gin.Context) {
			c.JSON(http.StatusOK, ver)
		})
	}
}
