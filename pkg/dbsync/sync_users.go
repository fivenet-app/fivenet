package dbsync

import (
	"context"
	"strconv"
	"strings"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/sync"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	pbsync "github.com/fivenet-app/fivenet/gen/go/proto/services/sync"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/multierr"
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

	if s.cfg.Tables.UserLicenses.Enabled {
		// Retrieve user' licenses
		errs := multierr.Combine()
		var err error
		for k := range users {
			identifier := ""
			if users[k].Identifier != nil {
				identifier = *users[k].Identifier
			}

			users[k].Licenses, err = s.retrieveLicenses(ctx, users[k].UserId, identifier)
			if err != nil {
				errs = multierr.Append(errs, err)
			}
		}

		if errs != nil {
			return errs
		}
	}

	if s.cli != nil {
		if _, err := s.cli.SendData(ctx, &pbsync.SendDataRequest{
			Data: &pbsync.SendDataRequest_Users{
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

func (s *usersSync) retrieveLicenses(ctx context.Context, userId int32, identifier string) ([]*users.License, error) {
	sQuery := s.cfg.Tables.UserLicenses
	query := prepareStringQuery(sQuery, s.state, 0, 100)

	args := []interface{}{}
	if strings.Contains(query, "$userId") {
		query = strings.ReplaceAll(query, "$userId", strconv.Itoa(int(userId)))
		args = append(args, userId)
	} else if strings.Contains(query, "$identifier") {
		query = strings.ReplaceAll(query, "$identifier", identifier)
		args = append(args, identifier)
	}

	licenses := []*users.License{}
	if _, err := qrm.Query(ctx, s.db, query, args, &licenses); err != nil {
		return nil, err
	}

	return licenses, nil
}

// Sync an individual user/char info
func (s *usersSync) SyncUser(ctx context.Context, userId int32) error {
	sQuery := s.cfg.Tables.Users
	query := prepareStringQuery(sQuery, s.state, 0, 1)
	_ = query

	// TODO

	return nil
}
