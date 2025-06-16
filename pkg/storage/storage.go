package storage

import (
	"fmt"

	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
)

// storageFactories maps storage types to their factory functions for dependency injection.
var storageFactories = map[config.StorageType]func(p Params) (IStorage, error){}

// Module is the Fx module for storage, providing the New constructor for dependency injection.
var Module = fx.Module("storage",
	fx.Provide(New),
)

// Params contains dependencies for constructing a storage backend.
type Params struct {
	fx.In

	// LC is the Fx lifecycle for managing start/stop hooks.
	LC fx.Lifecycle

	// TP is the OpenTelemetry tracer provider for tracing storage operations.
	TP *tracesdk.TracerProvider
	// Cfg is the application configuration.
	Cfg *config.Config
}

// New creates a new IStorage instance using the configured storage type and parameters.
// Returns an error if the storage type is invalid or the factory fails.
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
