package centrum

import (
	"context"

	database "github.com/galexrt/fivenet/gen/go/proto/resources/common/database"
	"github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	errorscentrum "github.com/galexrt/fivenet/gen/go/proto/services/centrum/errors"
	eventscentrum "github.com/galexrt/fivenet/gen/go/proto/services/centrum/events"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"google.golang.org/protobuf/proto"
)

var (
	tUnits      = table.FivenetCentrumUnits.AS("unit")
	tUnitStatus = table.FivenetCentrumUnitsStatus.AS("unitstatus")
	tUsers      = table.Users.AS("usershort")
)

func (s *Server) ListUnits(ctx context.Context, req *ListUnitsRequest) (*ListUnitsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: CentrumService_ServiceDesc.ServiceName,
		Method:  "ListUnits",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.auditer.Log(auditEntry, req)

	resp := &ListUnitsResponse{
		Units: []*dispatch.Unit{},
	}

	units, ok := s.state.ListUnits(userInfo.Job)
	if !ok {
		return nil, errorscentrum.ErrModeForbidsAction
	}

	for i := 0; i < len(units); i++ {
		if len(req.Status) > 0 {
			for _, status := range req.Status {
				if units[i].Status != nil && units[i].Status.Status == status {
					resp.Units = append(resp.Units, units[i])
				}
			}
		} else {
			resp.Units = append(resp.Units, units[i])
		}
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_VIEWED)

	return resp, nil
}

func (s *Server) CreateOrUpdateUnit(ctx context.Context, req *CreateOrUpdateUnitRequest) (*CreateOrUpdateUnitResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: CentrumService_ServiceDesc.ServiceName,
		Method:  "CreateOrUpdateUnit",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.auditer.Log(auditEntry, req)

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errorscentrum.ErrFailedQuery
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	// No unit id set
	if req.Unit.Id <= 0 {
		tUnits := table.FivenetCentrumUnits
		stmt := tUnits.
			INSERT(
				tUnits.Job,
				tUnits.Name,
				tUnits.Initials,
				tUnits.Color,
				tUnits.Description,
			).
			VALUES(
				userInfo.Job,
				req.Unit.Name,
				req.Unit.Initials,
				req.Unit.Color,
				req.Unit.Description,
			)

		result, err := stmt.ExecContext(ctx, tx)
		if err != nil {
			return nil, errorscentrum.ErrFailedQuery
		}

		lastId, err := result.LastInsertId()
		if err != nil {
			return nil, errorscentrum.ErrFailedQuery
		}

		req.Unit.Id = uint64(lastId)

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)
	} else {
		description := ""
		if req.Unit.Description != nil {
			description = *req.Unit.Description
		}

		stmt := tUnits.
			UPDATE(
				tUnits.Name,
				tUnits.Initials,
				tUnits.Color,
				tUnits.Description,
			).
			SET(
				tUnits.Name.SET(jet.String(req.Unit.Name)),
				tUnits.Initials.SET(jet.String(req.Unit.Initials)),
				tUnits.Color.SET(jet.String(req.Unit.Color)),
				tUnits.Description.SET(jet.String(description)),
			).
			WHERE(jet.AND(
				tUnits.Job.EQ(jet.String(userInfo.Job)),
				tUnits.ID.EQ(jet.Uint64(req.Unit.Id)),
			))

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return nil, errorscentrum.ErrFailedQuery
		}

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		auditEntry.State = int16(rector.EventType_EVENT_TYPE_ERRORED)
		return nil, errorscentrum.ErrFailedQuery
	}

	// Load new/updated unit from database
	if err := s.state.LoadUnits(ctx, req.Unit.Id); err != nil {
		return nil, errorscentrum.ErrFailedQuery
	}

	unit, ok := s.state.GetUnit(userInfo.Job, req.Unit.Id)
	if !ok {
		return nil, errorscentrum.ErrFailedQuery
	}

	var x, y *float64
	var postal *string
	marker, ok := s.tracker.GetUserById(userInfo.UserId)
	if ok {
		x = &marker.Info.X
		y = &marker.Info.Y
		postal = marker.Info.Postal
	}

	// A new unit shouldn't have a status, so we make sure it has one
	if unit.Status == nil {
		if err := s.state.UpdateUnitStatus(ctx, userInfo.Job, unit, &dispatch.UnitStatus{
			UnitId:    unit.Id,
			Status:    dispatch.StatusUnit_STATUS_UNIT_UNKNOWN,
			UserId:    &userInfo.UserId,
			CreatorId: &userInfo.UserId,
			X:         x,
			Y:         y,
			Postal:    postal,
		}); err != nil {
			return nil, errorscentrum.ErrFailedQuery
		}
	}

	data, err := proto.Marshal(unit)
	if err != nil {
		return nil, errorscentrum.ErrFailedQuery
	}
	s.events.JS.PublishAsync(eventscentrum.BuildSubject(eventscentrum.TopicUnit, eventscentrum.TypeUnitUpdated, userInfo.Job, unit.Id), data)

	if err := s.state.LoadUnits(ctx, unit.Id); err != nil {
		return nil, errorscentrum.ErrFailedQuery
	}

	return &CreateOrUpdateUnitResponse{
		Unit: unit,
	}, nil
}

