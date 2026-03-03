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
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/cache"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/protoutils"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/multierr"
	"go.uber.org/zap"
)

// userHashCacheTTL Cache user hash for 6 hours, which should be sufficient to cover the next few sync cycles and avoid memory bloat from caching too many hashes for long periods.
const userHashCacheTTL = 6 * time.Hour

type UsersSync struct {
	*Syncer

	logger        *zap.Logger
	state         *dbsyncconfig.TableSyncState
	saveUpdatedAt bool

	hashes *cache.LRUCache[int32, uint64]
}

func NewUsersSync(s *Syncer, state *dbsyncconfig.TableSyncState, saveUpdatedAt bool) *UsersSync {
	var hashes *cache.LRUCache[int32, uint64]
	if saveUpdatedAt {
		// Cache up to 500 user hashes to avoid memory bloat, as this is only used to compare against
		// the most recent hash for each user and not all historical hashes
		hashes = cache.NewLRUCache[int32, uint64](500)

		// Ensure a sane last check value is set for the "update" sync to work immediately
		if state.GetLastCheck() == nil {
			initialLastCheck := time.Now().Add(-15 * time.Minute)
			state.SetLastCheck(&initialLastCheck)
		}
	}

	logger := s.logger.With(
		zap.String("syncer", "users"),
		zap.Bool("resync", !saveUpdatedAt),
	)

	return &UsersSync{
		Syncer: s,

		logger:        logger,
		state:         state,
		saveUpdatedAt: saveUpdatedAt,

		hashes: hashes,
	}
}

func (s *UsersSync) Sync(ctx context.Context) (int64, string, *time.Time, error) {
	limit := s.cfg.Limits.Users

	var total int64
	lastID := "0"
	var lastUpdatedAt *time.Time

	for {
		fetched, sent, cursorID, cursorTime, err := s.syncOnce(ctx)
		if err != nil {
			return total, lastID, lastUpdatedAt, err
		}

		total += sent
		if cursorID != "" {
			lastID = cursorID
		}
		if cursorTime != nil {
			lastUpdatedAt = cursorTime
		}

		// Nothing left for this cycle.
		if fetched < limit {
			break
		}
	}

	return total, lastID, lastUpdatedAt, nil
}

func (s *UsersSync) Resync(ctx context.Context) (int64, string, *time.Time, error) {
	// Full resync mode paginates only by user id.
	if !s.saveUpdatedAt {
		s.state.SetLastCheck(nil)
	}

	_, sent, cursorID, cursorTime, err := s.syncOnce(ctx)
	return sent, cursorID, cursorTime, err
}

