package centrum

import (
	"context"

	database "github.com/galexrt/fivenet/gen/go/proto/resources/common/database"
	"github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	errorscentrum "github.com/galexrt/fivenet/gen/go/proto/services/centrum/errors"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"google.golang.org/protobuf/proto"
)

var (
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

	resp.Units = s.state.FilterUnits(userInfo.Job, req.Status, nil)
	if resp.Units == nil {
		return nil, errorscentrum.ErrModeForbidsAction
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

	var unit *dispatch.Unit
	var err error
	// No unit id set
	if req.Unit.Id <= 0 {
		unit, err = s.state.CreateUnit(ctx, userInfo.Job, req.Unit)
		if err != nil {
			return nil, err
		}

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)
	} else {
		unit, err = s.state.UpdateUnit(ctx, userInfo.Job, req.Unit)
		if err != nil {
			return nil, err
		}

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)
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

	if err := s.state.DeleteUnit(ctx, userInfo.Job, req.UnitId); err != nil {
		return nil, errorscentrum.ErrFailedQuery
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)

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
	if marker, ok := s.tracker.GetUserById(userInfo.UserId); ok {
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

	if err := s.state.UpdateUnitAssignments(ctx, userInfo.Job, &userInfo.UserId, unit, req.ToAdd, req.ToRemove); err != nil {
		return nil, errorscentrum.ErrFailedQuery
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

	currentUnitId, _ := s.state.GetUnitIDForUserID(userInfo.UserId)

	resp := &JoinUnitResponse{}
	// User tries to join his own unit
	if req.UnitId != nil && *req.UnitId == currentUnitId {
		return resp, nil
	}

	currentUnit, _ := s.state.GetUnit(userInfo.Job, currentUnitId)

	// User joins unit
	if req.UnitId != nil && *req.UnitId > 0 {
		// Remove user from his current unit
		if currentUnit != nil {
			if err := s.state.UpdateUnitAssignments(ctx, userInfo.Job, &userInfo.UserId, currentUnit, nil, []int32{userInfo.UserId}); err != nil {
				return nil, errorscentrum.ErrFailedQuery
			}
		}

		newUnit, ok := s.state.GetUnit(userInfo.Job, *req.UnitId)
		if !ok {
			return nil, errorscentrum.ErrFailedQuery
		}

		if err := s.state.UpdateUnitAssignments(ctx, userInfo.Job, &userInfo.UserId, newUnit, []int32{userInfo.UserId}, nil); err != nil {
			return nil, errorscentrum.ErrFailedQuery
		}

		resp.Unit = newUnit
	} else {
		// User leaves his current unit (if he is in an unit)
		if currentUnit != nil {
			if err := s.state.UpdateUnitAssignments(ctx, userInfo.Job, &userInfo.UserId, currentUnit, nil, []int32{userInfo.UserId}); err != nil {
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
			tUsers.Identifier,
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Sex,
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
