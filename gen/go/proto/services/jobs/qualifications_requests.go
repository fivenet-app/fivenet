package jobs

import (
	"context"
	"errors"

	database "github.com/galexrt/fivenet/gen/go/proto/resources/common/database"
	jobs "github.com/galexrt/fivenet/gen/go/proto/resources/jobs"
	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	errorsjobs "github.com/galexrt/fivenet/gen/go/proto/services/jobs/errors"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	"github.com/galexrt/fivenet/pkg/grpc/errswrap"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var (
	tQualiRequests = table.FivenetJobsQualificationsRequests.AS("qualificationrequest")
	tApprover      = tUser.AS("approver")
)

func (s *Server) ListQualificationRequests(ctx context.Context, req *ListQualificationRequestsRequest) (*ListQualificationRequestsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	condition := tQualiRequests.UserID.EQ(jet.Int32(userInfo.UserId))

	if req.QualificationId != nil {
		ok, err := s.checkIfUserHasAccessToQuali(ctx, *req.QualificationId, userInfo, jobs.AccessLevel_ACCESS_LEVEL_EDIT)
		if err != nil {
			return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
		}
		if !ok {
			return nil, errorsjobs.ErrFailedQuery
		}

		condition = condition.AND(tQualiRequests.QualificationID.EQ(jet.Uint64(*req.QualificationId)))
	} else {
		condition = condition.AND(jet.AND(
			tQuali.DeletedAt.IS_NULL(),
			jet.OR(
				tQuali.CreatorID.EQ(jet.Int32(userInfo.UserId)),
				jet.AND(
					tQJobAccess.Access.IS_NOT_NULL(),
					tQJobAccess.Access.NOT_EQ(jet.Int32(int32(jobs.AccessLevel_ACCESS_LEVEL_BLOCKED))),
				),
			),
		))
	}

	countStmt := tQualiRequests.
		SELECT(
			jet.COUNT(tQualiRequests.ID).AS("datacount.totalcount"),
		).
		FROM(tQualiRequests).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(count.TotalCount, QualificationsPageSize)
	resp := &ListQualificationRequestsResponse{
		Pagination: pag,
		Requests:   []*jobs.QualificationRequest{},
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

	stmt := tQualiRequests.
		SELECT(
			tQualiRequests.ID,
			tQualiRequests.CreatedAt,
			tQualiRequests.QualificationID,
			tQualiRequests.UserID,
			tQualiRequests.UserComment,
			tQualiRequests.Approved,
			tQualiRequests.ApprovedAt,
			tQualiRequests.ApproverComment,
			tQualiRequests.ApproverID,
			tQualiRequests.ApproverJob,
			tCreator.ID,
			tCreator.Identifier,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Dateofbirth,
		).
		FROM(
			tQualiRequests.
				LEFT_JOIN(tCreator,
					tQuali.CreatorID.EQ(tCreator.ID),
				).
				LEFT_JOIN(tQJobAccess,
					tQJobAccess.QualificationID.EQ(tQuali.ID).
						AND(tQJobAccess.Job.EQ(jet.String(userInfo.Job))).
						AND(tQJobAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
				),
		).
		WHERE(condition).
		OFFSET(req.Pagination.Offset).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Requests); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
	}

	return resp, nil
}

func (s *Server) CreateOrUpdateQualificationRequest(ctx context.Context, req *CreateOrUpdateQualificationRequestRequest) (*CreateOrUpdateQualificationRequestResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: JobsQualificationsService_ServiceDesc.ServiceName,
		Method:  "CreateOrUpdateQualificationRequest",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	ok, err := s.checkIfUserHasAccessToQuali(ctx, req.Request.QualificationId, userInfo, jobs.AccessLevel_ACCESS_LEVEL_GRADE)
	if err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}
	if !ok {
		return nil, errorsjobs.ErrFailedQuery
	}

	if req.Request.QualificationId <= 0 {
		stmt := tQualiRequests.
			INSERT(
				tQualiRequests.QualificationID,
				tQualiRequests.UserID,
				tQualiRequests.UserComment,
				tQualiRequests.Approved,
				tQualiRequests.ApprovedAt,
				tQualiRequests.ApproverComment,
				tQualiRequests.ApproverID,
				tQualiRequests.ApproverJob,
			).
			VALUES(
				req.Request.QualificationId,
				req.Request.UserId,
				req.Request.UserComment,
				req.Request.Approved,
				req.Request.ApprovedAt,
				req.Request.ApproverComment,
				req.Request.ApproverId,
				req.Request.ApproverJob,
			)

		res, err := stmt.ExecContext(ctx, s.db)
		if err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}

		lastId, err := res.LastInsertId()
		if err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}

		req.Request.Id = uint64(lastId)

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)
	} else {
		stmt := tQualiRequests.
			UPDATE(
				tQualiRequests.QualificationID,
				tQualiRequests.UserID,
				tQualiRequests.UserID,
				tQualiRequests.UserComment,
				tQualiRequests.Approved,
				tQualiRequests.ApprovedAt,
				tQualiRequests.ApproverComment,
				tQualiRequests.ApproverID,
				tQualiRequests.ApproverJob,
			).
			SET(
				tQualiRequests.QualificationID.SET(jet.Uint64(req.Request.QualificationId)),
				tQualiRequests.UserID.SET(jet.Int32(req.Request.UserId)),
				tQualiRequests.UserComment.SET(jet.String(*req.Request.UserComment)),
				tQualiRequests.Approved.SET(jet.Bool(*req.Request.Approved)),
				tQualiRequests.ApprovedAt.SET(jet.TimestampT(req.Request.ApprovedAt.AsTime())),
				tQualiRequests.ApproverComment.SET(jet.String(*req.Request.ApproverComment)),
				tQualiRequests.ApproverID.SET(jet.Int32(*req.Request.ApproverId)),
				tQualiRequests.ApproverJob.SET(jet.String(*req.Request.ApproverJob)),
			).
			WHERE(jet.AND(
				tQualiRequests.QualificationID.EQ(jet.Uint64(req.Request.QualificationId)),
				tQualiRequests.UserID.EQ(jet.Int32(userInfo.UserId)),
			))

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)
	}

	request, err := s.getQualificationRequest(ctx, tQualiRequests.ID.EQ(jet.Uint64(req.Request.QualificationId)).
		AND(tQualiRequests.UserID.EQ(jet.Int32(userInfo.UserId))), userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	return &CreateOrUpdateQualificationRequestResponse{
		Request: request,
	}, nil
}

