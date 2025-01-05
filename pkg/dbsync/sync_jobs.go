package dbsync

import (
	"context"
	"database/sql"

	pbsync "github.com/fivenet-app/fivenet/gen/go/proto/services/sync"
	"github.com/fivenet-app/fivenet/pkg/config"
	"go.uber.org/zap"
)

func init() {
	syncerFactories["jobs"] = NewJobsSync
}

type jobsSync struct {
	logger *zap.Logger
	db     *sql.DB

	cfg *config.DBSync

	cli pbsync.SyncServiceClient
}

func NewJobsSync(logger *zap.Logger, db *sql.DB, cfg *config.DBSync, cli pbsync.SyncServiceClient) (ISyncer, error) {
	return &jobsSync{
		logger: logger,
		db:     db,
		cfg:    cfg,

		cli: cli,
	}, nil
}

func (s *jobsSync) Sync(ctx context.Context) (*TableSyncState, error) {
	// TODO

	return nil, nil
}
