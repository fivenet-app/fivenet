package livemap

import (
	"context"
	"errors"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/livemap"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/permissions"
	pblivemap "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/livemap"
	permslivemap "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/livemap/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/access"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2025/pkg/grpc/interceptors/audit"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorslivemap "github.com/fivenet-app/fivenet/v2025/services/livemap/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

var tMarkers = table.FivenetCentrumMarkers.AS("marker_marker")

func (s *Server) CreateOrUpdateMarker(
	ctx context.Context,
	req *pblivemap.CreateOrUpdateMarkerRequest,
) (*pblivemap.CreateOrUpdateMarkerResponse, error) {
	if req.GetMarker() != nil && req.GetMarker().GetId() > 0 {
		logging.InjectFields(
			ctx,
			logging.Fields{"fivenet.livemap.marker_id", req.GetMarker().GetId()},
		)
	}

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	if req.Marker.Postal == nil || *req.Marker.Postal == "" {
		if postal, ok := s.postals.Closest(req.Marker.X, req.Marker.Y); postal != nil && ok {
			req.Marker.Postal = postal.Code
		}
	}

	// No marker id set
	if req.GetMarker().GetId() <= 0 {
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
				req.GetMarker().GetExpiresAt(),
				userInfo.GetJob(),
				req.GetMarker().GetName(),
				req.GetMarker().Description,
				req.GetMarker().GetX(),
				req.GetMarker().GetY(),
				req.GetMarker().Postal,
				req.GetMarker().Color,
				req.GetMarker().GetType(),
				req.GetMarker().GetData(),
				userInfo.GetUserId(),
			)

		res, err := stmt.ExecContext(ctx, s.db)
		if err != nil {
			return nil, errswrap.NewError(err, errorslivemap.ErrMarkerFailed)
		}

		lastId, err := res.LastInsertId()
		if err != nil {
			return nil, errswrap.NewError(err, errorslivemap.ErrMarkerFailed)
		}

		req.Marker.Id = lastId

		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_CREATED)
	} else {
		fields, err := s.ps.AttrStringList(userInfo, permslivemap.LivemapServicePerm, permslivemap.LivemapServiceCreateOrUpdateMarkerPerm, permslivemap.LivemapServiceCreateOrUpdateMarkerAccessPermField)
		if err != nil {
			return nil, errswrap.NewError(err, errorslivemap.ErrMarkerFailed)
		}

		marker, err := s.getMarker(ctx, req.GetMarker().GetId())
		if err != nil {
			return nil, errswrap.NewError(err, errorslivemap.ErrMarkerFailed)
		}

		if !access.CheckIfHasOwnJobAccess(fields, userInfo, marker.GetCreator().GetJob(), marker.GetCreator()) {
			return nil, errorslivemap.ErrMarkerDenied
		}

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
				req.GetMarker().GetExpiresAt(),
				req.GetMarker().GetName(),
				req.GetMarker().Description,
				req.GetMarker().GetX(),
				req.GetMarker().GetY(),
				req.GetMarker().Postal,
				req.GetMarker().Color,
				req.GetMarker().GetType(),
				req.GetMarker().GetData(),
			).
			WHERE(mysql.AND(
				tMarkers.Job.EQ(mysql.String(userInfo.GetJob())),
				tMarkers.ID.EQ(mysql.Int64(req.GetMarker().GetId())),
			))

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, errswrap.NewError(err, errorslivemap.ErrMarkerFailed)
		}

		grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)
	}

	marker, err := s.getMarker(ctx, req.GetMarker().GetId())
	if err != nil {
		return nil, errswrap.NewError(err, errorslivemap.ErrMarkerFailed)
	}

	if err := s.sendUpdateEvent(ctx, MarkerTopic, MarkerUpdate, marker.GetJob(), marker); err != nil {
		return nil, errswrap.NewError(err, errorslivemap.ErrMarkerFailed)
	}

	return &pblivemap.CreateOrUpdateMarkerResponse{
		Marker: marker,
	}, nil
}

