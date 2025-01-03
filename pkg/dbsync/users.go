package dbsync

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/zap"
)

func init() {
	syncerFactories["users"] = NewUsersSync
}

type usersSync struct {
	logger *zap.Logger
	db     *sql.DB

	cfg *config.DBSync
}

func NewUsersSync(logger *zap.Logger, db *sql.DB, cfg *config.DBSync) (ISyncer, error) {
	return &usersSync{
		logger: logger,
		db:     db,
		cfg:    cfg,
	}, nil
}

func (s *usersSync) Sync(ctx context.Context) (*TableSyncState, error) {
	if !s.cfg.Tables.Users.Enabled {
		return nil, nil
	}

	offset := 0
	limit := 100

	sQuery := s.cfg.Tables.Users
	query := prepareStringQuery(sQuery.Query, offset, limit)

	users := []*users.User{}
	if _, err := qrm.Query(ctx, s.db, query, []interface{}{}, &users); err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return &TableSyncState{}, nil
	}

	// If less users than limit are returned, we probably have reached the "end" of the table
	// and need to reset the offset to 0
	if len(users) < limit {
		offset = 0
	}

	lastUserId := strconv.Itoa(int(users[len(users)-1].UserId))

	return &TableSyncState{
		IDField: s.cfg.Tables.Users.IDField,
		Offset:  uint64(limit + offset),
		LastID:  &lastUserId,
	}, nil
}
