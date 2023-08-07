package centrum

import (
	"context"

	database "github.com/galexrt/fivenet/gen/go/proto/resources/common/database"
	dispatch "github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/galexrt/fivenet/pkg/utils/dbutils"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"google.golang.org/protobuf/proto"
)

var (
	tDispatch       = table.FivenetCentrumDispatches.AS("dispatch")
	tDispatchStatus = table.FivenetCentrumDispatchesStatus.AS("dispatchstatus")
	tDispatchUnit   = table.FivenetCentrumDispatchesAsgmts.AS("dispatchassignment")
)

// TODO does it make sense to distinguish between "All" and "Assigned" dispatches here?
// A unit user would only get to see their assigned dispatches

func (s *Server) ListDispatches(ctx context.Context, req *ListDispatchesRequest) (*ListDispatchesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: CentrumService_ServiceDesc.ServiceName,
		Method:  "ListDispatches",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EVENT_TYPE_ERRORED),
	}
	defer s.a.Log(auditEntry, req)

	resp := &ListDispatchesResponse{}

	condition := jet.AND(
		tDispatch.Job.EQ(jet.String(userInfo.Job)),
		jet.OR(
			tDispatchStatus.ID.IS_NULL(),
			tDispatchStatus.ID.EQ(
				jet.RawInt("SELECT MAX(`dispatchstatus`.`id`) FROM `fivenet_centrum_dispatches_status` AS `dispatchstatus` WHERE `dispatchstatus`.`dispatch_id` = `dispatch`.`id`"),
			),
		),
	)

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
			tDispatch.Anon,
			tDispatch.UserID,
			tDispatchStatus.ID,
			tDispatchStatus.CreatedAt,
			tDispatchStatus.UnitID,
			tDispatchStatus.Status,
			tDispatchStatus.Reason,
			tDispatchStatus.Code,
			tDispatchStatus.UserID,
			tDispatchUnit.UnitID,
			tDispatchUnit.DispatchID,
			tDispatchUnit.CreatedAt,
			tDispatchUnit.ExpiresAt,
		).
		FROM(
			tDispatch.
				LEFT_JOIN(tDispatchStatus,
					tDispatchStatus.DispatchID.EQ(tDispatch.ID),
				).
				LEFT_JOIN(tDispatchUnit,
					tDispatchUnit.DispatchID.EQ(tDispatch.ID),
				),
		).
		WHERE(condition).
		ORDER_BY(
			tDispatch.ID.DESC(),
			tDispatchStatus.Status.ASC(),
		).
		LIMIT(150)

	if err := stmt.QueryContext(ctx, s.db, &resp.Dispatches); err != nil {
		return nil, err
	}

	auditEntry.State = int16(rector.EVENT_TYPE_VIEWED)

	return resp, nil
}

func (s *Server) CreateDispatch(ctx context.Context, req *CreateDispatchRequest) (*CreateDispatchResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: CentrumService_ServiceDesc.ServiceName,
		Method:  "CreateDispatch",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EVENT_TYPE_ERRORED),
	}
	defer s.a.Log(auditEntry, req)

	req.Dispatch.UserId = &userInfo.UserId
	dsp, err := s.createDispatch(ctx, req.Dispatch)
	if err != nil {
		return nil, err
	}

	resp := &CreateDispatchResponse{
		Dispatch: dsp,
	}

	data, err := proto.Marshal(resp.Dispatch)
	if err != nil {
		return nil, err
	}
	s.events.JS.Publish(s.buildSubject(TopicDispatch, TypeDispatchUpdated, userInfo, 0), data)

	data, err = proto.Marshal(resp.Dispatch.Status)
	if err != nil {
		return nil, err
	}
	for _, u := range dsp.Units {
		s.events.JS.Publish(s.buildSubject(TopicDispatch, TypeDispatchStatus, userInfo, u.UnitId), data)
	}

	auditEntry.State = int16(rector.EVENT_TYPE_CREATED)

	return resp, nil
}

