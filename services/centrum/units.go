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
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorscentrum "github.com/fivenet-app/fivenet/v2025/services/centrum/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
)

var (
	tUnitStatus     = table.FivenetCentrumUnitsStatus.AS("unit_status")
	tUserProps      = table.FivenetUserProps
	tUnits          = table.FivenetCentrumUnits.AS("unit")
	tColleagueProps = table.FivenetJobColleagueProps.AS("colleague_props")
)

func (s *Server) ListUnits(
	ctx context.Context,
	req *pbcentrum.ListUnitsRequest,
) (*pbcentrum.ListUnitsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbcentrum.CentrumService_ServiceDesc.ServiceName,
		Method:  "ListUnits",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	resp := &pbcentrum.ListUnitsResponse{
		Units: []*centrum.Unit{},
	}

	resp.Units = s.units.Filter(ctx, []string{userInfo.GetJob()}, req.GetStatus(), nil, nil)
	if resp.Units == nil {
		return nil, errorscentrum.ErrFailedQuery
	}

	// Resolve qualifications access for the user
	for _, unit := range resp.GetUnits() {
		if unit.GetAccess() != nil && len(unit.GetAccess().GetQualifications()) > 0 {
			qualificationsAccess, err := s.units.GetAccess().Qualifications.List(
				ctx,
				s.db,
				unit.GetId(),
			)
			if err != nil {
				return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
			}
			unit.Access.Qualifications = qualificationsAccess
		}
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_VIEWED

	return resp, nil
}

func (s *Server) CreateOrUpdateUnit(
	ctx context.Context,
	req *pbcentrum.CreateOrUpdateUnitRequest,
) (*pbcentrum.CreateOrUpdateUnitResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbcentrum.CentrumService_ServiceDesc.ServiceName,
		Method:  "CreateOrUpdateUnit",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	req.Unit.Job = userInfo.GetJob()

	var unit *centrum.Unit
	var err error
	// No unit id set
	if req.GetUnit().GetId() <= 0 {
		unit, err = s.units.CreateUnit(ctx, userInfo.GetJob(), req.GetUnit())
		if err != nil {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}

		auditEntry.State = audit.EventType_EVENT_TYPE_CREATED
	} else {
		// Check that the unit belongs to the user's job
		u, err := s.units.Get(ctx, req.GetUnit().GetId())
		if err != nil {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}
		if u.GetJob() != userInfo.GetJob() {
			return nil, errorscentrum.ErrUnitPermDenied
		}

		unit, err = s.units.Update(ctx, req.GetUnit())
		if err != nil {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}

		auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED
	}

	return &pbcentrum.CreateOrUpdateUnitResponse{
		Unit: unit,
	}, nil
}

func (s *Server) DeleteUnit(
	ctx context.Context,
	req *pbcentrum.DeleteUnitRequest,
) (*pbcentrum.DeleteUnitResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.centrum.unit_id", req.GetUnitId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbcentrum.CentrumService_ServiceDesc.ServiceName,
		Method:  "DeleteUnit",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	// Check that the unit belongs to the user's job
	unit, err := s.units.Get(ctx, req.GetUnitId())
	if err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}
	if unit.GetJob() != userInfo.GetJob() {
		return nil, errorscentrum.ErrUnitPermDenied
	}

	resp := &pbcentrum.DeleteUnitResponse{}

	if err := s.units.Delete(ctx, req.GetUnitId()); err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_DELETED

	return resp, nil
}

