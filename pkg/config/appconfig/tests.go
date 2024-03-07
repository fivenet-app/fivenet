package appconfig

import (
	"context"
	"sync/atomic"

	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	"github.com/galexrt/fivenet/pkg/utils"
)

type TestConfig struct {
	cfg atomic.Pointer[Cfg]

	broker *utils.Broker[*Cfg]
}

func NewTest(ctx context.Context) (IConfig, error) {
	cfg := &TestConfig{
		cfg:    atomic.Pointer[rector.AppConfig]{},
		broker: utils.NewBroker[*Cfg](),
	}

	c := &Cfg{}
	c.Default()

	cfg.Set(c)

	go cfg.broker.Start(ctx)

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
