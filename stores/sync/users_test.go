package syncstore

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/notifications"
	notificationsclientview "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/notifications/clientview"
	notificationsevents "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/notifications/events"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/settings"
	syncdata "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/sync/data"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/config/appconfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type labelEnricher struct{}

func (labelEnricher) EnrichJobInfo(user common.IJobInfo) {
	user.SetJobLabel(user.GetJob())
	user.SetJobGradeLabel(fmt.Sprintf("Rank %d", user.GetJobGrade()))
}

func (labelEnricher) EnrichJobInfoNoFallback(common.IJobInfo) {}

func (labelEnricher) EnrichJobName(common.IJobName) {}

func (labelEnricher) GetJobByName(string) *jobs.Job { return nil }

func (labelEnricher) GetJobGrade(string, int32) (*jobs.Job, *jobs.JobGrade) {
	return nil, nil
}

type recordingNotifi struct {
	events []*notificationsevents.UserEvent
}

func (n *recordingNotifi) NotifyUser(context.Context, *notifications.Notification) error {
	return nil
}

func (n *recordingNotifi) SendObjectEvent(
	context.Context,
	*notificationsclientview.ObjectEvent,
) error {
	return nil
}

func (n *recordingNotifi) SendAccountEvent(
	_ context.Context,
	_ int64,
	event *notificationsevents.UserEvent,
) error {
	n.events = append(n.events, event)
	return nil
}

func (n *recordingNotifi) SendUserEvent(
	_ context.Context,
	_ int32,
	event *notificationsevents.UserEvent,
) error {
	n.events = append(n.events, event)
	return nil
}

func (n *recordingNotifi) SendSystemEvent(context.Context, *notificationsevents.SystemEvent) error {
	return nil
}

func TestCompareJobs(t *testing.T) {
	t.Parallel()
	job := func(name string, grade int32, primary bool) *users.UserJob {
		return &users.UserJob{Job: name, Grade: grade, IsPrimary: primary}
	}

	jobNames := func(jobs []*users.UserJob) []string {
		names := make([]string, 0, len(jobs))
		for _, j := range jobs {
			names = append(names, j.GetJob())
		}
		return names
	}

	jobMap := func(jobs []*users.UserJob) map[string]*users.UserJob {
		m := make(map[string]*users.UserJob, len(jobs))
		for _, j := range jobs {
			m[j.GetJob()] = j
		}
		return m
	}

	tests := []struct {
		name     string
		current  []*users.UserJob
		incoming []*users.UserJob
		add      []string
		update   []string
		remove   []string
	}{
		{
			name:     "add new jobs when none exist",
			current:  nil,
			incoming: []*users.UserJob{job("police", 3, true), job("ems", 1, false)},
			add:      []string{"police", "ems"},
			update:   nil,
			remove:   nil,
		},
		{
			name:     "update when grade changes",
			current:  []*users.UserJob{job("police", 1, true)},
			incoming: []*users.UserJob{job("police", 2, true)},
			add:      nil,
			update:   []string{"police"},
			remove:   nil,
		},
		{
			name:     "update when primary flag changes",
			current:  []*users.UserJob{job("ems", 1, false)},
			incoming: []*users.UserJob{job("ems", 1, true)},
			add:      nil,
			update:   []string{"ems"},
			remove:   nil,
		},
		{
			name:     "remove missing jobs",
			current:  []*users.UserJob{job("police", 3, true)},
			incoming: []*users.UserJob{},
			add:      nil,
			update:   nil,
			remove:   []string{"police"},
		},
		{
			name:     "mixed add update remove",
			current:  []*users.UserJob{job("police", 2, true), job("ems", 1, false)},
			incoming: []*users.UserJob{job("police", 3, true), job("fire", 1, false)},
			add:      []string{"fire"},
			update:   []string{"police"},
			remove:   []string{"ems"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			toAdd, toUpdate, toRemove := compareJobs(tt.current, tt.incoming)

			assert.ElementsMatch(t, tt.add, jobNames(toAdd))
			assert.ElementsMatch(t, tt.update, jobNames(toUpdate))
			assert.ElementsMatch(t, tt.remove, jobNames(toRemove))

			incomingByName := jobMap(tt.incoming)
			currentByName := jobMap(tt.current)

			for _, j := range toAdd {
				assert.Same(t, incomingByName[j.GetJob()], j)
			}

			for _, j := range toUpdate {
				assert.Same(t, incomingByName[j.GetJob()], j)
			}

			for _, j := range toRemove {
				assert.Same(t, currentByName[j.GetJob()], j)
			}
		})
	}
}

