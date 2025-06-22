package centrum

import (
	"context"
	"errors"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	centrum "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	database "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	pbcentrum "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/centrum"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorscentrum "github.com/fivenet-app/fivenet/v2025/services/centrum/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var (
	tUnitStatus     = table.FivenetCentrumUnitsStatus.AS("unit_status")
	tUserProps      = table.FivenetUserProps
	tUnits          = table.FivenetCentrumUnits.AS("unit")
	tColleagueProps = table.FivenetJobColleagueProps.AS("colleague_props")
)

func (s *Server) ListUnits(ctx context.Context, req *pbcentrum.ListUnitsRequest) (*pbcentrum.ListUnitsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbcentrum.CentrumService_ServiceDesc.ServiceName,
		Method:  "ListUnits",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	resp := &pbcentrum.ListUnitsResponse{
		Units: []*centrum.Unit{},
	}

	resp.Units = s.units.Filter(ctx, []string{userInfo.Job}, req.Status, nil, nil)
	if resp.Units == nil {
		return nil, errorscentrum.ErrFailedQuery
	}

	// Resolve qualifications access for the user
	for _, unit := range resp.Units {
		if unit.Access != nil && len(unit.Access.Qualifications) > 0 {
			qualificationsAccess, err := s.units.GetAccess().Qualifications.List(ctx, s.db, unit.Id)
			if err != nil {
				return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
			}
			unit.Access.Qualifications = qualificationsAccess
		}
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_VIEWED

	return resp, nil
}

func (s *Server) CreateOrUpdateUnit(ctx context.Context, req *pbcentrum.CreateOrUpdateUnitRequest) (*pbcentrum.CreateOrUpdateUnitResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbcentrum.CentrumService_ServiceDesc.ServiceName,
		Method:  "CreateOrUpdateUnit",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	req.Unit.Job = userInfo.Job

	var unit *centrum.Unit
	var err error
	// No unit id set
	if req.Unit.Id <= 0 {
		unit, err = s.units.CreateUnit(ctx, userInfo.Job, req.Unit)
		if err != nil {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}

		auditEntry.State = audit.EventType_EVENT_TYPE_CREATED
	} else {
		// Check that the unit belongs to the user's job
		u, err := s.units.Get(ctx, req.Unit.Id)
		if err != nil {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}
		if u.Job != userInfo.Job {
			return nil, errorscentrum.ErrUnitPermDenied
		}

		unit, err = s.units.Update(ctx, req.Unit)
		if err != nil {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}

		auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED
	}

	return &pbcentrum.CreateOrUpdateUnitResponse{
		Unit: unit,
	}, nil
}

func (s *Server) DeleteUnit(ctx context.Context, req *pbcentrum.DeleteUnitRequest) (*pbcentrum.DeleteUnitResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.centrum.unit_id", int64(req.UnitId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbcentrum.CentrumService_ServiceDesc.ServiceName,
		Method:  "DeleteUnit",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	// Check that the unit belongs to the user's job
	unit, err := s.units.Get(ctx, req.UnitId)
	if err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}
	if unit.Job != userInfo.Job {
		return nil, errorscentrum.ErrUnitPermDenied
	}

	resp := &pbcentrum.DeleteUnitResponse{}

	if err := s.units.Delete(ctx, req.UnitId); err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_DELETED

	return resp, nil
}

func (s *Server) UpdateUnitStatus(ctx context.Context, req *pbcentrum.UpdateUnitStatusRequest) (*pbcentrum.UpdateUnitStatusResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.centrum.unit_id", int64(req.UnitId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbcentrum.CentrumService_ServiceDesc.ServiceName,
		Method:  "UpdateUnitStatus",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	// Check that the user has access to the unit
	unit, err := s.units.Get(ctx, req.UnitId)
	if err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	// Only check unit access when not empty
	if unit.Access != nil && !unit.Access.IsEmpty() {
		// Make sure requestor is not a dispatcher
		if !s.helpers.CheckIfUserIsDispatcher(ctx, userInfo.Job, userInfo.UserId) {
			check, err := s.units.GetAccess().CanUserAccessTarget(ctx, unit.Id, userInfo, centrum.UnitAccessLevel_UNIT_ACCESS_LEVEL_JOIN)
			if err != nil {
				return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
			}
			if !check {
				return nil, errorscentrum.ErrUnitPermDenied
			}
		}
	} else if unit.Job != userInfo.Job {
		return nil, errorscentrum.ErrUnitPermDenied
	}

	if !s.helpers.CheckIfUserPartOfUnit(ctx, userInfo.Job, userInfo.UserId, unit, true) {
		return nil, errorscentrum.ErrNotPartOfUnit
	}

	if _, err := s.units.UpdateStatus(ctx, unit.Id, &centrum.UnitStatus{
		CreatedAt:  timestamp.Now(),
		UnitId:     unit.Id,
		Unit:       unit,
		Status:     req.Status,
		Reason:     req.Reason,
		Code:       req.Code,
		UserId:     &userInfo.UserId,
		CreatorId:  &userInfo.UserId,
		CreatorJob: &userInfo.Job,
	}); err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_CREATED

	return &pbcentrum.UpdateUnitStatusResponse{}, nil
}

