package auth

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreListOAuth2Connections(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	store := New(db, &config.CustomDB{})

	expectedQuery := regexp.QuoteMeta(`FROM fivenet_accounts_oauth2 AS oauth2_account`) +
		`(?s).*` + regexp.QuoteMeta(`oauth2_account.account_id = ?`) +
		`(?s).*` + regexp.QuoteMeta(`LIMIT ?;`)

	mock.ExpectQuery(expectedQuery).
		WillReturnRows(sqlmock.NewRows([]string{
			"oauth2_account.account_id",
			"oauth2_account.created_at",
			"oauth2_account.providername",
			"oauth2_account.external_id",
			"oauth2_account.username",
			"oauth2_account.avatar",
		}).AddRow(int64(3), time.Now(), "discord", "ext-1", "user-3", "avatar.png"))

	conns, err := store.ListOAuth2Connections(t.Context(), 3)
	require.NoError(t, err)
	require.Len(t, conns, 1)
	assert.Equal(t, "discord", conns[0].GetProviderName())
	assert.Equal(t, "ext-1", conns[0].GetExternalId())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreDeleteSocialLogin(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	store := New(db, &config.CustomDB{})

	expectedQuery := regexp.QuoteMeta(`DELETE FROM fivenet_accounts_oauth2`) +
		`(?s).*` + regexp.QuoteMeta(`fivenet_accounts_oauth2.account_id = ?`) +
		`(?s).*` + regexp.QuoteMeta(`fivenet_accounts_oauth2.provider = ?`) +
		`(?s).*` + regexp.QuoteMeta(`LIMIT ?;`)

	mock.ExpectExec(expectedQuery).
		WillReturnResult(sqlmock.NewResult(0, 1))

	require.NoError(t, store.DeleteSocialLogin(t.Context(), 3, "discord"))
	require.NoError(t, mock.ExpectationsWereMet())
}
