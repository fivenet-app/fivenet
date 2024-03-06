package appconfig

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/galexrt/fivenet/query/fivenet/table"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/nats-io/nats.go"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var (
	tConfig = table.FivenetConfig.AS("appconfig")
)

type IConfig interface {
	Get() *rector.AppConfig
	Set(val *rector.AppConfig)
	Update(*rector.AppConfig) error

	Subscribe() chan *rector.AppConfig
	Unsubscribe(c chan *rector.AppConfig)
}

var Module = fx.Module("appconfig",
	fx.Provide(
		New,
	),
)

type Config struct {
	IConfig

	ctx    context.Context
	logger *zap.Logger
	db     *sql.DB
	js     nats.JetStreamContext

	jsSub *nats.Subscription

	cfg atomic.Pointer[rector.AppConfig]

	broker *utils.Broker[*rector.AppConfig]
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger *zap.Logger
	JS     nats.JetStreamContext
	DB     *sql.DB
}

func New(p Params) (*Config, error) {
	ctx, cancel := context.WithCancel(context.Background())

	c := &Config{
		ctx:    ctx,
		logger: p.Logger,
		db:     p.DB,
		js:     p.JS,

		cfg: atomic.Pointer[rector.AppConfig]{},

		broker: utils.NewBroker[*rector.AppConfig](ctx),
	}

	p.LC.Append(fx.StartHook(func(ctx context.Context) error {
		if err := c.updateConfigFromDB(ctx); err != nil {
			return err
		}

		return c.registerEvents(ctx)
	}))

	p.LC.Append(fx.StopHook(func(ctx context.Context) error {
		if c.jsSub != nil {
			return c.jsSub.Unsubscribe()
		}

		cancel()

		return nil
	}))

	return c, nil
}

func (c *Config) Get() *rector.AppConfig {
	return c.cfg.Load()
}

func (c *Config) Set(val *rector.AppConfig) {
	c.cfg.Store(val)
}

func (c *Config) Update(val *rector.AppConfig) error {
	if _, err := c.js.Publish(fmt.Sprintf("%s.%s", BaseSubject, UpdateSubject), nil); err != nil {
		return err
	}

	c.Set(val)

	return nil
}

func (c *Config) Subscribe() chan *rector.AppConfig {
	return c.broker.Subscribe()
}

func (c *Config) Unsubscribe(ch chan *rector.AppConfig) {
	c.broker.Unsubscribe(ch)
}

func (c *Config) updateConfigFromDB(ctx context.Context) error {
	cfg, err := c.LoadFromDB(c.ctx)
	if err != nil {
		return err
	}

	c.Set(cfg)

	return nil
}

func (c *Config) LoadFromDB(ctx context.Context) (*rector.AppConfig, error) {
	stmt := tConfig.
		SELECT(
			tConfig.AppConfig.AS("appconfig"),
		).
		FROM(tConfig).
		LIMIT(1)

	dest := struct {
		AppConfig *rector.AppConfig
	}{
		AppConfig: &rector.AppConfig{},
	}
	if err := stmt.QueryContext(ctx, c.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	dest.AppConfig.Default()

	return dest.AppConfig, nil
}