func (s *Server) DeleteUnit(ctx context.Context, req *DeleteUnitRequest) (*DeleteUnitResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: CentrumService_ServiceDesc.ServiceName,
		Method:  "DeleteUnit",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.auditer.Log(auditEntry, req)

	resp := &DeleteUnitResponse{}

	unit, ok := s.state.GetUnit(userInfo.Job, req.UnitId)
	if !ok {
		return resp, nil
	}

	stmt := tUnits.
		DELETE().
		WHERE(jet.AND(
			tUnits.Job.EQ(jet.String(userInfo.Job)),
			tUnits.ID.EQ(jet.Uint64(req.UnitId)),
		)).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	data, err := proto.Marshal(unit)
	if err != nil {
		return nil, err
	}
	s.events.JS.PublishAsync(eventscentrum.BuildSubject(eventscentrum.TopicUnit, eventscentrum.TypeUnitDeleted, userInfo.Job, req.UnitId), data)

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)

	units, ok := s.state.Units.Load(userInfo.Job)
	if ok {
		units.Delete(req.UnitId)
	}

	return resp, nil
}

func (s *Server) UpdateUnitStatus(ctx context.Context, req *UpdateUnitStatusRequest) (*UpdateUnitStatusResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: CentrumService_ServiceDesc.ServiceName,
		Method:  "UpdateUnitStatus",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.auditer.Log(auditEntry, req)

	unit, ok := s.state.GetUnit(userInfo.Job, req.UnitId)
	if !ok {
		return nil, errorscentrum.ErrFailedQuery
	}

	if !s.state.CheckIfUserPartOfUnit(userInfo.Job, userInfo.UserId, unit, true) {
		return nil, errorscentrum.ErrNotPartOfUnit
	}

	var x, y *float64
	var postal *string
	marker, ok := s.tracker.GetUserById(userInfo.UserId)
	if ok {
		x = &marker.Info.X
		y = &marker.Info.Y
		postal = marker.Info.Postal
	}

	if err := s.state.UpdateUnitStatus(ctx, userInfo.Job, unit, &dispatch.UnitStatus{
		UnitId:    unit.Id,
		Status:    req.Status,
		Reason:    req.Reason,
		Code:      req.Code,
		UserId:    &userInfo.UserId,
		X:         x,
		Y:         y,
		Postal:    postal,
		CreatorId: &userInfo.UserId,
	}); err != nil {
		return nil, errorscentrum.ErrFailedQuery
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)

	return &UpdateUnitStatusResponse{}, nil
}

