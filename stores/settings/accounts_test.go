package settingsstore

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	pbsettings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/settings"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreListAccountsAppliesFiltersAndFallbackSort(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	license := "ABC"
	disabled := true
	username := "alice"
	group := "staff"
	pageSize := int64(10)
	req := &pbsettings.ListAccountsRequest{
		Pagination:   &database.PaginationRequest{PageSize: &pageSize},
		License:      &license,
		OnlyDisabled: &disabled,
		Username:     &username,
		Group:        &group,
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(account.id) AS "data_count.total" FROM fivenet_accounts AS account WHERE`) + `(?s).*` + regexp.QuoteMeta(`account.license LIKE ?`) + `(?s).*` + regexp.QuoteMeta(`account.enabled = ?`) + `(?s).*` + regexp.QuoteMeta(`account.username LIKE ?`) + `(?s).*` + regexp.QuoteMeta("JSON_CONTAINS(account.`groups`, ?)")).
		WillReturnRows(sqlmock.NewRows([]string{"data_count.total"}).AddRow(int64(1)))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT account.id AS "account.id" FROM fivenet_accounts AS account WHERE`) + `(?s).*` + regexp.QuoteMeta(`ORDER BY account.created_at DESC`) + `(?s).*` + regexp.QuoteMeta(`LIMIT ? OFFSET ?;`)).
		WillReturnRows(sqlmock.NewRows([]string{"account.id"}).AddRow(int64(7)))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT account.id AS "account.id"`) + `(?s).*` + regexp.QuoteMeta(`FROM fivenet_accounts AS account LEFT JOIN fivenet_accounts_oauth2 AS oauth2account ON`) + `(?s).*` + regexp.QuoteMeta(`WHERE account.id IN (?)`)).
		WillReturnRows(sqlmock.NewRows([]string{
			"account.id",
			"account.created_at",
			"account.updated_at",
			"account.deleted_at",
			"account.enabled",
			"account.username",
			"account.license",
			"account.groups",
			"account.last_char",
			"oauth2account.account_id",
			"oauth2account.created_at",
			"oauth2account.providername",
			"oauth2account.external_id",
			"oauth2account.username",
			"oauth2account.avatar",
		}).AddRow(
			int64(7),
			time.Now(),
			time.Now(),
			nil,
			false,
			"alice",
			"ABC",
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
			nil,
		))

	resp, err := store.ListAccounts(t.Context(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.GetAccounts(), 1)
	assert.Equal(t, int64(7), resp.GetAccounts()[0].GetId())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreUpdateAccountReturnsUpdatedAccount(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)
	enabled := true
	lastChar := int32(9)

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE fivenet_accounts AS account SET enabled = ?, last_char = ? WHERE account.id = ? LIMIT ?;`)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT account.id AS "account.id"`) + `(?s).*` + regexp.QuoteMeta(`FROM fivenet_accounts AS account WHERE ( (account.enabled IS TRUE) AND (account.deleted_at IS NULL) AND (account.id = ?) ) LIMIT ?;`)).
		WillReturnRows(sqlmock.NewRows([]string{
			"account.id",
			"account.created_at",
			"account.updated_at",
			"account.deleted_at",
			"account.enabled",
			"account.username",
			"account.license",
			"account.groups",
			"account.last_char",
		}).AddRow(
			int64(7),
			time.Now(),
			time.Now(),
			nil,
			true,
			"alice",
			"ABC",
			nil,
			nil,
		))

	resp, err := store.UpdateAccount(
		t.Context(),
		&pbsettings.UpdateAccountRequest{Id: 7, Enabled: &enabled, LastChar: &lastChar},
	)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.GetAccount())
	assert.Equal(t, int64(7), resp.GetAccount().GetId())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreDeleteAccountSetsDeletedAt(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE fivenet_accounts SET deleted_at = NULL WHERE fivenet_accounts.id = ? LIMIT ?;`)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	resp, err := store.DeleteAccount(t.Context(), 7, nil)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Nil(t, resp.GetDeletedAt())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreGetAccountByIDReturnsNilForMissingAccount(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_accounts AS account`) + `(?s).*` + regexp.QuoteMeta(`account.id = ?`) + `(?s).*` + regexp.QuoteMeta(`LIMIT ?`)).
		WillReturnRows(sqlmock.NewRows([]string{"account.id"}))

	acc, err := store.GetAccountByID(t.Context(), 7)
	require.NoError(t, err)
	assert.Nil(t, acc)
	require.NoError(t, mock.ExpectationsWereMet())
}
