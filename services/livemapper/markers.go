package livemapper

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/livemap"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/rector"
	pblivemapper "github.com/fivenet-app/fivenet/gen/go/proto/services/livemapper"
	permslivemapper "github.com/fivenet-app/fivenet/gen/go/proto/services/livemapper/perms"
	"github.com/fivenet-app/fivenet/pkg/access"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/pkg/perms"
	"github.com/fivenet-app/fivenet/pkg/utils/dbutils/tables"
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
	if req.Marker != nil && req.Marker.Info != nil && req.Marker.Info.Id < 1 {
		trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.livemapper.marker.id", int64(req.Marker.Info.Id)))
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
	if req.Marker.Info.Id <= 0 {
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
				tMarkers.Icon,
				tMarkers.MarkerType,
				tMarkers.MarkerData,
				tMarkers.CreatorID,
			).
			VALUES(
				req.Marker.ExpiresAt,
				userInfo.Job,
				req.Marker.Info.Name,
				req.Marker.Info.Description,
				req.Marker.Info.X,
				req.Marker.Info.Y,
				req.Marker.Info.Postal,
				req.Marker.Info.Color,
				req.Marker.Info.Icon,
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

		req.Marker.Info.Id = uint64(lastId)

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)
	} else {
		fieldsAttr, err := s.ps.Attr(userInfo, permslivemapper.LivemapperServicePerm, permslivemapper.LivemapperServiceCreateOrUpdateMarkerPerm, permslivemapper.LivemapperServiceCreateOrUpdateMarkerAccessPermField)
		if err != nil {
			return nil, errswrap.NewError(err, errorslivemapper.ErrMarkerFailed)
		}
		var fields perms.StringList
		if fieldsAttr != nil {
			fields = fieldsAttr.([]string)
		}

		marker, err := s.getMarker(ctx, req.Marker.Info.Id)
		if err != nil {
			return nil, errswrap.NewError(err, errorslivemapper.ErrMarkerFailed)
		}

		if !access.CheckIfHasAccess(fields, userInfo, marker.Creator.Job, marker.Creator) {
			return nil, errorslivemapper.ErrMarkerDenied
		}

		stmt := tMarkers.
			UPDATE(
				tMarkers.ExpiresAt,
				tMarkers.Job,
				tMarkers.Name,
				tMarkers.Description,
				tMarkers.X,
				tMarkers.Y,
				tMarkers.Postal,
				tMarkers.Color,
				tMarkers.Icon,
				tMarkers.MarkerType,
				tMarkers.MarkerData,
			).
			SET(
				req.Marker.ExpiresAt,
				userInfo.Job,
				req.Marker.Info.Name,
				req.Marker.Info.Description,
				req.Marker.Info.X,
				req.Marker.Info.Y,
				req.Marker.Info.Postal,
				req.Marker.Info.Color,
				req.Marker.Info.Icon,
				req.Marker.Type,
				req.Marker.Data,
			).
			WHERE(jet.AND(
				tMarkers.Job.EQ(jet.String(userInfo.Job)),
				tMarkers.ID.EQ(jet.Uint64(req.Marker.Info.Id)),
			))

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, errswrap.NewError(err, errorslivemapper.ErrMarkerFailed)
		}

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)
	}

	marker, err := s.getMarker(ctx, req.Marker.Info.Id)
	if err != nil {
		return nil, errswrap.NewError(err, errorslivemapper.ErrMarkerFailed)
	}

	if err := s.sendUpdateEvent(ctx, MarkerUpdate, marker); err != nil {
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

	fieldsAttr, err := s.ps.Attr(userInfo, permslivemapper.LivemapperServicePerm, permslivemapper.LivemapperServiceDeleteMarkerPerm, permslivemapper.LivemapperServiceDeleteMarkerAccessPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorslivemapper.ErrMarkerFailed)
	}
	var fields perms.StringList
	if fieldsAttr != nil {
		fields = fieldsAttr.([]string)
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
		DELETE().
		WHERE(
			tMarkers.ID.EQ(jet.Uint64(req.Id)),
		).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorslivemapper.ErrMarkerFailed)
	}

	if err := s.sendUpdateEvent(ctx, MarkerUpdate, marker); err != nil {
		return nil, errswrap.NewError(err, errorslivemapper.ErrMarkerFailed)
	}

	return &pblivemapper.DeleteMarkerResponse{}, nil
}

func (s *Server) getMarker(ctx context.Context, id uint64) (*livemap.MarkerMarker, error) {
	tUsers := tables.Users().AS("user_short")

	stmt := tMarkers.
		SELECT(
			tMarkers.ID.AS("markerinfo.id"),
			tMarkers.ExpiresAt,
			tMarkers.Job.AS("markerinfo.job"),
			tMarkers.Name.AS("markerinfo.name"),
			tMarkers.Description.AS("markerinfo.description"),
			tMarkers.X.AS("markerinfo.x"),
			tMarkers.Y.AS("markerinfo.y"),
			tMarkers.Postal.AS("markerinfo.postal"),
			tMarkers.Color.AS("markerinfo.color"),
			tMarkers.Icon.AS("markerinfo.icon"),
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

	if dest.Info != nil {
		s.enricher.EnrichJobName(dest.Info)
	}

	return &dest, nil
}

func (s *Server) getMarkerMarkers(jobs []string) ([]*livemap.MarkerMarker, error) {
	ds := []*livemap.MarkerMarker{}

	for _, job := range jobs {
		markers, ok := s.markersCache.Load(job)
		if !ok {
			continue
		}

		ds = append(ds, markers...)
	}

	return ds, nil
}

func (s *Server) refreshMarkers(ctx context.Context) error {
	tUsers := tables.Users().AS("user_short")

	stmt := tMarkers.
		SELECT(
			tMarkers.ID.AS("markerinfo.id"),
			tMarkers.ExpiresAt,
			tMarkers.Job.AS("markerinfo.job"),
			tMarkers.Name.AS("markerinfo.name"),
			tMarkers.Description.AS("markerinfo.description"),
			tMarkers.X.AS("markerinfo.x"),
			tMarkers.Y.AS("markerinfo.y"),
			tMarkers.Postal.AS("markerinfo.postal"),
			tMarkers.Color.AS("markerinfo.color"),
			tMarkers.Icon.AS("markerinfo.icon"),
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
		WHERE(jet.OR(
			tMarkers.ExpiresAt.IS_NULL(),
			tMarkers.ExpiresAt.GT_EQ(jet.NOW()),
		))

	var dest []*livemap.MarkerMarker
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
	}

	markers := map[string][]*livemap.MarkerMarker{}
	for _, job := range s.appCfg.Get().UserTracker.GetLivemapJobs() {
		markers[job] = []*livemap.MarkerMarker{}
	}

	for _, m := range dest {
		if m.Info != nil {
			s.enricher.EnrichJobName(m.Info)
		}

		if _, ok := markers[m.Info.Job]; !ok {
			markers[m.Info.Job] = []*livemap.MarkerMarker{}
		}

		markers[m.Info.Job] = append(markers[m.Info.Job], m)
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

	return nil
}