func (s *Server) AssignUnit(ctx context.Context, req *AssignUnitRequest) (*AssignUnitResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: CentrumService_ServiceDesc.ServiceName,
		Method:  "AssignUnit",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.auditer.Log(auditEntry, req)

	unit, ok := s.state.GetUnit(userInfo.Job, req.UnitId)
	if !ok {
		return nil, errorscentrum.ErrFailedQuery
	}
	if unit.Job != userInfo.Job {
		return nil, errorscentrum.ErrFailedQuery
	}

	if err := s.state.UpdateUnitAssignments(ctx, userInfo, unit, req.ToAdd, req.ToRemove); err != nil {
		return nil, err
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return &AssignUnitResponse{}, nil
}

func (s *Server) JoinUnit(ctx context.Context, req *JoinUnitRequest) (*JoinUnitResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Check if user is on duty
	if _, ok := s.tracker.GetUserByJobAndID(userInfo.Job, userInfo.UserId); !ok {
		return nil, errorscentrum.ErrNotOnDuty
	}

	unitId, _ := s.state.GetUnitIDForUserID(userInfo.UserId)

	resp := &JoinUnitResponse{}
	// User tries to join his own unit
	if req.UnitId != nil && *req.UnitId == unitId {
		return resp, nil
	}

	currentUnit, _ := s.state.GetUnit(userInfo.Job, unitId)

	// User joins unit
	if req.UnitId != nil && *req.UnitId > 0 {
		// Remove user from his current unit
		if unitId > 0 && currentUnit != nil {
			if err := s.state.UpdateUnitAssignments(ctx, userInfo, currentUnit, nil, []int32{userInfo.UserId}); err != nil {
				return nil, errorscentrum.ErrFailedQuery
			}
		}

		newUnit, ok := s.state.GetUnit(userInfo.Job, *req.UnitId)
		if !ok {
			return nil, errorscentrum.ErrFailedQuery
		}

		if err := s.state.UpdateUnitAssignments(ctx, userInfo, newUnit, []int32{userInfo.UserId}, nil); err != nil {
			return nil, errorscentrum.ErrFailedQuery
		}

		resp.Unit = newUnit
	} else {
		// User leaves his current unit (if he is in an unit)
		if currentUnit != nil {
			if err := s.state.UpdateUnitAssignments(ctx, userInfo, currentUnit, nil, []int32{userInfo.UserId}); err != nil {
				return nil, errorscentrum.ErrFailedQuery
			}
		}
	}

	return resp, nil
}

func (s *Server) ListUnitActivity(ctx context.Context, req *ListUnitActivityRequest) (*ListUnitActivityResponse, error) {
	countStmt := tUnitStatus.
		SELECT(
			jet.COUNT(jet.DISTINCT(tUnitStatus.ID)).AS("datacount.totalcount"),
		).
		FROM(tUnitStatus).
		WHERE(
			tUnitStatus.UnitID.EQ(jet.Uint64(req.Id)),
		)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		return nil, errorscentrum.ErrFailedQuery
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(10)
	resp := &ListUnitActivityResponse{
		Pagination: pag,
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

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
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Dateofbirth,
			tUsers.PhoneNumber,
		).
		FROM(
			tUnitStatus.
				LEFT_JOIN(tUsers,
					tUsers.ID.EQ(tUnitStatus.UserID),
				),
		).
		WHERE(
			tUnitStatus.UnitID.EQ(jet.Uint64(req.Id)),
		).
		ORDER_BY(tUnitStatus.ID.DESC()).
		OFFSET(req.Pagination.Offset).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Activity); err != nil {
		return nil, err
	}

	for i := 0; i < len(resp.Activity); i++ {
		if resp.Activity[i].UnitId > 0 && resp.Activity[i].User != nil {
			unit, ok := s.state.GetUnit(resp.Activity[i].User.Job, resp.Activity[i].UnitId)
			if ok && unit != nil {
				newUnit := proto.Clone(unit)
				resp.Activity[i].Unit = newUnit.(*dispatch.Unit)
			}
		}
	}

	resp.Pagination.Update(count.TotalCount, len(resp.Activity))

	return resp, nil
}
