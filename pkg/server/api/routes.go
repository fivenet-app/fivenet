package api

import (
	"context"
	"net/http"
	"sync/atomic"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/clientconfig"
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
	clientCfg        *atomic.Pointer[clientconfig.ClientConfig]
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
		clientCfg:        &atomic.Pointer[clientconfig.ClientConfig]{},
		discordInviteUrl: &atomic.Pointer[string]{},
	}

	providers := clientconfig.BuildProviderList(p.Config)

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
			c.String(
				http.StatusOK,
				"Your local site data should be cleared now, please go back to the FiveNet homepage yourself.",
			)
		})

		g.GET("/version", func(c *gin.Context) {
			c.JSON(http.StatusOK, versionInfo)
		})

		g.GET("/discord/invite-bot", func(c *gin.Context) {
			c.JSON(http.StatusOK, r.discordInviteUrl.Load())
		})
	}
}
