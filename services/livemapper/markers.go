package livemapper

import (
	"context"
	"errors"
	"time"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/livemap"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/permissions"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/rector"
	pblivemapper "github.com/fivenet-app/fivenet/gen/go/proto/services/livemapper"
	permslivemapper "github.com/fivenet-app/fivenet/gen/go/proto/services/livemapper/perms"
	"github.com/fivenet-app/fivenet/pkg/access"
	"github.com/fivenet-app/fivenet/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/query/fivenet/model"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	errorslivemapper "github.com/fivenet-app/fivenet/services/livemapper/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var tMarkers = table.FivenetCentrumMarkers.AS("markermarker")

func (s *Server) CreateOrUpdateMarker(ctx context.Context, req *pblivemapper.CreateOrUpdateMarkerRequest) (*pblivemapper.CreateOrUpdateMarkerResponse, error) {
	if req.Marker != nil && req.Marker.Id > 0 {
		trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.livemapper.marker.id", int64(req.Marker.Id)))
	}

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pblivemapper.LivemapperService_ServiceDesc.ServiceName,
		Method:  "CreateOrUpdateMarker",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	// No marker id set
	if req.Marker.Id <= 0 {
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
				req.Marker.ExpiresAt,
				userInfo.Job,
				req.Marker.Name,
				req.Marker.Description,
				req.Marker.X,
				req.Marker.Y,
				req.Marker.Postal,
				req.Marker.Color,
				req.Marker.Type,
				req.Marker.Data,
				userInfo.UserId,
			)

		res, err := stmt.ExecContext(ctx, s.db)
		if err != nil {
			return nil, errswrap.NewError(err, errorslivemapper.ErrMarkerFailed)
		}

		lastId, err := res.LastInsertId()
		if err != nil {
			return nil, errswrap.NewError(err, errorslivemapper.ErrMarkerFailed)
		}

		req.Marker.Id = uint64(lastId)

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)
	} else {
		fields, err := s.ps.AttrStringList(userInfo, permslivemapper.LivemapperServicePerm, permslivemapper.LivemapperServiceCreateOrUpdateMarkerPerm, permslivemapper.LivemapperServiceCreateOrUpdateMarkerAccessPermField)
		if err != nil {
			return nil, errswrap.NewError(err, errorslivemapper.ErrMarkerFailed)
		}

		marker, err := s.getMarker(ctx, req.Marker.Id)
		if err != nil {
			return nil, errswrap.NewError(err, errorslivemapper.ErrMarkerFailed)
		}

		if !access.CheckIfHasAccess(fields, userInfo, marker.Creator.Job, marker.Creator) {
			return nil, errorslivemapper.ErrMarkerDenied
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
				req.Marker.ExpiresAt,
				req.Marker.Name,
				req.Marker.Description,
				req.Marker.X,
				req.Marker.Y,
				req.Marker.Postal,
				req.Marker.Color,
				req.Marker.Type,
				req.Marker.Data,
			).
			WHERE(jet.AND(
				tMarkers.Job.EQ(jet.String(userInfo.Job)),
				tMarkers.ID.EQ(jet.Uint64(req.Marker.Id)),
			))

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, errswrap.NewError(err, errorslivemapper.ErrMarkerFailed)
		}

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)
	}

	marker, err := s.getMarker(ctx, req.Marker.Id)
	if err != nil {
		return nil, errswrap.NewError(err, errorslivemapper.ErrMarkerFailed)
	}

	if err := s.sendUpdateEvent(ctx, MarkerTopic, MarkerUpdate, marker.Job, marker); err != nil {
		return nil, errswrap.NewError(err, errorslivemapper.ErrMarkerFailed)
	}

	return &pblivemapper.CreateOrUpdateMarkerResponse{
		Marker: marker,
	}, nil
}

