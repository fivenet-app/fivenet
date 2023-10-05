package jobs

import (
	"context"
	"errors"

	database "github.com/galexrt/fivenet/gen/go/proto/resources/common/database"
	jobs "github.com/galexrt/fivenet/gen/go/proto/resources/jobs"
	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var (
	tRequests    = table.FivenetJobsRequests.AS("request")
	tReqTypes    = table.FivenetJobsRequestsTypes.AS("request_type")
	tReqComments = table.FivenetJobsRequestsComments.AS("request_")
)

func (s *Server) RequestsListEntries(ctx context.Context, req *RequestsListEntriesRequest) (*RequestsListEntriesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	condition := tRequests.Job.EQ(jet.String(userInfo.Job))

	countStmt := tRequests.
		SELECT(jet.COUNT(tRequests.ID).AS("datacount.totalcount")).
		FROM(tRequests).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		return nil, ErrFailedQuery
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(10)
	resp := &RequestsListEntriesResponse{
		Pagination: pag,
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

	tCreator := tUser.AS("creator")
	tApprover := tUser.AS("approver")
	stmt := tRequests.
		SELECT(
			tRequests.ID,
			tRequests.CreatedAt,
			tRequests.UpdatedAt,
			tRequests.DeletedAt,
			tRequests.Job,
			tRequests.TypeID,
			tRequests.Title,
			tRequests.Message,
			tRequests.Status,
			tRequests.CreatorID,
			tRequests.Approved,
			tRequests.ApproverID,
			tRequests.BeginsAt,
			tRequests.EndsAt,
		).
		FROM(
			tRequests.
				INNER_JOIN(tCreator,
					tCreator.ID.EQ(tRequests.CreatorID),
				).
				LEFT_JOIN(tApprover,
					tApprover.ID.EQ(tRequests.ApproverID),
				),
		).
		WHERE(jet.AND(
			tRequests.Job.EQ(jet.String(userInfo.Job)),
			tRequests.DeletedAt.IS_NULL(),
		)).
		OFFSET(req.Pagination.Offset).
		LIMIT(limit)

	var dest jobs.Request
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, ErrFailedQuery
		}
	}

	return resp, nil
}

func (s *Server) RequestsCreateEntry(ctx context.Context, req *RequestsCreateEntryRequest) (*RequestsCreateEntryResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: JobsService_ServiceDesc.ServiceName,
		Method:  "RequestsCreateEntry",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.auditer.Log(auditEntry, req)

	// TODO

	return &RequestsCreateEntryResponse{}, nil
}

func (s *Server) RequestsUpdateEntry(ctx context.Context, req *RequestsUpdateEntryRequest) (*RequestsUpdateEntryResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: JobsService_ServiceDesc.ServiceName,
		Method:  "RequestsUpdateEntry",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.auditer.Log(auditEntry, req)

	// TODO

	return &RequestsUpdateEntryResponse{}, nil
}

func (s *Server) RequestsDeleteEntry(ctx context.Context, req *RequestsDeleteEntryRequest) (*RequestsDeleteEntryResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: JobsService_ServiceDesc.ServiceName,
		Method:  "RequestsDeleteEntry",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.auditer.Log(auditEntry, req)

	stmt := tRequests.
		DELETE().
		WHERE(jet.AND(
			tRequests.Job.EQ(jet.String(userInfo.Job)),
			tRequests.ID.EQ(jet.Uint64(req.Id)),
		)).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, ErrFailedQuery
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)

	return &RequestsDeleteEntryResponse{}, nil
}

func (s *Server) RequestsListTypes(ctx context.Context, req *RequestsListTypesRequest) (*RequestsListTypesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	stmt := tReqTypes.
		SELECT(
			tReqTypes.ID,
			tReqTypes.CreatedAt,
			tReqTypes.UpdatedAt,
			tReqTypes.DeletedAt,
			tReqTypes.Job,
			tReqTypes.Name,
			tReqTypes.Description,
		).
		FROM(
			tReqTypes,
		).
		WHERE(jet.AND(
			tRequests.Job.EQ(jet.String(userInfo.Job)),
			tRequests.DeletedAt.IS_NULL(),
		)).
		LIMIT(15)

	var dest []*jobs.RequestType
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, ErrFailedQuery
		}
	}

	return &RequestsListTypesResponse{
		Types: dest,
	}, nil
}

func (s *Server) RequestsCreateOrUpdateType(ctx context.Context, req *RequestsCreateOrUpdateTypeRequest) (*RequestsCreateOrUpdateTypeResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: JobsService_ServiceDesc.ServiceName,
		Method:  "RequestsCreateOrUpdateType",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
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
	if req.RequestType.Id <= 0 {
		stmt := tReqTypes.
			INSERT(
				tReqTypes.Job,
				tReqTypes.Name,
				tReqTypes.Description,
				tReqTypes.Weight,
			).
			VALUES(
				userInfo.Job,
				req.RequestType.Name,
				req.RequestType.Description,
				req.RequestType.Weight,
			)

		result, err := stmt.ExecContext(ctx, tx)
		if err != nil {
			return nil, err
		}

		lastId, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		req.RequestType.Id = uint64(lastId)

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)
	} else {
		stmt := tReqTypes.
			UPDATE(
				tReqTypes.Name,
				tReqTypes.Description,
				tReqTypes.Weight,
			).
			SET(
				req.RequestType.Name,
				req.RequestType.Description,
				req.RequestType.Weight,
			).
			WHERE(jet.AND(
				tReqTypes.Job.EQ(jet.String(userInfo.Job)),
				tReqTypes.ID.EQ(jet.Uint64(req.RequestType.Id)),
			))

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return nil, err
		}

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		auditEntry.State = int16(rector.EventType_EVENT_TYPE_ERRORED)
		return nil, ErrFailedQuery
	}

	requestType, err := s.getRequestType(ctx, userInfo.Job, req.RequestType.Id)
	if err != nil {
		return nil, ErrFailedQuery
	}

	return &RequestsCreateOrUpdateTypeResponse{
		RequestType: requestType,
	}, nil
}

func (s *Server) RequestsDeleteType(ctx context.Context, req *RequestsDeleteTypeRequest) (*RequestsDeleteTypeResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: JobsService_ServiceDesc.ServiceName,
		Method:  "RequestsDeleteType",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.auditer.Log(auditEntry, req)

	stmt := tReqTypes.
		DELETE().
		WHERE(jet.AND(
			tReqTypes.Job.EQ(jet.String(userInfo.Job)),
			tReqTypes.ID.EQ(jet.Uint64(req.Id)),
		)).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, ErrFailedQuery
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)

	return &RequestsDeleteTypeResponse{}, nil
}

func (s *Server) RequestsPostComment(ctx context.Context, req *RequestsPostCommentRequest) (*RequestsPostCommentResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: JobsService_ServiceDesc.ServiceName,
		Method:  "RequestsPostComment",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.auditer.Log(auditEntry, req)

	// TODO

	return &RequestsPostCommentResponse{}, nil
}

func (s *Server) RequestsDeleteComment(ctx context.Context, req *RequestsDeleteCommentRequest) (*RequestsDeleteCommentResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: JobsService_ServiceDesc.ServiceName,
		Method:  "RequestsDeleteComment",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.auditer.Log(auditEntry, req)

	// TODO

	return &RequestsDeleteCommentResponse{}, nil
}
