package mailer

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreGetEmailSettings(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)

	expectedQuery := regexp.QuoteMeta(
		`FROM fivenet_mailer_settings AS email_settings LEFT JOIN fivenet_mailer_settings_blocked ON`,
	) +
		`(?s).*` + regexp.QuoteMeta(
		`email_settings.email_id = ?`,
	) +
		`(?s).*` + regexp.QuoteMeta(
		`LIMIT ?;`,
	)
	mock.ExpectQuery(expectedQuery).
		WithArgs(int64(7), int64(25)).
		WillReturnRows(sqlmock.NewRows([]string{
			"email_settings.email_id",
			"email_settings.signature",
			"email_settings.blocked_emails",
		}).
			AddRow(int64(7), nil, "one@example.com").
			AddRow(int64(7), nil, "two@example.com"),
		)

	settings, err := store.GetEmailSettings(t.Context(), db, 7)
	require.NoError(t, err)
	require.NotNil(t, settings)
	assert.Equal(t, int64(7), settings.GetEmailId())
	assert.Nil(t, settings.GetSignature())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreUpsertEmailSettingsSignature(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)

	expectedQuery := regexp.QuoteMeta(`INSERT INTO fivenet_mailer_settings`) +
		`(?s).*` + regexp.QuoteMeta(`ON DUPLICATE KEY UPDATE`)
	mock.ExpectExec(expectedQuery).
		WithArgs(int64(7), nil, "VALUES(`signature`)").
		WillReturnResult(sqlmock.NewResult(0, 1))

	require.NoError(t, store.UpsertEmailSettingsSignature(t.Context(), db, 7, nil))
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreAddBlockedEmails(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)

	expectedQuery := regexp.QuoteMeta(`INSERT INTO fivenet_mailer_settings_blocked`) +
		`(?s).*` + regexp.QuoteMeta(`VALUES (?, ?), (?, ?)`)
	mock.ExpectExec(expectedQuery).
		WithArgs(int64(7), "a@example.com", int64(7), "b@example.com").
		WillReturnResult(sqlmock.NewResult(0, 2))

	require.NoError(
		t,
		store.AddBlockedEmails(t.Context(), db, 7, []string{"a@example.com", "b@example.com"}),
	)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreDeleteBlockedEmails(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)

	expectedQuery := regexp.QuoteMeta(`DELETE FROM fivenet_mailer_settings_blocked`) +
		`(?s).*` + regexp.QuoteMeta(`email_id = ?`) +
		`(?s).*` + regexp.QuoteMeta(`target_email IN (?, ?)`)
	mock.ExpectExec(expectedQuery).
		WithArgs(int64(7), "a@example.com", "b@example.com", int64(2)).
		WillReturnResult(sqlmock.NewResult(0, 2))

	require.NoError(
		t,
		store.DeleteBlockedEmails(t.Context(), db, 7, []string{"a@example.com", "b@example.com"}),
	)
	require.NoError(t, mock.ExpectationsWereMet())
}
