package livemapper

import (
	"context"
	"slices"

	"github.com/galexrt/fivenet/gen/go/proto/resources/livemap"
	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	permslivemapper "github.com/galexrt/fivenet/gen/go/proto/services/livemapper/perms"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	"github.com/galexrt/fivenet/pkg/grpc/errswrap"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrMarkerFailed = status.Error(codes.Internal, "errors.LivemapperService.ErrMarkerFailed")
	ErrMarkerDenied = status.Error(codes.PermissionDenied, "errors.LivemapperService.ErrMarkerDenied")
)

var (
	tUsers   = table.Users.AS("usershort")
	tMarkers = table.FivenetCentrumMarkers.AS("marker")
)

func (s *Server) checkIfHasAccessToMarker(levels []string, marker *livemap.Marker, userInfo *userinfo.UserInfo) bool {
	if userInfo.SuperUser {
		return true
	}

	if marker.Creator == nil {
		return true
	}

	creator := marker.Creator

	if creator.Job != userInfo.Job {
		return false
	}

	if len(levels) == 0 {
		return creator.UserId == userInfo.UserId
	}

	if slices.Contains(levels, "Any") {
		return true
	}
	if slices.Contains(levels, "Lower_Rank") {
		if creator.JobGrade < userInfo.JobGrade {
			return true
		}
	}
	if slices.Contains(levels, "Same_Rank") {
		if creator.JobGrade <= userInfo.JobGrade {
			return true
		}
	}
	if slices.Contains(levels, "Own") {
		if creator.UserId == userInfo.UserId {
			return true
		}
	}

	return false
}

func (s *Server) CreateOrUpdateMarker(ctx context.Context, req *CreateOrUpdateMarkerRequest) (*CreateOrUpdateMarkerResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: LivemapperService_ServiceDesc.ServiceName,
		Method:  "CreateOrUpdateMarker",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.auditer.Log(auditEntry, req)

	// No marker id set
	if req.Marker.Info.Id <= 0 {
		tMarkers := table.FivenetCentrumMarkers
		stmt := tMarkers.INSERT(
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
		).VALUES(
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
			return nil, errswrap.NewError(ErrMarkerFailed, err)
		}

		lastId, err := res.LastInsertId()
		if err != nil {
			return nil, errswrap.NewError(ErrMarkerFailed, err)
		}

		req.Marker.Info.Id = uint64(lastId)

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)
	} else {
		fieldsAttr, err := s.ps.Attr(userInfo, permslivemapper.LivemapperServicePerm, permslivemapper.LivemapperServiceCreateOrUpdateMarkerPerm, permslivemapper.LivemapperServiceCreateOrUpdateMarkerAccessPermField)
		if err != nil {
			return nil, errswrap.NewError(ErrMarkerFailed, err)
		}
		var fields perms.StringList
		if fieldsAttr != nil {
			fields = fieldsAttr.([]string)
		}

		marker, err := s.getMarker(ctx, req.Marker.Info.Id)
		if err != nil {
			return nil, errswrap.NewError(ErrMarkerFailed, err)
		}

		if !s.checkIfHasAccessToMarker(fields, marker, userInfo) {
			return nil, errswrap.NewError(ErrMarkerFailed, err)
		}

		stmt := tMarkers.
			UPDATE(
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
			return nil, errswrap.NewError(ErrMarkerFailed, err)
		}

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)
	}

	marker, err := s.getMarker(ctx, req.Marker.Info.Id)
	if err != nil {
		return nil, errswrap.NewError(ErrMarkerFailed, err)
	}

	return &CreateOrUpdateMarkerResponse{
		Marker: marker,
	}, nil
}

func (s *Server) DeleteMarker(ctx context.Context, req *DeleteMarkerRequest) (*DeleteMarkerResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: LivemapperService_ServiceDesc.ServiceName,
		Method:  "DeleteMarker",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.auditer.Log(auditEntry, req)

	fieldsAttr, err := s.ps.Attr(userInfo, permslivemapper.LivemapperServicePerm, permslivemapper.LivemapperServiceCreateOrUpdateMarkerPerm, permslivemapper.LivemapperServiceCreateOrUpdateMarkerAccessPermField)
	if err != nil {
		return nil, errswrap.NewError(ErrMarkerFailed, err)
	}
	var fields perms.StringList
	if fieldsAttr != nil {
		fields = fieldsAttr.([]string)
	}

	marker, err := s.getMarker(ctx, req.Id)
	if err != nil {
		return nil, errswrap.NewError(ErrMarkerFailed, err)
	}

	if !s.checkIfHasAccessToMarker(fields, marker, userInfo) {
		return nil, errswrap.NewError(ErrMarkerFailed, err)
	}

	stmt := tMarkers.
		DELETE().
		WHERE(
			tMarkers.ID.EQ(jet.Uint64(req.Id)),
		).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(ErrMarkerFailed, err)
	}

	return &DeleteMarkerResponse{}, nil
}

func (s *Server) getMarker(ctx context.Context, id uint64) (*livemap.Marker, error) {
	stmt := tMarkers.
		SELECT(
			tMarkers.ID.AS("markerinfo.id"),
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
			tUsers.Identifier,
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

	var dest livemap.Marker
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, err
	}

	return &dest, nil
}

func (s *Server) getMarkerMarkers(jobs []string) ([]*livemap.Marker, error) {
	ds := []*livemap.Marker{}

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
	stmt := tMarkers.
		SELECT(
			tMarkers.ID.AS("markerinfo.id"),
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
			tUsers.Identifier,
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
			tMarkers.CreatedAt.GT_EQ(
				jet.CURRENT_TIMESTAMP().SUB(jet.INTERVAL(60, jet.MINUTE)),
			),
		)

	var dest []*livemap.Marker
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return err
	}

	markers := map[string][]*livemap.Marker{}
	for _, job := range s.trackedJobs {
		markers[job] = []*livemap.Marker{}
	}
	for _, m := range dest {
		if _, ok := markers[m.Info.Job]; !ok {
			markers[m.Info.Job] = []*livemap.Marker{}
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

	return nil
}
