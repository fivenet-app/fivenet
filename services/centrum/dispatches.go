package centrum

import (
	"context"
	"errors"
	"slices"
	"time"

	centrum "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	database "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/rector"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	pbcentrum "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/centrum"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/model"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorscentrum "github.com/fivenet-app/fivenet/v2025/services/centrum/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var (
	tDispatch       = table.FivenetCentrumDispatches.AS("dispatch")
	tDispatchStatus = table.FivenetCentrumDispatchesStatus.AS("dispatchstatus")
)

func (s *Server) ListDispatches(ctx context.Context, req *pbcentrum.ListDispatchesRequest) (*pbcentrum.ListDispatchesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbcentrum.CentrumService_ServiceDesc.ServiceName,
		Method:  "ListDispatches",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	condition := jet.AND(tDispatch.Job.EQ(jet.String(userInfo.Job)).
		AND(
			tDispatchStatus.ID.IS_NULL().OR(
				tDispatchStatus.ID.EQ(
					jet.RawInt("SELECT MAX(`dispatchstatus`.`id`) FROM `fivenet_centrum_dispatches_status` AS `dispatchstatus` WHERE `dispatchstatus`.`dispatch_id` = `dispatch`.`id`"),
				),
			),
		),
	)

	if len(req.Status) > 0 {
		statuses := make([]jet.Expression, len(req.Status))
		for i := range req.Status {
			statuses[i] = jet.Int16(int16(*req.Status[i].Enum()))
		}

		condition = condition.AND(tDispatchStatus.Status.IN(statuses...))
	}
	if len(req.NotStatus) > 0 {
		statuses := make([]jet.Expression, len(req.NotStatus))
		for i := range req.NotStatus {
			statuses[i] = jet.Int16(int16(*req.NotStatus[i].Enum()))
		}

		condition = condition.AND(tDispatchStatus.Status.NOT_IN(statuses...))
	}

	if len(req.Ids) > 0 {
		ids := make([]jet.Expression, len(req.Ids))
		for i := range req.Ids {
			ids[i] = jet.Uint64(req.Ids[i])
		}

		condition = condition.AND(tDispatch.ID.IN(ids...))
	}

	if req.Postal != nil && *req.Postal != "" {
		condition = condition.AND(tDispatch.Postal.EQ(jet.String(*req.Postal)))
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
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(count.TotalCount, 16)
	resp := &pbcentrum.ListDispatchesResponse{
		Pagination: pag,
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

	tUsers := tables.Users().AS("colleague")

	stmt := tDispatch.
		SELECT(
			tDispatch.ID,
			tDispatch.CreatedAt,
			tDispatch.UpdatedAt,
			tDispatch.Job,
			tDispatch.Message,
			tDispatch.Description,
			tDispatch.Attributes,
			tDispatch.References,
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
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}
	}

	resp.Pagination.Update(len(resp.Dispatches))

	publicJobs := s.appCfg.Get().JobInfo.PublicJobs
	for i := range resp.Dispatches {
		var err error
		resp.Dispatches[i].Units, err = s.state.LoadDispatchAssignments(ctx, resp.Dispatches[i].Job, resp.Dispatches[i].Id)
		if err != nil {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}

		if resp.Dispatches[i].CreatorId != nil {
			resp.Dispatches[i].Creator, err = s.state.RetrieveUserById(ctx, *resp.Dispatches[i].CreatorId)
			if err != nil {
				return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
			}

			if resp.Dispatches[i].Creator != nil {
				// Clear dispatch creator's job info if it isn't a public job
				if !slices.Contains(publicJobs, resp.Dispatches[i].Creator.Job) {
					resp.Dispatches[i].Creator.Job = ""
				}
				resp.Dispatches[i].Creator.JobGrade = 0
			}
		}
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_VIEWED)

	return resp, nil
}

func (s *Server) GetDispatch(ctx context.Context, req *pbcentrum.GetDispatchRequest) (*pbcentrum.GetDispatchResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.centrum.dispatch_id", int64(req.Id)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbcentrum.CentrumService_ServiceDesc.ServiceName,
		Method:  "GetDispatch",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	condition := jet.AND(tDispatch.Job.EQ(jet.String(userInfo.Job)).
		AND(
			tDispatchStatus.ID.IS_NULL().OR(
				tDispatchStatus.ID.EQ(
					jet.RawInt("SELECT MAX(`dispatchstatus`.`id`) FROM `fivenet_centrum_dispatches_status` AS `dispatchstatus` WHERE `dispatchstatus`.`dispatch_id` = `dispatch`.`id`"),
				),
			))).
		AND(tDispatch.ID.EQ(jet.Uint64(req.Id)))

	resp := &pbcentrum.GetDispatchResponse{
		Dispatch: &centrum.Dispatch{},
	}

	tUsers := tables.Users().AS("colleague")

	stmt := tDispatch.
		SELECT(
			tDispatch.ID,
			tDispatch.CreatedAt,
			tDispatch.UpdatedAt,
			tDispatch.Job,
			tDispatch.Message,
			tDispatch.Description,
			tDispatch.Attributes,
			tDispatch.References,
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
		ORDER_BY(
			tDispatch.ID.DESC(),
		).
		LIMIT(1)

	if err := stmt.QueryContext(ctx, s.db, resp.Dispatch); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}
	}

	if resp.Dispatch == nil || resp.Dispatch.Id <= 0 {
		return &pbcentrum.GetDispatchResponse{}, nil
	}

	var err error
	resp.Dispatch.Units, err = s.state.LoadDispatchAssignments(ctx, resp.Dispatch.Job, resp.Dispatch.Id)
	if err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	if resp.Dispatch.CreatorId != nil {
		creator, err := s.state.RetrieveUserById(ctx, *resp.Dispatch.CreatorId)
		if err != nil {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}

		if creator != nil {
			resp.Dispatch.Creator = creator
			// Clear dispatch creator's job info if not a visible job
			if !slices.Contains(s.appCfg.Get().JobInfo.PublicJobs, creator.Job) {
				resp.Dispatch.Creator.Job = ""
			}
			resp.Dispatch.Creator.JobGrade = 0
		}
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_VIEWED)

	return resp, nil
}

