package notifi

import (
	"context"
	"errors"
	"testing"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/notifications"
	notificationsstore "github.com/fivenet-app/fivenet/v2026/stores/notifications"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

type notificationStoreStub struct {
	delivery *notifications.NotificationDelivery
}

func (s notificationStoreStub) Count(context.Context, notificationsstore.ListQuery) (int64, error) {
	return 0, nil
}

func (s notificationStoreStub) List(
	context.Context,
	notificationsstore.ListQuery,
) ([]*notifications.Notification, error) {
	return nil, nil
}

func (s notificationStoreStub) MarkNotifications(
	context.Context,
	notificationsstore.MarkQuery,
) (int64, error) {
	return 0, nil
}

func (s notificationStoreStub) UpdateNotificationState(
	context.Context,
	notificationsstore.StateQuery,
) (int64, error) {
	return 0, nil
}

func (s notificationStoreStub) CountUnread(context.Context, int32) (int64, error) {
	return 0, nil
}

func (s notificationStoreStub) ListPreferences(
	context.Context,
	int32,
) ([]*notifications.NotificationPreference, error) {
	return nil, nil
}

func (s notificationStoreStub) UpsertPreference(
	context.Context,
	int32,
	*notifications.NotificationPreference,
) error {
	return nil
}

func (s notificationStoreStub) DeletePreference(
	context.Context,
	int32,
	notifications.NotificationCategory,
	notifications.NotificationKind,
) error {
	return nil
}

func (s notificationStoreStub) ResolveDelivery(
	context.Context,
	int32,
	notifications.NotificationCategory,
	notifications.NotificationKind,
) (*notifications.NotificationDelivery, error) {
	return s.delivery, nil
}

func TestNoopPublish(t *testing.T) {
	t.Parallel()

	assert.NoError(t, noopPublish(t.Context()))
}

func TestNewUserNotification(t *testing.T) {
	t.Parallel()

	actorID := int32(4)
	entityType := "documents.document"
	entityID := int64(9)
	notification := NewUserNotification(UserNotificationParams{
		UserID:      3,
		Type:        notifications.NotificationType_NOTIFICATION_TYPE_WARNING,
		Category:    notifications.NotificationCategory_NOTIFICATION_CATEGORY_DOCUMENT,
		Kind:        notifications.NotificationKind_NOTIFICATION_KIND_DOCUMENT_REQUEST_CREATED,
		ActorUserID: &actorID,
		EntityType:  &entityType,
		EntityID:    &entityID,
		Title:       &common.I18NItem{Key: "title"},
		Content:     &common.I18NItem{Key: "content"},
		Data:        &notifications.Data{Link: &notifications.Link{To: "/documents/9"}},
	})

	assert.Equal(t, int32(3), notification.GetUserId())
	assert.Equal(
		t,
		notifications.NotificationType_NOTIFICATION_TYPE_WARNING,
		notification.GetType(),
	)
	assert.Equal(
		t,
		notifications.NotificationCategory_NOTIFICATION_CATEGORY_DOCUMENT,
		notification.GetCategory(),
	)
	assert.Equal(
		t,
		notifications.NotificationKind_NOTIFICATION_KIND_DOCUMENT_REQUEST_CREATED,
		notification.GetKind(),
	)
	assert.Equal(t, int32(4), notification.GetActorUserId())
	assert.Equal(t, "documents.document", notification.GetEntityType())
	assert.Equal(t, int64(9), notification.GetEntityId())
	assert.Equal(t, "title", notification.GetTitle().GetKey())
	assert.Equal(t, "/documents/9", notification.GetData().GetLink().GetTo())
}

func TestPrepareUserNotificationReturnsNoopWhenDeliveryDisabled(t *testing.T) {
	t.Parallel()

	notifier := &Notifi{
		preferences: notificationStoreStub{delivery: &notifications.NotificationDelivery{}},
	}
	notification := &notifications.Notification{UserId: 3}
	publish, err := notifier.PrepareUserNotification(t.Context(), nil, notification)

	require.NoError(t, err)
	require.NotNil(t, publish)
	require.NoError(t, publish(t.Context()))
	assert.NotNil(t, notification.GetCreatedAt())
	assert.Equal(t, int64(0), notification.GetId())
}

func TestPublishAfterCommit(t *testing.T) {
	t.Parallel()

	core, logs := observer.New(zap.WarnLevel)
	logger := zap.New(core)
	called := false

	PublishAfterCommit(
		t.Context(),
		logger,
		"test_notification",
		nil,
		func(context.Context) error {
			called = true
			return nil
		},
		func(context.Context) error { return errors.New("publish failed") },
	)

	assert.True(t, called)
	entries := logs.All()
	assert.Len(t, entries, 1)
	assert.Equal(t, "failed to publish notification", entries[0].Message)
	assert.Equal(t, "test_notification", entries[0].ContextMap()["notification_operation"])
}
