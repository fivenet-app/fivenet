package cmd

import (
	"testing"

	resourcesnotifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/notifications"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseNotificationType(t *testing.T) {
	t.Parallel()

	typ, err := parseNotificationType("info")
	require.NoError(t, err)
	assert.Equal(t, resourcesnotifications.NotificationType_NOTIFICATION_TYPE_INFO, typ)

	typ, err = parseNotificationType("SUCCESS")
	require.NoError(t, err)
	assert.Equal(t, resourcesnotifications.NotificationType_NOTIFICATION_TYPE_SUCCESS, typ)
}

func TestParseNotificationCategory(t *testing.T) {
	t.Parallel()

	category, err := parseNotificationCategory("general")
	require.NoError(t, err)
	assert.Equal(t, resourcesnotifications.NotificationCategory_NOTIFICATION_CATEGORY_GENERAL, category)

	category, err = parseNotificationCategory("CALENDAR")
	require.NoError(t, err)
	assert.Equal(t, resourcesnotifications.NotificationCategory_NOTIFICATION_CATEGORY_CALENDAR, category)
}

func TestBuildNotification(t *testing.T) {
	t.Parallel()

	not, err := (&NotificationSendCmd{
		UserID:     42,
		Type:       "warning",
		Category:   "document",
		TitleKey:   "notifications.system.test_notification.title",
		ContentKey: "notifications.system.test_notification.content",
		Index:      3,
	}).buildNotification()

	require.NoError(t, err)
	assert.Equal(t, int32(42), not.GetUserId())
	assert.Equal(t, resourcesnotifications.NotificationType_NOTIFICATION_TYPE_WARNING, not.GetType())
	assert.Equal(t, resourcesnotifications.NotificationCategory_NOTIFICATION_CATEGORY_DOCUMENT, not.GetCategory())
	require.NotNil(t, not.GetTitle())
	assert.Equal(t, "notifications.system.test_notification.title", not.GetTitle().GetKey())
	assert.Equal(t, "3", not.GetTitle().GetParameters()["index"])
	require.NotNil(t, not.GetContent())
	assert.Equal(t, "notifications.system.test_notification.content", not.GetContent().GetKey())
	assert.Equal(t, "NOTIFICATION_TYPE_WARNING", not.GetContent().GetParameters()["type"])
}
