package eventsmigrations

import (
	"context"
	"errors"
	"fmt"

	"github.com/nats-io/nats.go/jetstream"
)

func init() {
	registerMigration(Migration{ID: "006_recreate_centrum_dispatches_kv", Fn: migrate006})
}

func migrate006(ctx context.Context, js JetStream) error {
	// Dispatch projections changed from one full resource per key to separate
	// spec, assignments, and status components. The SQL database remains
	// authoritative, so recreate the bucket and let startup repopulate it.
	if err := js.DeleteKeyValue(ctx, "centrum_dispatches"); err != nil &&
		!errors.Is(err, jetstream.ErrBucketNotFound) {
		return fmt.Errorf("failed to recreate centrum dispatches KV bucket: %w", err)
	}

	return nil
}
