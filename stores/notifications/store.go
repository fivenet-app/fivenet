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
	CountUnread(ctx context.Context, userID int32) (int64, error)
}

type Store struct {
	db *sql.DB
}

type ListQuery struct {
	UserID     int32
	UnreadOnly bool
	Categories []resourcesnotifications.NotificationCategory
	Offset     int64
	Limit      int64
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
