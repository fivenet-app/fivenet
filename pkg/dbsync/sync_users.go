package dbsync

import (
	"context"
	"strconv"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/sync"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	pbsync "github.com/fivenet-app/fivenet/gen/go/proto/services/sync"
	"github.com/go-jet/jet/v2/qrm"
)

type usersSync struct {
	*syncer

	state *TableSyncState
}

func NewUsersSync(s *syncer, state *TableSyncState) (ISyncer, error) {
	return &usersSync{
		syncer: s,
		state:  state,
	}, nil
}

func (s *usersSync) Sync(ctx context.Context) error {
	if !s.cfg.Tables.Users.Enabled {
		return nil
	}

	limit := 500
	var offset uint64
	if s.state != nil && s.state.Offset > 0 {
		offset = s.state.Offset
	}

	sQuery := s.cfg.Tables.Users
	query := prepareStringQuery(sQuery, s.state, offset, limit)

	users := []*users.User{}
	if _, err := qrm.Query(ctx, s.db, query, []interface{}{}, &users); err != nil {
		return err
	}

	if len(users) == 0 {
		return nil
	}

	// TODO retrieve licenses and vehicles for each user selected

	if s.cli != nil {
		if _, err := s.cli.SyncData(ctx, &pbsync.SyncDataRequest{
			Data: &pbsync.SyncDataRequest_Users{
				Users: &sync.DataUsers{
					Users: users,
				},
			},
		}); err != nil {
			return err
		}
	}

	// If less users than limit are returned, we probably have reached the "end" of the table
	// and need to reset the offset to 0
	if len(users) < limit {
		offset = 0
	}

	lastUserId := strconv.Itoa(int(users[len(users)-1].UserId))
	s.state.Set(
		s.cfg.Tables.Users.IDField,
		uint64(limit)+offset,
		&lastUserId,
	)

	return nil
}
