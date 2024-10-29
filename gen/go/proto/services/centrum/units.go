package centrum

import (
	"context"
	"errors"
	"time"

	centrum "github.com/fivenet-app/fivenet/gen/go/proto/resources/centrum"
	database "github.com/fivenet-app/fivenet/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/rector"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/timestamp"
	errorscentrum "github.com/fivenet-app/fivenet/gen/go/proto/services/centrum/errors"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/pkg/utils"
	"github.com/fivenet-app/fivenet/query/fivenet/model"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var (
	tUnitStatus = table.FivenetCentrumUnitsStatus.AS("unitstatus")
	tUsers      = table.Users.AS("usershort")
	tUserProps  = table.FivenetUserProps
	tUnits      = table.FivenetCentrumUnits.AS("unit")
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
	defer s.aud.Log(auditEntry, req)

	resp := &ListUnitsResponse{
		Units: []*centrum.Unit{},
	}

	resp.Units = s.state.FilterUnits(ctx, userInfo.Job, req.Status, nil, nil)
	if resp.Units == nil {
		return nil, errorscentrum.ErrFailedQuery
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
	defer s.aud.Log(auditEntry, req)

	req.Unit.Job = userInfo.Job

	var unit *centrum.Unit
	var err error
	// No unit id set
	if req.Unit.Id <= 0 {
		unit, err = s.state.CreateUnit(ctx, userInfo.Job, req.Unit)
		if err != nil {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)
	} else {
		unit, err = s.state.UpdateUnit(ctx, userInfo.Job, req.Unit)
		if err != nil {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)
	}

	return &CreateOrUpdateUnitResponse{
		Unit: unit,
	}, nil
}

func (s *Server) DeleteUnit(ctx context.Context, req *DeleteUnitRequest) (*DeleteUnitResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.centrum.unit_id", int64(req.UnitId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: CentrumService_ServiceDesc.ServiceName,
		Method:  "DeleteUnit",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	resp := &DeleteUnitResponse{}

	if err := s.state.DeleteUnit(ctx, userInfo.Job, req.UnitId); err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)

	return resp, nil
}

func (s *Server) UpdateUnitStatus(ctx context.Context, req *UpdateUnitStatusRequest) (*UpdateUnitStatusResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.centrum.unit_id", int64(req.UnitId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: CentrumService_ServiceDesc.ServiceName,
		Method:  "UpdateUnitStatus",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	unit, err := s.state.GetUnit(ctx, userInfo.Job, req.UnitId)
	if err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	if !s.state.CheckIfUserPartOfUnit(ctx, userInfo.Job, userInfo.UserId, unit, true) {
		return nil, errorscentrum.ErrNotPartOfUnit
	}

	if _, err := s.state.UpdateUnitStatus(ctx, userInfo.Job, unit.Id, &centrum.UnitStatus{
		CreatedAt: timestamp.Now(),
		UnitId:    unit.Id,
		Unit:      unit,
		Status:    req.Status,
		Reason:    req.Reason,
		Code:      req.Code,
		UserId:    &userInfo.UserId,
		CreatorId: &userInfo.UserId,
	}); err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)

	return &UpdateUnitStatusResponse{}, nil
}

