package centrum

import (
	"context"
	"errors"
	"time"

	centrum "github.com/galexrt/fivenet/gen/go/proto/resources/centrum"
	database "github.com/galexrt/fivenet/gen/go/proto/resources/common/database"
	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	errorscentrum "github.com/galexrt/fivenet/gen/go/proto/services/centrum/errors"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/paulmach/orb"
)

var (
	tDispatch       = table.FivenetCentrumDispatches.AS("dispatch")
	tDispatchStatus = table.FivenetCentrumDispatchesStatus.AS("dispatchstatus")
)

func (s *Server) ListDispatches(ctx context.Context, req *ListDispatchesRequest) (*ListDispatchesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: CentrumService_ServiceDesc.ServiceName,
		Method:  "ListDispatches",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.auditer.Log(auditEntry, req)

	condition := tDispatch.Job.EQ(jet.String(userInfo.Job)).
		AND(
			tDispatchStatus.ID.IS_NULL().OR(
				tDispatchStatus.ID.EQ(
					jet.RawInt("SELECT MAX(`dispatchstatus`.`id`) FROM `fivenet_centrum_dispatches_status` AS `dispatchstatus` WHERE `dispatchstatus`.`dispatch_id` = `dispatch`.`id`"),
				),
			))

	if len(req.Status) > 0 {
		statuses := make([]jet.Expression, len(req.Status))
		for i := 0; i < len(req.Status); i++ {
			statuses[i] = jet.Int16(int16(*req.Status[i].Enum()))
		}

		condition = condition.AND(tDispatchStatus.Status.IN(statuses...))
	}
	if len(req.NotStatus) > 0 {
		statuses := make([]jet.Expression, len(req.NotStatus))
		for i := 0; i < len(req.NotStatus); i++ {
			statuses[i] = jet.Int16(int16(*req.NotStatus[i].Enum()))
		}

		condition = condition.AND(tDispatchStatus.Status.NOT_IN(statuses...))
	}

	countStmt := tDispatch.
		SELECT(
			jet.COUNT(tDispatch.ID).AS("datacount.totalcount"),
		).
		FROM(
			tDispatch.
				LEFT_JOIN(tDispatchStatus,
					tDispatchStatus.DispatchID.EQ(tDispatch.ID),
				),
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		return nil, errorscentrum.ErrFailedQuery
	}

	pag, limit := req.Pagination.GetResponse()
	resp := &ListDispatchesResponse{
		Pagination: pag,
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

	stmt := tDispatch.
		SELECT(
			tDispatch.ID,
			tDispatch.CreatedAt,
			tDispatch.UpdatedAt,
			tDispatch.Job,
			tDispatch.Message,
			tDispatch.Description,
			tDispatch.Attributes,
			tDispatch.X,
			tDispatch.Y,
			tDispatch.Postal,
			tDispatch.Anon,
			tDispatch.CreatorID,
			tDispatchStatus.ID,
			tDispatchStatus.CreatedAt,
			tDispatchStatus.DispatchID,
			tDispatchStatus.UnitID,
			tDispatchStatus.Status,
			tDispatchStatus.Reason,
			tDispatchStatus.Code,
			tDispatchStatus.UserID,
			tDispatchStatus.X,
			tDispatchStatus.Y,
			tDispatchStatus.Postal,
			tUsers.Identifier,
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Sex,
			tUsers.Dateofbirth,
			tUsers.PhoneNumber,
		).
		FROM(tDispatch.
			LEFT_JOIN(tDispatchStatus,
				tDispatchStatus.DispatchID.EQ(tDispatch.ID),
			).
			LEFT_JOIN(tUsers,
				tUsers.ID.EQ(tDispatchStatus.UserID),
			)).
		WHERE(condition).
		OFFSET(req.Pagination.Offset).
		ORDER_BY(
			tDispatch.ID.DESC(),
		).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Dispatches); err != nil {
		return nil, errorscentrum.ErrFailedQuery
	}

	resp.Pagination.Update(count.TotalCount, len(resp.Dispatches))

	for i := 0; i < len(resp.Dispatches); i++ {
		var err error
		resp.Dispatches[i].Units, err = s.state.LoadDispatchAssignments(ctx, resp.Dispatches[i].Job, resp.Dispatches[i].Id)
		if err != nil {
			return nil, errorscentrum.ErrFailedQuery
		}

		if resp.Dispatches[i].CreatorId != nil {
			resp.Dispatches[i].Creator, err = s.state.ResolveUserById(ctx, *resp.Dispatches[i].CreatorId)
			if err != nil {
				return nil, err
			}

			// Clear dispatch creator's job info if not a visible job
			if !utils.InSlice(s.publicJobs, resp.Dispatches[i].Creator.Job) {
				resp.Dispatches[i].Creator.Job = ""
			}
			resp.Dispatches[i].Creator.JobGrade = 0
		}
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_VIEWED)

	return resp, nil
}

