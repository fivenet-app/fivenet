package storage

import (
	"fmt"

	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
)

var storageFactories = map[config.StorageType]func(p Params) (IStorage, error){}

var Module = fx.Module("storage",
	fx.Provide(New),
)

type Params struct {
	fx.In

	LC fx.Lifecycle

	TP  *tracesdk.TracerProvider
	Cfg *config.Config
}

func New(p Params) (IStorage, error) {
	fn, ok := storageFactories[p.Cfg.Storage.Type]
	if !ok {
		return nil, fmt.Errorf("invalid storage '%s' factory given", p.Cfg.Storage.Type)
	}

	st, err := fn(p)
	if err != nil {
		return nil, err
	}

	return st, nil
}
