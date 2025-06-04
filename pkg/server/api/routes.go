package api

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"sync/atomic"

	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2025/pkg/version"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Routes struct {
	logger *zap.Logger

	cfg              *config.Config
	clientCfg        *atomic.Pointer[ClientConfig]
	discordInviteUrl *atomic.Pointer[string]
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger    *zap.Logger
	Config    *config.Config
	AppConfig appconfig.IConfig
}

func New(p Params) *Routes {
	r := &Routes{
		logger:           p.Logger,
		cfg:              p.Config,
		clientCfg:        &atomic.Pointer[ClientConfig]{},
		discordInviteUrl: &atomic.Pointer[string]{},
	}

	providers := make([]*ProviderConfig, len(p.Config.OAuth2.Providers))

	for i, p := range p.Config.OAuth2.Providers {
		providers[i] = &ProviderConfig{
			Name:     p.Name,
			Label:    p.Label,
			Icon:     p.Icon,
			Homepage: &p.Homepage,
		}
	}

	ctx, cancel := context.WithCancel(context.Background())

	p.LC.Append(fx.StartHook(func(_ context.Context) error {
		r.handleAppConfigUpdate(providers, p.AppConfig.Get())

		// Handle app config updates
		go func() {
			configUpdateCh := p.AppConfig.Subscribe()
			for {
				select {
				case <-ctx.Done():
					p.AppConfig.Unsubscribe(configUpdateCh)
					return

				case cfg := <-configUpdateCh:
					if cfg == nil {
						continue
					}

					r.handleAppConfigUpdate(providers, cfg)
				}
			}
		}()

		return nil
	}))

	p.LC.Append(fx.StopHook(func() error {
		cancel()

		return nil
	}))

	return r
}

func (r *Routes) handleAppConfigUpdate(providers []*ProviderConfig, appCfg *appconfig.Cfg) {
	clientCfg := r.buildClientConfig(providers, appCfg)
	r.clientCfg.Store(clientCfg)

	if appCfg.Discord.BotId == nil || *appCfg.Discord.BotId == "" {
		r.discordInviteUrl.Store(nil)
		return
	}

	clientId := appCfg.Discord.BotId
	permissions := strconv.FormatInt(appCfg.Discord.BotPermissions, 10)
	redirectUri, err := url.JoinPath(r.cfg.HTTP.PublicURL, "/settings/props")
	if err != nil {
		r.logger.Warn("failed to build redirect URI for discord invite", zap.Error(err))
		return
	}
	redirectUri = redirectUri + "?tab=discord#"
	scopes := "bot identify"

	u, err := url.Parse("https://discord.com/oauth2/authorize")
	if err != nil {
		r.logger.Warn("failed to build discord invite URL", zap.Error(err))
		return
	}
	q := u.Query()
	q.Set("client_id", *clientId)
	q.Set("permissions", permissions)
	q.Set("scope", scopes)
	q.Set("redirect_uri", redirectUri)
	q.Set("response_type", "code")
	u.RawQuery = q.Encode()

	inviteUrl := u.String()
	r.discordInviteUrl.Store(&inviteUrl)
}

func (r *Routes) buildClientConfig(providers []*ProviderConfig, appCfg *appconfig.Cfg) *ClientConfig {
	clientCfg := &ClientConfig{
		Version: version.Version,

		DefaultLocale: appCfg.DefaultLocale,

		Login: LoginConfig{
			SignupEnabled: appCfg.Auth.SignupEnabled,
			LastCharLock:  appCfg.Auth.LastCharLock,
			Providers:     providers,
		},
		Discord: Discord{
			BotEnabled: appCfg.Discord.BotId != nil && *appCfg.Discord.BotId != "",
		},
		Website: Website{
			Links: Links{
				Imprint:       appCfg.Website.Links.Imprint,
				PrivacyPolicy: appCfg.Website.Links.PrivacyPolicy,
			},
			StatsPage: appCfg.Website.StatsPage,
		},
		FeatureGates: FeatureGates{
			ImageProxy: r.cfg.ImageProxy.Enabled,
		},
		Game: Game{
			UnemployedJobName: appCfg.JobInfo.UnemployedJob.Name,
			StartJobGrade:     r.cfg.Game.StartJobGrade,
		},
		System: System{
			BannerMessageEnabled: appCfg.System.BannerMessageEnabled,
			OTLP: OTLPFrontend{
				Enabled: r.cfg.OTLP.Enabled,
				URL:     r.cfg.OTLP.Frontend.URL,
				Headers: r.cfg.OTLP.Frontend.Headers,
			},
		},
	}

	if appCfg.System.BannerMessage != nil {
		clientCfg.System.BannerMessage = &BannerMessage{
			Id:        appCfg.System.BannerMessage.Id,
			Title:     appCfg.System.BannerMessage.Title,
			Icon:      appCfg.System.BannerMessage.Icon,
			Color:     appCfg.System.BannerMessage.Color,
			CreatedAt: appCfg.System.BannerMessage.CreatedAt,
			ExpiresAt: appCfg.System.BannerMessage.ExpiresAt,
		}
	}

	return clientCfg
}

var versionInfo = &Version{
	Version: version.Version,
}

func (r *Routes) RegisterHTTP(e *gin.Engine) {
	// API Base
	g := e.Group("/api")
	{
		g.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, PingResponse)
		})

		g.POST("/config", func(c *gin.Context) {
			c.JSON(http.StatusOK, r.clientCfg.Load())
		})

		g.GET("/clear-site-data", func(c *gin.Context) {
			c.Header("Clear-Site-Data", "\"cache\", \"cookies\", \"storage\"")
			c.String(http.StatusOK, "Your local site data should be cleared now, please go back to the FiveNet homepage yourself.")
		})

		g.GET("/version", func(c *gin.Context) {
			c.JSON(http.StatusOK, versionInfo)
		})

		g.GET("/discord/invite-bot", func(c *gin.Context) {
			c.JSON(http.StatusOK, r.discordInviteUrl.Load())
		})
	}
}