func (s *Server) AssignUnit(ctx context.Context, req *pbcentrum.AssignUnitRequest) (*pbcentrum.AssignUnitResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.centrum.unit_id", int64(req.UnitId)))
	trace.SpanFromContext(ctx).SetAttributes(attribute.IntSlice("fivenet.centrum.users.to_add", utils.SliceInt32ToInt(req.ToAdd)))
	trace.SpanFromContext(ctx).SetAttributes(attribute.IntSlice("fivenet.centrum.users.to_remove", utils.SliceInt32ToInt(req.ToRemove)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbcentrum.CentrumService_ServiceDesc.ServiceName,
		Method:  "AssignUnit",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	unit, err := s.units.Get(ctx, req.UnitId)
	if err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}
	if unit.Job != userInfo.Job {
		return nil, errorscentrum.ErrUnitPermDenied
	}

	// Only check unit access when not empty
	if unit.Access != nil && !unit.Access.IsEmpty() {
		// Make sure requestor is not a dispatcher
		if !s.helpers.CheckIfUserIsDispatcher(ctx, userInfo.Job, userInfo.UserId) {
			check, err := s.units.GetAccess().CanUserAccessTarget(ctx, unit.Id, userInfo, centrum.UnitAccessLevel_UNIT_ACCESS_LEVEL_JOIN)
			if err != nil {
				return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
			}
			if !check {
				return nil, errorscentrum.ErrUnitPermDenied
			}
		}
	}

	if err := s.units.UpdateUnitAssignments(ctx, userInfo.Job, &userInfo.UserId, unit.Id, req.ToAdd, req.ToRemove); err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED

	return &pbcentrum.AssignUnitResponse{}, nil
}

func (s *Server) JoinUnit(ctx context.Context, req *pbcentrum.JoinUnitRequest) (*pbcentrum.JoinUnitResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if req.UnitId != nil {
		trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.centrum.unit_id", int64(*req.UnitId)))
	}

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Check if user is on duty, if not make sure to unset any unit id
	if um, ok := s.tracker.GetUserMarkerById(userInfo.UserId); !ok || um.Hidden {
		if err := s.tracker.SetUserMappingForUser(ctx, userInfo.UserId, nil); err != nil {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}

		return nil, errorscentrum.ErrNotOnDuty
	}

	currentUnitMapping, _ := s.tracker.GetUserMapping(userInfo.UserId)

	var currentUnit *centrum.Unit
	if currentUnitMapping != nil && currentUnitMapping.UnitId != nil && *currentUnitMapping.UnitId > 0 {
		var err error
		currentUnit, err = s.units.Get(ctx, *currentUnitMapping.UnitId)
		if err != nil && !errors.Is(err, jetstream.ErrKeyNotFound) {
			return nil, errorscentrum.ErrNotOnDuty
		}
	}

	resp := &pbcentrum.JoinUnitResponse{}
	// User tries to join his own unit
	if req.UnitId != nil && currentUnitMapping != nil && currentUnitMapping.UnitId != nil && *req.UnitId == *currentUnitMapping.UnitId {
		resp.Unit = currentUnit
		return resp, nil
	}

	// User joins unit
	if req.UnitId != nil && *req.UnitId > 0 {
		s.logger.Debug("user joining unit", zap.String("job", userInfo.Job), zap.Int32("user_id", userInfo.UserId), zap.Uint64p("current_unit_id", currentUnitMapping.UnitId), zap.Uint64p("unit_id", req.UnitId))

		// Remove user from his current unit
		if currentUnit != nil {
			if err := s.units.UpdateUnitAssignments(ctx, userInfo.Job, &userInfo.UserId, currentUnit.Id, nil, []int32{userInfo.UserId}); err != nil {
				return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
			}
		}

		newUnit, err := s.units.Get(ctx, *req.UnitId)
		if err != nil {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}

		if newUnit.Job != userInfo.Job {
			return nil, errorscentrum.ErrUnitPermDenied
		}

		// Only check unit access when not empty
		if newUnit.Access != nil && !newUnit.Access.IsEmpty() {
			// Make sure requestor is not a dispatcher
			if !s.helpers.CheckIfUserIsDispatcher(ctx, userInfo.Job, userInfo.UserId) {
				check, err := s.units.GetAccess().CanUserAccessTarget(ctx, newUnit.Id, userInfo, centrum.UnitAccessLevel_UNIT_ACCESS_LEVEL_JOIN)
				if err != nil {
					return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
				}
				if !check {
					return nil, errorscentrum.ErrUnitPermDenied
				}
			}
		}

		if err := s.units.UpdateUnitAssignments(ctx, userInfo.Job, &userInfo.UserId, newUnit.Id, []int32{userInfo.UserId}, nil); err != nil {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}

		resp.Unit = newUnit
	} else {
		s.logger.Debug("user leaving unit", zap.Uint64p("current_unit_id", currentUnitMapping.UnitId), zap.Uint64p("unit_id", req.UnitId))
		// User leaves his current unit (if he is in an unit)
		if currentUnit != nil {
			if err := s.units.UpdateUnitAssignments(ctx, userInfo.Job, &userInfo.UserId, currentUnit.Id, nil, []int32{userInfo.UserId}); err != nil {
				return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
			}
		}
	}

	return resp, nil
}

func (s *Server) ListUnitActivity(ctx context.Context, req *pbcentrum.ListUnitActivityRequest) (*pbcentrum.ListUnitActivityResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.centrum.unit_id", int64(req.Id)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	countStmt := tUnitStatus.
		SELECT(
			jet.COUNT(jet.DISTINCT(tUnitStatus.ID)).AS("data_count.total"),
		).
		FROM(
			tUnitStatus.
				INNER_JOIN(tUnits,
					tUnits.ID.EQ(tUnitStatus.UnitID),
				),
		).
		WHERE(jet.AND(
			tUnitStatus.UnitID.EQ(jet.Uint64(req.Id)),
			tUnits.Job.EQ(jet.String(userInfo.Job)),
		))

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(count.Total, 10)
	resp := &pbcentrum.ListUnitActivityResponse{
		Pagination: pag,
	}
	if count.Total <= 0 {
		return resp, nil
	}

	tColleague := tables.User().AS("colleague")
	tAvatar := table.FivenetFiles.AS("avatar")

	stmt := tUnitStatus.
		SELECT(
			tUnitStatus.ID,
			tUnitStatus.CreatedAt,
			tUnitStatus.UnitID,
			tUnitStatus.Status,
			tUnitStatus.Reason,
			tUnitStatus.Code,
			tUnitStatus.UserID,
			tUnitStatus.CreatorID,
			tUnitStatus.X,
			tUnitStatus.Y,
			tUnitStatus.Postal,
			tUnitStatus.CreatorJob,
			tColleague.ID,
			tColleague.Firstname,
			tColleague.Lastname,
			tColleague.Job,
			tColleague.JobGrade,
			tColleague.Sex,
			tColleague.Dateofbirth,
			tColleague.PhoneNumber,
			tColleagueProps.UserID,
			tColleagueProps.Job,
			tColleagueProps.NamePrefix,
			tColleagueProps.NameSuffix,
			tUserProps.AvatarFileID.AS("colleague.avatar_file_id"),
			tAvatar.FilePath.AS("colleague.avatar"),
		).
		FROM(
			tUnitStatus.
				LEFT_JOIN(tColleague,
					tColleague.ID.EQ(tUnitStatus.UserID),
				).
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tUnitStatus.UserID).
						AND(tColleague.Job.EQ(jet.String(userInfo.Job))),
				).
				LEFT_JOIN(tColleagueProps,
					tColleagueProps.UserID.EQ(tColleague.ID).
						AND(tColleagueProps.Job.EQ(tColleague.Job)),
				).
				LEFT_JOIN(tAvatar,
					tAvatar.ID.EQ(tUserProps.AvatarFileID),
				),
		).
		WHERE(
			tUnitStatus.UnitID.EQ(jet.Uint64(req.Id)),
		).
		ORDER_BY(tUnitStatus.ID.DESC()).
		OFFSET(req.Pagination.Offset).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Activity); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.Activity {
		if resp.Activity[i].UnitId > 0 && resp.Activity[i].User != nil {
			unit, err := s.units.Get(ctx, resp.Activity[i].UnitId)
			if err != nil {
				return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
			}

			if unit.Job != userInfo.Job {
				return nil, errorscentrum.ErrUnitPermDenied
			}

			resp.Activity[i].Unit = unit
		}

		if resp.Activity[i].User != nil {
			jobInfoFn(resp.Activity[i].User)
		}
	}

	resp.Pagination.Update(len(resp.Activity))

	return resp, nil
}
