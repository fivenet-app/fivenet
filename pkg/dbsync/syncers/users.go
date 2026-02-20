package syncers

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	syncdata "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/sync/data"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users"
	userslicenses "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/licenses"
	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
	dbsyncconfig "github.com/fivenet-app/fivenet/v2026/pkg/dbsync/config"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

type UsersSync struct {
	*Syncer

	state         *dbsyncconfig.TableSyncState
	saveLastCheck bool
}

func NewUsersSync(s *Syncer, state *dbsyncconfig.TableSyncState, saveLastCheck bool) *UsersSync {
	return &UsersSync{
		Syncer: s,
		state:  state,

		saveLastCheck: saveLastCheck,
	}
}

func (s *UsersSync) Sync(ctx context.Context) (int64, int64, string, error) {
	// Ensure last check is nil when we don't want to save it
	if !s.saveLastCheck {
		s.state.SetLastCheck(nil)
	}

	limit := int64(150)
	offset := s.state.GetOffset()
	s.logger.Debug("usersSync", zap.Int64("offset", offset))

	s.resetLastCheckIfNotSynced()

	sQuery := s.cfg.Tables.Users
	q := sQuery.GetQuery(s.state, offset, limit)

	us, err := s.fetchUsers(ctx, q)
	if err != nil {
		return 0, offset, "0", err
	}

	count := int64(len(us))
	s.logger.Debug("usersSync", zap.Int64("len", count))
	if len(us) == 0 {
		s.logger.Debug("no users found to sync, resetting state offset")
		s.state.Set(0, nil)
		s.resetLastCheckIfNotSynced()
		return 0, offset, "0", nil
	}

	offset, err = s.updateSyncState(count, offset, limit)
	if err != nil {
		return 0, offset, "0", err
	}

	us = s.applyFiltersAndTransformations(us, sQuery)

	if err := s.retrieveAndAttachLicenses(ctx, us); err != nil {
		return 0, offset, "0", err
	}
	if err := s.retrieveAndAttachJobs(ctx, us); err != nil {
		return 0, offset, "0", err
	}
	if err := s.retrieveAndAttachPhoneNumbers(ctx, us); err != nil {
		return 0, offset, "0", err
	}

	if err := s.sendData(ctx, &pbsync.SendDataRequest{
		Data: &pbsync.SendDataRequest_Users{
			Users: &syncdata.DataUsers{
				Users: us,
			},
		},
	}); err != nil {
		return 0, offset, "0", err
	}

	s.logger.Debug("usersSync", zap.Bool("syncedUp", s.state.GetSyncedUp()))

	lastId := int64(us[count-1].GetUserId())
	lastUserId := strconv.FormatInt(lastId, 10)
	s.state.Set(offset+limit, &lastUserId)

	return count, offset, lastUserId, nil
}

func (s *UsersSync) resetLastCheckIfNotSynced() {
	if !s.state.GetSyncedUp() {
		s.state.SetLastCheck(nil)
	}
}

func (s *UsersSync) fetchUsers(ctx context.Context, query string) ([]*syncdata.DataUser, error) {
	s.logger.Debug("accounts sync query", zap.String("query", query))

	us := []*syncdata.DataUser{}
	if _, err := qrm.Query(ctx, s.db, query, []any{}, &us); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, fmt.Errorf("failed to query users table. %w", err)
		}
	}

	return us, nil
}

func (s *UsersSync) updateSyncState(usersCount int64, offset, limit int64) (int64, error) {
	if usersCount < limit {
		offset = 0
		s.state.SetSyncedUp(true)
	}
	return offset, nil
}

func (s *UsersSync) applyFiltersAndTransformations(
	us []*syncdata.DataUser,
	sQuery dbsyncconfig.UsersTable,
) []*syncdata.DataUser {
	if s.cfg.Tables.Users.IgnoreEmptyName {
		us = slices.DeleteFunc(us, func(in *syncdata.DataUser) bool {
			return in == nil || (in.GetFirstname() == "" && in.GetLastname() == "")
		})
	}

	hasFilters := len(sQuery.Filters.Jobs) > 0

	foundNullUserId := false
	for i, u := range slices.Backward(us) {
		if u.GetUserId() <= 0 {
			foundNullUserId = true
			s.logger.Debug(
				"user with null/zero id found",
				zap.String("identifier", u.GetIdentifier()),
			)
			continue
		}

		if s.cfg.Tables.Users.ValueMapping != nil {
			s.applyValueMapping(u)
		}

		if hasFilters {
			if s.applyFilters(u, sQuery) {
				// Remove "skipped" user
				us = slices.Delete(us, i, i+1)
				continue
			}
		}

		s.splitNamesIfRequired(u)
		s.parseDateOfBirth(u)
		s.cleanupUserJob(u)

		s.logger.Debug(
			"user data",
			zap.Int32("user_id", u.GetUserId()),
			zap.String("job", u.GetJob()),
			zap.Int32("job_grade", u.GetJobGrade()),
			zap.Int("jobs_len", len(u.GetJobs())),
		)
	}

	if foundNullUserId {
		s.logger.Warn(
			"some queried users have null or zero id, which have been skipped during processing",
		)
	}

	return us
}