func TestComparePhoneNumbers(t *testing.T) {
	t.Parallel()
	phone := func(number string, primary bool) *users.PhoneNumber {
		return &users.PhoneNumber{Number: number, IsPrimary: primary}
	}

	numbers := func(list []*users.PhoneNumber) []string {
		res := make([]string, 0, len(list))
		for _, p := range list {
			res = append(res, p.GetNumber())
		}
		return res
	}

	incomingMap := func(list []*users.PhoneNumber) map[string]*users.PhoneNumber {
		m := make(map[string]*users.PhoneNumber, len(list))
		for _, p := range list {
			m[p.GetNumber()] = p
		}
		return m
	}

	currentMap := func(list []*users.PhoneNumber) map[string]*users.PhoneNumber {
		m := make(map[string]*users.PhoneNumber, len(list))
		for _, p := range list {
			m[p.GetNumber()] = p
		}
		return m
	}

	tests := []struct {
		name     string
		current  []*users.PhoneNumber
		incoming []*users.PhoneNumber
		add      []string
		update   []string
		remove   []string
	}{
		{
			name:     "add new numbers when none exist",
			current:  nil,
			incoming: []*users.PhoneNumber{phone("111", true), phone("222", false)},
			add:      []string{"111", "222"},
			update:   nil,
			remove:   nil,
		},
		{
			name:     "update primary flag on existing number",
			current:  []*users.PhoneNumber{phone("333", false)},
			incoming: []*users.PhoneNumber{phone("333", true)},
			add:      nil,
			update:   []string{"333"},
			remove:   nil,
		},
		{
			name:     "remove missing numbers",
			current:  []*users.PhoneNumber{phone("444", false), phone("555", false)},
			incoming: []*users.PhoneNumber{phone("444", false)},
			add:      nil,
			update:   nil,
			remove:   []string{"555"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			toAdd, toUpdate, toRemove := comparePhoneNumbers(tt.current, tt.incoming)

			assert.ElementsMatch(t, tt.add, numbers(toAdd))
			assert.ElementsMatch(t, tt.update, numbers(toUpdate))
			assert.ElementsMatch(t, tt.remove, numbers(toRemove))

			incomingByNumber := incomingMap(tt.incoming)
			currentByNumber := currentMap(tt.current)

			for _, p := range toAdd {
				assert.Same(t, incomingByNumber[p.GetNumber()], p)
			}

			for _, p := range toUpdate {
				assert.Same(t, incomingByNumber[p.GetNumber()], p)
			}

			for _, p := range toRemove {
				assert.Same(t, currentByNumber[p.GetNumber()], p)
			}
		})
	}

	t.Run("switches primary and demotes old one", func(t *testing.T) {
		t.Parallel()
		currentPrimary := phone("666", true)
		currentOther := phone("777", false)
		current := []*users.PhoneNumber{currentPrimary, currentOther}

		incomingPrimary := phone("777", true)
		incomingOld := phone("666", false)
		incoming := []*users.PhoneNumber{incomingOld, incomingPrimary}

		toAdd, toUpdate, toRemove := comparePhoneNumbers(current, incoming)

		assert.Empty(t, toAdd)
		assert.Empty(t, toRemove)
		assert.ElementsMatch(t, []string{"666", "777", "666"}, numbers(toUpdate))

		incomingByNumber := incomingMap(incoming)

		seenCurrentDemotion := false
		seenIncomingNewPrimary := false
		for _, p := range toUpdate {
			if p == currentPrimary {
				seenCurrentDemotion = true
			}
			if p == incomingByNumber["777"] {
				seenIncomingNewPrimary = true
			}
		}

		assert.True(
			t,
			seenCurrentDemotion,
			"old primary should be demoted via current slice pointer",
		)
		assert.True(t, seenIncomingNewPrimary, "new primary should come from incoming slice")
		assert.False(
			t,
			currentPrimary.GetIsPrimary(),
			"old primary must be marked non-primary after compare",
		)
	})
}

func TestCleanupUserPhoneNumbersDefaultsToSinglePrimaryPhone(t *testing.T) {
	t.Parallel()

	store, _ := newTestStore(t)

	user := &syncdata.DataUser{
		PhoneNumber: func() *string { v := "111"; return &v }(),
	}

	store.cleanupUserPhoneNumbers(user)

	require.Len(t, user.GetPhoneNumbers(), 1)
	assert.Equal(t, "111", user.GetPhoneNumbers()[0].GetNumber())
	assert.True(t, user.GetPhoneNumbers()[0].GetIsPrimary())
}