func (s *UsersSync) syncOnce(
	ctx context.Context,
) (fetched int64, sent int64, cursorID string, cursorTime *time.Time, err error) {
	limit := s.cfg.Limits.Users
	sQuery := s.cfg.Tables.Users
	q := sQuery.GetQuery(s.state, 0, limit)

	us, err := s.fetchUsers(ctx, q)
	if err != nil {
		return 0, 0, "0", nil, err
	}

	fetchedCount := int64(len(us))
	s.logger.Debug("usersSync", zap.Int64("len", fetchedCount))
	if len(us) == 0 {
		if !s.saveUpdatedAt {
			s.logger.Debug("no users found during full resync, resetting cursor")
			s.state.ResetCursor()
		}
		return 0, 0, "0", nil, nil
	}

	cursorTime, cursorLastID := s.cursorFromUsersResults(us)
	if s.saveUpdatedAt && cursorTime == nil {
		return 0, 0, "0", nil, errors.New(
			"users result is missing updated_at, cannot persist cursor timestamp",
		)
	}
	cursorIDValue := ""
	if cursorLastID != nil {
		cursorIDValue = *cursorLastID
	}

	us = s.applyFiltersAndTransformations(us, sQuery)

	if err := s.retrieveAndAttachLicenses(ctx, us); err != nil {
		return 0, 0, "", nil, err
	}
	if err := s.retrieveAndAttachJobs(ctx, us); err != nil {
		return 0, 0, "", nil, err
	}
	if err := s.retrieveAndAttachPhoneNumbers(ctx, us); err != nil {
		return 0, 0, "", nil, err
	}

	if s.hashes != nil {
		for i, u := range slices.Backward(us) {
			// Get hash of user data to compare with existing hash and skip sending if data is the same (treat as not updated)
			_, hash, err := protoutils.JSONAndHash(u)
			if err != nil {
				s.logger.Warn(
					"failed to compute user data hash, skipping hash check and treating as new/updated user",
					zap.Int32("user_id", u.GetUserId()),
					zap.String("identifier", u.GetIdentifier()),
					zap.Error(err),
				)
			}

			if existingHash, ok := s.hashes.Get(u.GetUserId()); ok {
				if existingHash == hash {
					s.logger.Debug(
						"user data hash is the same as existing entry, skipping update for user",
						zap.Int32("user_id", u.GetUserId()),
						zap.String("identifier", u.GetIdentifier()),
					)
					// Remove "skipped" user
					us = slices.Delete(us, i, i+1)
					continue
				}
			} else {
				s.hashes.Put(u.GetUserId(), hash, userHashCacheTTL)
			}
		}
	}

	// Sync users to FiveNet server (if there are any left after hash check)
	if len(us) > 0 {
		if err := s.sendData(ctx, &pbsync.SendDataRequest{
			Data: &pbsync.SendDataRequest_Users{
				Users: &syncdata.DataUsers{
					Users: us,
				},
			},
		}); err != nil {
			return 0, 0, "", nil, err
		}
	}

	s.persistCursor(fetchedCount, limit, cursorTime, cursorLastID)

	return fetchedCount, int64(len(us)), cursorIDValue, cursorTime, nil
}

func (s *UsersSync) cursorFromUsersResults(
	us []*syncdata.DataUser,
) (*time.Time, *string) {
	if len(us) == 0 {
		return nil, nil
	}

	last := us[len(us)-1]
	lastID := strconv.FormatInt(int64(last.GetUserId()), 10)

	ts := last.GetUpdatedAt()
	if ts == nil || ts.GetTimestamp() == nil {
		return nil, &lastID
	}

	t := ts.GetTimestamp().AsTime()
	return &t, &lastID
}

func (s *UsersSync) persistCursor(
	fetchedCount int64,
	limit int64,
	cursorTime *time.Time,
	lastID *string,
) {
	if s.saveUpdatedAt {
		s.state.SetCursor(cursorTime, lastID)
		return
	}

	// Full resync mode restarts from beginning after reaching table end.
	if fetchedCount < limit {
		s.state.ResetCursor()
		return
	}

	s.state.SetCursor(nil, lastID)
}

func (s *UsersSync) fetchUsers(ctx context.Context, query string) ([]*syncdata.DataUser, error) {
	s.logger.Debug("users sync query", zap.String("query", query))

	us := []*syncdata.DataUser{}
	if _, err := qrm.Query(ctx, s.db, query, []any{}, &us); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, fmt.Errorf("failed to query users table. %w", err)
		}
	}

	return us, nil
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
	for idx, user := range slices.Backward(us) {
		if user.GetUserId() <= 0 {
			foundNullUserId = true
			s.logger.Debug(
				"user with null/zero id found",
				zap.String("identifier", user.GetIdentifier()),
			)
			continue
		}

		if s.cfg.Tables.Users.ValueMapping != nil {
			s.applyValueMapping(user)
		}

		if hasFilters {
			if s.applyFilters(user, sQuery) {
				// Remove "skipped" user
				us = slices.Delete(us, idx, idx+1)
				continue
			}
		}

		s.splitNamesIfRequired(user)
		s.parseDateOfBirth(user)
		s.cleanupUserJob(user)
		s.cleanupUserPhoneNumbers(user)
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

		return
	}

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

