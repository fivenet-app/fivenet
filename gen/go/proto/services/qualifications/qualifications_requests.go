package qualifications

import (
	"context"
	"errors"
	"time"

	database "github.com/galexrt/fivenet/gen/go/proto/resources/common/database"
	"github.com/galexrt/fivenet/gen/go/proto/resources/qualifications"
	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	errorsqualifications "github.com/galexrt/fivenet/gen/go/proto/services/jobs/errors"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	"github.com/galexrt/fivenet/pkg/grpc/errswrap"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var (
	tQualiRequests = table.FivenetQualificationsRequests.AS("qualificationrequest")
	tApprover      = tUser.AS("approver")
)

func (s *Server) ListQualificationRequests(ctx context.Context, req *ListQualificationRequestsRequest) (*ListQualificationRequestsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	condition := tQualiRequests.UserID.EQ(jet.Int32(userInfo.UserId))

	if req.QualificationId != nil {
		ok, err := s.checkIfUserHasAccessToQuali(ctx, *req.QualificationId, userInfo, qualifications.AccessLevel_ACCESS_LEVEL_EDIT)
		if err != nil {
			return nil, errswrap.NewError(errorsqualifications.ErrFailedQuery, err)
		}
		if !ok {
			return nil, errorsqualifications.ErrFailedQuery
		}

		condition = condition.AND(tQualiRequests.QualificationID.EQ(jet.Uint64(*req.QualificationId)))
	} else {
		condition = condition.AND(jet.AND(
			tQuali.DeletedAt.IS_NULL(),
			jet.OR(
				tQuali.CreatorID.EQ(jet.Int32(userInfo.UserId)),
				jet.AND(
					tQJobAccess.Access.IS_NOT_NULL(),
					tQJobAccess.Access.NOT_EQ(jet.Int32(int32(qualifications.AccessLevel_ACCESS_LEVEL_BLOCKED))),
				),
			),
		))
	}

	countStmt := tQualiRequests.
		SELECT(
			jet.COUNT(tQualiRequests.QualificationID).AS("datacount.totalcount"),
		).
		FROM(
			tQualiRequests.
				INNER_JOIN(tQuali,
					tQuali.ID.EQ(tQualiRequests.QualificationID),
				).
				LEFT_JOIN(tQJobAccess,
					tQJobAccess.QualificationID.EQ(tQuali.ID).
						AND(tQJobAccess.Job.EQ(jet.String(userInfo.Job))).
						AND(tQJobAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
				),
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		return nil, errswrap.NewError(errorsqualifications.ErrFailedQuery, err)
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(count.TotalCount, QualificationsPageSize)
	resp := &ListQualificationRequestsResponse{
		Pagination: pag,
		Requests:   []*qualifications.QualificationRequest{},
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

	stmt := tQualiRequests.
		SELECT(
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
				INNER_JOIN(tQuali,
					tQuali.ID.EQ(tQualiRequests.QualificationID),
				).
				LEFT_JOIN(tCreator,
					tQualiRequests.UserID.EQ(tCreator.ID),
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
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
	}

	return resp, nil
}

func (s *Server) CreateOrUpdateQualificationRequest(ctx context.Context, req *CreateOrUpdateQualificationRequestRequest) (*CreateOrUpdateQualificationRequestResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: QualificationsService_ServiceDesc.ServiceName,
		Method:  "CreateOrUpdateQualificationRequest",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	ok, err := s.checkIfUserHasAccessToQuali(ctx, req.Request.QualificationId, userInfo, qualifications.AccessLevel_ACCESS_LEVEL_GRADE)
	if err != nil {
		return nil, errswrap.NewError(errorsqualifications.ErrFailedQuery, err)
	}
	// If user can grade a qualification, they are treated as an "approver" of requests
	if ok {
		stmt := tQualiRequests.
			UPDATE(
				tQualiRequests.Approved,
				tQualiRequests.ApprovedAt,
				tQualiRequests.ApproverComment,
				tQualiRequests.ApproverID,
				tQualiRequests.ApproverJob,
			).
			SET(
				req.Request.Approved,
				time.Now(),
				req.Request.ApproverComment,
				userInfo.UserId,
				userInfo.Job,
			).
			WHERE(jet.AND(
				tQualiRequests.QualificationID.EQ(jet.Uint64(req.Request.QualificationId)),
				tQualiRequests.UserID.EQ(jet.Int32(req.Request.UserId)),
			))

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)
	} else {
		ok, err = s.checkIfUserHasAccessToQuali(ctx, req.Request.QualificationId, userInfo, qualifications.AccessLevel_ACCESS_LEVEL_REQUEST)
		if err != nil {
			return nil, errswrap.NewError(errorsqualifications.ErrFailedQuery, err)
		}
		if !ok {
			return nil, errorsqualifications.ErrFailedQuery
		}

		tQualiRequests := table.FivenetQualificationsRequests
		stmt := tQualiRequests.
			INSERT(
				tQualiRequests.QualificationID,
				tQualiRequests.UserID,
				tQualiRequests.UserComment,
			).
			VALUES(
				req.Request.QualificationId,
				userInfo.UserId,
				req.Request.UserComment,
			).
			ON_DUPLICATE_KEY_UPDATE(
				tQualiRequests.UserComment.SET(jet.StringExp(jet.Raw("VALUES(`user_comment`)"))),
			)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)
	}

	request, err := s.getQualificationRequest(ctx, req.Request.QualificationId, req.Request.UserId, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	return &CreateOrUpdateQualificationRequestResponse{
		Request: request,
	}, nil
}

func (s *Server) getQualificationRequest(ctx context.Context, requestId uint64, userId int32, userInfo *userinfo.UserInfo) (*qualifications.QualificationRequest, error) {
	var request qualifications.QualificationRequest

	stmt := tQualiRequests.
		SELECT(
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
			LEFT_JOIN(tUser,
				tUser.ID.EQ(tQualiRequests.UserID),
			).
			LEFT_JOIN(tApprover,
				tApprover.ID.EQ(tQualiRequests.ApproverID),
			),
		).
		WHERE(jet.AND(
			tQualiRequests.QualificationID.EQ(jet.Uint64(requestId)),
			tQualiRequests.UserID.EQ(jet.Int32(userId)),
		)).
		LIMIT(1)

	if err := stmt.QueryContext(ctx, s.db, &request); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if request.User != nil {
		s.enricher.EnrichJobInfoSafe(userInfo, request.User)
	}

	if request.Approver != nil {
		s.enricher.EnrichJobInfoSafe(userInfo, request.Approver)
	}

	return &request, nil
}

func (s *Server) DeleteQualificationReq(ctx context.Context, req *DeleteQualificationReqRequest) (*DeleteQualificationReqResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: QualificationsService_ServiceDesc.ServiceName,
		Method:  "DeleteQualificationReq",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	re, err := s.getQualificationRequest(ctx, req.QualificationId, req.UserId, userInfo)
	if err != nil {
		return nil, errswrap.NewError(errorsqualifications.ErrFailedQuery, err)
	}

	ok, err := s.checkIfUserHasAccessToQuali(ctx, re.QualificationId, userInfo, qualifications.AccessLevel_ACCESS_LEVEL_MANAGE)
	if err != nil {
		return nil, errswrap.NewError(errorsqualifications.ErrFailedQuery, err)
	}
	if !ok {
		return nil, errorsqualifications.ErrFailedQuery
	}

	stmt := tQualiRequests.
		UPDATE(
			tQualiRequests.DeletedAt,
		).
		SET(
			jet.CURRENT_TIMESTAMP(),
		).
		WHERE(jet.AND(
			tQualiRequests.QualificationID.EQ(jet.Uint64(re.QualificationId)),
			tQualiRequests.UserID.EQ(jet.Int32(re.UserId)),
		))

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)

	return &DeleteQualificationReqResponse{}, nil
}
