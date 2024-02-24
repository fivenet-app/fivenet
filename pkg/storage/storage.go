package storage

import (
	"fmt"

	"github.com/galexrt/fivenet/pkg/config"
	"go.uber.org/fx"
)

var storageFactories = map[string]func(cfg *config.Config) (IStorage, error){}

var Module = fx.Module("storage",
	fx.Provide(New),
)

func New(cfg *config.Config) (IStorage, error) {
	if !cfg.Storage.Enabled {
		return nil, nil
	}

	fn, ok := storageFactories[cfg.Storage.Type]
	if !ok {
		return nil, fmt.Errorf("invalid storage '%s' factory given", cfg.Storage.Type)
	}

	st, err := fn(cfg)
	if err != nil {
		return nil, err
	}

	return st, nil
}
