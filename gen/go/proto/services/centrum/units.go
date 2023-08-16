package centrum

import (
	"context"

	database "github.com/galexrt/fivenet/gen/go/proto/resources/common/database"
	"github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

var (
	tUnits      = table.FivenetCentrumUnits.AS("unit")
	tUnitStatus = table.FivenetCentrumUnitsStatus.AS("unitstatus")
	tUnitUser   = table.FivenetCentrumUnitsUsers.AS("unitassignment")
	tUsers      = table.Users.AS("usershort")
)

func (s *Server) ListUnits(ctx context.Context, req *ListUnitsRequest) (*ListUnitsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: CentrumService_ServiceDesc.ServiceName,
		Method:  "ListUnits",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EVENT_TYPE_ERRORED),
	}
	defer s.auditer.Log(auditEntry, req)

	resp := &ListUnitsResponse{
		Units: []*dispatch.Unit{},
	}

	if req.OwnOnly != nil && *req.OwnOnly {
		unitId, ok := s.getUnitIDForUserID(userInfo.UserId)
		if !ok {
			return nil, ErrFailedQuery
		}

		unit, ok := s.getUnit(userInfo.Job, unitId)
		if !ok {
			return nil, ErrFailedQuery
		}

		resp.Units = append(resp.Units, unit)
	} else {
		units, err := s.listUnits(userInfo.Job)
		if err != nil {
			return nil, err
		}

		for i := 0; i < len(units); i++ {
			if units[i].Status != nil {
				for _, status := range req.Status {
					if units[i].Status.Status == status {
						resp.Units = append(resp.Units, units[i])
					}
				}
			}
		}
	}

	auditEntry.State = int16(rector.EVENT_TYPE_VIEWED)

	return resp, nil
}

func (s *Server) CreateOrUpdateUnit(ctx context.Context, req *CreateOrUpdateUnitRequest) (*CreateOrUpdateUnitResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: CentrumService_ServiceDesc.ServiceName,
		Method:  "CreateOrUpdateUnit",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EVENT_TYPE_ERRORED),
	}
	defer s.auditer.Log(auditEntry, req)

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, ErrFailedQuery
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
			return nil, err
		}

		lastId, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		req.Unit.Id = uint64(lastId)

		if err := s.updateUnitStatus(ctx, userInfo, &dispatch.UnitStatus{
			UnitId: uint64(lastId),
		}); err != nil {
			return nil, err
		}

		auditEntry.State = int16(rector.EVENT_TYPE_CREATED)
	} else {
		stmt := tUnits.
			UPDATE(
				tUnits.Name,
				tUnits.Initials,
				tUnits.Color,
				tUnits.Description,
			).
			SET(
				req.Unit.Name,
				req.Unit.Initials,
				req.Unit.Color,
				req.Unit.Description,
			).
			WHERE(jet.AND(
				tUnits.Job.EQ(jet.String(userInfo.Job)),
				tUnits.ID.EQ(jet.Uint64(req.Unit.Id)),
			))

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return nil, err
		}

		auditEntry.State = int16(rector.EVENT_TYPE_UPDATED)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		auditEntry.State = int16(rector.EVENT_TYPE_ERRORED)
		return nil, ErrFailedQuery
	}

	// Load new/updated unit from database
	if err := s.loadUnits(ctx, req.Unit.Id); err != nil {
		return nil, err
	}

	unit, ok := s.getUnit(userInfo.Job, req.Unit.Id)
	if !ok {
		return nil, ErrFailedQuery
	}

	data, err := proto.Marshal(unit)
	if err != nil {
		return nil, err
	}
	s.events.JS.Publish(s.buildSubject(TopicUnit, TypeUnitUpdated, userInfo.Job, unit.Id), data)

	if err := s.loadUnits(ctx, unit.Id); err != nil {
		return nil, err
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
		State:   int16(rector.EVENT_TYPE_ERRORED),
	}
	defer s.auditer.Log(auditEntry, req)

	resp := &DeleteUnitResponse{}

	unit, ok := s.getUnit(userInfo.Job, req.UnitId)
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
	s.events.JS.Publish(s.buildSubject(TopicUnit, TypeUnitDeleted, userInfo.Job, req.UnitId), data)

	auditEntry.State = int16(rector.EVENT_TYPE_DELETED)

	units, ok := s.units.Load(userInfo.Job)
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
		State:   int16(rector.EVENT_TYPE_ERRORED),
	}
	defer s.auditer.Log(auditEntry, req)

	unit, ok := s.getUnit(userInfo.Job, req.UnitId)
	if !ok {
		return nil, ErrFailedQuery
	}

	can := s.ps.Can(userInfo, CentrumServicePerm, CentrumServiceDeleteUnitPerm)
	if !can {
		if !s.checkIfUserPartOfUnit(userInfo.UserId, unit) {
			return nil, ErrFailedQuery
		}
	}

	var x, y *float64
	marker, ok := s.tracker.GetUserById(userInfo.UserId)
	if ok {
		x = &marker.Marker.X
		y = &marker.Marker.Y
	}

	if err := s.updateUnitStatus(ctx, userInfo, &dispatch.UnitStatus{
		UnitId:    unit.Id,
		Status:    req.Status,
		Reason:    req.Reason,
		Code:      req.Code,
		UserId:    &userInfo.UserId,
		X:         x,
		Y:         y,
		CreatorId: &userInfo.UserId,
	}); err != nil {
		return nil, ErrFailedQuery
	}

	auditEntry.State = int16(rector.EVENT_TYPE_CREATED)

	return &UpdateUnitStatusResponse{}, nil
}