func TestCleanupUserPhoneNumbersPrefersIncomingPrimaryFlag(t *testing.T) {
	t.Parallel()

	store, _ := newTestStore(t)

	user := &syncdata.DataUser{
		PhoneNumber: func() *string { v := "222"; return &v }(),
		PhoneNumbers: []*users.PhoneNumber{
			{Number: "222", IsPrimary: false},
			{Number: "111", IsPrimary: true},
		},
	}

	store.cleanupUserPhoneNumbers(user)

	require.Len(t, user.GetPhoneNumbers(), 2)
	assert.Equal(t, "111", user.GetPhoneNumber())
	assert.True(t, user.GetPhoneNumbers()[0].GetIsPrimary())
	assert.False(t, user.GetPhoneNumbers()[1].GetIsPrimary())
}

func TestCleanupUserPhoneNumbersPrefersIncomingPrimaryOverLegacy(t *testing.T) {
	t.Parallel()

	store, _ := newTestStore(t)

	user := &syncdata.DataUser{
		PhoneNumbers: []*users.PhoneNumber{
			{Number: "111", IsPrimary: false},
			{Number: "222", IsPrimary: true},
		},
	}

	store.cleanupUserPhoneNumbers(user)

	require.Len(t, user.GetPhoneNumbers(), 2)
	assert.Equal(t, "222", user.GetPhoneNumber())
	assert.True(t, user.GetPhoneNumbers()[0].GetIsPrimary())
	assert.False(t, user.GetPhoneNumbers()[1].GetIsPrimary())
}

func TestCleanupUserPhoneNumbersFallsBackToLegacyWhenNoPrimary(t *testing.T) {
	t.Parallel()

	store, _ := newTestStore(t)

	user := &syncdata.DataUser{
		PhoneNumber: func() *string { v := "222"; return &v }(),
		PhoneNumbers: []*users.PhoneNumber{
			{Number: "111", IsPrimary: false},
			{Number: "222", IsPrimary: false},
		},
	}

	store.cleanupUserPhoneNumbers(user)

	require.Len(t, user.GetPhoneNumbers(), 2)
	assert.Equal(t, "222", user.GetPhoneNumber())
	assert.False(t, user.GetPhoneNumbers()[0].GetIsPrimary())
	assert.True(t, user.GetPhoneNumbers()[1].GetIsPrimary())
}

func TestCleanupUserPhoneNumbersKeepsOnlyOnePrimary(t *testing.T) {
	t.Parallel()

	store, _ := newTestStore(t)

	user := &syncdata.DataUser{
		PhoneNumber: func() *string { v := "111"; return &v }(),
		PhoneNumbers: []*users.PhoneNumber{
			{Number: "111", IsPrimary: true},
			{Number: "222", IsPrimary: true},
			{Number: "333", IsPrimary: false},
		},
	}

	store.cleanupUserPhoneNumbers(user)

	require.Len(t, user.GetPhoneNumbers(), 3)
	assert.True(t, user.GetPhoneNumbers()[0].GetIsPrimary())
	assert.False(t, user.GetPhoneNumbers()[1].GetIsPrimary())
	assert.False(t, user.GetPhoneNumbers()[2].GetIsPrimary())
}

func newTestStore(t *testing.T) (*Store, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	store := New(db, zap.NewNop(), &config.Config{}, &appconfig.TestConfig{}, nil, nil, nil, nil, nil).(*Store)
	t.Cleanup(func() {
		_ = db.Close()
	})

	return store, mock
}