func (s *Server) UpdateUnitStatus(
	ctx context.Context,
	req *pbcentrum.UpdateUnitStatusRequest,
) (*pbcentrum.UpdateUnitStatusResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.centrum.unit_id", req.GetUnitId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbcentrum.CentrumService_ServiceDesc.ServiceName,
		Method:  "UpdateUnitStatus",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	// Check that the user has access to the unit
	unit, err := s.units.Get(ctx, req.GetUnitId())
	if err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	// Only check unit access when not empty
	if unit.GetAccess() != nil && !unit.GetAccess().IsEmpty() {
		// Make sure requestor is not a dispatcher
		if !s.helpers.CheckIfUserIsDispatcher(ctx, userInfo.GetJob(), userInfo.GetUserId()) {
			check, err := s.units.GetAccess().
				CanUserAccessTarget(ctx, unit.GetId(), userInfo, centrum.UnitAccessLevel_UNIT_ACCESS_LEVEL_JOIN)
			if err != nil {
				return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
			}
			if !check {
				return nil, errorscentrum.ErrUnitPermDenied
			}
		}
	} else if unit.GetJob() != userInfo.GetJob() {
		return nil, errorscentrum.ErrUnitPermDenied
	}

	if !s.helpers.CheckIfUserPartOfUnit(ctx, userInfo.GetJob(), userInfo.GetUserId(), unit, true) {
		return nil, errorscentrum.ErrNotPartOfUnit
	}

	if _, err := s.units.UpdateStatus(ctx, unit.GetId(), &centrum.UnitStatus{
		CreatedAt:  timestamp.Now(),
		UnitId:     unit.GetId(),
		Unit:       unit,
		Status:     req.GetStatus(),
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

func (s *Server) AssignUnit(
	ctx context.Context,
	req *pbcentrum.AssignUnitRequest,
) (*pbcentrum.AssignUnitResponse, error) {
	logging.InjectFields(ctx, logging.Fields{
		"fivenet.centrum.unit_id", req.GetUnitId(),
		"fivenet.centrum.users.to_add", req.GetToAdd(),
		"fivenet.centrum.users.to_remove", req.GetToRemove(),
	})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbcentrum.CentrumService_ServiceDesc.ServiceName,
		Method:  "AssignUnit",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	unit, err := s.units.Get(ctx, req.GetUnitId())
	if err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}
	if unit.GetJob() != userInfo.GetJob() {
		return nil, errorscentrum.ErrUnitPermDenied
	}

	// Only check unit access when not empty
	if unit.GetAccess() != nil && !unit.GetAccess().IsEmpty() {
		// Make sure requestor is not a dispatcher
		if !s.helpers.CheckIfUserIsDispatcher(ctx, userInfo.GetJob(), userInfo.GetUserId()) {
			check, err := s.units.GetAccess().
				CanUserAccessTarget(ctx, unit.GetId(), userInfo, centrum.UnitAccessLevel_UNIT_ACCESS_LEVEL_JOIN)
			if err != nil {
				return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
			}
			if !check {
				return nil, errorscentrum.ErrUnitPermDenied
			}
		}
	}

	if err := s.units.UpdateUnitAssignments(ctx, userInfo.GetJob(), &userInfo.UserId, unit.GetId(), req.GetToAdd(), req.GetToRemove()); err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED

	return &pbcentrum.AssignUnitResponse{}, nil
}

func (s *Server) JoinUnit(
	ctx context.Context,
	req *pbcentrum.JoinUnitRequest,
) (*pbcentrum.JoinUnitResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if req.UnitId != nil {
		logging.InjectFields(ctx, logging.Fields{"fivenet.centrum.unit_id", req.GetUnitId()})
	}

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Check if user is on duty, if not make sure to unset any unit id
	if um, ok := s.tracker.GetUserMarkerById(userInfo.GetUserId()); !ok || um.GetHidden() {
		if err := s.tracker.SetUserMappingForUser(ctx, userInfo.GetUserId(), nil); err != nil {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}

		return nil, errorscentrum.ErrNotOnDuty
	}

	currentUnitMapping, _ := s.tracker.GetUserMapping(userInfo.GetUserId())

	var currentUnit *centrum.Unit
	if currentUnitMapping != nil && currentUnitMapping.UnitId != nil &&
		currentUnitMapping.GetUnitId() > 0 {
		var err error
		currentUnit, err = s.units.Get(ctx, currentUnitMapping.GetUnitId())
		if err != nil && !errors.Is(err, jetstream.ErrKeyNotFound) {
			return nil, errorscentrum.ErrNotOnDuty
		}
	}

	resp := &pbcentrum.JoinUnitResponse{}
	// User tries to join his own unit
	if req.UnitId != nil && currentUnitMapping != nil && currentUnitMapping.UnitId != nil &&
		req.GetUnitId() == currentUnitMapping.GetUnitId() {
		resp.Unit = currentUnit
		return resp, nil
	}

	// User joins unit
	if req.UnitId != nil && req.GetUnitId() > 0 {
		s.logger.Debug(
			"user joining unit",
			zap.String("job", userInfo.GetJob()),
			zap.Int32("user_id", userInfo.GetUserId()),
			zap.Int64("current_unit_id", currentUnitMapping.GetUnitId()),
			zap.Int64("unit_id", req.GetUnitId()),
		)

		// Remove user from his current unit
		if currentUnit != nil {
			if err := s.units.UpdateUnitAssignments(ctx, userInfo.GetJob(), &userInfo.UserId, currentUnit.GetId(), nil, []int32{userInfo.GetUserId()}); err != nil {
				return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
			}
		}

		newUnit, err := s.units.Get(ctx, req.GetUnitId())
		if err != nil {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}

		if newUnit.GetJob() != userInfo.GetJob() {
			return nil, errorscentrum.ErrUnitPermDenied
		}

		// Only check unit access when not empty
		if newUnit.GetAccess() != nil && !newUnit.GetAccess().IsEmpty() {
			// Make sure requestor is not a dispatcher
			if !s.helpers.CheckIfUserIsDispatcher(ctx, userInfo.GetJob(), userInfo.GetUserId()) {
				check, err := s.units.GetAccess().
					CanUserAccessTarget(ctx, newUnit.GetId(), userInfo, centrum.UnitAccessLevel_UNIT_ACCESS_LEVEL_JOIN)
				if err != nil {
					return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
				}
				if !check {
					return nil, errorscentrum.ErrUnitPermDenied
				}
			}
		}

		if err := s.units.UpdateUnitAssignments(ctx, userInfo.GetJob(), &userInfo.UserId, newUnit.GetId(), []int32{userInfo.GetUserId()}, nil); err != nil {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}

		resp.Unit = newUnit
	} else {
		s.logger.Debug("user leaving unit", zap.Int64("current_unit_id", currentUnitMapping.GetUnitId()), zap.Int64("unit_id", req.GetUnitId()))
		// User leaves his current unit (if he is in an unit)
		if currentUnit != nil {
			if err := s.units.UpdateUnitAssignments(ctx, userInfo.GetJob(), &userInfo.UserId, currentUnit.GetId(), nil, []int32{userInfo.GetUserId()}); err != nil {
				return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
			}
		}
	}

	return resp, nil
}

func (s *Server) ListUnitActivity(
	ctx context.Context,
	req *pbcentrum.ListUnitActivityRequest,
) (*pbcentrum.ListUnitActivityResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.centrum.unit_id", req.GetId()})

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
			tUnitStatus.UnitID.EQ(jet.Int64(req.GetId())),
			tUnits.Job.EQ(jet.String(userInfo.GetJob())),
		))

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}
	}

	pag, limit := req.GetPagination().GetResponseWithPageSize(count.Total, 10)
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
						AND(tColleague.Job.EQ(jet.String(userInfo.GetJob()))),
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
			tUnitStatus.UnitID.EQ(jet.Int64(req.GetId())),
		).
		ORDER_BY(tUnitStatus.ID.DESC()).
		OFFSET(req.GetPagination().GetOffset()).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Activity); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.GetActivity() {
		if resp.GetActivity()[i].GetUnitId() > 0 && resp.GetActivity()[i].GetUser() != nil {
			unit, err := s.units.Get(ctx, resp.GetActivity()[i].GetUnitId())
			if err != nil {
				return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
			}

			if unit.GetJob() != userInfo.GetJob() {
				return nil, errorscentrum.ErrUnitPermDenied
			}

			resp.Activity[i].Unit = unit
		}

		if resp.GetActivity()[i].GetUser() != nil {
			jobInfoFn(resp.GetActivity()[i].GetUser())
		}
	}

	resp.GetPagination().Update(len(resp.GetActivity()))

	return resp, nil
}