func (s *Server) createDispatch(ctx context.Context, d *dispatch.Dispatch) (*dispatch.Dispatch, error) {
	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	tDispatch := table.FivenetCentrumDispatches
	stmt := tDispatch.
		INSERT(
			tDispatch.Job,
			tDispatch.Message,
			tDispatch.Description,
			tDispatch.Attributes,
			tDispatch.X,
			tDispatch.Y,
			tDispatch.Anon,
			tDispatch.UserID,
		).
		VALUES(
			d.Job,
			d.Message,
			d.Description,
			d.Attributes,
			d.X,
			d.Y,
			d.Anon,
			d.UserId,
		)

	result, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return nil, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	if err := s.addDispatchStatus(ctx, tx, uint64(lastId), *d.UserId, dispatch.DISPATCH_STATUS_NEW); err != nil {
		return nil, err
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	dsp, err := s.getDispatchFromDB(ctx, s.db, uint64(lastId))
	if err != nil {
		return nil, err
	}

	return dsp, nil
}

func (s *Server) addDispatchStatus(ctx context.Context, tx qrm.DB, dispatchId uint64, userId int32, status dispatch.DISPATCH_STATUS) error {
	tDispatchStatus := table.FivenetCentrumDispatchesStatus
	stmt := tDispatchStatus.
		INSERT(
			tDispatchStatus.DispatchID,
			tDispatchStatus.Status,
			tDispatchStatus.UserID,
		).
		VALUES(
			dispatchId,
			status,
			userId,
		)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	return nil
}

func (s *Server) UpdateDispatch(ctx context.Context, req *UpdateDispatchRequest) (*UpdateDispatchResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: CentrumService_ServiceDesc.ServiceName,
		Method:  "UpdateDispatch",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EVENT_TYPE_ERRORED),
	}
	defer s.a.Log(auditEntry, req)

	resp := &UpdateDispatchResponse{}

	stmt := tDispatch.
		UPDATE(
			tDispatch.Job,
			tDispatch.Message,
			tDispatch.Description,
			tDispatch.Attributes,
			tDispatch.X,
			tDispatch.Y,
			tDispatch.Anon,
			tDispatch.UserID,
		).
		SET(
			userInfo.Job,
			req.Dispatch.Message,
			req.Dispatch.Description,
			req.Dispatch.Attributes,
			req.Dispatch.X,
			req.Dispatch.Y,
			req.Dispatch.Anon,
			userInfo.UserId,
		).
		WHERE(jet.AND(
			tDispatch.Job.EQ(jet.String(userInfo.Job)),
			tDispatch.ID.EQ(jet.Uint64(req.Dispatch.Id)),
		))

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	dsp, err := s.getDispatchFromDB(ctx, s.db, req.Dispatch.Id)
	if err != nil {
		return nil, ErrFailedQuery
	}

	resp.Dispatch = dsp

	auditEntry.State = int16(rector.EVENT_TYPE_UPDATED)

	return nil, nil
}

func (s *Server) TakeDispatch(ctx context.Context, req *TakeDispatchRequest) (*TakeDispatchResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: CentrumService_ServiceDesc.ServiceName,
		Method:  "TakeDispatch",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EVENT_TYPE_ERRORED),
	}
	defer s.a.Log(auditEntry, req)

	dsp, err := s.getDispatchFromDB(ctx, s.db, req.DispatchId)
	if err != nil {
		return nil, ErrFailedQuery
	}

	unitId, err := s.getUnitIDFromUserID(ctx, userInfo.UserId)
	if err != nil {
		return nil, err
	}

	tDispatchUnit := table.FivenetCentrumDispatchesAsgmts
	stmt := tDispatchUnit.
		INSERT(
			tDispatchUnit.DispatchID,
			tDispatchUnit.UnitID,
			tDispatchUnit.ExpiresAt,
		).
		VALUES(
			dsp.Id,
			unitId,
			jet.NULL,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tDispatchUnit.ExpiresAt.SET(jet.TimestampExp(jet.NULL)),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		if !dbutils.IsDuplicateError(err) {
			return nil, err
		}
	}

	if _, err := s.updateDispatchStatus(ctx, userInfo, dsp, &dispatch.DispatchStatus{
		DispatchId: req.DispatchId,
		Status:     dispatch.DISPATCH_STATUS_UNIT_ASSIGNED,
		UserId:     &userInfo.UserId,
		UnitId:     unitId,
	}); err != nil {
		return nil, err
	}

	auditEntry.State = int16(rector.EVENT_TYPE_UPDATED)

	return &TakeDispatchResponse{
		Dispatch: dsp,
	}, nil
}

