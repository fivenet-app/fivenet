package mailerstore

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreCountEmails(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)

	expectedQuery := regexp.QuoteMeta(`FROM fivenet_mailer_emails AS email`) +
		`(?s).*` + regexp.QuoteMeta(`WHERE ?`)
	mock.ExpectQuery(expectedQuery).
		WithArgs(true).
		WillReturnRows(sqlmock.NewRows([]string{"data_count.total"}).AddRow(int64(7)))

	total, err := store.CountEmails(t.Context(), db, mysql.Bool(true))
	require.NoError(t, err)
	assert.Equal(t, int64(7), total)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreListEmails(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)
	now := time.Unix(0, 0).UTC()

	countQuery := regexp.QuoteMeta(`SELECT COUNT(email.id) AS "data_count.total" FROM fivenet_mailer_emails AS email;`)
	mock.ExpectQuery(countQuery).
		WillReturnRows(sqlmock.NewRows([]string{"data_count.total"}).AddRow(int64(1)))

	listQuery := regexp.QuoteMeta(`FROM fivenet_mailer_emails AS email`) +
		`(?s).*` + regexp.QuoteMeta(
		`ORDER BY email.job ASC, email.label ASC LIMIT ? OFFSET ?;`,
	)
	mock.ExpectQuery(listQuery).
		WillReturnRows(sqlmock.NewRows([]string{
			"email.id",
			"email.created_at",
			"email.updated_at",
			"email.deleted_at",
			"email.deactivated",
			"email.job",
			"email.user_id",
			"email.email",
			"email.email_changed",
			"email.label",
		}).AddRow(
			int64(9),
			now,
			now,
			nil,
			false,
			"police",
			int32(3),
			"user@example.com",
			nil,
			"Primary",
		))

	pageSize := int64(20)
	pag := &database.PaginationRequest{Offset: 5, PageSize: &pageSize}
	emailsPag, emails, err := store.ListEmails(
		t.Context(),
		db,
		&userinfo.UserInfo{Superuser: true},
		pag,
		true,
	)
	require.NoError(t, err)
	require.NotNil(t, emailsPag)
	assert.Equal(t, int64(1), emailsPag.GetTotalCount())
	require.Len(t, emails, 1)
	assert.Equal(t, int64(9), emails[0].GetId())
	assert.Equal(t, "user@example.com", emails[0].GetEmail())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreListEmailsVisible(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)
	now := time.Unix(0, 0).UTC()

	countQuery := `(?s).*WITH user_subjects AS.*visible_sources AS.*winning_visibility AS.*data_count.total`
	mock.ExpectQuery(countQuery).
		WillReturnRows(sqlmock.NewRows([]string{"data_count.total"}).AddRow(int64(1)))

	listQuery := `(?s).*WITH user_subjects AS.*visible_sources AS.*winning_visibility AS.*` +
		regexp.QuoteMeta(`ORDER BY email.job ASC, email.label ASC LIMIT ? OFFSET ?;`)
	mock.ExpectQuery(listQuery).
		WillReturnRows(sqlmock.NewRows([]string{
			"email.id",
			"email.created_at",
			"email.updated_at",
			"email.deleted_at",
			"email.deactivated",
			"email.job",
			"email.user_id",
			"email.email",
			"email.email_changed",
			"email.label",
		}).AddRow(
			int64(12),
			now,
			now,
			nil,
			false,
			"police",
			int32(3),
			"user2@example.com",
			nil,
			"Secondary",
		))

	pageSize := int64(20)
	emailsPag, emails, err := store.ListEmails(
		t.Context(),
		db,
		&userinfo.UserInfo{UserId: 3, Job: "police", Superuser: true},
		&database.PaginationRequest{PageSize: &pageSize},
		false,
	)
	require.NoError(t, err)
	require.NotNil(t, emailsPag)
	assert.Equal(t, int64(1), emailsPag.GetTotalCount())
	require.Len(t, emails, 1)
	assert.Equal(t, int64(12), emails[0].GetId())
	assert.Equal(t, "user2@example.com", emails[0].GetEmail())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreGetEmail(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)
	now := time.Unix(0, 0).UTC()

	expectedQuery := regexp.QuoteMeta(`FROM fivenet_mailer_emails AS email`) +
		`(?s).*` + regexp.QuoteMeta(`email.id = ?`) +
		`(?s).*` + regexp.QuoteMeta(`LIMIT ?;`)
	mock.ExpectQuery(expectedQuery).
		WithArgs(int64(42), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"email.id",
			"email.created_at",
			"email.updated_at",
			"email.deleted_at",
			"email.deactivated",
			"email.job",
			"email.user_id",
			"email.email",
			"email.email_changed",
			"email.label",
		}).AddRow(
			int64(42),
			now,
			now,
			nil,
			true,
			"police",
			int32(3),
			"mail@example.com",
			nil,
			nil,
		))

	email, err := store.GetEmail(t.Context(), db, 42, false)
	require.NoError(t, err)
	require.NotNil(t, email)
	assert.Equal(t, int64(42), email.GetId())
	assert.Equal(t, "mail@example.com", email.GetEmail())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreGetUserShort(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)

	expectedQuery := regexp.QuoteMeta(`FROM fivenet_user AS user_short`) +
		`(?s).*` + regexp.QuoteMeta(`user_short.id = ?`) +
		`(?s).*` + regexp.QuoteMeta(`LIMIT ?;`)
	mock.ExpectQuery(expectedQuery).
		WithArgs(int32(3), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"user_short.firstname",
			"user_short.lastname",
			"user_short.dateofbirth",
		}).AddRow("Jane", "Doe", "01.01.2000"))

	user, err := store.GetUserShort(t.Context(), db, 3)
	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, "Jane", user.GetFirstname())
	assert.Equal(t, "01.01.2000", user.GetDateofbirth())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreListRecipientsByEmails(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)

	expectedQuery := regexp.QuoteMeta(`FROM fivenet_mailer_emails AS email`) +
		`(?s).*` + regexp.QuoteMeta(`email.email IN (?, ?)`) +
		`(?s).*` + regexp.QuoteMeta(`email.deleted_at IS NULL`) +
		`(?s).*` + regexp.QuoteMeta(`LIMIT ?;`)
	mock.ExpectQuery(expectedQuery).
		WithArgs("a@example.com", "b@example.com", int64(2)).
		WillReturnRows(sqlmock.NewRows([]string{
			"thread_recipient_email.email_id",
			"email.email",
			"email.deactivated",
		}).AddRow(int64(3), "a@example.com", false).AddRow(int64(4), "b@example.com", true))

	recipients, err := store.ListRecipientsByEmails(
		t.Context(),
		db,
		[]string{"a@example.com", "b@example.com"},
	)
	require.NoError(t, err)
	require.Len(t, recipients, 2)
	assert.Equal(t, int64(3), recipients[0].GetEmailId())
	assert.Equal(t, "a@example.com", recipients[0].GetEmail().GetEmail())
	assert.Equal(t, int64(4), recipients[1].GetEmailId())
	assert.Equal(t, "b@example.com", recipients[1].GetEmail().GetEmail())
	require.NoError(t, mock.ExpectationsWereMet())
}
