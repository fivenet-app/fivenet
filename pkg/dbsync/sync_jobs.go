package dbsync

import (
	"context"
	"database/sql"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/sync"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	pbsync "github.com/fivenet-app/fivenet/gen/go/proto/services/sync"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/go-jet/jet/v2/qrm"
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
	if !s.cfg.Tables.Jobs.Enabled {
		return nil, nil
	}

	offset := 0
	limit := 100

	sQuery := s.cfg.Tables.Jobs
	query := prepareStringQuery(sQuery.Query, offset, limit)

	jobs := []*users.Job{}
	if _, err := qrm.Query(ctx, s.db, query, []interface{}{}, &jobs); err != nil {
		return nil, err
	}

	if len(jobs) == 0 {
		return &TableSyncState{}, nil
	}

	if s.cli != nil {
		if _, err := s.cli.SyncData(ctx, &pbsync.SyncDataRequest{
			Data: &pbsync.SyncDataRequest_Jobs{
				Jobs: &sync.DataJobs{
					Jobs: jobs,
				},
			},
		}); err != nil {
			return nil, err
		}
	}

	// If less users than limit are returned, we probably have reached the "end" of the table
	// and need to reset the offset to 0
	if len(jobs) < limit {
		offset = 0
	}

	lastUserId := jobs[len(jobs)-1].Name

	return &TableSyncState{
		IDField: s.cfg.Tables.Jobs.IDField,
		Offset:  uint64(limit + offset),
		LastID:  &lastUserId,
	}, nil
}
