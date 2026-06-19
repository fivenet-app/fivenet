package mailerstore

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	mailerthreads "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/mailer/threads"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreGetThreadState(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)
	now := time.Unix(0, 0).UTC()

	expectedQuery := regexp.QuoteMeta(`FROM fivenet_mailer_threads_state AS thread_state`) +
		`(?s).*` + regexp.QuoteMeta(`thread_state.thread_id = ?`) +
		`(?s).*` + regexp.QuoteMeta(`thread_state.email_id = ?`)
	mock.ExpectQuery(expectedQuery).
		WithArgs(int64(42), int64(7)).
		WillReturnRows(sqlmock.NewRows([]string{
			"thread_state.thread_id",
			"thread_state.email_id",
			"thread_state.unread",
			"thread_state.last_read",
			"thread_state.important",
			"thread_state.favorite",
			"thread_state.muted",
			"thread_state.archived",
		}).AddRow(int64(42), int64(7), true, now, true, false, false, false))

	state, err := store.GetThreadState(t.Context(), db, 42, 7)
	require.NoError(t, err)
	require.NotNil(t, state)
	assert.Equal(t, int64(42), state.GetThreadId())
	assert.Equal(t, int64(7), state.GetEmailId())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreSetThreadState(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db)
	now := time.Unix(0, 0).UTC()
	state := &mailerthreads.ThreadState{
		ThreadId:  42,
		EmailId:   7,
		Unread:    func() *bool { v := true; return &v }(),
		LastRead:  timestamp.New(now),
		Important: func() *bool { v := true; return &v }(),
		Favorite:  func() *bool { v := false; return &v }(),
		Muted:     func() *bool { v := false; return &v }(),
		Archived:  func() *bool { v := false; return &v }(),
	}

	expectedQuery := regexp.QuoteMeta(`INSERT INTO fivenet_mailer_threads_state`) +
		`(?s).*` + regexp.QuoteMeta(`VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
	mock.ExpectExec(expectedQuery).
		WithArgs(int64(42), int64(7), true, now, true, false, false, false).
		WillReturnResult(sqlmock.NewResult(0, 1))

	require.NoError(t, store.SetThreadState(t.Context(), db, state))
	require.NoError(t, mock.ExpectationsWereMet())
}
