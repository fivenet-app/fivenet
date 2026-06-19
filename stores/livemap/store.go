package livemapstore

import (
	"context"
	"database/sql"

	livemapmarkers "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/livemap/markers"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
)

type IStore interface {
	CreateMarker(
		ctx context.Context,
		marker *livemapmarkers.MarkerMarker,
		creatorID int32,
		job string,
	) (int64, error)
	UpdateMarker(ctx context.Context, marker *livemapmarkers.MarkerMarker, job string) error
	SoftDeleteMarker(ctx context.Context, id int64) error
	GetMarker(ctx context.Context, id int64) (*livemapmarkers.MarkerMarker, error)
	ListActiveMarkers(ctx context.Context) ([]*livemapmarkers.MarkerMarker, error)
	ListDeletedMarkers(ctx context.Context) ([]*livemapmarkers.MarkerMarker, error)
}

type Store struct {
	db       *sql.DB
	customDB *config.CustomDB
}

func New(db *sql.DB, customDB *config.CustomDB) IStore {
	return &Store{
		db:       db,
		customDB: customDB,
	}
}
