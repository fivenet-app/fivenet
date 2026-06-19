package livemapstore

import (
	"context"
	"errors"

	livemapmarkers "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/livemap/markers"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) CreateMarker(
	ctx context.Context,
	marker *livemapmarkers.MarkerMarker,
	creatorID int32,
	job string,
) (int64, error) {
	tMarkers := table.FivenetCentrumMarkers
	stmt := tMarkers.
		INSERT(
			tMarkers.ExpiresAt,
			tMarkers.Job,
			tMarkers.Name,
			tMarkers.Description,
			tMarkers.X,
			tMarkers.Y,
			tMarkers.Postal,
			tMarkers.Color,
			tMarkers.MarkerType,
			tMarkers.MarkerData,
			tMarkers.CreatorID,
		).
		VALUES(
			marker.GetExpiresAt(),
			job,
			marker.GetName(),
			marker.Description,
			marker.GetX(),
			marker.GetY(),
			marker.Postal,
			marker.Color,
			marker.GetType(),
			marker.GetData(),
			creatorID,
		)

	res, err := stmt.ExecContext(ctx, s.db)
	if err != nil {
		return 0, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastID, nil
}

func (s *Store) UpdateMarker(
	ctx context.Context,
	marker *livemapmarkers.MarkerMarker,
	job string,
) error {
	tMarkers := table.FivenetCentrumMarkers
	stmt := tMarkers.
		UPDATE(
			tMarkers.ExpiresAt,
			tMarkers.Name,
			tMarkers.Description,
			tMarkers.X,
			tMarkers.Y,
			tMarkers.Postal,
			tMarkers.Color,
			tMarkers.MarkerType,
			tMarkers.MarkerData,
		).
		SET(
			marker.GetExpiresAt(),
			marker.GetName(),
			marker.Description,
			marker.GetX(),
			marker.GetY(),
			marker.Postal,
			marker.Color,
			marker.GetType(),
			marker.GetData(),
		).
		WHERE(mysql.AND(
			tMarkers.Job.EQ(mysql.String(job)),
			tMarkers.ID.EQ(mysql.Int64(marker.GetId())),
		)).
		LIMIT(1)

	_, err := stmt.ExecContext(ctx, s.db)
	return err
}

func (s *Store) SoftDeleteMarker(ctx context.Context, id int64) error {
	tMarkers := table.FivenetCentrumMarkers
	stmt := tMarkers.
		UPDATE(
			tMarkers.DeletedAt,
		).
		SET(
			tMarkers.DeletedAt.SET(mysql.CURRENT_TIMESTAMP()),
		).
		WHERE(
			tMarkers.ID.EQ(mysql.Int64(id)),
		).
		LIMIT(1)

	_, err := stmt.ExecContext(ctx, s.db)
	return err
}

func (s *Store) GetMarker(ctx context.Context, id int64) (*livemapmarkers.MarkerMarker, error) {
	tMarkers := table.FivenetCentrumMarkers.AS("marker_marker")
	tUsers := table.FivenetUser.AS("user_short")

	stmt := tMarkers.
		SELECT(
			tMarkers.ID,
			tMarkers.CreatedAt,
			tMarkers.UpdatedAt,
			tMarkers.DeletedAt,
			tMarkers.ExpiresAt,
			tMarkers.Job,
			tMarkers.Name,
			tMarkers.Description,
			tMarkers.X,
			tMarkers.Y,
			tMarkers.Postal,
			tMarkers.Color,
			tMarkers.MarkerType,
			tMarkers.MarkerData,
			tMarkers.CreatorID,
			tUsers.ID,
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.Sex,
			tUsers.PhoneNumber,
		).
		FROM(
			tMarkers.
				LEFT_JOIN(tUsers,
					tMarkers.CreatorID.EQ(tUsers.ID),
				),
		).
		WHERE(
			tMarkers.ID.EQ(mysql.Int64(id)),
		).
		LIMIT(1)

	var dest livemapmarkers.MarkerMarker
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, err
	}

	return &dest, nil
}

func (s *Store) ListActiveMarkers(ctx context.Context) ([]*livemapmarkers.MarkerMarker, error) {
	tMarkers := table.FivenetCentrumMarkers.AS("marker_marker")
	tUsers := table.FivenetUser.AS("user_short")

	stmt := tMarkers.
		SELECT(
			tMarkers.ID,
			tMarkers.CreatedAt,
			tMarkers.UpdatedAt,
			tMarkers.DeletedAt,
			tMarkers.ExpiresAt,
			tMarkers.Job,
			tMarkers.Name,
			tMarkers.Description,
			tMarkers.X,
			tMarkers.Y,
			tMarkers.Postal,
			tMarkers.Color,
			tMarkers.MarkerType,
			tMarkers.MarkerData,
			tMarkers.CreatorID,
			tUsers.ID,
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.Sex,
			tUsers.PhoneNumber,
		).
		FROM(
			tMarkers.
				LEFT_JOIN(tUsers,
					tMarkers.CreatorID.EQ(tUsers.ID),
				),
		).
		WHERE(mysql.AND(
			tMarkers.DeletedAt.IS_NULL(),
			mysql.OR(
				tMarkers.ExpiresAt.IS_NULL(),
				tMarkers.ExpiresAt.GT(mysql.CURRENT_TIMESTAMP()),
			),
		)).
		ORDER_BY(
			tMarkers.Job.ASC(),
			tMarkers.ID.ASC(),
		)

	var dest []*livemapmarkers.MarkerMarker
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return []*livemapmarkers.MarkerMarker{}, nil
		}
		return nil, err
	}

	return dest, nil
}

func (s *Store) ListDeletedMarkers(ctx context.Context) ([]*livemapmarkers.MarkerMarker, error) {
	tMarkers := table.FivenetCentrumMarkers.AS("marker_marker")

	stmt := tMarkers.
		SELECT(
			tMarkers.ID,
			tMarkers.Job,
		).
		FROM(
			tMarkers,
		).
		WHERE(mysql.OR(
			tMarkers.DeletedAt.IS_NOT_NULL(),
			tMarkers.ExpiresAt.LT_EQ(mysql.CURRENT_TIMESTAMP()),
		)).
		ORDER_BY(
			tMarkers.ID.ASC(),
		)

	var dest []*livemapmarkers.MarkerMarker
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return []*livemapmarkers.MarkerMarker{}, nil
		}
		return nil, err
	}

	return dest, nil
}
