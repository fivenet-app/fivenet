package api

import (
	"context"
	_ "embed"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/clientconfig"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	serveroauth2 "github.com/fivenet-app/fivenet/v2026/pkg/server/oauth2"
	"github.com/fivenet-app/fivenet/v2026/pkg/version"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

//go:embed clear-site-data.html
var clearSiteDataPage []byte

type Routes struct {
	logger *zap.Logger

	cfg              *config.Config
	providers        []*clientconfig.ProviderConfig
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
		providers:        clientconfig.BuildProviderList(p.Config),
		clientCfg:        &atomic.Pointer[clientconfig.ClientConfig]{},
		discordInviteUrl: &atomic.Pointer[string]{},
	}

	ctx, cancel := context.WithCancel(context.Background())

	p.LC.Append(fx.StartHook(func(_ context.Context) error {
		r.handleAppConfigUpdate(p.AppConfig.Get())

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

					r.handleAppConfigUpdate(cfg)
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

var cookiesToExpire = []string{
	auth.AccCookieName,
	auth.AuthedCookieName,
	serveroauth2.SessionKeyOAuth2State,
	// Legacy cookies
	"fivenet_token",
	"fivenet_usr",
}

func (r *Routes) RegisterHTTP(e *gin.Engine) {
	// API Base
	g := e.Group("/api")
	{
		g.Any("/", func(ctx *gin.Context) {
			ctx.AbortWithStatus(http.StatusBadRequest)
		})

		g.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, PingResponse)
		})

		g.POST("/config", func(c *gin.Context) {
			c.JSON(http.StatusOK, r.clientCfg.Load())
		})

		g.GET("/clear-site-data", func(c *gin.Context) {
			c.Header("Clear-Site-Data", "\"cache\", \"cookies\", \"storage\"")
			c.Header("Cache-Control", "no-store, max-age=0")
			c.Header("Pragma", "no-cache")
			for _, name := range cookiesToExpire {
				//nolint:gosec // `Same-Site: None` is required because otherwise users can't login in the in-game tablet (iframe).
				c.SetCookieData(&http.Cookie{
					Name:     name,
					Value:    "",
					Expires:  time.Unix(0, 0).UTC(),
					MaxAge:   -1,
					Domain:   r.cfg.HTTP.Sessions.Domain,
					Path:     "/",
					HttpOnly: true,
					Secure:   true,
					SameSite: http.SameSiteNoneMode,
				})
			}
			c.Data(http.StatusOK, "text/html; charset=utf-8", clearSiteDataPage)
		})

		g.GET("/version", func(c *gin.Context) {
			c.JSON(http.StatusOK, versionInfo)
		})

		g.GET("/discord/invite-bot", func(c *gin.Context) {
			c.JSON(http.StatusOK, r.discordInviteUrl.Load())
		})
	}
}
