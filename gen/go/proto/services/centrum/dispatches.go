package centrum

import (
	"context"

	database "github.com/galexrt/fivenet/gen/go/proto/resources/common/database"
	dispatch "github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/utils/dbutils"
	"github.com/galexrt/fivenet/pkg/utils/syncx"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
)

var (
	tDispatch       = table.FivenetCentrumDispatches.AS("dispatch")
	tDispatchStatus = table.FivenetCentrumDispatchesStatus.AS("dispatchstatus")
	tDispatchUnit   = table.FivenetCentrumDispatchesAsgmts.AS("dispatchassignment")
)

func (s *Server) loadDispatches(ctx context.Context) error {
	condition := tDispatchStatus.ID.IS_NULL().OR(
		tDispatchStatus.ID.EQ(
			jet.RawInt(`SELECT MAX(dispatchstatus.id) FROM fivenet_centrum_dispatches_status AS dispatchstatus WHERE dispatchstatus.dispatch_id = dispatch.id`),
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
			tDispatch.CreatedAt.ASC(),
			tDispatch.Job.ASC(),
			tDispatchStatus.Status.ASC(),
		).LIMIT(300)

	dispatches := []*dispatch.Dispatch{}
	if err := stmt.QueryContext(ctx, s.db, &dispatches); err != nil {
		return err
	}

	for i := 0; i < len(dispatches); i++ {
		jobDispatches, _ := s.dispatches.LoadOrStore(dispatches[i].Job, &syncx.Map[uint64, *dispatch.Dispatch]{})
		jobDispatches.Store(dispatches[i].Id, dispatches[i])
	}

	return nil
}

func (s *Server) ListDispatches(ctx context.Context, req *ListDispatchesRequest) (*ListDispatchesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: CentrumService_ServiceDesc.ServiceName,
		Method:  "ListDispatches",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EVENT_TYPE_ERRORED),
	}
	defer s.a.AddEntryWithData(auditEntry, req)

	resp := &ListDispatchesResponse{
		Dispatches: []*dispatch.Dispatch{},
	}

	val, ok := s.dispatches.Load(userInfo.Job)
	if val == nil || !ok {
		return resp, nil
	}

	dispatches := []*dispatch.Dispatch{}
	val.Range(func(key uint64, value *dispatch.Dispatch) bool {
		dispatches = append(dispatches, value)
		return true
	})

	resp.Dispatches = dispatches

	auditEntry.State = int16(rector.EVENT_TYPE_VIEWED)

	return resp, nil
}

