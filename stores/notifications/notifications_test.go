package notificationsstore

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	resourcesnotifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/notifications"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreCountAppliesNotificationFilters(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	query := ListQuery{
		UserID:     3,
		UnreadOnly: true,
		Categories: []resourcesnotifications.NotificationCategory{
			resourcesnotifications.NotificationCategory_NOTIFICATION_CATEGORY_GENERAL,
			resourcesnotifications.NotificationCategory_NOTIFICATION_CATEGORY_CALENDAR,
		},
	}

	expectedQuery := regexp.QuoteMeta(`FROM fivenet_notifications AS notification`) +
		`(?s).*` + regexp.QuoteMeta(`notification.user_id = ?`) +
		`(?s).*` + regexp.QuoteMeta(`notification.read_at IS NULL`) +
		`(?s).*` + regexp.QuoteMeta(`notification.category IN (?, ?)`)
	mock.ExpectQuery(expectedQuery).
		WithArgs(int32(3), int32(resourcesnotifications.NotificationCategory_NOTIFICATION_CATEGORY_GENERAL), int32(resourcesnotifications.NotificationCategory_NOTIFICATION_CATEGORY_CALENDAR)).
		WillReturnRows(sqlmock.NewRows([]string{"data_count.total"}).AddRow(int64(7)))

	total, err := store.Count(t.Context(), query)
	require.NoError(t, err)
	assert.Equal(t, int64(7), total)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreListAppliesNotificationFiltersAndPaging(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	query := ListQuery{
		UserID: 3,
		Offset: 10,
		Limit:  20,
		Categories: []resourcesnotifications.NotificationCategory{
			resourcesnotifications.NotificationCategory_NOTIFICATION_CATEGORY_GENERAL,
		},
	}

	expectedQuery := regexp.QuoteMeta(`ORDER BY notification.id DESC LIMIT ? OFFSET ?;`)
	mock.ExpectQuery(expectedQuery).
		WithArgs(int32(3), int32(resourcesnotifications.NotificationCategory_NOTIFICATION_CATEGORY_GENERAL), int64(20), int64(10)).
		WillReturnRows(sqlmock.NewRows([]string{"notification.id", "notification.user_id"}).AddRow(int64(42), int32(3)))

	notifications, err := store.List(t.Context(), query)
	require.NoError(t, err)
	require.Len(t, notifications, 1)
	assert.Equal(t, int64(42), notifications[0].GetId())
	assert.Equal(t, int32(3), notifications[0].GetUserId())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreMarkNotificationsMarksSelectedRows(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	query := MarkQuery{
		UserID: 3,
		IDs:    []int64{10, 11},
	}

	expectedQuery := regexp.QuoteMeta(
		`UPDATE fivenet_notifications SET read_at = CURRENT_TIMESTAMP WHERE`,
	) +
		`(?s).*` + regexp.QuoteMeta(
		`fivenet_notifications.user_id = ?`,
	) +
		`(?s).*` + regexp.QuoteMeta(
		`fivenet_notifications.read_at IS NULL`,
	) +
		`(?s).*` + regexp.QuoteMeta(
		`fivenet_notifications.id IN (?, ?)`,
	) +
		`(?s).*`
	mock.ExpectExec(expectedQuery).
		WithArgs(int32(3), int64(10), int64(11)).
		WillReturnResult(sqlmock.NewResult(0, 2))

	affected, err := store.MarkNotifications(t.Context(), query)
	require.NoError(t, err)
	assert.Equal(t, int64(2), affected)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreMarkNotificationsMarksUnread(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	query := MarkQuery{
		UserID: 3,
		All:    true,
		Unread: true,
	}

	expectedQuery := regexp.QuoteMeta(`UPDATE fivenet_notifications SET read_at = NULL WHERE`) +
		`(?s).*` + regexp.QuoteMeta(`fivenet_notifications.user_id = ?`) +
		`(?s).*` + regexp.QuoteMeta(`fivenet_notifications.read_at IS NULL`)
	mock.ExpectExec(expectedQuery).
		WithArgs(int32(3)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	affected, err := store.MarkNotifications(t.Context(), query)
	require.NoError(t, err)
	assert.Equal(t, int64(1), affected)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreMarkNotificationsSkipsWhenNothingSelected(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	affected, err := store.MarkNotifications(t.Context(), MarkQuery{UserID: 3})
	require.NoError(t, err)
	assert.Equal(t, int64(0), affected)
	require.NoError(t, mock.ExpectationsWereMet())
}
