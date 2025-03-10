package api

import (
	"context"
	"net/http"
	"sync/atomic"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/rector"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/pkg/version"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type Routes struct {
	logger *zap.Logger

	cfg       *config.Config
	clientCfg *atomic.Pointer[ClientConfig]
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
		logger:    p.Logger,
		cfg:       p.Config,
		clientCfg: &atomic.Pointer[ClientConfig]{},
	}

	providers := make([]*ProviderConfig, len(p.Config.OAuth2.Providers))

	for i, p := range p.Config.OAuth2.Providers {
		providers[i] = &ProviderConfig{
			Name:  p.Name,
			Label: p.Label,
			Icon:  p.Icon,
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
		Discord: Discord{},
		Website: Website{
			Links:     Links{},
			StatsPage: appCfg.Website.StatsPage,
		},
		FeatureGates: FeatureGates{
			ImageProxy: r.cfg.ImageProxy.Enabled,
		},
		Game: Game{
			UnemployedJobName: "unemployed",
			StartJobGrade:     0,
		},
		System: System{},
	}

	clientCfg.Discord.BotInviteURL = appCfg.Discord.InviteUrl

	if appCfg.Website.Links.Imprint != nil {
		clientCfg.Website.Links.Imprint = appCfg.Website.Links.Imprint
	}
	if appCfg.Website.Links.PrivacyPolicy != nil {
		clientCfg.Website.Links.PrivacyPolicy = appCfg.Website.Links.PrivacyPolicy
	}

	if appCfg.JobInfo.UnemployedJob != nil {
		clientCfg.Game.UnemployedJobName = appCfg.JobInfo.UnemployedJob.Name
	}
	clientCfg.Game.StartJobGrade = r.cfg.Game.StartJobGrade

	if appCfg.System.BannerMessage != nil {
		clientCfg.System.BannerMessage = proto.Clone(appCfg.System.BannerMessage).(*rector.BannerMessage)
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
	}
}
