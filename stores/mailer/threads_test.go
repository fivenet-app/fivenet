package mailerstore

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreCountThreads(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)
	query := ThreadListQuery{
		EmailIDs:  []int64{7},
		Archived:  nil,
		Superuser: false,
	}

	expectedQuery := regexp.QuoteMeta(
		`FROM fivenet_mailer_threads AS thread LEFT JOIN fivenet_mailer_threads_state AS thread_state ON`,
	) +
		`(?s).*` + regexp.QuoteMeta(
		`thread_state.email_id = ?`,
	) +
		`(?s).*` + regexp.QuoteMeta(
		`thread.deleted_at IS NULL`,
	) +
		`(?s).*` + regexp.QuoteMeta(
		`thread_state.archived IS NULL`,
	) +
		`(?s).*` + regexp.QuoteMeta(
		`thread_state.archived = ?`,
	) +
		`(?s).*` + regexp.QuoteMeta(
		`EXISTS ( SELECT ? FROM fivenet_mailer_threads_recipients WHERE`,
	) +
		`(?s).*` + regexp.QuoteMeta(
		`fivenet_mailer_threads_recipients.email_id IN (?)`,
	)
	mock.ExpectQuery(expectedQuery).
		WithArgs(int64(7), false, int64(1), int64(7)).
		WillReturnRows(sqlmock.NewRows([]string{"data_count.total"}).AddRow(int64(4)))

	total, err := store.CountThreads(t.Context(), db, query)
	require.NoError(t, err)
	assert.Equal(t, int64(4), total)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreListThreads(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)
	now := time.Unix(0, 0).UTC()
	query := ThreadListQuery{
		EmailIDs:  []int64{7},
		Archived:  nil,
		Superuser: false,
		Offset:    5,
		Limit:     20,
	}

	expectedQuery := `(?s).*` + regexp.QuoteMeta(
		`FROM fivenet_mailer_threads AS thread LEFT JOIN fivenet_mailer_threads_state AS thread_state ON`,
	) +
		`.*` + regexp.QuoteMeta(
		`EXISTS ( SELECT ? FROM fivenet_mailer_threads_recipients WHERE`,
	) +
		`.*` + regexp.QuoteMeta(
		`ORDER BY COALESCE(thread.updated_at, thread.created_at) DESC, thread.id DESC LIMIT ? OFFSET ? ) AS page INNER JOIN fivenet_mailer_threads AS thread ON`,
	) +
		`.*` + regexp.QuoteMeta(
		`LEFT JOIN fivenet_mailer_threads_state AS thread_state ON`,
	)
	mock.ExpectQuery(expectedQuery).
		WithArgs(int64(7), false, int64(1), int64(7), int64(20), int64(5), int64(7)).
		WillReturnRows(sqlmock.NewRows([]string{
			"thread.id",
			"thread.created_at",
			"thread.updated_at",
			"thread.deleted_at",
			"thread.title",
			"thread.creator_email_id",
			"thread.creator_email",
			"email.id",
			"email.email",
			"thread.creator_id",
			"thread_state.thread_id",
			"thread_state.email_id",
			"thread_state.unread",
			"thread_state.last_read",
			"thread_state.important",
			"thread_state.favorite",
			"thread_state.muted",
			"thread_state.archived",
		}).AddRow(
			int64(
				11,
			), now, now, nil, "Hello", int64(7), "sender@example.com", int64(7), "sender@example.com", int32(3), nil, nil, nil, nil, nil, nil, nil, nil,
		),
		)

	threads, err := store.ListThreads(t.Context(), db, query)
	require.NoError(t, err)
	require.Len(t, threads, 1)
	assert.Equal(t, int64(11), threads[0].GetId())
	assert.Equal(t, "Hello", threads[0].GetTitle())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreGetThread(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)
	now := time.Unix(0, 0).UTC()

	expectedQuery := regexp.QuoteMeta(
		`FROM fivenet_mailer_threads AS thread LEFT JOIN fivenet_mailer_threads_state AS thread_state ON`,
	) +
		`(?s).*` + regexp.QuoteMeta(
		`thread_state.email_id = ?`,
	) +
		`(?s).*` + regexp.QuoteMeta(
		`thread.id = ?`,
	) +
		`(?s).*` + regexp.QuoteMeta(
		`LIMIT ?;`,
	)
	mock.ExpectQuery(expectedQuery).
		WithArgs(int64(7), int64(42), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"thread.id",
			"thread.created_at",
			"thread.updated_at",
			"thread.deleted_at",
			"thread.title",
			"thread.creator_email_id",
			"thread.creator_email",
			"email.id",
			"email.email",
			"thread.creator_id",
			"thread_state.thread_id",
			"thread_state.email_id",
			"thread_state.unread",
			"thread_state.last_read",
			"thread_state.important",
			"thread_state.favorite",
			"thread_state.muted",
			"thread_state.archived",
		}).AddRow(
			int64(
				42,
			), now, now, nil, "Hello", int64(7), "sender@example.com", int64(7), "sender@example.com", int32(3), int64(42), int64(7), true, now, true, false, false, false,
		),
		)

	thread, err := store.GetThread(t.Context(), db, 42, 7)
	require.NoError(t, err)
	require.NotNil(t, thread)
	assert.Equal(t, int64(42), thread.GetId())
	assert.Equal(t, "Hello", thread.GetTitle())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreListThreadRecipients(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)
	now := time.Unix(0, 0).UTC()

	expectedQuery := regexp.QuoteMeta(
		`FROM fivenet_mailer_threads_recipients AS thread_recipient_email INNER JOIN fivenet_mailer_emails AS email ON`,
	) +
		`(?s).*` + regexp.QuoteMeta(
		`email.deactivated IS FALSE`,
	) +
		`(?s).*` + regexp.QuoteMeta(
		`thread_recipient_email.thread_id = ?`,
	) +
		`(?s).*` + regexp.QuoteMeta(
		`email.deleted_at IS NULL`,
	)
	mock.ExpectQuery(expectedQuery).
		WithArgs(int64(42)).
		WillReturnRows(sqlmock.NewRows([]string{
			"thread_recipient_email.id",
			"thread_recipient_email.thread_id",
			"thread_recipient_email.email_id",
			"email.id",
			"email.email",
			"thread_recipient_email.created_at",
		}).AddRow(int64(9), int64(42), int64(7), int64(7), "sender@example.com", now))

	recipients, err := store.ListThreadRecipients(t.Context(), db, 42)
	require.NoError(t, err)
	require.Len(t, recipients, 1)
	assert.Equal(t, int64(9), recipients[0].GetId())
	assert.Equal(t, int64(7), recipients[0].GetEmailId())
	require.NoError(t, mock.ExpectationsWereMet())
}
