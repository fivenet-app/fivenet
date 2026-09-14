package appconfig

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync/atomic"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/settings"
	serverconfig "github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/broker"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Cfg = settings.AppConfig

var tConfig = table.FivenetConfig.AS("app_config")

type IConfig interface {
	Get() *Cfg
	Set(val *Cfg)
	Update(ctx context.Context, val *Cfg) error
	Reload(ctx context.Context) (*Cfg, error)

	Subscribe() chan *Cfg
	Unsubscribe(ch chan *Cfg)
}

var Module = fx.Module("app_config",
	fx.Provide(
		New,
	),
)

type Config struct {
	IConfig

	logger       *zap.Logger
	db           *sql.DB
	tracer       trace.Tracer
	nc           *nats.Conn
	ncSub        *nats.Subscription
	serverConfig *serverconfig.Config

	jsCons jetstream.ConsumeContext

	cfg atomic.Pointer[Cfg]

	broker *broker.Broker[*Cfg]
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger       *zap.Logger
	NC           *nats.Conn
	TP           *tracesdk.TracerProvider
	DB           *sql.DB
	ServerConfig *serverconfig.Config
}

func New(p Params) (IConfig, error) {
	cfg := &Config{
		logger:       p.Logger.Named("appconfig"),
		db:           p.DB,
		tracer:       p.TP.Tracer("appconfig"),
		nc:           p.NC,
		serverConfig: p.ServerConfig,

		cfg: atomic.Pointer[Cfg]{},

		broker: broker.New[*Cfg](),
	}

	ctxCancel, cancel := context.WithCancel(context.Background())
	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		go cfg.broker.Start(ctxCancel)

		if _, err := cfg.updateConfigFromDB(ctxStartup); err != nil {
			return err
		}

		return cfg.registerSubscriptions(ctxCancel)
	}))

	p.LC.Append(fx.StopHook(func(ctx context.Context) error {
		cancel()

		if cfg.jsCons != nil {
			cfg.jsCons.Stop()
			cfg.jsCons = nil
		}

		return nil
	}))

	return cfg, nil
}

func (c *Config) Get() *Cfg {
	return c.cfg.Load()
}

func (c *Config) Set(val *Cfg) {
	c.cfg.Store(val)
}

func (c *Config) Update(ctx context.Context, val *Cfg) error {
	if current := c.Get(); current != nil && !val.HasSetupComplete() && current.HasSetupComplete() {
		val.SetSetupComplete(current.GetSetupComplete())
	}

	c.Set(val)

	// Send update message to inform components
	if err := c.nc.Publish(fmt.Sprintf("%s.%s", BaseSubject, UpdateSubject), nil); err != nil {
		return err
	}

	// Retrieve config and publish event to "self" (we don't want to rely on nats echo functionality)
	c.broker.Publish(c.Get())

	return nil
}

func (c *Config) Subscribe() chan *Cfg {
	return c.broker.Subscribe()
}

func (c *Config) Unsubscribe(ch chan *Cfg) {
	c.broker.Unsubscribe(ch)
}

func (c *Config) updateConfigFromDB(ctx context.Context) (*Cfg, error) {
	cfg, err := c.Reload(ctx)
	if err != nil {
		return nil, err
	}

	c.Set(cfg)

	return cfg, nil
}

func (c *Config) updateConfigInDB(ctx context.Context, cfg *Cfg) error {
	tConfig := table.FivenetConfig
	stmt := tConfig.
		INSERT(
			tConfig.Key,
			tConfig.AppConfig,
		).
		VALUES(
			1,
			cfg,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tConfig.AppConfig.SET(mysql.RawString("VALUES(`app_config`)")),
		)

	if _, err := stmt.ExecContext(ctx, c.db); err != nil {
		return err
	}

	return nil
}

// insertInitialConfig inserts an initial app config without replacing a value created by another replica.
func (c *Config) insertInitialConfig(ctx context.Context, cfg *Cfg) error {
	tConfig := table.FivenetConfig
	stmt := tConfig.
		INSERT(
			tConfig.Key,
			tConfig.AppConfig,
		).
		VALUES(
			1,
			cfg,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tConfig.Key.SET(tConfig.Key),
		)

	_, err := stmt.ExecContext(ctx, c.db)
	return err
}

func (c *Config) Reload(ctx context.Context) (*Cfg, error) {
	ctx, span := c.tracer.Start(ctx, "appconfig.reload")
	defer span.End()

	stmt := tConfig.
		SELECT(
			tConfig.AppConfig.AS("app_config"),
			tConfig.SetupComplete.AS("setup_complete"),
		).
		FROM(tConfig).
		LIMIT(1)

	dest := struct {
		AppConfig     *Cfg
		SetupComplete bool
	}{
		AppConfig: &Cfg{},
	}

	if err := stmt.QueryContext(ctx, c.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		} else {
			// No app config found in database? Insert into database.
			dest.AppConfig.Default()
			c.applyInitialConfig(dest.AppConfig)
			if err := c.insertInitialConfig(ctx, dest.AppConfig); err != nil {
				return nil, err
			}
			// Always load the persisted row: another replica may have created it first.
			return c.Reload(ctx)
		}
	}
	dest.AppConfig.Default()
	dest.AppConfig.Migrate()
	dest.AppConfig.SetSetupComplete(dest.SetupComplete)
	c.applyStartupOverride(dest.AppConfig)

	if slices.ContainsFunc(dest.AppConfig.Perms.GetDefault(), func(p *settings.Perm) bool {
		return !strings.Contains(p.GetCategory(), ".")
	}) {
		c.logger.Warn(
			"WARNING! You must update the default permissions in the app config to include the category prefix.",
		)
	}

	return dest.AppConfig, nil
}

func (c *Config) applyInitialConfig(cfg *Cfg) {
	if c.serverConfig == nil {
		return
	}

	applyWebsiteLinks(cfg, c.serverConfig.AppConfig.Initial.Website.Links)
}

func (c *Config) applyStartupOverride(cfg *Cfg) {
	if c.serverConfig == nil || !c.serverConfig.AppConfig.Override.Enabled {
		return
	}

	applyWebsiteLinks(cfg, c.serverConfig.AppConfig.Override.Website.Links)
}

func applyWebsiteLinks(cfg *Cfg, links *serverconfig.AppConfigLinks) {
	if cfg == nil || links == nil {
		return
	}

	if cfg.Website == nil {
		cfg.Website = &settings.Website{}
	}
	if cfg.Website.Links == nil {
		cfg.Website.Links = &settings.Links{}
	}

	if links.PrivacyPolicy != nil {
		cfg.Website.Links.SetPrivacyPolicy(*links.PrivacyPolicy)
	}
	if links.Imprint != nil {
		cfg.Website.Links.SetImprint(*links.Imprint)
	}
}
