package notificationsstore

import (
	"context"
	"database/sql"

	resourcesnotifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/notifications"
)

type IStore interface {
	Count(ctx context.Context, q ListQuery) (int64, error)
	List(ctx context.Context, q ListQuery) ([]*resourcesnotifications.Notification, error)
	MarkNotifications(ctx context.Context, q MarkQuery) (int64, error)
	UpdateNotificationState(ctx context.Context, q StateQuery) (int64, error)
	CountUnread(ctx context.Context, userID int32) (int64, error)
	ListPreferences(
		ctx context.Context,
		userID int32,
	) ([]*resourcesnotifications.NotificationPreference, error)
	UpsertPreference(
		ctx context.Context,
		userID int32,
		preference *resourcesnotifications.NotificationPreference,
	) error
	DeletePreference(
		ctx context.Context,
		userID int32,
		category resourcesnotifications.NotificationCategory,
		kind resourcesnotifications.NotificationKind,
	) error
	ResolveDelivery(
		ctx context.Context,
		userID int32,
		category resourcesnotifications.NotificationCategory,
		kind resourcesnotifications.NotificationKind,
	) (*resourcesnotifications.NotificationDelivery, error)
}

type Store struct {
	db *sql.DB
}

type ListQuery struct {
	UserID       int32
	UnreadOnly   bool
	Categories   []resourcesnotifications.NotificationCategory
	Kinds        []resourcesnotifications.NotificationKind
	Archived     bool
	ArchivedOnly bool
	Starred      bool
	Offset       int64
	Limit        int64
}

type StateQuery struct {
	UserID   int32
	IDs      []int64
	Unread   *bool
	Starred  *bool
	Archived *bool
}

type MarkQuery struct {
	UserID int32
	IDs    []int64
	All    bool
	Unread bool
}

func New(db *sql.DB) IStore {
	return &Store{db: db}
}
