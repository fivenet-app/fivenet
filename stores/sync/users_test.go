package syncstore

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

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

func newTestStore(t *testing.T) (*Store, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	store := New(db, zap.NewNop(), &config.Config{}, nil, nil, nil).(*Store)
	t.Cleanup(func() {
		_ = db.Close()
	})

	return store, mock
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
	require.NoError(t, store.syncUserAccount(t.Context(), tx, 11, "char1:license-42"))
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
	require.NoError(t, store.syncUserAccount(t.Context(), tx, 11, "char1:license-42"))
	require.NoError(t, tx.Commit())
	require.NoError(t, mock.ExpectationsWereMet())
}