func (s *Server) UpdateDispatchStatus(ctx context.Context, req *UpdateDispatchStatusRequest) (*UpdateDispatchStatusResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: CentrumService_ServiceDesc.ServiceName,
		Method:  "UpdateDispatchStatus",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EVENT_TYPE_ERRORED),
	}
	defer s.a.Log(auditEntry, req)

	dsp, err := s.getDispatchFromDB(ctx, s.db, req.DispatchId)
	if err != nil {
		return nil, ErrFailedQuery
	}

	ok, err := s.checkIfUserPartOfDispatchUnits(ctx, userInfo, dsp)
	if err != nil {
		return nil, ErrFailedQuery
	}
	if !ok && !userInfo.SuperUser {
		return nil, ErrFailedQuery
	}

	unitId, err := s.getUnitIDFromUserID(ctx, userInfo.UserId)
	if err != nil {
		return nil, err
	}

	if _, err := s.updateDispatchStatus(ctx, userInfo, dsp, &dispatch.DispatchStatus{
		DispatchId: dsp.Id,
		Status:     req.Status,
		Code:       req.Code,
		Reason:     req.Reason,
		UnitId:     unitId,
	}); err != nil {
		return nil, err
	}

	auditEntry.State = int16(rector.EVENT_TYPE_UPDATED)

	return &UpdateDispatchStatusResponse{}, nil
}

