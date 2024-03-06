package storage

import (
	"fmt"

	"github.com/galexrt/fivenet/pkg/config"
	"go.uber.org/fx"
)

var storageFactories = map[string]func(lc fx.Lifecycle, cfg *config.Config) (IStorage, error){}

var Module = fx.Module("storage",
	fx.Provide(New),
)

type Params struct {
	fx.In

	LC     fx.Lifecycle
	Config *config.Config
}

func New(p Params) (IStorage, error) {
	fn, ok := storageFactories[p.Config.Storage.Type]
	if !ok {
		return nil, fmt.Errorf("invalid storage '%s' factory given", p.Config.Storage.Type)
	}

	st, err := fn(p.LC, p.Config)
	if err != nil {
		return nil, err
	}

	return st, nil
}
