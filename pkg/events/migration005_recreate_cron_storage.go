package events

import (
	"context"
	"errors"
	"fmt"

	"github.com/nats-io/nats.go/jetstream"
)

func init() {
	registerMigration(Migration{ID: "005_recreate_cron_storage", Fn: migrate005})
}

func migrate005(ctx context.Context, js *JSWrapper) error {
	// Cron storage changed from file-backed to memory-backed for the schedule
	// stream, while the cron job KV changed in the opposite direction. NATS
	// does not allow changing storage type in place, so remove both resources
	// and let the cron registry recreate them with the current settings.
	if err := js.DeleteStream(ctx, "CRON_SCHEDULE"); err != nil &&
		!errors.Is(err, jetstream.ErrStreamNotFound) {
		return fmt.Errorf("failed to delete cron schedule stream. %w", err)
	}

	if err := js.DeleteKeyValue(ctx, "cron"); err != nil &&
		!errors.Is(err, jetstream.ErrBucketNotFound) {
		return fmt.Errorf("failed to delete cron job bucket. %w", err)
	}

	return nil
}
