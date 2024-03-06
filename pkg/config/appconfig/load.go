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

type Cfg = rector.AppConfig

var (
	tConfig = table.FivenetConfig.AS("appconfig")
)

type IConfig interface {
	Get() *Cfg
	Set(val *Cfg)
	Update(*Cfg) error

	Subscribe() chan *Cfg
	Unsubscribe(c chan *Cfg)
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

	cfg atomic.Pointer[Cfg]

	broker *utils.Broker[*Cfg]
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

		cfg: atomic.Pointer[Cfg]{},

		broker: utils.NewBroker[*Cfg](ctx),
	}

	p.LC.Append(fx.StartHook(func(ctx context.Context) error {
		go c.broker.Start()

		if _, err := c.updateConfigFromDB(ctx); err != nil {
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

func (c *Config) Get() *Cfg {
	return c.cfg.Load()
}

func (c *Config) Set(val *Cfg) {
	c.cfg.Store(val)
}

func (c *Config) Update(val *Cfg) error {
	c.Set(val)

	// Send update message to inform components
	if _, err := c.js.Publish(fmt.Sprintf("%s.%s", BaseSubject, UpdateSubject), nil); err != nil {
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
	cfg, err := c.LoadFromDB(c.ctx)
	if err != nil {
		return nil, err
	}

	c.Set(cfg)

	return cfg, nil
}

func (c *Config) LoadFromDB(ctx context.Context) (*Cfg, error) {
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

func LoadTest() (*Config, error) {
	cfg := &Config{
		cfg: atomic.Pointer[rector.AppConfig]{},
	}

	c := &Cfg{}
	c.Default()

	cfg.Set(c)

	return cfg, nil
}
