package centrum

import (
	"context"
	"fmt"

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
	defer s.a.Log(auditEntry, req)

	resp := &ListUnitsResponse{
		Units: []*dispatch.Unit{},
	}

	var err error
	resp.Units, err = s.listUnits(ctx, userInfo.Job)
	if err != nil {
		return nil, err
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
	defer s.a.Log(auditEntry, req)

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, ErrFailedQuery
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	resp := &CreateOrUpdateUnitResponse{}
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

		tUnitStatus := table.FivenetCentrumUnitsStatus
		stmt = tUnitStatus.
			INSERT(
				tUnitStatus.UnitID,
				tUnitStatus.Status,
				tUnitStatus.UserID,
			).
			VALUES(
				resp.Unit.Id,
				dispatch.UNIT_STATUS_UNKNOWN,
				userInfo.UserId,
			)

		res, err := stmt.ExecContext(ctx, tx)
		if err != nil {
			return nil, err
		}

		lastId, err = res.LastInsertId()
		if err != nil {
			return nil, err
		}

		status, err := s.getUnitStatus(ctx, uint64(lastId))
		if err != nil {
			return nil, err
		}

		data, err := proto.Marshal(status)
		if err != nil {
			return nil, err
		}
		s.events.JS.Publish(s.buildSubject(TopicUnit, TypeUnitStatus, userInfo, status.UnitId), data)

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

		unit, err := s.getUnitFromDB(ctx, tx, req.Unit.Id)
		if err != nil {
			return nil, ErrFailedQuery
		}

		resp.Unit = unit

		auditEntry.State = int16(rector.EVENT_TYPE_UPDATED)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		auditEntry.State = int16(rector.EVENT_TYPE_ERRORED)
		return nil, ErrFailedQuery
	}

	unit, err := s.getUnitFromDB(ctx, tx, uint64(req.Unit.Id))
	if err != nil {
		return nil, ErrFailedQuery
	}

	resp.Unit = unit

	data, err := proto.Marshal(resp.Unit)
	if err != nil {
		return nil, err
	}
	s.events.JS.Publish(s.buildSubject(TopicUnit, TypeUnitUpdated, userInfo, resp.Unit.Id), data)

	if err := s.loadUnits(ctx, resp.Unit.Id); err != nil {
		return nil, err
	}

	return resp, nil
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
	defer s.a.Log(auditEntry, req)

	resp := &DeleteUnitResponse{}

	unit, ok := s.getUnit(ctx, userInfo, req.UnitId)
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

	if err := s.units.Delete(fmt.Sprintf("%s/%d", userInfo.Job, req.UnitId)); err != nil {
		return nil, err
	}

	data, err := proto.Marshal(unit)
	if err != nil {
		return nil, err
	}
	s.events.JS.Publish(s.buildSubject(TopicUnit, TypeUnitDeleted, userInfo, req.UnitId), data)

	auditEntry.State = int16(rector.EVENT_TYPE_DELETED)

	if err := s.units.Delete(fmt.Sprintf("%s/%d", userInfo.Job, req.UnitId)); err != nil {
		return nil, err
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
	defer s.a.Log(auditEntry, req)

	unit, ok := s.getUnit(ctx, userInfo, req.UnitId)
	if !ok {
		return nil, ErrFailedQuery
	}

	can := s.p.Can(userInfo, CentrumServicePerm, CentrumServiceDeleteUnitPerm)
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

	if _, err := s.updateUnitStatus(ctx, userInfo, unit, &dispatch.UnitStatus{
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
	defer s.a.Log(auditEntry, req)

	unit, ok := s.getUnit(ctx, userInfo, req.UnitId)
	if !ok {
		return nil, ErrFailedQuery
	}
	if unit.Job != userInfo.Job {
		return nil, ErrFailedQuery
	}

	if err := s.updateDispatchUnitAssignments(ctx, userInfo, unit, req.ToAdd, req.ToRemove); err != nil {
		return nil, err
	}

	data, err := proto.Marshal(unit)
	if err != nil {
		return nil, err
	}
	s.events.JS.Publish(s.buildSubject(TopicUnit, TypeUnitUserAssigned, userInfo, unit.Id), data)
	s.events.JS.Publish(s.buildSubject(TopicUnit, TypeUnitUserAssigned, userInfo, 0), data)

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

		unitId, err := s.getUnitIDForUserID(ctx, userInfo.UserId)
		if err != nil {
			return nil, ErrFailedQuery
		}
		if unitId > 0 {
			return nil, ErrAlreadyInUnit
		}
	}

	unit, ok := s.getUnit(ctx, userInfo, req.UnitId)
	if !ok {
		return nil, ErrFailedQuery
	}

	resp := &JoinUnitResponse{}
	if req.Leave != nil && !*req.Leave {
		if err := s.updateDispatchUnitAssignments(ctx, userInfo, unit, []int32{userInfo.UserId}, nil); err != nil {
			return nil, err
		}

		unit, ok := s.getUnit(ctx, userInfo, req.UnitId)
		if !ok {
			return nil, ErrFailedQuery
		}

		resp.Unit = unit
	} else {
		if err := s.updateDispatchUnitAssignments(ctx, userInfo, unit, nil, []int32{userInfo.UserId}); err != nil {
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
