package dbsync

import (
	"context"
	"database/sql"

	pbsync "github.com/fivenet-app/fivenet/gen/go/proto/services/sync"
	"github.com/fivenet-app/fivenet/pkg/config"
	"go.uber.org/zap"
)

type ISyncer interface {
	Sync(ctx context.Context) (*TableSyncState, error)
}

type SyncerFactory = func(logger *zap.Logger, db *sql.DB, cfg *config.DBSync, cli pbsync.SyncServiceClient) (ISyncer, error)

var syncerFactories = map[string]SyncerFactory{}
