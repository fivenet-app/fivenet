package dbsync

import (
	"context"
	"database/sql"

	pbsync "github.com/fivenet-app/fivenet/gen/go/proto/services/sync"
	"github.com/fivenet-app/fivenet/pkg/config"
	"go.uber.org/zap"
)

type ISyncer interface {
	Sync(ctx context.Context) error
}

type SyncerFactory = func(s *syncer, state *TableSyncState) (ISyncer, error)

type syncer struct {
	logger *zap.Logger
	db     *sql.DB

	cfg *config.DBSync

	cli pbsync.SyncServiceClient
}
