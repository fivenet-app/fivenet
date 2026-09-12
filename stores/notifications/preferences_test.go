package notificationsstore

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	resourcesnotifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/notifications"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreListPreferences(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)
	mock.ExpectQuery(`(?s)SELECT .*category AS "NotificationPreference.category".*kind AS "NotificationPreference.kind".*inbox_enabled AS "NotificationPreference.inbox_enabled".*FROM fivenet_user_notification_preferences WHERE .*user_id = \?.*ORDER BY .*category.*kind`).
		WithArgs(int32(3), int64(64)).
		WillReturnRows(sqlmock.NewRows([]string{"NotificationPreference.category", "NotificationPreference.kind", "NotificationPreference.inbox_enabled", "NotificationPreference.toast_enabled", "NotificationPreference.sound_enabled"}).
			AddRow(int32(resourcesnotifications.NotificationCategory_NOTIFICATION_CATEGORY_JOBS), int32(0), false, nil, true),
		)

	preferences, err := store.ListPreferences(t.Context(), 3)
	require.NoError(t, err)
	require.Len(t, preferences, 1)
	assert.False(t, preferences[0].GetInboxEnabled())
	assert.False(t, preferences[0].HasToastEnabled())
	assert.True(t, preferences[0].GetSoundEnabled())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreResolveDeliveryUsesGlobalCategoryAndKindOverrides(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)
	mock.ExpectQuery(`(?s)SELECT COALESCE\(kind_preference.inbox_enabled, category_preference.inbox_enabled, global_preference.inbox_enabled, \?\) AS "NotificationDelivery.inbox_enabled".*FROM \( SELECT \? AS "user_id" \) AS target LEFT JOIN fivenet_user_notification_preferences AS global_preference.*category_preference.*kind_preference`).
		WithArgs(
			true, true, true,
			int32(3),
			int32(0), int32(0),
			int32(resourcesnotifications.NotificationCategory_NOTIFICATION_CATEGORY_JOBS), int32(0),
			int32(resourcesnotifications.NotificationCategory_NOTIFICATION_CATEGORY_JOBS),
			int32(
				resourcesnotifications.NotificationKind_NOTIFICATION_KIND_JOBS_GROUP_LEADERSHIP_ADDED,
			),
		).
		WillReturnRows(sqlmock.NewRows([]string{
			"NotificationDelivery.inbox_enabled",
			"NotificationDelivery.toast_enabled",
			"NotificationDelivery.sound_enabled",
		}).AddRow(false, false, true))

	delivery, err := store.ResolveDelivery(
		t.Context(),
		3,
		resourcesnotifications.NotificationCategory_NOTIFICATION_CATEGORY_JOBS,
		resourcesnotifications.NotificationKind_NOTIFICATION_KIND_JOBS_GROUP_LEADERSHIP_ADDED,
	)
	require.NoError(t, err)
	assert.False(t, delivery.GetInboxEnabled())
	assert.False(t, delivery.GetToastEnabled())
	assert.True(t, delivery.GetSoundEnabled())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreUpsertPreference(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)
	preference := &resourcesnotifications.NotificationPreference{
		Category:     resourcesnotifications.NotificationCategory_NOTIFICATION_CATEGORY_DOCUMENT,
		InboxEnabled: new(false),
		SoundEnabled: new(true),
	}
	mock.ExpectExec(`(?s)INSERT INTO fivenet_user_notification_preferences \(user_id, category, kind, inbox_enabled, toast_enabled, sound_enabled\) VALUES \(\?, \?, \?, \?, \?, \?\) ON DUPLICATE KEY UPDATE.*`).
		WithArgs(
			int32(3),
			int32(resourcesnotifications.NotificationCategory_NOTIFICATION_CATEGORY_DOCUMENT),
			int32(resourcesnotifications.NotificationKind_NOTIFICATION_KIND_UNSPECIFIED),
			preference.InboxEnabled,
			nil,
			preference.SoundEnabled,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	require.NoError(t, store.UpsertPreference(t.Context(), 3, preference))
	require.NoError(t, mock.ExpectationsWereMet())
}