func (s *UsersSync) cleanupUserJob(user *syncdata.DataUser) {
	if len(user.Jobs) == 0 {
		// If no jobs are set, create one from the user job field
		user.Jobs = []*users.UserJob{
			{
				Job:       user.GetJob(),
				Grade:     user.GetJobGrade(),
				IsPrimary: true,
			},
		}
	} else {
		// Sort the user's jobs by is primary and then alphabetically to ensure consistent order
		slices.SortFunc(user.GetJobs(), func(a *users.UserJob, b *users.UserJob) int {
			if a.GetIsPrimary() && !b.GetIsPrimary() {
				return -1
			}
			if !a.GetIsPrimary() && b.GetIsPrimary() {
				return 1
			}
			return strings.Compare(a.GetJob(), b.GetJob())
		})

		foundPrimary := false
		primaryJob := user.GetJob()
		for _, job := range user.GetJobs() {
			if job.GetJob() == primaryJob {
				// Make sure the "primary" job (user's job field if set) is marked as primary
				foundPrimary = true
				job.IsPrimary = true
			} else {
				job.IsPrimary = false
			}
		}

		// If not ensure user has at least one primary job set
		if !foundPrimary {
			user.Jobs[0].IsPrimary = true
		}
	}
}

func (s *UsersSync) applyValueMapping(user *syncdata.DataUser) {
	if user.Sex != nil && !s.cfg.Tables.Users.ValueMapping.Sex.IsEmpty() {
		s.cfg.Tables.Users.ValueMapping.Sex.Process(user.Sex)
	}
}

func (s *UsersSync) applyFilters(
	us *syncdata.DataUser,
	sQuery dbsyncconfig.UsersTable,
) bool {
	for _, filter := range sQuery.Filters.Jobs {
		if filter.CompiledPattern.MatchString(us.GetJob()) {
			switch filter.Action {
			case dbsyncconfig.FilterActionDrop:
				return true

			case dbsyncconfig.FilterActionReplace:
				us.Job = filter.CompiledPattern.ReplaceAllString(
					us.GetJob(),
					filter.Replacement,
				)

			default:
				s.logger.Warn("unknown filter action", zap.String("action", string(filter.Action)))
			}
		}
	}
	return false
}

func (s *UsersSync) splitNamesIfRequired(user *syncdata.DataUser) {
	if s.cfg.Tables.Users.SplitName && user.GetLastname() == "" {
		ss := strings.Split(user.GetFirstname(), " ")
		if len(ss) > 1 {
			user.Lastname = &ss[len(ss)-1]
			user.Firstname = strings.Replace(user.GetFirstname(), " "+user.GetLastname(), "", 1)
		}
	}
}

func (s *UsersSync) parseDateOfBirth(user *syncdata.DataUser) {
	for _, format := range s.cfg.Tables.Users.DateOfBirth.Formats {
		parsedTime, err := time.Parse(format, user.GetDateofbirth())
		if err == nil {
			user.Dateofbirth = parsedTime.Format(s.cfg.Tables.Users.DateOfBirth.OutputFormat)
			break
		}
	}
}

func (s *UsersSync) retrieveAndAttachLicenses(ctx context.Context, us []*syncdata.DataUser) error {
	if !s.cfg.Tables.UserLicenses.Enabled {
		return nil
	}

	errs := multierr.Combine()
	for k := range us {
		licenses, err := s.retrieveLicenses(ctx, us[k].GetUserId(), us[k].GetIdentifier())
		if err != nil {
			errs = multierr.Append(
				errs,
				fmt.Errorf(
					"failed to retrieve user %d (%s) licenses. %w",
					us[k].GetUserId(),
					us[k].GetIdentifier(),
					err,
				),
			)
		}
		us[k].Licenses = licenses
	}

	return errs
}

func (s *UsersSync) retrieveLicenses(
	ctx context.Context,
	userId int32,
	identifier string,
) ([]*userslicenses.License, error) {
	q := s.cfg.Tables.UserLicenses.GetQuery(0, 100)
	s.logger.Debug("users licenses sync query", zap.String("query", q))

	args := []any{}
	if strings.Contains(q, "$userId") {
		count := strings.Count(q, "$userId")
		q = strings.ReplaceAll(q, "$userId", "?")
		for range count {
			args = append(args, userId)
		}
	} else if strings.Contains(q, "$identifier") {
		count := strings.Count(q, "$identifier")
		q = strings.ReplaceAll(q, "$identifier", "?")
		for range count {
			args = append(args, identifier)
		}
	}

	licenses := []*userslicenses.License{}
	if _, err := qrm.Query(ctx, s.db, q, args, &licenses); err != nil {
		return nil, err
	}

	return licenses, nil
}

