package notificationsstore

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreCountUnread(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	expectedQuery := regexp.QuoteMeta(`FROM fivenet_notifications`) +
		`(?s).*` + regexp.QuoteMeta(`fivenet_notifications.user_id = ?`) +
		`(?s).*` + regexp.QuoteMeta(`fivenet_notifications.read_at IS NULL`)
	mock.ExpectQuery(expectedQuery).
		WithArgs(int32(3)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(int64(4)))

	total, err := store.CountUnread(t.Context(), 3)
	require.NoError(t, err)
	assert.Equal(t, int64(4), total)
	require.NoError(t, mock.ExpectationsWereMet())
}