func (s *Server) AssignDispatch(ctx context.Context, req *AssignDispatchRequest) (*AssignDispatchResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: CentrumService_ServiceDesc.ServiceName,
		Method:  "AssignDispatch",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EVENT_TYPE_ERRORED),
	}
	defer s.a.Log(auditEntry, req)

	dsp, err := s.getDispatchFromDB(ctx, s.db, req.DispatchId)
	if err != nil {
		return nil, ErrFailedQuery
	}

	if dsp.Job != userInfo.Job {
		return nil, ErrFailedQuery
	}

	addIds := make([]jet.Expression, len(req.ToAdd))
	for i := 0; i < len(req.ToAdd); i++ {
		addIds[i] = jet.Uint64(req.ToAdd[i])
	}
	removeIds := make([]jet.Expression, len(req.ToRemove))
	for i := 0; i < len(req.ToRemove); i++ {
		removeIds[i] = jet.Uint64(req.ToRemove[i])
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, ErrFailedQuery
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	tDispatchUnit := table.FivenetCentrumDispatchesAsgmts
	if len(removeIds) > 0 {
		stmt := tDispatchUnit.
			DELETE().
			WHERE(jet.AND(
				tDispatchUnit.DispatchID.EQ(jet.Uint64(dsp.Id)),
				tDispatchUnit.UnitID.IN(removeIds...),
			))

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return nil, err
		}

		for i := 0; i < len(dsp.Units); i++ {
			for k := 0; k < len(req.ToRemove); k++ {
				if dsp.Units[i].UnitId == req.ToRemove[k] {
					dsp.Units = utils.RemoveFromSlice(dsp.Units, i)

					if _, err := s.updateDispatchStatus(ctx, userInfo, dsp, &dispatch.DispatchStatus{
						DispatchId: dsp.Id,
						UnitId:     req.ToRemove[k],
						UserId:     &userInfo.UserId,
						Status:     dispatch.DISPATCH_STATUS_UNIT_UNASSIGNED,
					}); err != nil {
						return nil, ErrFailedQuery
					}

					continue
				}
			}
		}
	}

	if len(addIds) > 0 {
		for _, id := range addIds {
			stmt := tDispatchUnit.
				INSERT(
					tDispatchUnit.DispatchID,
					tDispatchUnit.UnitID,
					tDispatchUnit.ExpiresAt,
				).
				VALUES(
					dsp.Id,
					id,
					jet.CURRENT_TIMESTAMP().ADD(jet.INTERVAL(13, jet.SECOND)),
				).
				ON_DUPLICATE_KEY_UPDATE(
					tDispatchUnit.ExpiresAt.SET(jet.RawTimestamp("VALUES(`expires_at`)")),
				)

			if _, err := stmt.ExecContext(ctx, tx); err != nil {
				if !dbutils.IsDuplicateError(err) {
					return nil, err
				}
			}
		}

		assignments := []*dispatch.DispatchAssignment{}
		needsUpdate := []uint64{}
		for k := 0; k < len(req.ToAdd); k++ {
			found := false
			for i := 0; i < len(dsp.Units); i++ {
				if dsp.Units[i].UnitId == req.ToAdd[k] {
					found = true
					break
				}
			}

			unit, ok := s.getUnit(ctx, userInfo, req.ToAdd[k])
			if !ok {
				return nil, ErrFailedQuery
			}

			assignments = append(assignments, &dispatch.DispatchAssignment{
				UnitId:     unit.Id,
				DispatchId: dsp.Id,
				Unit:       unit,
			})

			if !found {
				needsUpdate = append(needsUpdate, unit.Id)
			}
		}
		dsp.Units = assignments

		for _, unitId := range needsUpdate {
			if _, err := s.updateDispatchStatus(ctx, userInfo, dsp, &dispatch.DispatchStatus{
				DispatchId: dsp.Id,
				UnitId:     unitId,
				UserId:     &userInfo.UserId,
				Status:     dispatch.DISPATCH_STATUS_UNIT_ASSIGNED,
			}); err != nil {
				return nil, ErrFailedQuery
			}
		}
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, ErrFailedQuery
	}

	data, err := proto.Marshal(dsp)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(req.ToRemove); i++ {
		s.events.JS.Publish(s.buildSubject(TopicDispatch, TypeDispatchUnassigned, userInfo, req.ToRemove[i]), data)
	}
	for i := 0; i < len(req.ToAdd); i++ {
		s.events.JS.Publish(s.buildSubject(TopicDispatch, TypeDispatchAssigned, userInfo, req.ToAdd[i]), data)
	}

	auditEntry.State = int16(rector.EVENT_TYPE_UPDATED)

	if len(dsp.Units) <= 0 {
		if _, err := s.updateDispatchStatus(ctx, userInfo, dsp, &dispatch.DispatchStatus{
			DispatchId: dsp.Id,
			Status:     dispatch.DISPATCH_STATUS_UNASSIGNED,
			UnitId:     0,
			UserId:     &userInfo.UserId,
		}); err != nil {
			return nil, err
		}
	}

	return &AssignDispatchResponse{}, nil
}

func (s *Server) ListDispatchActivity(ctx context.Context, req *ListActivityRequest) (*ListDispatchActivityResponse, error) {
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
		return nil, ErrFailedQuery
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
			tUser.Firstname,
			tUser.Lastname,
			tUser.Dateofbirth,
			tUser.Job,
		).
		FROM(
			tDispatchStatus.
				LEFT_JOIN(tUser,
					tUser.ID.EQ(tDispatchStatus.UserID),
				),
		).
		WHERE(
			tDispatchStatus.DispatchID.EQ(jet.Uint64(req.Id)),
		).
		ORDER_BY(tDispatchStatus.ID.DESC()).
		OFFSET(req.Pagination.Offset).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Activity); err != nil {
		return nil, err
	}

	resp.Pagination.Update(count.TotalCount, len(resp.Activity))

	return resp, nil
}
