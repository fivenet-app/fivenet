package centrum

import (
	"context"
	"time"

	database "github.com/galexrt/fivenet/gen/go/proto/resources/common/database"
	dispatch "github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	"github.com/galexrt/fivenet/gen/go/proto/resources/timestamp"
	errorscentrum "github.com/galexrt/fivenet/gen/go/proto/services/centrum/errors"
	eventscentrum "github.com/galexrt/fivenet/gen/go/proto/services/centrum/events"
	centrumutils "github.com/galexrt/fivenet/gen/go/proto/services/centrum/utils"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/galexrt/fivenet/pkg/utils/dbutils"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"google.golang.org/protobuf/proto"
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
			tUsers.ID,
			tUsers.Identifier,
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.Sex,
			tUsers.Dateofbirth,
			tUsers.PhoneNumber,
		).
		OPTIMIZER_HINTS(jet.OptimizerHint("idx_users_firstname_lastname_fulltext")).
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
			tDispatch.ID.ASC(),
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

			// Alawys clear dispatch creator's job info
			resp.Dispatches[i].Creator.Job = ""
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

	oldDsp, ok := s.state.GetDispatch(userInfo.Job, req.Dispatch.Id)
	if !ok {
		return nil, errorscentrum.ErrFailedQuery
	}
	if oldDsp.X != req.Dispatch.X || oldDsp.Y != req.Dispatch.Y {
		s.state.DispatchLocations[oldDsp.Job].Remove(oldDsp, nil)
	}

	if err := s.state.UpdateDispatch(ctx, userInfo, req.Dispatch); err != nil {
		return nil, errorscentrum.ErrFailedQuery
	}

	// Load dispatch into cache
	if err := s.state.LoadDispatches(ctx, req.Dispatch.Id); err != nil {
		return nil, errorscentrum.ErrFailedQuery
	}

	dsp, ok := s.state.GetDispatch(userInfo.Job, req.Dispatch.Id)
	if !ok {
		return nil, errorscentrum.ErrFailedQuery
	}

	data, err := proto.Marshal(dsp)
	if err != nil {
		return nil, errorscentrum.ErrFailedQuery
	}
	s.events.JS.PublishAsync(eventscentrum.BuildSubject(eventscentrum.TopicDispatch, eventscentrum.TypeDispatchUpdated, userInfo.Job, 0), data)

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

	unitId, ok := s.state.GetUnitIDForUserID(userInfo.UserId)
	if !ok {
		return nil, errorscentrum.ErrFailedQuery
	}
	unit, ok := s.state.GetUnit(userInfo.Job, unitId)
	if !ok {
		return nil, errorscentrum.ErrFailedQuery
	}

	settings := s.state.GetSettings(userInfo.Job)
	var x, y *float64
	var postal *string
	marker, ok := s.tracker.GetUserById(userInfo.UserId)
	if ok {
		x = &marker.Info.X
		y = &marker.Info.Y
		postal = marker.Info.Postal
	}

	for _, dispatchId := range req.DispatchIds {
		dsp, ok := s.state.GetDispatch(userInfo.Job, dispatchId)
		if !ok {
			return nil, errorscentrum.ErrFailedQuery
		}

		// If the dispatch center is in central command mode, units can't self assign dispatches
		if settings.Mode == dispatch.CentrumMode_CENTRUM_MODE_CENTRAL_COMMAND {
			if !utils.InSliceFunc(dsp.Units, func(in *dispatch.DispatchAssignment) bool {
				return in.UnitId == unitId
			}) {
				return nil, errorscentrum.ErrModeForbidsAction
			}
		}

		// If dispatch is completed, disallow to accept the dispatch
		if dsp.Status != nil && centrumutils.IsStatusDispatchComplete(dsp.Status.Status) {
			return nil, errorscentrum.ErrDispatchAlreadyCompleted
		}

		status := dispatch.StatusDispatch_STATUS_DISPATCH_UNSPECIFIED

		tDispatchUnit := table.FivenetCentrumDispatchesAsgmts
		// Dispatch accepted
		if req.Resp == TakeDispatchResp_TAKE_DISPATCH_RESP_ACCEPTED {
			status = dispatch.StatusDispatch_STATUS_DISPATCH_UNIT_ACCEPTED

			stmt := tDispatchUnit.
				INSERT(
					tDispatchUnit.DispatchID,
					tDispatchUnit.UnitID,
					tDispatchUnit.ExpiresAt,
				).
				VALUES(
					dsp.Id,
					unit.Id,
					jet.NULL,
				).
				ON_DUPLICATE_KEY_UPDATE(
					tDispatchUnit.ExpiresAt.SET(jet.TimestampExp(jet.NULL)),
				)

			if _, err := stmt.ExecContext(ctx, s.db); err != nil {
				if !dbutils.IsDuplicateError(err) {
					return nil, errorscentrum.ErrFailedQuery
				}
			}

			found := false
			accepted := true
			// Set unit expires at to nil
			for _, ua := range dsp.Units {
				if ua.UnitId == unit.Id {
					found = true
					// If there's no expiration time the unit has been directly assigned
					if ua.ExpiresAt == nil {
						accepted = false
					}
					ua.ExpiresAt = nil
					break
				}
			}

			if !found {
				dsp.Units = append(dsp.Units, &dispatch.DispatchAssignment{
					DispatchId: dsp.Id,
					UnitId:     unit.Id,
					Unit:       unit,
					CreatedAt:  timestamp.Now(),
				})
			}

			if accepted {
				// Set unit to busy when unit accepts a dispatch
				if unit.Status == nil || unit.Status.Status != dispatch.StatusUnit_STATUS_UNIT_BUSY {
					if err := s.state.UpdateUnitStatus(ctx, userInfo.Job, unit, &dispatch.UnitStatus{
						UnitId:    unit.Id,
						Status:    dispatch.StatusUnit_STATUS_UNIT_BUSY,
						UserId:    &userInfo.UserId,
						CreatorId: &userInfo.UserId,
						X:         x,
						Y:         y,
						Postal:    postal,
					}); err != nil {
						return nil, errorscentrum.ErrFailedQuery
					}
				}
			}
		} else {
			// Dispatch declined
			status = dispatch.StatusDispatch_STATUS_DISPATCH_UNIT_DECLINED

			stmt := tDispatchUnit.
				DELETE().
				WHERE(jet.AND(
					tDispatchUnit.DispatchID.EQ(jet.Uint64(dsp.Id)),
					tDispatchUnit.UnitID.EQ(jet.Uint64(unit.Id)),
				)).
				LIMIT(1)

			if _, err := stmt.ExecContext(ctx, s.db); err != nil {
				if !dbutils.IsDuplicateError(err) {
					return nil, errorscentrum.ErrFailedQuery
				}
			}

			// Remove the unit's assignment
			for k, u := range dsp.Units {
				if u.UnitId == unit.Id {
					dsp.Units = utils.RemoveFromSlice(dsp.Units, k)
					break
				}
			}
		}

		if err := s.state.UpdateDispatchStatus(ctx, userInfo.Job, dsp, &dispatch.DispatchStatus{
			DispatchId: dispatchId,
			Status:     status,
			UnitId:     &unitId,
			UserId:     &userInfo.UserId,
			X:          x,
			Y:          y,
			Postal:     postal,
		}); err != nil {
			return nil, errorscentrum.ErrFailedQuery
		}
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

	dsp, ok := s.state.GetDispatch(userInfo.Job, req.DispatchId)
	if !ok {
		return nil, errorscentrum.ErrFailedQuery
	}

	if !s.state.CheckIfUserIsPartOfDispatch(userInfo, dsp, true) && !userInfo.SuperUser {
		return nil, errorscentrum.ErrNotPartOfDispatch
	}

	var statusUnitId *uint64
	unitId, ok := s.state.GetUnitIDForUserID(userInfo.UserId)
	if !ok {
		if !s.state.CheckIfUserIsDisponent(userInfo.Job, userInfo.UserId) {
			return nil, errorscentrum.ErrNotPartOfDispatch
		}
	} else {
		statusUnitId = &unitId
	}

	var x, y *float64
	var postal *string
	marker, ok := s.tracker.GetUserById(userInfo.UserId)
	if ok {
		x = &marker.Info.X
		y = &marker.Info.Y
		postal = marker.Info.Postal
	}

	if err := s.state.UpdateDispatchStatus(ctx, userInfo.Job, dsp, &dispatch.DispatchStatus{
		DispatchId: dsp.Id,
		UnitId:     statusUnitId,
		Status:     req.Status,
		Code:       req.Code,
		Reason:     req.Reason,
		UserId:     &userInfo.UserId,
		X:          x,
		Y:          y,
		Postal:     postal,
	}); err != nil {
		return nil, errorscentrum.ErrFailedQuery
	}

	if req.Status == dispatch.StatusDispatch_STATUS_DISPATCH_EN_ROUTE ||
		req.Status == dispatch.StatusDispatch_STATUS_DISPATCH_ON_SCENE ||
		req.Status == dispatch.StatusDispatch_STATUS_DISPATCH_NEED_ASSISTANCE {
		unit, ok := s.state.GetUnit(userInfo.Job, unitId)
		if ok && unit != nil {
			// Set unit to busy when unit accepts a dispatch
			if unit.Status == nil || unit.Status.Status != dispatch.StatusUnit_STATUS_UNIT_BUSY {
				if err := s.state.UpdateUnitStatus(ctx, userInfo.Job, unit, &dispatch.UnitStatus{
					UnitId:    unit.Id,
					Status:    dispatch.StatusUnit_STATUS_UNIT_BUSY,
					UserId:    &userInfo.UserId,
					CreatorId: &userInfo.UserId,
					X:         x,
					Y:         y,
					Postal:    postal,
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

	dsp, ok := s.state.GetDispatch(userInfo.Job, req.DispatchId)
	if !ok {
		return nil, errorscentrum.ErrFailedQuery
	}

	if dsp.Job != userInfo.Job {
		return nil, errorscentrum.ErrFailedQuery
	}

	expiresAt := time.Time{}
	if req.Forced == nil || !*req.Forced {
		expiresAt = s.state.DispatchAssignmentExpirationTime()
	}

	if err := s.state.UpdateDispatchAssignments(ctx, userInfo.Job, &userInfo.UserId, dsp, req.ToAdd, req.ToRemove, expiresAt); err != nil {
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
		FROM(tDispatchStatus).
		WHERE(
			tDispatchStatus.DispatchID.EQ(jet.Uint64(req.Id)),
		)

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
		return nil, errorscentrum.ErrFailedQuery
	}
	for i := 0; i < len(resp.Activity); i++ {
		if resp.Activity[i].UnitId != nil && *resp.Activity[i].UnitId > 0 {
			resp.Activity[i].Unit, _ = s.state.GetUnit(userInfo.Job, *resp.Activity[i].UnitId)
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

	if err := s.state.DeleteDispatch(ctx, userInfo.Job, req.Id); err != nil {
		return nil, errorscentrum.ErrFailedQuery
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)

	return &DeleteDispatchResponse{}, nil
}