func TestHandleUserJobsPublishesPrimaryJobChange(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)
	rec := &recordingNotifi{}
	store.notifi = rec
	store.enricher = labelEnricher{}
	accountID := int64(42)

	mock.ExpectBegin()
	mock.ExpectQuery(`(?s)SELECT .*fivenet_user_jobs.*WHERE .*user_id = \?.*ORDER BY .*`).
		WithArgs(int64(11)).
		WillReturnRows(sqlmock.NewRows([]string{"user_job.job", "user_job.grade", "user_job.is_primary"}).
			AddRow("police", int32(1), true).
			AddRow("ems", int32(1), false))
	mock.ExpectExec(`(?s)INSERT INTO .*fivenet_user_jobs.*ON DUPLICATE KEY UPDATE.*`).
		WithArgs(int64(11), "sheriff", int32(1), true).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`(?s)DELETE FROM .*fivenet_user_jobs.*job IN \(\?\).*LIMIT \?.*`).
		WithArgs(int64(11), "police", int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	user := &syncdata.DataUser{
		UserId:   11,
		Job:      "sheriff",
		JobGrade: 1,
		Jobs: []*users.UserJob{
			{Job: "sheriff", Grade: 1, IsPrimary: true},
			{Job: "ems", Grade: 1, IsPrimary: false},
		},
		Firstname: "Test",
	}

	tx, err := store.db.Begin()
	require.NoError(t, err)

	jobChange, err := store.handleUserJobs(t.Context(), tx, user)
	require.NoError(t, err)
	require.NoError(t, tx.Commit())

	store.publishUserInfoChanged(t.Context(), &accountID, user.GetUserId(), jobChange)

	require.Len(t, rec.events, 1)
	evt := rec.events[0].GetUserInfoChanged()
	require.NotNil(t, evt)
	assert.Equal(t, int64(42), evt.GetAccountId())
	assert.Equal(t, int32(11), evt.GetUserId())
	assert.Equal(t, "sheriff", evt.GetNewJob())
	assert.Equal(t, "sheriff", evt.GetNewJobLabel())
	assert.Equal(t, int32(1), evt.GetNewJobGrade())
	assert.Equal(t, "Rank 1", evt.GetNewJobGradeLabel())
	require.NotNil(t, evt.GetChangedAt())

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestHandleUserJobsPublishesPrimaryGradeChange(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)
	rec := &recordingNotifi{}
	store.notifi = rec
	store.enricher = labelEnricher{}
	accountID := int64(42)

	mock.ExpectBegin()
	mock.ExpectQuery(`(?s)SELECT .*fivenet_user_jobs.*WHERE .*user_id = \?.*ORDER BY .*`).
		WithArgs(int64(11)).
		WillReturnRows(sqlmock.NewRows([]string{"user_job.job", "user_job.grade", "user_job.is_primary"}).
			AddRow("police", int32(1), true))
	mock.ExpectExec(`(?s)INSERT INTO .*fivenet_user_jobs.*ON DUPLICATE KEY UPDATE.*`).
		WithArgs(int64(11), "police", int32(2), true).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	user := &syncdata.DataUser{
		UserId:   11,
		Job:      "police",
		JobGrade: 2,
		Jobs:     []*users.UserJob{{Job: "police", Grade: 2, IsPrimary: true}},
	}

	tx, err := store.db.Begin()
	require.NoError(t, err)

	jobChange, err := store.handleUserJobs(t.Context(), tx, user)
	require.NoError(t, err)
	require.NoError(t, tx.Commit())

	store.publishUserInfoChanged(t.Context(), &accountID, user.GetUserId(), jobChange)

	require.Len(t, rec.events, 1)
	evt := rec.events[0].GetUserInfoChanged()
	require.NotNil(t, evt)
	assert.Equal(t, "police", evt.GetNewJob())
	assert.Equal(t, int32(2), evt.GetNewJobGrade())
	assert.Equal(t, "Rank 2", evt.GetNewJobGradeLabel())

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestHandleUserJobsIgnoresSecondaryJobChurn(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)
	rec := &recordingNotifi{}
	store.notifi = rec
	store.enricher = labelEnricher{}
	accountID := int64(42)

	mock.ExpectBegin()
	mock.ExpectQuery(`(?s)SELECT .*fivenet_user_jobs.*WHERE .*user_id = \?.*ORDER BY .*`).
		WithArgs(int64(11)).
		WillReturnRows(sqlmock.NewRows([]string{"user_job.job", "user_job.grade", "user_job.is_primary"}).
			AddRow("ems", int32(1), false).
			AddRow("police", int32(1), true))
	mock.ExpectExec(`(?s)INSERT INTO .*fivenet_user_jobs.*ON DUPLICATE KEY UPDATE.*`).
		WithArgs(int64(11), "medic", int32(1), false).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(`(?s)DELETE FROM .*fivenet_user_jobs.*job IN \(\?\).*LIMIT \?.*`).
		WithArgs(int64(11), "ems", int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	user := &syncdata.DataUser{
		UserId:   11,
		Job:      "police",
		JobGrade: 1,
		Jobs: []*users.UserJob{
			{Job: "police", Grade: 1, IsPrimary: true},
			{Job: "medic", Grade: 1, IsPrimary: false},
		},
	}

	tx, err := store.db.Begin()
	require.NoError(t, err)

	jobChange, err := store.handleUserJobs(t.Context(), tx, user)
	require.NoError(t, err)
	require.NoError(t, tx.Commit())

	store.publishUserInfoChanged(t.Context(), &accountID, user.GetUserId(), jobChange)

	assert.Empty(t, rec.events)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestHandleUserJobsNoopDoesNotPublish(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)
	rec := &recordingNotifi{}
	store.notifi = rec
	store.enricher = labelEnricher{}
	accountID := int64(42)

	mock.ExpectBegin()
	mock.ExpectQuery(`(?s)SELECT .*fivenet_user_jobs.*WHERE .*user_id = \?.*ORDER BY .*`).
		WithArgs(int64(11)).
		WillReturnRows(sqlmock.NewRows([]string{"user_job.job", "user_job.grade", "user_job.is_primary"}).
			AddRow("police", int32(1), true))
	mock.ExpectCommit()

	user := &syncdata.DataUser{
		UserId:   11,
		Job:      "police",
		JobGrade: 1,
		Jobs:     []*users.UserJob{{Job: "police", Grade: 1, IsPrimary: true}},
	}

	tx, err := store.db.Begin()
	require.NoError(t, err)

	jobChange, err := store.handleUserJobs(t.Context(), tx, user)
	require.NoError(t, err)
	require.NoError(t, tx.Commit())

	store.publishUserInfoChanged(t.Context(), &accountID, user.GetUserId(), jobChange)

	assert.Empty(t, rec.events)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestHandleUserJobsDefaultsToUnemployedWhenNoJobs(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)
	rec := &recordingNotifi{}
	store.notifi = rec
	store.enricher = labelEnricher{}
	accountID := int64(42)

	cfg := &appconfig.Cfg{}
	cfg.Default()
	cfg.JobInfo.UnemployedJob = &settings.UnemployedJob{
		Name:  "civilian",
		Grade: 7,
	}
	store.appCfg.Set(cfg)

	mock.ExpectBegin()
	mock.ExpectQuery(`(?s)SELECT .*fivenet_user_jobs.*WHERE .*user_id = \?.*ORDER BY .*`).
		WithArgs(int64(11)).
		WillReturnRows(sqlmock.NewRows([]string{"user_job.job", "user_job.grade", "user_job.is_primary"}))
	mock.ExpectExec(`(?s)INSERT INTO .*fivenet_user_jobs.*ON DUPLICATE KEY UPDATE.*`).
		WithArgs(int64(11), "civilian", int32(7), true).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	user := &syncdata.DataUser{
		UserId:   11,
		Job:      "civilian",
		JobGrade: 7,
		Jobs: []*users.UserJob{
			{
				Job:       "civilian",
				Grade:     7,
				IsPrimary: true,
			},
		},
	}

	tx, err := store.db.Begin()
	require.NoError(t, err)

	jobChange, err := store.handleUserJobs(t.Context(), tx, user)
	require.NoError(t, err)
	require.NoError(t, tx.Commit())

	store.publishUserInfoChanged(t.Context(), &accountID, user.GetUserId(), jobChange)

	require.Len(t, rec.events, 1)
	evt := rec.events[0].GetUserInfoChanged()
	require.NotNil(t, evt)
	assert.Equal(t, "civilian", evt.GetNewJob())
	assert.Equal(t, int32(7), evt.GetNewJobGrade())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSyncUserAccountUpsertsMapping(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT fivenet_accounts.id AS "id" FROM fivenet_accounts WHERE fivenet_accounts.license = ? LIMIT ?;`)).
		WithArgs("license-42", int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(42)))
	mock.ExpectExec(`(?s)INSERT INTO .*fivenet_user_accounts.*ON DUPLICATE KEY UPDATE.*account_id = .*VALUES.*`).
		WithArgs(int32(11), int64(42)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	tx, err := store.db.Begin()
	require.NoError(t, err)
	accountID, err := store.syncUserAccount(t.Context(), tx, 11, "char1:license-42")
	require.NoError(t, err)
	require.NotNil(t, accountID)
	require.Equal(t, int64(42), *accountID)
	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSyncUserAccountDeletesMappingWhenUnresolved(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT fivenet_accounts.id AS "id" FROM fivenet_accounts WHERE fivenet_accounts.license = ? LIMIT ?;`)).
		WithArgs("license-42", int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	mock.ExpectExec(
		`(?s)DELETE FROM .*fivenet_user_accounts.*user_id = \?.*`+
			`(?s).*`+regexp.QuoteMeta(`LIMIT ?;`),
	).
		WithArgs(int32(11), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	tx, err := store.db.Begin()
	require.NoError(t, err)
	accountID, err := store.syncUserAccount(t.Context(), tx, 11, "char1:license-42")
	require.NoError(t, err)
	require.Nil(t, accountID)
	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}
