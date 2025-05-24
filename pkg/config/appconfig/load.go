package appconfig

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync/atomic"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/settings"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/broker"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
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

	logger *zap.Logger
	db     *sql.DB
	tracer trace.Tracer
	nc     *nats.Conn
	ncSub  *nats.Subscription

	jsCons jetstream.ConsumeContext

	cfg atomic.Pointer[Cfg]

	broker *broker.Broker[*Cfg]
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger *zap.Logger
	NC     *nats.Conn
	TP     *tracesdk.TracerProvider
	DB     *sql.DB
}

func New(p Params) (IConfig, error) {
	cfg := &Config{
		logger: p.Logger,
		db:     p.DB,
		tracer: p.TP.Tracer("appconfig"),
		nc:     p.NC,

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
			tConfig.AppConfig.SET(jet.RawString("VALUES(`app_config`)")),
		)

	if _, err := stmt.ExecContext(ctx, c.db); err != nil {
		return err
	}

	return nil
}

func (c *Config) Reload(ctx context.Context) (*Cfg, error) {
	stmt := tConfig.
		SELECT(
			tConfig.AppConfig.AS("app_config"),
		).
		FROM(tConfig).
		LIMIT(1)

	dest := struct {
		AppConfig *Cfg
	}{
		AppConfig: &Cfg{},
	}

	if err := stmt.QueryContext(ctx, c.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		} else {
			dest.AppConfig.Default()
			if err := c.updateConfigInDB(ctx, dest.AppConfig); err != nil {
				return nil, err
			}
		}
	}
	dest.AppConfig.Default()

	if slices.ContainsFunc(dest.AppConfig.Perms.Default, func(p *settings.Perm) bool {
		return !strings.Contains(p.Category, ".")
	}) {
		c.logger.Error("WARNING! You must update the default permissions in the app config to include the category prefix.")
	}

	return dest.AppConfig, nil
}