func (s *Server) CreateDispatch(ctx context.Context, req *CreateDispatchRequest) (*CreateDispatchResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: CentrumService_ServiceDesc.ServiceName,
		Method:  "CreateOrUpdateUnit",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EVENT_TYPE_ERRORED),
	}
	defer s.a.AddEntryWithData(auditEntry, req)

	resp := &CreateDispatchResponse{}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, ErrFailedQuery
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
			userInfo.Job,
			req.Dispatch.Message,
			req.Dispatch.Description,
			req.Dispatch.Attributes,
			req.Dispatch.X,
			req.Dispatch.Y,
			req.Dispatch.Anon,
			userInfo.UserId,
		)

	result, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return nil, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	req.Dispatch.Id = uint64(lastId)

	dsp, err := s.getDispatchFromDB(ctx, tx, uint64(lastId))
	if err != nil {
		return nil, ErrFailedQuery
	}

	resp.Dispatch = dsp

	tDispatchStatus := table.FivenetCentrumDispatchesStatus
	stmt = tDispatchStatus.
		INSERT(
			tDispatchStatus.DispatchID,
			tDispatchStatus.Status,
			tDispatchStatus.UserID,
		).
		VALUES(
			uint64(lastId),
			dispatch.DISPATCH_STATUS_NEW,
			userInfo.UserId,
		)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return nil, err
	}

	val, ok := s.dispatches.Load(userInfo.Job)
	if val != nil && ok {
		val.Store(dsp.Id, dsp)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, ErrFailedQuery
	}

	auditEntry.State = int16(rector.EVENT_TYPE_CREATED)

	return resp, nil
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
	defer s.a.AddEntryWithData(auditEntry, req)

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

	val, ok := s.dispatches.Load(userInfo.Job)
	if val != nil && ok {
		val.Store(dsp.Id, dsp)
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
	defer s.a.AddEntryWithData(auditEntry, req)

	// TODO

	auditEntry.State = int16(rector.EVENT_TYPE_UPDATED)

	return nil, nil
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
	defer s.a.AddEntryWithData(auditEntry, req)

	dsp, err := s.getDispatchFromDB(ctx, s.db, req.DispatchId)
	if err != nil {
		return nil, ErrFailedQuery
	}

	found := false
	for i := 0; i < len(dsp.Units); i++ {
		unit, ok := s.getUnit(ctx, userInfo, dsp.Units[i].UnitId)
		if !ok {
			return nil, ErrFailedQuery
		}
		for i := 0; i < len(unit.Users); i++ {
			if unit.Users[i].UserId == userInfo.UserId {
				found = true
				break
			}
		}
	}
	if !found {
		return nil, ErrFailedQuery
	}

	stmt := tDispatchStatus.
		INSERT(
			tDispatchStatus.UnitID,
			tDispatchStatus.Status,
			tDispatchStatus.Reason,
			tDispatchStatus.Code,
			tDispatchStatus.UserID,
		).
		VALUES(
			req.DispatchId,
			req.Status,
			req.Reason,
			req.Code,
			userInfo.UserId,
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
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
	defer s.a.AddEntryWithData(auditEntry, req)

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

	if len(removeIds) > 0 {
		stmt := tUnitUser.
			DELETE().
			WHERE(jet.AND(
				tUnitUser.UnitID.EQ(jet.Uint64(dsp.Id)),
				tUnitUser.UserID.IN(removeIds...),
			))

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return nil, err
		}

		for i := 0; i < len(dsp.Units); i++ {
			for k := 0; k < len(req.ToRemove); k++ {
				if dsp.Units[i].UnitId == req.ToRemove[k] {
					break
				}
			}
		}
	}

	if len(addIds) > 0 {
		for _, id := range addIds {
			stmt := tUnitUser.
				INSERT(
					tUnitUser.UnitID,
					tUnitUser.UserID,
				).
				VALUES(
					dsp.Id,
					id,
				)

			if _, err := stmt.ExecContext(ctx, tx); err != nil {
				if !dbutils.IsDuplicateError(err) {
					return nil, err
				}
			}

		}

		found := []uint64{}
		for k := 0; k < len(req.ToAdd); k++ {
			for i := 0; i < len(dsp.Units); i++ {
				if dsp.Units[i].UnitId == req.ToAdd[k] {
					found = append(found, req.ToAdd[k])
				}
			}
		}

		assignments := []*dispatch.DispatchAssignment{}
		for _, dId := range found {
			unit, ok := s.getUnit(ctx, userInfo, dId)
			if !ok {
				return nil, ErrFailedQuery
			}

			assignments = append(assignments, &dispatch.DispatchAssignment{
				UnitId:     dId,
				DispatchId: dsp.Id,
				Unit:       unit,
			})
		}
		dsp.Units = assignments
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, ErrFailedQuery
	}

	auditEntry.State = int16(rector.EVENT_TYPE_UPDATED)

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

func (s *Server) Stream(req *StreamRequest, srv CentrumService_StreamServer) error {
	signalCh := s.broker.Subscribe()
	defer s.broker.Unsubscribe(signalCh)

	resp := &StreamResponse{}

	for {
		select {
		case <-srv.Context().Done():
			return nil
		case data := <-signalCh:
			switch data.(type) {
			case StreamResponse_UnitChange:
				resp.Change = data.(*StreamResponse_UnitChange)
			case StreamResponse_UnitStatus:
				resp.Change = data.(*StreamResponse_UnitStatus)
			case StreamResponse_DispatchChange:
				resp.Change = data.(*StreamResponse_DispatchChange)
			case StreamResponse_DispatchStatus:
				resp.Change = data.(*StreamResponse_DispatchStatus)

			case StreamResponse_DispatchAssigned:
				resp.Change = data.(*StreamResponse_DispatchAssigned)
			}
		}

		if err := srv.Send(resp); err != nil {
			return err
		}
	}
}