func (s *Server) DeleteMarker(ctx context.Context, req *pblivemapper.DeleteMarkerRequest) (*pblivemapper.DeleteMarkerResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.livemapper.marker.id", int64(req.Id)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pblivemapper.LivemapperService_ServiceDesc.ServiceName,
		Method:  "DeleteMarker",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	fields, err := s.ps.AttrStringList(userInfo, permslivemapper.LivemapperServicePerm, permslivemapper.LivemapperServiceDeleteMarkerPerm, permslivemapper.LivemapperServiceDeleteMarkerAccessPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorslivemapper.ErrMarkerFailed)
	}

	marker, err := s.getMarker(ctx, req.Id)
	if err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorslivemapper.ErrMarkerFailed)
		}

		return &pblivemapper.DeleteMarkerResponse{}, nil
	}

	if !access.CheckIfHasAccess(fields, userInfo, marker.Creator.Job, marker.Creator) {
		return nil, errorslivemapper.ErrMarkerDenied
	}

	stmt := tMarkers.
		UPDATE(
			tMarkers.DeletedAt,
		).
		SET(
			tMarkers.DeletedAt.SET(jet.CURRENT_TIMESTAMP()),
		).
		WHERE(
			tMarkers.ID.EQ(jet.Uint64(req.Id)),
		).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorslivemapper.ErrMarkerFailed)
	}

	if err := s.sendUpdateEvent(ctx, MarkerTopic, MarkerDelete, marker.Job, marker); err != nil {
		return nil, errswrap.NewError(err, errorslivemapper.ErrMarkerFailed)
	}

	return &pblivemapper.DeleteMarkerResponse{}, nil
}

func (s *Server) getMarker(ctx context.Context, id uint64) (*livemap.MarkerMarker, error) {
	tUsers := tables.Users().AS("user_short")

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
			tMarkers.ID.EQ(jet.Uint64(id)),
		).
		LIMIT(1)

	var dest livemap.MarkerMarker
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, err
	}

	s.enricher.EnrichJobName(&dest)

	return &dest, nil
}

func (s *Server) getMarkerMarkers(jobs *permissions.StringList, updatedAt time.Time) ([]*livemap.MarkerMarker, []uint64, error) {
	updated := []*livemap.MarkerMarker{}
	deleted := []uint64{}

	for _, job := range jobs.Strings {
		markers, _ := s.markersCache.Load(job)

		for _, marker := range markers {
			if updatedAt.IsZero() || marker.UpdatedAt != nil && updatedAt.Sub(marker.UpdatedAt.AsTime()) < 0 {
				// Make sure marker isn't expired if expiresAt is set
				if marker.ExpiresAt == nil || time.Since(marker.ExpiresAt.AsTime()) < 0 {
					updated = append(updated, marker)
				} else {
					// Just to be sure in regards to cleaning up the client side, add marker id to deleted list
					deleted = append(deleted, marker.Id)
				}
			}
		}

		// Load the deleted markers list
		deletedMarkers, _ := s.markersDeletedCache.Load(job)
		deleted = append(deleted, deletedMarkers...)
	}

	return updated, deleted, nil
}

func (s *Server) refreshMarkers(ctx context.Context) error {
	tUsers := tables.Users().AS("user_short")

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
		WHERE(jet.AND(
			tMarkers.DeletedAt.IS_NULL(),
			jet.OR(
				tMarkers.ExpiresAt.IS_NULL(),
				tMarkers.ExpiresAt.GT(jet.CURRENT_TIMESTAMP()),
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

		if _, ok := markers[m.Job]; !ok {
			markers[m.Job] = []*livemap.MarkerMarker{}
		}

		markers[m.Job] = append(markers[m.Job], m)
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
	deletedMarkers := map[string][]uint64{}

	stmt := tMarkers.
		SELECT(
			tMarkers.ID,
		).
		FROM(
			tMarkers,
		).
		WHERE(jet.OR(
			tMarkers.DeletedAt.IS_NOT_NULL(),
			tMarkers.ExpiresAt.LT_EQ(jet.CURRENT_TIMESTAMP()),
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
		if _, ok := deletedMarkers[m.Job]; !ok {
			deletedMarkers[m.Job] = []uint64{}
		}

		deletedMarkers[m.Job] = append(deletedMarkers[m.Job], m.Id)
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
