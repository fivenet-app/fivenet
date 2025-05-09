package appconfig

import (
	"context"
	"sync/atomic"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/rector"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/broker"
	"go.uber.org/fx"
)

type TestConfig struct {
	cfg atomic.Pointer[Cfg]

	broker *broker.Broker[*Cfg]
}

var TestModule = fx.Module("appconfig_test",
	fx.Provide(
		NewTest,
	),
)

type TestParams struct {
	fx.In

	LC fx.Lifecycle
}

func NewTest(p TestParams) (IConfig, error) {
	cfg := &TestConfig{
		cfg:    atomic.Pointer[rector.AppConfig]{},
		broker: broker.New[*Cfg](),
	}

	c := &Cfg{}
	c.Default()

	cfg.Set(c)

	ctx, cancel := context.WithCancel(context.Background())
	p.LC.Append(fx.StartHook(func(_ context.Context) error {
		go cfg.broker.Start(ctx)

		return nil
	}))

	p.LC.Append(fx.StartHook(func(ctx context.Context) error {
		cancel()

		return nil
	}))

	return cfg, nil
}

func (c *TestConfig) Get() *Cfg {
	return c.cfg.Load()
}

func (c *TestConfig) Set(val *Cfg) {
	c.cfg.Store(val)
}

func (c *TestConfig) Update(ctx context.Context, val *Cfg) error {
	c.Set(val)

	return nil
}

func (c *TestConfig) Reload(ctx context.Context) (*Cfg, error) {
	// Noop during testing
	return c.Get(), nil
}

func (c *TestConfig) Subscribe() chan *Cfg {
	return c.broker.Subscribe()
}

func (c *TestConfig) Unsubscribe(ch chan *Cfg) {
	c.broker.Unsubscribe(ch)
}
