package centrum

import (
	"context"

	database "github.com/galexrt/fivenet/gen/go/proto/resources/common/database"
	dispatch "github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
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
			tDispatch.Active,
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
		jobDispatches, _ := s.dispatches.LoadOrStore(*dispatches[i].Job, &syncx.Map[uint64, *dispatch.Dispatch]{})
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

	// TODO

	return nil, nil
}

func (s *Server) UpdateDispatch(ctx context.Context, req *UpdateDispatchRequest) (*UpdateDispatchResponse, error) {

	// TODO

	return nil, nil
}

func (s *Server) TakeDispatch(ctx context.Context, req *TakeDispatchRequest) (*TakeDispatchResponse, error) {

	// TODO

	return nil, nil
}

func (s *Server) UpdateDispatchStatus(ctx context.Context, req *UpdateDispatchStatusRequest) (*UpdateDispatchStatusResponse, error) {

	// TODO

	return nil, nil
}

func (s *Server) AssignDispatch(ctx context.Context, req *AssignDispatchRequest) (*AssignDispatchResponse, error) {

	// TODO

	return nil, nil
}

func (s *Server) ListDispatchActivity(ctx context.Context, req *ListActivityRequest) (*ListDispatchActivityResponse, error) {
	countStmt := tDispatchStatus.
		SELECT(
			jet.COUNT(jet.DISTINCT(tDispatchStatus.ID)).AS("datacount.totalcount"),
		).
		FROM(tDispatchStatus).
		WHERE(
			tDispatchStatus.UnitID.EQ(jet.Uint64(req.Id)),
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
		).
		FROM(tDispatchStatus).
		WHERE(
			tDispatchStatus.DispatchID.EQ(jet.Uint64(req.Id)),
		).
		OFFSET(req.Pagination.Offset).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Activity); err != nil {
		return nil, err
	}

	resp.Pagination.Update(count.TotalCount, len(resp.Activity))

	return resp, nil
}

func (s *Server) Stream(req *CentrumStreamRequest, srv CentrumService_StreamServer) error {

	// TODO

	return nil
}
