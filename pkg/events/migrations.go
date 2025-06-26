package events

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"

	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
)

var (
	mu                   sync.Mutex
	registeredMigrations = []Migration{}
)

// registerMigration adds a migration to the registry; called in each migration file's init()
func registerMigration(m Migration) {
	mu.Lock()
	defer mu.Unlock()

	registeredMigrations = append(registeredMigrations, m)
}

// migrationBucket is the KV bucket used to track applied migrations.
const migrationBucket = "NATS_MIGRATIONS"

// Migration encapsulates an id and a function that applies the migration.
type Migration struct {
	// ID should be a monotonically increasing string, e.g. "001", "002"
	ID string
	// Fn performs the migration against the provided JetStream context.
	Fn func(ctx context.Context, js *JSWrapper) error
}

// runMigrations executes any pending migrations.  On a fresh install (no version recorded),
// it records the latest migration ID and does nothing else.
func runMigrations(ctx context.Context, logger *zap.Logger, js *JSWrapper) error {
	// 1) Ensure the migration tracking bucket exists
	kv, err := js.CreateOrUpdateKeyValue(ctx, jetstream.KeyValueConfig{
		Bucket:      migrationBucket,
		Description: "Tracks applied migrations",
	})
	if err != nil {
		return fmt.Errorf("failed to ensure migration KV bucket. %w", err)
	}

	// Fetch the recorded latest version
	rec, err := kv.Get(ctx, "latest")
	if err == jetstream.ErrKeyNotFound {
		// Fresh install: record highest and exit
		latest := registeredMigrations[len(registeredMigrations)-1].ID
		if _, err := kv.Put(ctx, "latest", []byte(latest)); err != nil {
			return fmt.Errorf("failed to record initial migration %s. %w", latest, err)
		}
		logger.Info("fresh install: recorded latest migration", zap.String("version", latest))
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to fetch recorded migration. %w", err)
	}

	current := string(rec.Value())

	// Sort migrations by ID
	sort.Slice(registeredMigrations, func(i, j int) bool {
		return registeredMigrations[i].ID < registeredMigrations[j].ID
	})

	// Apply migrations with ID greater than current
	for _, m := range registeredMigrations {
		if m.ID <= current {
			continue
		}
		logger.Info("applying migration", zap.String("id", m.ID))
		if err := m.Fn(ctx, js); err != nil {
			return fmt.Errorf("migration %s failed. %w", m.ID, err)
		}
		// record this version
		if _, err := kv.Put(ctx, "latest", []byte(m.ID)); err != nil {
			return fmt.Errorf("failed to record migration %s. %w", m.ID, err)
		}
		logger.Info("migration applied", zap.String("id", m.ID))
	}

	return nil
}

func init() {
	registerMigration(Migration{ID: "001_remove_user_locations", Fn: migrate001})
}

func migrate001(ctx context.Context, js *JSWrapper) error {
	// Remove any existing "_snapshot" mapping from the tracker buckets, as it is not used anymore
	for _, bucket := range []string{
		"userloc",
		"user_mappings",
		"userloc_by_id",
	} {
		kv, err := js.KeyValue(ctx, bucket)
		if err == nil {
			_ = kv.Delete(ctx, "_snapshot")
		} else if !errors.Is(err, jetstream.ErrBucketNotFound) {
			return fmt.Errorf("failed to get user locations bucket. %w", err)
		}
	}

	// Remove user_locations bucket if it exists
	if err := js.DeleteKeyValue(ctx, "user_locations"); err != nil {
		if !errors.Is(err, jetstream.ErrBucketNotFound) {
			return fmt.Errorf("failed to delete user_locations bucket. %w", err)
		}
	}

	// Remove buckets which didn't had their LimitMarkerTTL set correctly
	for _, bucket := range []string{
		"userinfo_poll_ttl",
		"leader_election",
	} {
		if err := js.DeleteKeyValue(ctx, bucket); err != nil {
			if !errors.Is(err, jetstream.ErrBucketNotFound) {
				return fmt.Errorf("failed to delete userloc bucket. %w", err)
			}
		}
	}

	// Remove _owner key from cron jobs bucket
	kv, err := js.KeyValue(ctx, "cron")
	if err == nil {
		_ = kv.Delete(ctx, "_owner")
		_ = kv.Delete(ctx, "LOCK._owner")
	} else if !errors.Is(err, jetstream.ErrBucketNotFound) {
		return fmt.Errorf("failed to get user locations bucket. %w", err)
	}

	if err := js.DeleteStream(ctx, "USERINFO"); err != nil {
		if !errors.Is(err, jetstream.ErrStreamNotFound) {
			return fmt.Errorf("failed to delete userinfo stream. %w", err)
		}
	}

	return nil
}