func (s *Server) AssignUnit(ctx context.Context, req *AssignUnitRequest) (*AssignUnitResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: CentrumService_ServiceDesc.ServiceName,
		Method:  "AssignUnit",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EVENT_TYPE_ERRORED),
	}
	defer s.auditer.Log(auditEntry, req)

	unit, ok := s.getUnit(userInfo.Job, req.UnitId)
	if !ok {
		return nil, ErrFailedQuery
	}
	if unit.Job != userInfo.Job {
		return nil, ErrFailedQuery
	}

	if err := s.updateUnitAssignments(ctx, userInfo, unit, req.ToAdd, req.ToRemove); err != nil {
		return nil, err
	}

	data, err := proto.Marshal(unit)
	if err != nil {
		return nil, err
	}
	s.events.JS.Publish(s.buildSubject(TopicUnit, TypeUnitUserAssigned, userInfo.Job, 0), data)
	s.events.JS.Publish(s.buildSubject(TopicUnit, TypeUnitUserAssigned, userInfo.Job, unit.Id), data)

	auditEntry.State = int16(rector.EVENT_TYPE_UPDATED)

	return &AssignUnitResponse{}, nil
}

func (s *Server) JoinUnit(ctx context.Context, req *JoinUnitRequest) (*JoinUnitResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Only check if the user is joining an unit
	if req.Leave != nil && !*req.Leave {
		_, ok := s.tracker.GetUserByJobAndID(userInfo.Job, userInfo.UserId)
		if !ok {
			return nil, status.Error(codes.InvalidArgument, "You are not on duty!")
		}

		unitId, ok := s.getUnitIDForUserID(userInfo.UserId)
		if !ok {
			return nil, ErrFailedQuery
		}
		if unitId > 0 {
			return nil, ErrAlreadyInUnit
		}
	}

	unit, ok := s.getUnit(userInfo.Job, req.UnitId)
	if !ok {
		return nil, ErrFailedQuery
	}

	resp := &JoinUnitResponse{}
	if req.Leave != nil && !*req.Leave {
		if err := s.updateUnitAssignments(ctx, userInfo, unit, []int32{userInfo.UserId}, nil); err != nil {
			return nil, err
		}

		unit, ok := s.getUnit(userInfo.Job, req.UnitId)
		if !ok {
			return nil, ErrFailedQuery
		}

		resp.Unit = unit
	} else {
		if err := s.updateUnitAssignments(ctx, userInfo, unit, nil, []int32{userInfo.UserId}); err != nil {
			return nil, err
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
		return nil, ErrFailedQuery
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
			tUnitStatus.X,
			tUnitStatus.Y,
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.Dateofbirth,
			tUsers.Job,
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

	resp.Pagination.Update(count.TotalCount, len(resp.Activity))

	return resp, nil
}
