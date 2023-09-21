package livemapper

import (
	"context"

	"github.com/galexrt/fivenet/gen/go/proto/resources/livemap"
	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/pkg/utils"
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
	tMarkers = table.FivenetCentrumMarkers
)

func (s *Server) checkIfHasAccessToMarker(levels []string, marker *livemap.Marker, userInfo *userinfo.UserInfo) bool {
	if userInfo.SuperUser {
		return true
	}

	if marker.Creator == nil {
		return true
	}

	creator := marker.Creator

	if len(levels) == 0 {
		return creator.UserId == userInfo.UserId
	}

	if utils.InSlice(levels, "Lower_Rank") {
		if creator.JobGrade < userInfo.JobGrade {
			return true
		}
	}
	if utils.InSlice(levels, "Same_Rank") {
		if creator.JobGrade <= userInfo.JobGrade {
			return true
		}
	}
	if utils.InSlice(levels, "Own") {
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
		)

		res, err := stmt.ExecContext(ctx, s.db)
		if err != nil {
			return nil, err
		}

		lastId, err := res.LastInsertId()
		if err != nil {
			return nil, err
		}

		req.Marker.Info.Id = uint64(lastId)

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)
	} else {
		fieldsAttr, err := s.ps.Attr(userInfo, LivemapperServicePerm, LivemapperServiceCreateOrUpdateMarkerPerm, LivemapperServiceCreateOrUpdateMarkerAccessPermField)
		if err != nil {
			return nil, ErrMarkerFailed
		}
		var fields perms.StringList
		if fieldsAttr != nil {
			fields = fieldsAttr.([]string)
		}

		marker, err := s.getMarker(ctx, req.Marker.Info.Id)
		if err != nil {
			return nil, ErrMarkerFailed
		}

		if !s.checkIfHasAccessToMarker(fields, marker, userInfo) {
			return nil, ErrMarkerFailed
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
			return nil, err
		}

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)
	}

	marker, err := s.getMarker(ctx, req.Marker.Info.Id)
	if err != nil {
		return nil, ErrMarkerFailed
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

	fieldsAttr, err := s.ps.Attr(userInfo, LivemapperServicePerm, LivemapperServiceCreateOrUpdateMarkerPerm, LivemapperServiceCreateOrUpdateMarkerAccessPermField)
	if err != nil {
		return nil, ErrMarkerFailed
	}
	var fields perms.StringList
	if fieldsAttr != nil {
		fields = fieldsAttr.([]string)
	}

	marker, err := s.getMarker(ctx, req.Id)
	if err != nil {
		return nil, ErrMarkerFailed
	}

	if !s.checkIfHasAccessToMarker(fields, marker, userInfo) {
		return nil, ErrMarkerFailed
	}

	stmt := tMarkers.
		DELETE().
		WHERE(
			tMarkers.ID.EQ(jet.Uint64(req.Id)),
		).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	return &DeleteMarkerResponse{}, nil
}

func (s *Server) getMarker(ctx context.Context, id uint64) (*livemap.Marker, error) {
	tUsers := tUsers.AS("creator")
	tMarkers := tMarkers.AS("marker")
	stmt := tMarkers.
		SELECT(
			tMarkers.ID.AS("markerinfo.id"),
			tMarkers.Job.AS("markerinfo.job"),
			tMarkers.Name.AS("markerinfo.name"),
			tMarkers.Description.AS("markerinfo.description"),
			tMarkers.X.AS("markerinfo.x"),
			tMarkers.Y.AS("markerinfo.y"),
			tMarkers.Postal.AS("markerinfo.postal"),
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