func (s *Server) CreateDispatch(ctx context.Context, req *CreateDispatchRequest) (*CreateDispatchResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: CentrumService_ServiceDesc.ServiceName,
		Method:  "CreateDispatch",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.auditer.Log(auditEntry, req)

	req.Dispatch.Job = userInfo.Job
	req.Dispatch.CreatorId = &userInfo.UserId

	dsp, err := s.state.CreateDispatch(ctx, req.Dispatch)
	if err != nil {
		return nil, errorscentrum.ErrFailedQuery
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)

	return &CreateDispatchResponse{
		Dispatch: dsp,
	}, nil
}

func (s *Server) UpdateDispatch(ctx context.Context, req *UpdateDispatchRequest) (*UpdateDispatchResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: CentrumService_ServiceDesc.ServiceName,
		Method:  "UpdateDispatch",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.auditer.Log(auditEntry, req)

	oldDsp, err := s.state.GetDispatch(userInfo.Job, req.Dispatch.Id)
	if oldDsp == nil || err != nil {
		return nil, errorscentrum.ErrFailedQuery
	}
	if oldDsp.X != req.Dispatch.X || oldDsp.Y != req.Dispatch.Y {
		s.state.GetDispatchLocations(oldDsp.Job).Remove(oldDsp, func(p orb.Pointer) bool {
			return p.(*centrum.Dispatch).Id == oldDsp.Id
		})
	}

	if _, err := s.state.UpdateDispatch(ctx, userInfo.Job, &userInfo.UserId, req.Dispatch, true); err != nil {
		return nil, errorscentrum.ErrFailedQuery
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return &UpdateDispatchResponse{}, nil
}

func (s *Server) TakeDispatch(ctx context.Context, req *TakeDispatchRequest) (*TakeDispatchResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: CentrumService_ServiceDesc.ServiceName,
		Method:  "TakeDispatch",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.auditer.Log(auditEntry, req)

	unitId, ok := s.state.GetUserUnitID(userInfo.UserId)
	if !ok {
		return nil, errorscentrum.ErrFailedQuery
	}

	if err := s.state.TakeDispatch(ctx, userInfo.Job, userInfo.UserId, unitId, req.Resp, req.DispatchIds); err != nil {
		return nil, err
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return &TakeDispatchResponse{}, nil
}

func (s *Server) UpdateDispatchStatus(ctx context.Context, req *UpdateDispatchStatusRequest) (*UpdateDispatchStatusResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: CentrumService_ServiceDesc.ServiceName,
		Method:  "UpdateDispatchStatus",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.auditer.Log(auditEntry, req)

	dsp, err := s.state.GetDispatch(userInfo.Job, req.DispatchId)
	if err != nil {
		return nil, errorscentrum.ErrFailedQuery
	}

	if !s.state.CheckIfUserIsPartOfDispatch(userInfo, dsp, true) && !userInfo.SuperUser {
		return nil, errorscentrum.ErrNotPartOfDispatch
	}

	var statusUnitId *uint64
	unitId, ok := s.state.GetUserUnitID(userInfo.UserId)
	if !ok {
		if !s.state.CheckIfUserIsDisponent(userInfo.Job, userInfo.UserId) {
			return nil, errorscentrum.ErrNotPartOfDispatch
		}
	} else {
		statusUnitId = &unitId
	}

	if _, err := s.state.UpdateDispatchStatus(ctx, userInfo.Job, dsp.Id, &centrum.DispatchStatus{
		DispatchId: dsp.Id,
		UnitId:     statusUnitId,
		Status:     req.Status,
		Code:       req.Code,
		Reason:     req.Reason,
		UserId:     &userInfo.UserId,
	}); err != nil {
		return nil, errorscentrum.ErrFailedQuery
	}

	if req.Status == centrum.StatusDispatch_STATUS_DISPATCH_EN_ROUTE ||
		req.Status == centrum.StatusDispatch_STATUS_DISPATCH_ON_SCENE ||
		req.Status == centrum.StatusDispatch_STATUS_DISPATCH_NEED_ASSISTANCE {
		if unit, err := s.state.GetUnit(userInfo.Job, unitId); err == nil {
			// Set unit to busy when unit accepts a dispatch
			if unit.Status == nil || unit.Status.Status != centrum.StatusUnit_STATUS_UNIT_BUSY {
				if _, err := s.state.UpdateUnitStatus(ctx, userInfo.Job, unitId, &centrum.UnitStatus{
					UnitId:    unit.Id,
					Status:    centrum.StatusUnit_STATUS_UNIT_BUSY,
					UserId:    &userInfo.UserId,
					CreatorId: &userInfo.UserId,
				}); err != nil {
					return nil, errorscentrum.ErrFailedQuery
				}
			}
		}
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return &UpdateDispatchStatusResponse{}, nil
}

func (s *Server) AssignDispatch(ctx context.Context, req *AssignDispatchRequest) (*AssignDispatchResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: CentrumService_ServiceDesc.ServiceName,
		Method:  "AssignDispatch",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.auditer.Log(auditEntry, req)

	dsp, err := s.state.GetDispatch(userInfo.Job, req.DispatchId)
	if err != nil {
		return nil, errorscentrum.ErrFailedQuery
	}

	if dsp.Job != userInfo.Job {
		return nil, errorscentrum.ErrFailedQuery
	}

	expiresAt := time.Time{}
	if req.Forced == nil || !*req.Forced {
		expiresAt = s.state.DispatchAssignmentExpirationTime()
	}

	if err := s.state.UpdateDispatchAssignments(ctx, userInfo.Job, &userInfo.UserId, dsp.Id, req.ToAdd, req.ToRemove, expiresAt); err != nil {
		return nil, errorscentrum.ErrFailedQuery
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return &AssignDispatchResponse{}, nil
}

func (s *Server) ListDispatchActivity(ctx context.Context, req *ListDispatchActivityRequest) (*ListDispatchActivityResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	countStmt := tDispatchStatus.
		SELECT(
			jet.COUNT(jet.DISTINCT(tDispatchStatus.ID)).AS("datacount.totalcount"),
		).
		FROM(
			tDispatchStatus.
				INNER_JOIN(tDispatch,
					tDispatch.ID.EQ(tDispatchStatus.UnitID),
				),
		).
		WHERE(jet.AND(
			tDispatchStatus.DispatchID.EQ(jet.Uint64(req.Id)),
			tDispatch.Job.EQ(jet.String(userInfo.Job)),
		))

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		return nil, errorscentrum.ErrFailedQuery
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(10)
	resp := &ListDispatchActivityResponse{
		Pagination: pag,
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

	stmt := tDispatchStatus.
		SELECT(
			tDispatchStatus.ID,
			tDispatchStatus.CreatedAt,
			tDispatchStatus.UnitID,
			tDispatchStatus.Status,
			tDispatchStatus.Reason,
			tDispatchStatus.Code,
			tDispatchStatus.UserID,
			tDispatchStatus.X,
			tDispatchStatus.Y,
			tDispatchStatus.Postal,
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Dateofbirth,
			tUsers.PhoneNumber,
		).
		FROM(
			tDispatchStatus.
				LEFT_JOIN(tUsers,
					tUsers.ID.EQ(tDispatchStatus.UserID),
				),
		).
		WHERE(
			tDispatchStatus.DispatchID.EQ(jet.Uint64(req.Id)),
		).
		ORDER_BY(tDispatchStatus.ID.DESC()).
		OFFSET(req.Pagination.Offset).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Activity); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errorscentrum.ErrFailedQuery
		}
	}
	for i := 0; i < len(resp.Activity); i++ {
		if resp.Activity[i].UnitId != nil && *resp.Activity[i].UnitId > 0 {
			var err error
			resp.Activity[i].Unit, err = s.state.GetUnit(userInfo.Job, *resp.Activity[i].UnitId)
			if err != nil {
				return nil, errorscentrum.ErrFailedQuery
			}
		}
	}

	resp.Pagination.Update(count.TotalCount, len(resp.Activity))

	return resp, nil
}

func (s *Server) DeleteDispatch(ctx context.Context, req *DeleteDispatchRequest) (*DeleteDispatchResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: CentrumService_ServiceDesc.ServiceName,
		Method:  "DeleteDispatch",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.auditer.Log(auditEntry, req)

	if err := s.state.DeleteDispatch(ctx, userInfo.Job, req.Id, true); err != nil {
		return nil, errorscentrum.ErrFailedQuery
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)

	return &DeleteDispatchResponse{}, nil
}
