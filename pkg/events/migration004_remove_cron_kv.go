package events

import (
	"context"
	"errors"
	"fmt"

	"github.com/nats-io/nats.go/jetstream"
)

func init() {
	registerMigration(Migration{ID: "004_remove_cron_kv", Fn: migrate004})
}

func migrate004(ctx context.Context, js *JSWrapper) error {
	// The cron KV storage changed from file-backed to memory-backed. NATS does
	// not allow changing the storage type of an existing bucket, so remove the
	// old bucket and let the cron registry recreate it with the new settings.
	if err := js.DeleteKeyValue(ctx, "cron"); err != nil &&
		!errors.Is(err, jetstream.ErrBucketNotFound) {
		return fmt.Errorf("failed to delete cron job bucket. %w", err)
	}

	return nil
}