func (s *UsersSync) retrieveAndAttachJobs(ctx context.Context, us []*syncdata.DataUser) error {
	if !s.cfg.Tables.UserJobs.Enabled {
		return nil
	}

	errs := multierr.Combine()
	for k := range us {
		jobs, err := s.retrieveJobs(ctx, us[k].GetUserId(), us[k].GetIdentifier())
		if err != nil {
			errs = multierr.Append(
				errs,
				fmt.Errorf(
					"failed to retrieve users %d (%s) jobs. %w",
					us[k].GetUserId(),
					us[k].GetIdentifier(),
					err,
				),
			)
		}
		us[k].Jobs = jobs
	}

	return errs
}

func (s *UsersSync) retrieveJobs(
	ctx context.Context,
	userId int32,
	identifier string,
) ([]*users.UserJob, error) {
	q := s.cfg.Tables.UserJobs.GetQuery(0, 10)
	s.logger.Debug("users jobs sync query", zap.String("query", q))

	args := []any{}
	if strings.Contains(q, "$userId") {
		count := strings.Count(q, "$userId")
		q = strings.ReplaceAll(q, "$userId", "?")
		for range count {
			args = append(args, userId)
		}
	} else if strings.Contains(q, "$identifier") {
		count := strings.Count(q, "$identifier")
		q = strings.ReplaceAll(q, "$identifier", "?")
		for range count {
			args = append(args, identifier)
		}
	}

	jobs := []*users.UserJob{}
	if _, err := qrm.Query(ctx, s.db, q, args, &jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}

func (s *UsersSync) retrieveAndAttachPhoneNumbers(
	ctx context.Context,
	us []*syncdata.DataUser,
) error {
	if !s.cfg.Tables.UserPhoneNumbers.Enabled {
		return nil
	}

	errs := multierr.Combine()
	for k := range us {
		phoneNumbers, err := s.retrievePhoneNumbers(ctx, us[k].GetUserId(), us[k].GetIdentifier())
		if err != nil {
			errs = multierr.Append(
				errs,
				fmt.Errorf(
					"failed to retrieve users %d (%s) jobs. %w",
					us[k].GetUserId(),
					us[k].GetIdentifier(),
					err,
				),
			)
		}
		us[k].PhoneNumbers = phoneNumbers
	}

	return errs
}

func (s *UsersSync) retrievePhoneNumbers(
	ctx context.Context,
	userId int32,
	identifier string,
) ([]*users.PhoneNumber, error) {
	q := s.cfg.Tables.UserPhoneNumbers.GetQuery(0, 10)
	s.logger.Debug("users phone numbers sync query", zap.String("query", q))

	args := []any{}
	if strings.Contains(q, "$userId") {
		count := strings.Count(q, "$userId")
		q = strings.ReplaceAll(q, "$userId", "?")
		for range count {
			args = append(args, userId)
		}
	} else if strings.Contains(q, "$identifier") {
		count := strings.Count(q, "$identifier")
		q = strings.ReplaceAll(q, "$identifier", "?")
		for range count {
			args = append(args, identifier)
		}
	}

	phoneNumbers := []*users.PhoneNumber{}
	if _, err := qrm.Query(ctx, s.db, q, args, &phoneNumbers); err != nil {
		return nil, err
	}

	return phoneNumbers, nil
}

// Sync an individual user's info.
func (s *UsersSync) SyncUser(ctx context.Context, userId int32) error {
	wheres := []string{}
	if userId != 0 {
		wheres = append(wheres, fmt.Sprintf("%#q = %d", s.cfg.Tables.Users.Columns.ID, userId))
	}
	q := s.cfg.Tables.Users.GetQuery(s.state, 0, 1, wheres...)
	s.logger.Debug("users sync query", zap.String("query", q))

	user := &syncdata.DataUser{}
	if _, err := qrm.Query(ctx, s.db, q, []any{}, user); err != nil {
		return err
	}
	us := []*syncdata.DataUser{user}

	if err := s.retrieveAndAttachJobs(ctx, us); err != nil {
		return err
	}
	if err := s.retrieveAndAttachLicenses(ctx, us); err != nil {
		return err
	}
	if err := s.retrieveAndAttachPhoneNumbers(ctx, us); err != nil {
		return err
	}

	s.splitNamesIfRequired(user)
	s.parseDateOfBirth(user)
	s.cleanupUserJob(user)
	us = s.applyFiltersAndTransformations(us, s.cfg.Tables.Users)

	if len(us) > 0 && s.cli != nil {
		if err := s.sendData(ctx, &pbsync.SendDataRequest{
			Data: &pbsync.SendDataRequest_Users{
				Users: &syncdata.DataUsers{
					Users: us,
				},
			},
		}); err != nil {
			return err
		}
	}

	s.logger.Debug(
		"user data",
		zap.Int32("user_id", user.GetUserId()),
		zap.String("job", user.GetJob()),
		zap.Int32("job_grade", user.GetJobGrade()),
		zap.Int("jobs_len", len(user.GetJobs())),
	)

	return nil
}
