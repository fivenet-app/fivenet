package livemapstore

import (
	"context"
	"database/sql"

	livemapmarkers "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/livemap/markers"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/pkg/events"
	"go.uber.org/fx"
)

type IStore interface {
	CreateMarker(
		ctx context.Context,
		marker *livemapmarkers.MarkerMarker,
		creatorID int32,
		job string,
	) (int64, error)
	UpdateMarker(ctx context.Context, marker *livemapmarkers.MarkerMarker, job string) error
	DeleteMarker(ctx context.Context, id int64, deletedAt *timestamp.Timestamp) error
	GetMarker(ctx context.Context, id int64) (*livemapmarkers.MarkerMarker, error)
	ListActiveMarkers(ctx context.Context) ([]*livemapmarkers.MarkerMarker, error)
	ListDeletedMarkers(ctx context.Context) ([]*livemapmarkers.MarkerMarker, error)
}

const (
	markerEventsSubject events.Subject = "livemap"
	markerEventsTopic   events.Topic   = "marker"
	markerEventUpdate   events.Type    = "update"
	markerEventDelete   events.Type    = "delete"
)

type Params struct {
	fx.In

	DB *sql.DB
	JS *events.JSWrapper `optional:"true"`
}

type Store struct {
	db *sql.DB
	js *events.JSWrapper
}

func New(p Params) IStore {
	return &Store{
		db: p.DB,
		js: p.JS,
	}
}

func (s *Store) publishMarkerEvent(
	ctx context.Context,
	tType events.Type,
	marker *livemapmarkers.MarkerMarker,
) error {
	if s.js == nil {
		return nil
	}

	_, err := s.js.PublishProto(
		ctx,
		buildMarkerSubject(tType, marker.GetJob()),
		marker,
	)
	return err
}

func buildMarkerSubject(tType events.Type, job string) string {
	return string(
		markerEventsSubject,
	) + "." + string(
		markerEventsTopic,
	) + "." + string(
		tType,
	) + "." + job
}
