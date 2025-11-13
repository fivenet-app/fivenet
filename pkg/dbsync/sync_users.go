package dbsync

import (
	"context"
	"errors"
	"fmt"
	"slices"
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

	limit := int64(500)
	offset := s.getInitialOffset()
	s.logger.Debug("usersSync", zap.Int64("offset", offset))

	s.resetLastCheckIfNotSynced()

	sQuery := s.cfg.Tables.Users
	q := sQuery.GetQuery(s.state, offset, limit)

	us, err := s.fetchUsers(ctx, q)
	if err != nil {
		return err
	}

	s.logger.Debug("usersSync", zap.Int("len", len(us)))

	if len(us) == 0 {
		s.logger.Debug("no users found to sync, resetting state offset")
		s.state.Set(0, nil)
		return nil
	}

	offset, err = s.updateSyncState(us, offset, limit)
	if err != nil {
		return err
	}

	s.applyFiltersAndTransformations(us, sQuery)

	if err := s.retrieveAndAttachLicenses(ctx, us); err != nil {
		return err
	}

	if err := s.sendData(ctx, &pbsync.SendDataRequest{
		Data: &pbsync.SendDataRequest_Users{
			Users: &sync.DataUsers{
				Users: us,
			},
		},
	}); err != nil {
		return err
	}

	s.logger.Debug("usersSync", zap.Bool("syncedUp", s.state.SyncedUp))

	lastUserId := strconv.FormatInt(int64(us[len(us)-1].GetUserId()), 10)
	s.state.Set(offset+limit, &lastUserId)

	return nil
}

func (s *usersSync) getInitialOffset() int64 {
	if s.state != nil && s.state.Offset > 0 {
		return s.state.Offset
	}
	return 0
}

func (s *usersSync) resetLastCheckIfNotSynced() {
	if !s.state.SyncedUp {
		s.state.LastCheck = nil
	}
}

func (s *usersSync) fetchUsers(ctx context.Context, query string) ([]*users.User, error) {
	us := []*users.User{}
	if _, err := qrm.Query(ctx, s.db, query, []any{}, &us); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, fmt.Errorf("failed to query users table. %w", err)
		}
	}
	return us, nil
}

func (s *usersSync) updateSyncState(us []*users.User, offset, limit int64) (int64, error) {
	if int64(len(us)) < limit {
		offset = 0
		s.state.SyncedUp = true
	}
	return offset, nil
}

func (s *usersSync) applyFiltersAndTransformations(us []*users.User, sQuery UsersTable) {
	if s.cfg.Tables.Users.IgnoreEmptyName {
		us = slices.DeleteFunc(us, func(in *users.User) bool {
			return in == nil || (in.GetFirstname() == "" && in.GetLastname() == "")
		})
	}

	hasFilters := len(sQuery.Filters.Jobs) > 0

	for k := 0; k < len(us); {
		if s.cfg.Tables.Users.ValueMapping != nil {
			s.applyValueMapping(us[k])
		}

		if hasFilters {
			if s.applyFilters(us, k, sQuery) {
				continue
			}
		}

		s.splitNamesIfRequired(us[k])
		s.parseDateOfBirth(us[k])
		k++
	}
}

func (s *usersSync) applyValueMapping(user *users.User) {
	if user.Sex != nil && !s.cfg.Tables.Users.ValueMapping.Sex.IsEmpty() {
		s.cfg.Tables.Users.ValueMapping.Sex.Process(user.Sex)
	}
}

func (s *usersSync) applyFilters(us []*users.User, k int, sQuery UsersTable) bool {
	for _, filter := range sQuery.Filters.Jobs {
		if filter.compiledPattern.MatchString(us[k].GetJob()) {
			switch filter.Action {
			case FilterActionDrop:
				us = append(us[:k], us[k+1:]...)
				return true
			case FilterActionReplace:
				us[k].Job = filter.compiledPattern.ReplaceAllString(
					us[k].GetJob(),
					filter.Replacement,
				)
			default:
				s.logger.Warn("unknown filter action", zap.String("action", string(filter.Action)))
			}
		}
	}
	return false
}

func (s *usersSync) splitNamesIfRequired(user *users.User) {
	if s.cfg.Tables.Users.SplitName && user.GetLastname() == "" {
		ss := strings.Split(user.GetFirstname(), " ")
		if len(ss) > 1 {
			user.Lastname = ss[len(ss)-1]
			user.Firstname = strings.Replace(user.GetFirstname(), " "+user.GetLastname(), "", 1)
		}
	}
}

func (s *usersSync) parseDateOfBirth(user *users.User) {
	for _, format := range s.cfg.Tables.Users.DateOfBirth.Formats {
		parsedTime, err := time.Parse(format, user.GetDateofbirth())
		if err == nil {
			user.Dateofbirth = parsedTime.Format(s.cfg.Tables.Users.DateOfBirth.OutputFormat)
			break
		}
	}
}

func (s *usersSync) retrieveAndAttachLicenses(ctx context.Context, us []*users.User) error {
	if !s.cfg.Tables.CitizensLicenses.Enabled {
		return nil
	}

	errs := multierr.Combine()
	for k := range us {
		identifier := ""
		if us[k].Identifier != nil {
			identifier = us[k].GetIdentifier()
		}

		licenses, err := s.retrieveLicenses(ctx, us[k].GetUserId(), identifier)
		if err != nil {
			errs = multierr.Append(
				errs,
				fmt.Errorf("failed to retrieve users %s licenses. %w", identifier, err),
			)
		}
		us[k].Licenses = licenses
	}

	return errs
}

func (s *usersSync) retrieveLicenses(
	ctx context.Context,
	userId int32,
	identifier string,
) ([]*users.License, error) {
	q := s.cfg.Tables.CitizensLicenses.GetQuery(s.state, 0, 100)

	args := []any{}
	if strings.Contains(q, "$userId") {
		q = strings.ReplaceAll(q, "$userId", strconv.FormatInt(int64(userId), 10))
		args = append(args, userId)
	} else if strings.Contains(q, "$identifier") {
		q = strings.ReplaceAll(q, "$identifier", identifier)
		args = append(args, identifier)
	}

	licenses := []*users.License{}
	if _, err := qrm.Query(ctx, s.db, q, args, &licenses); err != nil {
		return nil, err
	}

	return licenses, nil
}

// Sync an individual user's info.
func (s *usersSync) SyncUser(ctx context.Context, userId int32) error {
	wheres := []string{}
	if userId != 0 {
		wheres = append(wheres, fmt.Sprintf("`%s` = %d", s.cfg.Tables.Users.Columns.ID, userId))
	}
	q := s.cfg.Tables.Users.GetQuery(s.state, 0, 1, wheres...)

	user := &users.User{}
	if _, err := qrm.Query(ctx, s.db, q, []any{}, &user); err != nil {
		return err
	}

	if s.cfg.Tables.CitizensLicenses.Enabled {
		// Retrieve user's licenses
		identifier := ""
		if user.Identifier != nil {
			identifier = user.GetIdentifier()
		}

		var err error
		user.Licenses, err = s.retrieveLicenses(ctx, user.GetUserId(), identifier)
		if err != nil {
			return err
		}
	}

	// Split names if only one field is used by the framework
	if s.cfg.Tables.Users.SplitName {
		if user.GetLastname() == "" {
			ss := strings.Split(user.GetFirstname(), " ")
			user.Lastname = ss[len(ss)-1]

			user.Firstname = strings.Replace(user.GetFirstname(), " "+user.GetLastname(), "", 1)
		}
	}

	if s.cli != nil {
		if err := s.sendData(ctx, &pbsync.SendDataRequest{
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
