package mailerstore

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreCountThreadMessages(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)

	expectedQuery := regexp.QuoteMeta(`FROM fivenet_mailer_messages AS message`) +
		`(?s).*` + regexp.QuoteMeta(`message.deleted_at IS NULL`) +
		`(?s).*` + regexp.QuoteMeta(`message.thread_id = ?`)
	mock.ExpectQuery(expectedQuery).
		WithArgs(int64(42)).
		WillReturnRows(sqlmock.NewRows([]string{"data_count.total"}).AddRow(int64(3)))

	total, err := store.CountThreadMessages(t.Context(), db, 42)
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreListThreadMessages(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)
	now := time.Unix(0, 0).UTC()

	expectedQuery := regexp.QuoteMeta(
		`FROM fivenet_mailer_messages AS message LEFT JOIN fivenet_mailer_emails AS email ON`,
	) +
		`(?s).*` + regexp.QuoteMeta(
		`message.thread_id = ?`,
	) +
		`(?s).*` + regexp.QuoteMeta(
		`LIMIT ? OFFSET ?;`,
	)
	mock.ExpectQuery(expectedQuery).
		WithArgs(int64(42), int64(20), int64(5)).
		WillReturnRows(sqlmock.NewRows([]string{
			"message.id",
			"message.thread_id",
			"message.sender_id",
			"message.created_at",
			"message.updated_at",
			"message.deleted_at",
			"message.title",
			"message.content",
			"message.data",
			"message.creator_id",
			"sender.id",
			"sender.email",
		}).AddRow(
			int64(
				9,
			), int64(42), int64(7), now, now, nil, "Title", nil, nil, int32(3), int64(7), "sender@example.com",
		),
		)

	messages, err := store.ListThreadMessages(
		t.Context(),
		db,
		MessageListQuery{ThreadID: 42, Offset: 5, Limit: 20},
	)
	require.NoError(t, err)
	require.Len(t, messages, 1)
	assert.Equal(t, int64(9), messages[0].GetId())
	assert.Equal(t, "Title", messages[0].GetTitle())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreGetMessage(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)
	now := time.Unix(0, 0).UTC()

	expectedQuery := regexp.QuoteMeta(`FROM fivenet_mailer_messages AS message`) +
		`(?s).*` + regexp.QuoteMeta(`message.id = ?`) +
		`(?s).*` + regexp.QuoteMeta(`LIMIT ?;`)
	mock.ExpectQuery(expectedQuery).
		WithArgs(int64(9), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"message.id",
			"message.thread_id",
			"message.sender_id",
			"message.created_at",
			"message.updated_at",
			"message.deleted_at",
			"message.title",
			"message.content",
			"message.data",
			"message.creator_id",
			"sender.id",
			"sender.email",
		}).AddRow(
			int64(
				9,
			), int64(42), int64(7), now, now, nil, "Title", nil, nil, int32(3), int64(7), "sender@example.com",
		),
		)

	message, err := store.GetMessage(t.Context(), db, 9)
	require.NoError(t, err)
	require.NotNil(t, message)
	assert.Equal(t, int64(9), message.GetId())
	assert.Equal(t, "Title", message.GetTitle())
	require.NoError(t, mock.ExpectationsWereMet())
}