func (s *Server) AssignUnit(ctx context.Context, req *AssignUnitRequest) (*AssignUnitResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.centrum.unit_id", int64(req.UnitId)))
	trace.SpanFromContext(ctx).SetAttributes(attribute.IntSlice("fivenet.centrum.users.to_add", utils.SliceInt32ToInt(req.ToAdd)))
	trace.SpanFromContext(ctx).SetAttributes(attribute.IntSlice("fivenet.centrum.users.to_remove", utils.SliceInt32ToInt(req.ToRemove)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: CentrumService_ServiceDesc.ServiceName,
		Method:  "AssignUnit",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	unit, err := s.state.GetUnit(ctx, userInfo.Job, req.UnitId)
	if err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}
	if unit.Job != userInfo.Job {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	if err := s.state.UpdateUnitAssignments(ctx, userInfo.Job, &userInfo.UserId, unit.Id, req.ToAdd, req.ToRemove); err != nil {
		return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	return &AssignUnitResponse{}, nil
}

func (s *Server) JoinUnit(ctx context.Context, req *JoinUnitRequest) (*JoinUnitResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if req.UnitId != nil {
		trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.centrum.unit_id", int64(*req.UnitId)))
	}

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Check if user is on duty, if not make sure to unset any unit id
	if um, ok := s.tracker.GetUserById(userInfo.UserId); !ok || um.Hidden {
		if err := s.state.UnsetUnitIDForUser(ctx, userInfo.UserId); err != nil {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}

		return nil, errorscentrum.ErrNotOnDuty
	}

	currentUnitId, _ := s.state.GetUserUnitID(ctx, userInfo.UserId)

	resp := &JoinUnitResponse{}
	// User tries to join his own unit
	if req.UnitId != nil && *req.UnitId == currentUnitId {
		return resp, nil
	}

	currentUnit, err := s.state.GetUnit(ctx, userInfo.Job, currentUnitId)
	if err != nil && !errors.Is(err, jetstream.ErrKeyNotFound) {
		return nil, errorscentrum.ErrNotOnDuty
	}

	// User joins unit
	if req.UnitId != nil && *req.UnitId > 0 {
		s.logger.Debug("user joining unit", zap.String("job", userInfo.Job), zap.Int32("user_id", userInfo.UserId), zap.Uint64("current_unit_id", currentUnitId), zap.Uint64p("unit_id", req.UnitId))
		// Remove user from his current unit
		if currentUnit != nil {
			if err := s.state.UpdateUnitAssignments(ctx, userInfo.Job, &userInfo.UserId, currentUnit.Id, nil, []int32{userInfo.UserId}); err != nil {
				return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
			}
		}

		newUnit, err := s.state.GetUnit(ctx, userInfo.Job, *req.UnitId)
		if err != nil {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}

		if err := s.state.UpdateUnitAssignments(ctx, userInfo.Job, &userInfo.UserId, newUnit.Id, []int32{userInfo.UserId}, nil); err != nil {
			return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}

		resp.Unit = newUnit
	} else {
		s.logger.Debug("user leaving unit", zap.Uint64("current_unit_id", currentUnitId), zap.Uint64p("unit_id", req.UnitId))
		// User leaves his current unit (if he is in an unit)
		if currentUnit != nil {
			if err := s.state.UpdateUnitAssignments(ctx, userInfo.Job, &userInfo.UserId, currentUnit.Id, nil, []int32{userInfo.UserId}); err != nil {
				return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
			}
		}
	}

	return resp, nil
}

func (s *Server) ListUnitActivity(ctx context.Context, req *ListUnitActivityRequest) (*ListUnitActivityResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.centrum.unit_id", int64(req.Id)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	countStmt := tUnitStatus.
		SELECT(
			jet.COUNT(jet.DISTINCT(tUnitStatus.ID)).AS("datacount.totalcount"),
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

	pag, limit := req.Pagination.GetResponseWithPageSize(count.TotalCount, 10)
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
			tUsers.ID,
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Sex,
			tUsers.Dateofbirth,
			tUsers.PhoneNumber,
			tUserProps.Avatar.AS("usershort.avatar"),
		).
		FROM(
			tUnitStatus.
				LEFT_JOIN(tUsers,
					tUsers.ID.EQ(tUnitStatus.UserID),
				).
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tUnitStatus.UserID).
						AND(tUsers.Job.EQ(jet.String(userInfo.Job))),
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
	for i := 0; i < len(resp.Activity); i++ {
		if resp.Activity[i].UnitId > 0 && resp.Activity[i].User != nil {
			unit, err := s.state.GetUnit(ctx, userInfo.Job, resp.Activity[i].UnitId)
			if err != nil {
				return nil, errswrap.NewError(err, errorscentrum.ErrFailedQuery)
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