func (s *Server) getQualificationRequest(ctx context.Context, condition jet.BoolExpression, userInfo *userinfo.UserInfo) (*jobs.QualificationRequest, error) {
	var request jobs.QualificationRequest

	stmt := tQualiRequests.
		SELECT(
			tQualiRequests.ID,
			tQualiRequests.CreatedAt,
			tQualiRequests.DeletedAt,
			tQualiRequests.QualificationID,
			tQualiRequests.UserID,
			tUser.ID,
			tUser.Identifier,
			tUser.Job,
			tUser.JobGrade,
			tUser.Firstname,
			tUser.Lastname,
			tUser.Dateofbirth,
			tQualiRequests.UserComment,
			tQualiRequests.Approved,
			tQualiRequests.ApprovedAt,
			tQualiRequests.ApproverComment,
			tQualiRequests.ApproverID,
			tQualiRequests.ApproverJob,
			tApprover.ID,
			tApprover.Identifier,
			tApprover.Job,
			tApprover.JobGrade,
			tApprover.Firstname,
			tApprover.Lastname,
			tApprover.Dateofbirth,
		).
		FROM(tQualiRequests.
			LEFT_JOIN(tCreator,
				tCreator.ID.EQ(tQualiRequests.UserID),
			).
			LEFT_JOIN(tApprover,
				tApprover.ID.EQ(tQualiRequests.ApproverID),
			),
		).
		WHERE(condition).
		LIMIT(1)

	if err := stmt.QueryContext(ctx, s.db, &request); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
		}
	}

	if request.User != nil {
		s.enricher.EnrichJobInfoSafe(userInfo, request.User)
	}

	if request.Approver != nil {
		s.enricher.EnrichJobInfoSafe(userInfo, request.Approver)
	}

	return nil, nil
}
