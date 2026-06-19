package authstore

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestStore(t *testing.T) (IStore, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	return New(db, &config.CustomDB{}), mock
}

func TestStoreGetLoginAccountByUsername(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	expectedQuery := regexp.QuoteMeta(`FROM fivenet_accounts`) +
		`(?s).*` + regexp.QuoteMeta(`fivenet_accounts.username = ?`) +
		`(?s).*` + regexp.QuoteMeta(`fivenet_accounts.reg_token IS NULL`) +
		`(?s).*` + regexp.QuoteMeta(`fivenet_accounts.password IS NOT NULL`) +
		`(?s).*` + regexp.QuoteMeta(`LIMIT ?;`)

	mock.ExpectQuery(expectedQuery).
		WillReturnRows(sqlmock.NewRows([]string{
			"fivenet_accounts.id",
			"fivenet_accounts.password",
			"fivenet_accounts.last_char",
			"fivenet_accounts.license",
			"fivenet_accounts.username",
		}).AddRow(int64(3), "password", int64(7), "license-3", "user-3"))

	acc, err := store.GetLoginAccountByUsername(t.Context(), "user-3")
	require.NoError(t, err)
	require.NotNil(t, acc)
	assert.Equal(t, int64(3), acc.ID)
	require.NotNil(t, acc.Password)
	assert.Equal(t, "password", *acc.Password)
	require.NotNil(t, acc.LastChar)
	assert.Equal(t, int32(7), *acc.LastChar)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreGetPasswordResetAccountByRegToken(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	expectedQuery := regexp.QuoteMeta(`FROM fivenet_accounts`) +
		`(?s).*` + regexp.QuoteMeta(`fivenet_accounts.reg_token = ?`) +
		`(?s).*` + regexp.QuoteMeta(`fivenet_accounts.username IS NOT NULL`) +
		`(?s).*` + regexp.QuoteMeta(`fivenet_accounts.password IS NULL`) +
		`(?s).*` + regexp.QuoteMeta(`LIMIT ?;`)

	mock.ExpectQuery(expectedQuery).
		WillReturnRows(sqlmock.NewRows([]string{
			"fivenet_accounts.id",
			"fivenet_accounts.password",
			"fivenet_accounts.username",
			"fivenet_accounts.license",
		}).AddRow(int64(4), nil, "user-4", "license-4"))

	acc, err := store.GetPasswordResetAccountByRegToken(t.Context(), "reg-token")
	require.NoError(t, err)
	require.NotNil(t, acc)
	assert.Equal(t, int64(4), acc.ID)
	assert.Nil(t, acc.Password)
	require.NotNil(t, acc.Username)
	assert.Equal(t, "user-4", *acc.Username)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreActivateAndUpdatePassword(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	activateQuery := regexp.QuoteMeta(`UPDATE fivenet_accounts SET`) +
		`(?s).*` + regexp.QuoteMeta(`username = ?`) +
		`(?s).*` + regexp.QuoteMeta(`password = ?`) +
		`(?s).*` + regexp.QuoteMeta(`reg_token = NULL`) +
		`(?s).*` + regexp.QuoteMeta(`fivenet_accounts.id = ?`) +
		`(?s).*` + regexp.QuoteMeta(`fivenet_accounts.reg_token = ?`) +
		`(?s).*` + regexp.QuoteMeta(`LIMIT ?;`)
	mock.ExpectExec(activateQuery).
		WillReturnResult(sqlmock.NewResult(1, 1))

	require.NoError(t, store.ActivateAccount(t.Context(), 9, "reg-token", "new-user", "hashed"))

	updateQuery := regexp.QuoteMeta(`UPDATE fivenet_accounts SET`) +
		`(?s).*` + regexp.QuoteMeta(`password = ?`) +
		`(?s).*` + regexp.QuoteMeta(`fivenet_accounts.id = ?`) +
		`(?s).*` + regexp.QuoteMeta(`LIMIT ?;`)
	mock.ExpectExec(updateQuery).
		WillReturnResult(sqlmock.NewResult(1, 1))

	require.NoError(t, store.UpdatePassword(t.Context(), 9, "new-hash"))
	require.NoError(t, mock.ExpectationsWereMet())
}