func (s *Server) DeleteMarker(
	ctx context.Context,
	req *pblivemap.DeleteMarkerRequest,
) (*pblivemap.DeleteMarkerResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.livemap.marker_id", req.GetId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	fields, err := s.ps.AttrStringList(
		userInfo,
		permslivemap.LivemapServicePerm,
		permslivemap.LivemapServiceDeleteMarkerPerm,
		permslivemap.LivemapServiceDeleteMarkerAccessPermField,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorslivemap.ErrMarkerFailed)
	}

	marker, err := s.getMarker(ctx, req.GetId())
	if err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorslivemap.ErrMarkerFailed)
		}

		return &pblivemap.DeleteMarkerResponse{}, nil
	}

	if !access.CheckIfHasOwnJobAccess(
		fields,
		userInfo,
		marker.GetCreator().GetJob(),
		marker.GetCreator(),
	) {
		return nil, errorslivemap.ErrMarkerDenied
	}

	stmt := tMarkers.
		UPDATE(
			tMarkers.DeletedAt,
		).
		SET(
			tMarkers.DeletedAt.SET(mysql.CURRENT_TIMESTAMP()),
		).
		WHERE(
			tMarkers.ID.EQ(mysql.Int64(req.GetId())),
		).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorslivemap.ErrMarkerFailed)
	}

	if err := s.sendUpdateEvent(ctx, MarkerTopic, MarkerDelete, marker.GetJob(), marker); err != nil {
		return nil, errswrap.NewError(err, errorslivemap.ErrMarkerFailed)
	}

	return &pblivemap.DeleteMarkerResponse{}, nil
}

func (s *Server) getMarker(ctx context.Context, id int64) (*livemap.MarkerMarker, error) {
	tUsers := tables.User().AS("user_short")

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

	var dest livemap.MarkerMarker
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, err
	}

	s.enricher.EnrichJobName(&dest)

	return &dest, nil
}

func (s *Server) getMarkerMarkers(
	jobs *permissions.StringList,
) ([]*livemap.MarkerMarker, []int64) {
	updated := []*livemap.MarkerMarker{}
	deleted := []int64{}

	for _, job := range jobs.GetStrings() {
		markers, _ := s.markersCache.Load(job)

		for _, marker := range markers {
			// Make sure marker isn't expired if expiresAt is set
			if marker.GetExpiresAt() == nil || time.Since(marker.GetExpiresAt().AsTime()) < 0 {
				updated = append(updated, marker)
			} else {
				// Just to be sure in regards to cleaning up the client side, add marker id to deleted list
				deleted = append(deleted, marker.GetId())
			}
		}

		// Load the deleted markers list
		deletedMarkers, _ := s.markersDeletedCache.Load(job)
		deleted = append(deleted, deletedMarkers...)
	}

	return updated, deleted
}

func (s *Server) refreshMarkers(ctx context.Context) error {
	tUsers := tables.User().AS("user_short")

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
			tMarkers.ID.ASC(),
		)

	var dest []*livemap.MarkerMarker
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
	}

	markers := map[string][]*livemap.MarkerMarker{}
	for _, m := range dest {
		s.enricher.EnrichJobName(m)

		if _, ok := markers[m.GetJob()]; !ok {
			markers[m.GetJob()] = []*livemap.MarkerMarker{}
		}

		markers[m.GetJob()] = append(markers[m.GetJob()], m)
	}

	for job, ms := range markers {
		if len(ms) == 0 {
			s.markersCache.Delete(job)
		} else {
			s.markersCache.Store(job, ms)
		}
	}

	s.markersCache.Range(func(key string, value []*livemap.MarkerMarker) bool {
		if _, ok := markers[key]; !ok {
			s.markersCache.Delete(key)
		}

		return true
	})

	return s.refreshDeletedMarkers(ctx)
}

func (s *Server) refreshDeletedMarkers(ctx context.Context) error {
	deletedMarkers := map[string][]int64{}

	stmt := tMarkers.
		SELECT(
			tMarkers.ID,
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

	var dest []*livemap.MarkerMarker
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
	}

	for _, m := range dest {
		if _, ok := deletedMarkers[m.GetJob()]; !ok {
			deletedMarkers[m.GetJob()] = []int64{}
		}

		deletedMarkers[m.GetJob()] = append(deletedMarkers[m.GetJob()], m.GetId())
	}

	for job, ms := range deletedMarkers {
		if len(ms) == 0 {
			s.markersDeletedCache.Delete(job)
		} else {
			s.markersDeletedCache.Store(job, ms)
		}
	}

	return nil
}