func (s *UsersSync) cleanupUserPhoneNumbers(user *syncdata.DataUser) {
	if len(user.GetPhoneNumbers()) == 0 && user.GetPhoneNumber() == "" {
		return
	}

	// If no phone numbers are set, add the user's phone number field as the primary phone number (if not empty)
	if len(user.GetPhoneNumbers()) == 0 && user.GetPhoneNumber() != "" {
		user.PhoneNumbers = []*users.PhoneNumber{
			{
				Number:    user.GetPhoneNumber(),
				IsPrimary: true,
			},
		}
		return
	} else {
		primaryNumber := user.GetPhoneNumbers()[0].GetNumber()
		user.PhoneNumber = &primaryNumber
	}

	// Sort the user's phone numbers by is primary and then alphabetically to ensure consistent order
	slices.SortFunc(user.GetPhoneNumbers(), func(a *users.PhoneNumber, b *users.PhoneNumber) int {
		if a.GetIsPrimary() && !b.GetIsPrimary() {
			return -1
		}
		if !a.GetIsPrimary() && b.GetIsPrimary() {
			return 1
		}
		return strings.Compare(a.GetNumber(), b.GetNumber())
	})

	foundPrimary := false
	primaryNumber := user.GetPhoneNumber()
	for _, number := range user.GetPhoneNumbers() {
		if number.GetNumber() == primaryNumber {
			// Make sure the "primary" phone number (user's phone number field if set) is marked as primary
			foundPrimary = true
			number.IsPrimary = true
		} else {
			number.IsPrimary = false
		}
	}

	// If not ensure user has at least one primary phone number set
	if !foundPrimary {
		user.PhoneNumbers[0].IsPrimary = true
	}
}

func (s *UsersSync) applyValueMapping(user *syncdata.DataUser) {
	if !s.cfg.Tables.Users.ValueMapping.Sex.IsEmpty() {
		if sex := s.cfg.Tables.Users.ValueMapping.Sex.Process(user.GetSex()); sex != "" {
			user.Sex = &sex
		}
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
	q := s.cfg.Tables.Users.GetQuery(nil, 0, 1, wheres...)
	s.logger.Debug("users sync query", zap.String("query", q))

	user := &syncdata.DataUser{}
	if _, err := qrm.Query(ctx, s.db, q, []any{}, user); err != nil {
		return fmt.Errorf("failed to query single user %d. %w", userId, err)
	}
	us := []*syncdata.DataUser{user}

	if err := s.retrieveAndAttachJobs(ctx, us); err != nil {
		return fmt.Errorf("failed to retrieve and attach jobs for user %d. %w", userId, err)
	}
	if err := s.retrieveAndAttachLicenses(ctx, us); err != nil {
		return fmt.Errorf("failed to retrieve and attach licenses for user %d. %w", userId, err)
	}
	if err := s.retrieveAndAttachPhoneNumbers(ctx, us); err != nil {
		return fmt.Errorf(
			"failed to retrieve and attach phone numbers for user %d. %w",
			userId,
			err,
		)
	}

	s.splitNamesIfRequired(user)
	s.parseDateOfBirth(user)
	s.cleanupUserJob(user)
	s.cleanupUserPhoneNumbers(user)
	us = s.applyFiltersAndTransformations(us, s.cfg.Tables.Users)

	if len(us) > 0 && s.cli != nil {
		if err := s.sendData(ctx, &pbsync.SendDataRequest{
			Data: &pbsync.SendDataRequest_Users{
				Users: &syncdata.DataUsers{
					Users: us,
				},
			},
		}); err != nil {
			return fmt.Errorf("failed to send user data for user %d. %w", userId, err)
		}
	}

	s.logger.Debug(
		"sync single user data",
		zap.Int32("user_id", user.GetUserId()),
		zap.String("job", user.GetJob()),
		zap.Int32("job_grade", user.GetJobGrade()),
		zap.Int("jobs_len", len(user.GetJobs())),
	)

	return nil
}
