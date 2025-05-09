package dbsync

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/sync"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/users"
	pbsync "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/sync"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type usersSync struct {
	*syncer

	state *TableSyncState
}

func newUsersSync(s *syncer, state *TableSyncState) *usersSync {
	return &usersSync{
		syncer: s,
		state:  state,
	}
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

	// Ensure to zero the last check time if the data hasn't fully synced yet
	if !s.state.SyncedUp {
		s.state.LastCheck = nil
	}

	sQuery := s.cfg.Tables.Users.DBSyncTable
	query := prepareStringQuery(sQuery, s.state, offset, limit)

	users := []*users.User{}
	if _, err := qrm.Query(ctx, s.db, query, []any{}, &users); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return fmt.Errorf("failed to query users table. %w", err)
		}
	}

	s.logger.Debug("usersSync", zap.Any("users", users))

	if len(users) == 0 {
		s.state.Set(0, nil)
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
				errs = multierr.Append(errs, fmt.Errorf("failed to retrieve users %s licenses. %w", identifier, err))
			}
		}

		if errs != nil {
			return errs
		}
	}

	for k := range users {
		// Value mapping logic
		if s.cfg.Tables.Users.ValueMapping != nil {
			if users[k].Sex != nil && !s.cfg.Tables.Users.ValueMapping.Sex.IsEmpty() {
				s.cfg.Tables.Users.ValueMapping.Sex.Process(users[k].Sex)
			}
		}

		// Split names if only one field is used by the source data structure and only if we get 2 names out of it
		if s.cfg.Tables.Users.SplitName {
			if users[k].Lastname == "" {
				ss := strings.Split(users[k].Firstname, " ")
				if len(ss) > 1 {
					users[k].Lastname = ss[len(ss)-1]

					users[k].Firstname = strings.Replace(users[k].Firstname, " "+users[k].Lastname, "", 1)
				}
			}
		}

		// Attempt to parse date of birth via list of input formats
		for _, format := range s.cfg.Tables.Users.DateOfBirth.Formats {
			parsedTime, err := time.Parse(format, users[k].Dateofbirth)
			if err != nil {
				continue
			}

			// Format dates to the output format so all are the same if parseable
			users[k].Dateofbirth = parsedTime.Format(s.cfg.Tables.Users.DateOfBirth.OutputFormat)
			break
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
			return fmt.Errorf("failed to send users data to server. %w", err)
		}
	}

	// If less users than limit are returned, we probably have reached the "end" of the table
	// and need to reset the offset to 0. That means we are "synced up" and can start the normal
	// sync loop of checking the "updatedAt" date.
	if len(users) < limit {
		offset = 0
		s.state.SyncedUp = true
	}

	lastUserId := strconv.Itoa(int(users[len(users)-1].UserId))
	s.state.Set(uint64(limit)+offset, &lastUserId)

	return nil
}

func (s *usersSync) retrieveLicenses(ctx context.Context, userId int32, identifier string) ([]*users.License, error) {
	sQuery := s.cfg.Tables.UserLicenses
	query := prepareStringQuery(sQuery, s.state, 0, 100)

	args := []any{}
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
	sQuery := s.cfg.Tables.Users.DBSyncTable
	query := prepareStringQuery(sQuery, s.state, 0, 1)

	user := &users.User{}
	if _, err := qrm.Query(ctx, s.db, query, []any{}, &user); err != nil {
		return err
	}

	if s.cfg.Tables.UserLicenses.Enabled {
		// Retrieve user's licenses
		identifier := ""
		if user.Identifier != nil {
			identifier = *user.Identifier
		}

		var err error
		user.Licenses, err = s.retrieveLicenses(ctx, user.UserId, identifier)
		if err != nil {
			return err
		}
	}

	// Split names if only one field is used by the framework
	if s.cfg.Tables.Users.SplitName {
		if user.Lastname == "" {
			ss := strings.Split(user.Firstname, " ")
			user.Lastname = ss[len(ss)-1]

			user.Firstname = strings.Replace(user.Firstname, " "+user.Lastname, "", 1)
		}
	}

	if s.cli != nil {
		if _, err := s.cli.SendData(ctx, &pbsync.SendDataRequest{
			Data: &pbsync.SendDataRequest_Users{
				Users: &sync.DataUsers{
					Users: []*users.User{user},
				},
			},
		}); err != nil {
			return err
		}
	}

	return nil
}
