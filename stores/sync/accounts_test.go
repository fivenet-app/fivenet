package syncstore

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	accounts "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/accounts"
	activity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/sync/activity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleAccountUpdatePublishesGroupChange(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)
	rec := &recordingNotifi{}
	store.notifi = rec

	mock.ExpectBegin()
	mock.ExpectQuery(
		regexp.QuoteMeta(
			`SELECT fivenet_accounts.id AS "id", fivenet_accounts.license AS "license", fivenet_accounts.`+"`groups`"+` AS "groups" FROM fivenet_accounts WHERE fivenet_accounts.license = ? LIMIT ?;`,
		),
	).
		WithArgs("license-42", int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "license", "groups"}).
			AddRow(int64(42), "license-42", []byte(`["old"]`)))
	mock.ExpectExec(`(?s)UPDATE .*fivenet_accounts.*SET .*groups.*WHERE .*license = \?.*LIMIT \?.*`).
		WithArgs(sqlmock.AnyArg(), "license-42", int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := store.handleAccountUpdate(t.Context(), &activity.AccountUpdate{
		License: "license-42",
		Groups:  &accounts.AccountGroups{Groups: []string{"supporter", "donator"}},
	})
	require.NoError(t, err)

	require.Len(t, rec.events, 1)
	evt := rec.events[0].GetUserGroupsChanged()
	require.NotNil(t, evt)
	assert.Equal(t, int64(42), evt.GetAccountId())
	require.NotNil(t, evt.GetNewGroups())
	assert.Equal(t, []string{"supporter", "donator"}, evt.GetNewGroups().GetGroups())
	assert.False(t, evt.GetCanBeSuperuser())

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestLoadAccountGroupState(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	mock.ExpectBegin()
	mock.ExpectQuery(
		regexp.QuoteMeta(
			`SELECT fivenet_accounts.id AS "id", fivenet_accounts.license AS "license", fivenet_accounts.`+"`groups`"+` AS "groups" FROM fivenet_accounts WHERE fivenet_accounts.license = ? LIMIT ?;`,
		),
	).
		WithArgs("license-42", int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "license", "groups"}).
			AddRow(int64(42), "license-42", []byte(`["old"]`)))
	mock.ExpectRollback()

	tx, err := store.db.Begin()
	require.NoError(t, err)

	state, err := store.loadAccountGroupState(t.Context(), tx, "license-42")
	require.NoError(t, err)
	require.NotNil(t, state)
	assert.Equal(t, int64(42), state.ID)
	assert.Equal(t, []string{"old"}, state.Groups.GetGroups())

	require.NoError(t, tx.Rollback())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestHandleAccountUpdateNoopDoesNotPublish(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)
	rec := &recordingNotifi{}
	store.notifi = rec

	mock.ExpectBegin()
	mock.ExpectQuery(
		regexp.QuoteMeta(
			`SELECT fivenet_accounts.id AS "id", fivenet_accounts.license AS "license", fivenet_accounts.`+"`groups`"+` AS "groups" FROM fivenet_accounts WHERE fivenet_accounts.license = ? LIMIT ?;`,
		),
	).
		WithArgs("license-42", int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "license", "groups"}).
			AddRow(int64(42), "license-42", []byte(`["supporter","donator"]`)))
	mock.ExpectExec(`(?s)UPDATE .*fivenet_accounts.*SET .*groups.*WHERE .*license = \?.*LIMIT \?.*`).
		WithArgs(sqlmock.AnyArg(), "license-42", int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := store.handleAccountUpdate(t.Context(), &activity.AccountUpdate{
		License: "license-42",
		Groups:  &accounts.AccountGroups{Groups: []string{"supporter", "donator"}},
	})
	require.NoError(t, err)

	assert.Empty(t, rec.events)
	require.NoError(t, mock.ExpectationsWereMet())
}
