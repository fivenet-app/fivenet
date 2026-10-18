package eventsmigrations

import (
	"context"
	"errors"
	"fmt"

	"github.com/nats-io/nats.go/jetstream"
)

func init() {
	registerMigration(Migration{ID: "007_recreate_centrum_units_kv", Fn: migrate007})
}

func migrate007(ctx context.Context, js JetStream) error {
	// Unit projections changed from one full resource per key to independent
	// spec, status, and membership components. SQL is authoritative, so remove
	// the incompatible bucket before Centrum repopulates it during startup.
	if err := js.DeleteKeyValue(ctx, "centrum_units"); err != nil &&
		!errors.Is(err, jetstream.ErrBucketNotFound) {
		return fmt.Errorf("failed to recreate centrum units KV bucket: %w", err)
	}

	return nil
}
