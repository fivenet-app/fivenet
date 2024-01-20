package jobs

import (
	"context"
	"errors"

	jobs "github.com/galexrt/fivenet/gen/go/proto/resources/jobs"
	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	errorsjobs "github.com/galexrt/fivenet/gen/go/proto/services/jobs/errors"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/grpc/errswrap"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) ListRequestTypes(ctx context.Context, req *ListRequestTypesRequest) (*ListRequestTypesResponse, error) {
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
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
		}
	}

	return &ListRequestTypesResponse{
		Types: dest,
	}, nil
}

func (s *Server) CreateOrUpdateRequestType(ctx context.Context, req *CreateOrUpdateRequestTypeRequest) (*CreateOrUpdateRequestTypeResponse, error) {
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
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	tReqTypes := table.FivenetJobsRequestsTypes
	// No request type id set
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
			return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
		}

		lastId, err := result.LastInsertId()
		if err != nil {
			return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
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
			return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
		}

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		auditEntry.State = int16(rector.EventType_EVENT_TYPE_ERRORED)
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}

	requestType, err := s.getRequestType(ctx, userInfo.Job, req.RequestType.Id)
	if err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}

	return &CreateOrUpdateRequestTypeResponse{
		RequestType: requestType,
	}, nil
}

func (s *Server) RequestsDeleteType(ctx context.Context, req *DeleteRequestTypeRequest) (*DeleteRequestTypeResponse, error) {
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
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)

	return &DeleteRequestTypeResponse{}, nil
}