func (s *Server) CreateDispatch(ctx context.Context, req *pbcentrum.CreateDispatchRequest) (*pbcentrum.CreateDispatchResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbcentrum.CentrumService_ServiceDesc.ServiceName,
		Method:  "CreateDispatch",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	req.Dispatch.Job = userInfo.Job
	req.Dispatch.CreatorId = &userInfo.UserId

	req.Dispatch.CreatedAt = timestamp.Now()
	dsp, err := s.state.CreateDispatch(ctx, req.Dispatch)
	if err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)

	return &pbcentrum.CreateDispatchResponse{
		Dispatch: dsp,
	}, nil
}

func (s *Server) UpdateDispatch(ctx context.Context, req *pbcentrum.UpdateDispatchRequest) (*pbcentrum.UpdateDispatchResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.centrum.dispatch_id", int64(req.Dispatch.Id)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbcentrum.CentrumService_ServiceDesc.ServiceName,
		Method:  "UpdateDispatch",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	if _, err := s.state.UpdateDispatch(ctx, userInfo.Job, &userInfo.UserId, req.Dispatch, true); err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return &pbcentrum.UpdateDispatchResponse{}, nil
}

func (s *Server) TakeDispatch(ctx context.Context, req *pbcentrum.TakeDispatchRequest) (*pbcentrum.TakeDispatchResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64Slice("fivenet.centrum.dispatch_ids", utils.SliceUint64ToInt64(req.DispatchIds)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbcentrum.CentrumService_ServiceDesc.ServiceName,
		Method:  "TakeDispatch",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	unitId, ok := s.state.GetUserUnitID(ctx, userInfo.UserId)
	if !ok {
		return nil, errorscentrum.ErrFailedQuery
	}

	if err := s.state.TakeDispatch(ctx, userInfo.Job, userInfo.UserId, unitId, req.Resp, req.DispatchIds); err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return &pbcentrum.TakeDispatchResponse{}, nil
}

func (s *Server) UpdateDispatchStatus(ctx context.Context, req *pbcentrum.UpdateDispatchStatusRequest) (*pbcentrum.UpdateDispatchStatusResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbcentrum.CentrumService_ServiceDesc.ServiceName,
		Method:  "UpdateDispatchStatus",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	dsp, err := s.state.GetDispatch(ctx, userInfo.Job, req.DispatchId)
	if err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	if !s.state.CheckIfUserIsPartOfDispatch(ctx, userInfo, dsp, true) && !userInfo.SuperUser {
		return nil, errorscentrum.ErrNotPartOfDispatch
	}

	var statusUnitId *uint64
	unitId, ok := s.state.GetUserUnitID(ctx, userInfo.UserId)
	if !ok {
		if !s.state.CheckIfUserIsDisponent(ctx, userInfo.Job, userInfo.UserId) {
			return nil, errorscentrum.ErrNotPartOfDispatch
		}
	} else {
		statusUnitId = &unitId
	}

	if _, err := s.state.UpdateDispatchStatus(ctx, userInfo.Job, dsp.Id, &centrum.DispatchStatus{
		CreatedAt:  timestamp.Now(),
		DispatchId: dsp.Id,
		UnitId:     statusUnitId,
		Status:     req.Status,
		Code:       req.Code,
		Reason:     req.Reason,
		UserId:     &userInfo.UserId,
	}); err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	if req.Status == centrum.StatusDispatch_STATUS_DISPATCH_EN_ROUTE ||
		req.Status == centrum.StatusDispatch_STATUS_DISPATCH_ON_SCENE ||
		req.Status == centrum.StatusDispatch_STATUS_DISPATCH_NEED_ASSISTANCE {
		if unit, err := s.state.GetUnit(ctx, userInfo.Job, unitId); err == nil {
			// Set unit to busy when unit accepts a dispatch
			if unit.Status == nil || unit.Status.Status != centrum.StatusUnit_STATUS_UNIT_BUSY {
				if _, err := s.state.UpdateUnitStatus(ctx, userInfo.Job, unitId, &centrum.UnitStatus{
					CreatedAt: timestamp.Now(),
					UnitId:    unit.Id,
					Status:    centrum.StatusUnit_STATUS_UNIT_BUSY,
					UserId:    &userInfo.UserId,
					CreatorId: &userInfo.UserId,
				}); err != nil {
					return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
				}
			}
		}
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return &pbcentrum.UpdateDispatchStatusResponse{}, nil
}

func (s *Server) AssignDispatch(ctx context.Context, req *pbcentrum.AssignDispatchRequest) (*pbcentrum.AssignDispatchResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.centrum.dispatch_id", int64(req.DispatchId)))
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64Slice("fivenet.centrum.units.to_add", utils.SliceUint64ToInt64(req.ToAdd)))
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64Slice("fivenet.centrum.units.to_remove", utils.SliceUint64ToInt64(req.ToRemove)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbcentrum.CentrumService_ServiceDesc.ServiceName,
		Method:  "AssignDispatch",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	dsp, err := s.state.GetDispatch(ctx, userInfo.Job, req.DispatchId)
	if err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	if dsp.Job != userInfo.Job {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	expiresAt := time.Time{}
	if req.Forced == nil || !*req.Forced {
		expiresAt = s.state.DispatchAssignmentExpirationTime()
	}

	if err := s.state.UpdateDispatchAssignments(ctx, userInfo.Job, &userInfo.UserId, dsp.Id, req.ToAdd, req.ToRemove, expiresAt); err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return &pbcentrum.AssignDispatchResponse{}, nil
}

func (s *Server) ListDispatchActivity(ctx context.Context, req *pbcentrum.ListDispatchActivityRequest) (*pbcentrum.ListDispatchActivityResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	countStmt := tDispatchStatus.
		SELECT(
			jet.COUNT(jet.DISTINCT(tDispatchStatus.ID)).AS("datacount.totalcount"),
		).
		FROM(
			tDispatchStatus.
				INNER_JOIN(tDispatch,
					tDispatch.ID.EQ(tDispatchStatus.DispatchID),
				),
		).
		WHERE(jet.AND(
			tDispatchStatus.DispatchID.EQ(jet.Uint64(req.Id)),
			tDispatch.Job.EQ(jet.String(userInfo.Job)),
		))

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(count.TotalCount, 10)
	resp := &pbcentrum.ListDispatchActivityResponse{
		Pagination: pag,
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

	tUsers := tables.Users().AS("colleague")

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
			tUsers.ID,
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Sex,
			tUsers.Dateofbirth,
			tUsers.PhoneNumber,
			tJobsUserProps.UserID,
			tJobsUserProps.Job,
			tJobsUserProps.NamePrefix,
			tJobsUserProps.NameSuffix,
			tUserProps.Avatar.AS("colleague.avatar"),
		).
		FROM(
			tDispatchStatus.
				LEFT_JOIN(tUsers,
					tUsers.ID.EQ(tDispatchStatus.UserID),
				).
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tDispatchStatus.UserID).
						AND(tUsers.Job.EQ(jet.String(userInfo.Job))),
				).
				LEFT_JOIN(tJobsUserProps,
					tJobsUserProps.UserID.EQ(tUsers.ID).
						AND(tJobsUserProps.Job.EQ(tUsers.Job)),
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
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.Activity {
		if resp.Activity[i].UnitId != nil && *resp.Activity[i].UnitId > 0 {
			var err error
			resp.Activity[i].Unit, err = s.state.GetUnit(ctx, userInfo.Job, *resp.Activity[i].UnitId)
			if err != nil {
				return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
			}
		}

		if resp.Activity[i].User != nil {
			jobInfoFn(resp.Activity[i].User)
		}
	}

	resp.Pagination.Update(len(resp.Activity))

	return resp, nil
}

func (s *Server) DeleteDispatch(ctx context.Context, req *pbcentrum.DeleteDispatchRequest) (*pbcentrum.DeleteDispatchResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbcentrum.CentrumService_ServiceDesc.ServiceName,
		Method:  "DeleteDispatch",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	if err := s.state.DeleteDispatch(ctx, userInfo.Job, req.Id, true); err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)

	return &pbcentrum.DeleteDispatchResponse{}, nil
}
