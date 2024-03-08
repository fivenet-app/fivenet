package appconfig

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	"github.com/galexrt/fivenet/pkg/utils/broker"
	"github.com/galexrt/fivenet/query/fivenet/table"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Cfg = rector.AppConfig

var (
	tConfig = table.FivenetConfig.AS("appconfig")
)

type IConfig interface {
	Get() *Cfg
	Set(val *Cfg)
	Update(ctx context.Context, val *Cfg) error
	Reload(ctx context.Context) (*Cfg, error)

	Subscribe() chan *Cfg
	Unsubscribe(ch chan *Cfg)
}

var Module = fx.Module("appconfig",
	fx.Provide(
		New,
	),
)

type Config struct {
	IConfig

	logger *zap.Logger
	db     *sql.DB
	js     jetstream.JetStream

	jsCons jetstream.ConsumeContext

	cfg atomic.Pointer[Cfg]

	broker *broker.Broker[*Cfg]
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger *zap.Logger
	JS     jetstream.JetStream
	DB     *sql.DB
}

func New(p Params) (IConfig, error) {
	ctx, cancel := context.WithCancel(context.Background())

	cfg := &Config{
		logger: p.Logger,
		db:     p.DB,
		js:     p.JS,

		cfg: atomic.Pointer[Cfg]{},

		broker: broker.New[*Cfg](),
	}

	p.LC.Append(fx.StartHook(func(c context.Context) error {
		go cfg.broker.Start(ctx)

		if _, err := cfg.updateConfigFromDB(c); err != nil {
			return err
		}

		return cfg.registerSubscriptions(c)
	}))

	p.LC.Append(fx.StopHook(func(ctx context.Context) error {
		cancel()

		cfg.jsCons.Stop()

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
	if _, err := c.js.Publish(ctx, fmt.Sprintf("%s.%s", BaseSubject, UpdateSubject), nil); err != nil {
		return err
	}

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

func (c *Config) Reload(ctx context.Context) (*Cfg, error) {
	stmt := tConfig.
		SELECT(
			tConfig.AppConfig.AS("appconfig"),
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
		}
	}

	dest.AppConfig.Default()

	return dest.AppConfig, nil
}
