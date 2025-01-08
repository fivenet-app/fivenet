package qualifications

import (
	"context"
	"errors"
	"fmt"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common"
	database "github.com/fivenet-app/fivenet/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/notifications"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/qualifications"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/rector"
	errorsqualifications "github.com/fivenet-app/fivenet/gen/go/proto/services/qualifications/errors"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/pkg/utils/dbutils/tables"
	"github.com/fivenet-app/fivenet/query/fivenet/model"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var tQualiRequests = table.FivenetQualificationsRequests.AS("qualificationrequest")

func (s *Server) ListQualificationRequests(ctx context.Context, req *ListQualificationRequestsRequest) (*ListQualificationRequestsResponse, error) {
	if req.QualificationId != nil {
		trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.qualifications.id", int64(*req.QualificationId)))
	}
	if req.UserId != nil {
		trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.qualifications.user_id", int64(*req.UserId)))
	}

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	tQuali := tQuali.AS("qualificationshort")
	tUser := tables.Users().AS("user")
	tApprover := tUser.AS("approver")

	condition := tQualiRequests.DeletedAt.IS_NULL().
		AND(tQualiRequests.Status.NOT_EQ(jet.Int16(int16(qualifications.RequestStatus_REQUEST_STATUS_COMPLETED))))

	if req.QualificationId != nil {
		check, err := s.access.CanUserAccessTarget(ctx, *req.QualificationId, userInfo, qualifications.AccessLevel_ACCESS_LEVEL_GRADE)
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
		if !check {
			return nil, errorsqualifications.ErrFailedQuery
		}

		condition = condition.AND(tQualiRequests.QualificationID.EQ(jet.Uint64(*req.QualificationId)))
	} else {
		condition = condition.AND(jet.AND(
			tQuali.DeletedAt.IS_NULL(),
			jet.OR(
				tQualiRequests.UserID.EQ(jet.Int32(userInfo.UserId)),
				jet.AND(
					tQJobAccess.Access.IS_NOT_NULL(),
					jet.OR(
						tQJobAccess.Access.GT(jet.Int32(int32(qualifications.AccessLevel_ACCESS_LEVEL_GRADE))),
						jet.AND(
							tQJobAccess.Access.GT(jet.Int32(int32(qualifications.AccessLevel_ACCESS_LEVEL_BLOCKED))),
							tQualiRequests.UserID.EQ(jet.Int32(userInfo.UserId)),
						),
					),
				),
			),
		))
	}

	countColumn := jet.Expression(tQualiRequests.QualificationID)
	if req.UserId != nil {
		condition = condition.AND(tUser.Job.EQ(jet.String(userInfo.Job))).AND(tQualiRequests.UserID.EQ(jet.Int32(*req.UserId)))
	} else {
		if req.QualificationId == nil {
			condition = condition.AND(tUser.Job.EQ(jet.String(userInfo.Job))).AND(tQualiRequests.UserID.EQ(jet.Int32(userInfo.UserId)))
			countColumn = jet.DISTINCT(tQualiRequests.QualificationID)
		} else {
			countColumn = jet.DISTINCT(tQualiRequests.UserID)
		}
	}

	if len(req.Status) > 0 {
		statuses := []jet.Expression{}
		for i := 0; i < len(req.Status); i++ {
			statuses = append(statuses, jet.Int16(int16(req.Status[i])))
		}

		condition = condition.AND(tQualiRequests.Status.IN(statuses...))
	} else {
		condition = condition.AND(tQualiRequests.Status.NOT_EQ(jet.Int16(int16(qualifications.RequestStatus_REQUEST_STATUS_COMPLETED))))
	}

	countStmt := tQualiRequests.
		SELECT(
			jet.COUNT(countColumn).AS("datacount.totalcount"),
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
				).
				LEFT_JOIN(tUser,
					tQualiRequests.UserID.EQ(tUser.ID),
				),
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(count.TotalCount, QualificationsPageSize)
	resp := &ListQualificationRequestsResponse{
		Pagination: pag,
		Requests:   []*qualifications.QualificationRequest{},
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

	// Convert proto sort to db sorting
	orderBys := []jet.OrderByClause{}
	if req.Sort != nil {
		var column jet.Column
		switch req.Sort.Column {
		case "status":
			column = tQualiRequests.Status
		case "approvedAt":
			column = tQualiRequests.ApprovedAt
		case "createdAt":
			fallthrough
		default:
			column = tQualiRequests.CreatedAt
		}

		if req.Sort.Direction == database.AscSortDirection {
			orderBys = append(orderBys, column.ASC())
		} else {
			orderBys = append(orderBys, column.DESC())
		}
	} else {
		orderBys = append(orderBys, tQualiRequests.CreatedAt.DESC())
	}

	stmt := tQualiRequests.
		SELECT(
			tQualiRequests.CreatedAt,
			tQualiRequests.QualificationID,
			tQuali.ID,
			tQuali.CreatedAt,
			tQuali.UpdatedAt,
			tQuali.Job,
			tQuali.Closed,
			tQuali.Abbreviation,
			tQuali.Title,
			tQuali.Description,
			tQualiRequests.UserID,
			tUser.ID,
			tUser.Job,
			tUser.JobGrade,
			tUser.Firstname,
			tUser.Lastname,
			tUser.Dateofbirth,
			tUser.PhoneNumber,
			tQualiRequests.UserComment,
			tQualiRequests.Status,
			tQualiRequests.ApprovedAt,
			tQualiRequests.ApproverComment,
			tQualiRequests.ApproverID,
			tApprover.ID,
			tApprover.Job,
			tApprover.JobGrade,
			tApprover.Firstname,
			tApprover.Lastname,
			tApprover.Dateofbirth,
			tApprover.PhoneNumber,
			tQualiRequests.ApproverJob,
		).
		FROM(
			tQualiRequests.
				INNER_JOIN(tQuali,
					tQuali.ID.EQ(tQualiRequests.QualificationID),
				).
				LEFT_JOIN(tUser,
					tQualiRequests.UserID.EQ(tUser.ID),
				).
				LEFT_JOIN(tApprover,
					tQualiRequests.ApproverID.EQ(tApprover.ID),
				).
				LEFT_JOIN(tQJobAccess,
					tQJobAccess.QualificationID.EQ(tQuali.ID).
						AND(tQJobAccess.Job.EQ(jet.String(userInfo.Job))).
						AND(tQJobAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
				),
		).
		GROUP_BY(tQualiRequests.QualificationID, tQualiRequests.UserID).
		ORDER_BY(orderBys...).
		WHERE(condition).
		OFFSET(req.Pagination.Offset).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Requests); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := 0; i < len(resp.Requests); i++ {
		if resp.Requests[i].User != nil {
			jobInfoFn(resp.Requests[i].User)
		}

		if resp.Requests[i].Approver != nil {
			jobInfoFn(resp.Requests[i].Approver)
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

	canGrade, err := s.access.CanUserAccessTarget(ctx, req.Request.QualificationId, userInfo, qualifications.AccessLevel_ACCESS_LEVEL_GRADE)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	quali, err := s.getQualification(ctx, req.Request.QualificationId, tQuali.ID.EQ(jet.Uint64(req.Request.QualificationId)), userInfo, false)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	// If the qualification is closed and user is not a grade tutor
	if !canGrade && quali.Closed {
		return nil, errorsqualifications.ErrQualificationClosed
	}

	// If user can grade a qualification, they are treated as an "approver" of requests
	if canGrade && req.Request.UserId > 0 {
		stmt := tQualiRequests.
			UPDATE(
				tQualiRequests.Status,
				tQualiRequests.ApprovedAt,
				tQualiRequests.ApproverComment,
				tQualiRequests.ApproverID,
				tQualiRequests.ApproverJob,
			).
			SET(
				req.Request.Status,
				jet.CURRENT_TIMESTAMP(),
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

		request, err := s.getQualificationRequest(ctx, req.Request.QualificationId, req.Request.UserId, userInfo)
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}

		// Only send notification when it wasn't already in the same status
		if request == nil || request.Status == nil || request.Status.Enum() != req.Request.Status.Enum() {
			if err := s.notif.NotifyUser(ctx, &notifications.Notification{
				UserId: request.UserId,
				Title: &common.TranslateItem{
					Key: "notifications.qualifications.request_updated.title",
				},
				Content: &common.TranslateItem{
					Key:        "notifications.qualifications.request_updated.content",
					Parameters: map[string]string{"abbreviation": quali.Abbreviation, "title": quali.Title},
				},
				Category: notifications.NotificationCategory_NOTIFICATION_CATEGORY_GENERAL,
				Type:     notifications.NotificationType_NOTIFICATION_TYPE_INFO,
				Data: &notifications.Data{
					Link: &notifications.Link{
						To: fmt.Sprintf("/qualifications/%d", request.QualificationId),
					},
				},
			}); err != nil {
				return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
			}
		}

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)
	} else {
		canRequest, err := s.access.CanUserAccessTarget(ctx, req.Request.QualificationId, userInfo, qualifications.AccessLevel_ACCESS_LEVEL_REQUEST)
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
		if !canRequest {
			return nil, errorsqualifications.ErrFailedQuery
		}

		// Make sure the requirements of the qualification are fullfiled by the user, ErrRequirementsMissing
		reqsFullfilled, err := s.checkRequirementsMetForQualification(ctx, req.Request.QualificationId, userInfo.UserId)
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
		if !reqsFullfilled {
			return nil, errorsqualifications.ErrRequirementsMissing
		}

		request, err := s.getQualificationRequest(ctx, req.Request.QualificationId, userInfo.UserId, userInfo)
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}

		if request != nil &&
			(request.Status == nil || (*request.Status != qualifications.RequestStatus_REQUEST_STATUS_PENDING &&
				*request.Status != qualifications.RequestStatus_REQUEST_STATUS_COMPLETED)) {
			return nil, errorsqualifications.ErrFailedQuery
		}

		tQualiRequests := table.FivenetQualificationsRequests
		stmt := tQualiRequests.
			INSERT(
				tQualiRequests.QualificationID,
				tQualiRequests.UserID,
				tQualiRequests.UserComment,
				tQualiRequests.Status,
			).
			VALUES(
				req.Request.QualificationId,
				userInfo.UserId,
				req.Request.UserComment,
				qualifications.RequestStatus_REQUEST_STATUS_PENDING,
			).
			ON_DUPLICATE_KEY_UPDATE(
				tQualiRequests.DeletedAt.SET(jet.TimestampExp(jet.NULL)),
				tQualiRequests.UserComment.SET(jet.StringExp(jet.Raw("VALUES(`user_comment`)"))),
				tQualiRequests.Status.SET(jet.Int16(int16(qualifications.RequestStatus_REQUEST_STATUS_PENDING))),
				tQualiRequests.ApprovedAt.SET(jet.DateTimeExp(jet.NULL)),
				tQualiRequests.ApproverComment.SET(jet.StringExp(jet.NULL)),
				tQualiRequests.ApproverID.SET(jet.IntExp(jet.NULL)),
				tQualiRequests.ApproverJob.SET(jet.StringExp(jet.NULL)),
			)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)
	}

	request, err := s.getQualificationRequest(ctx, req.Request.QualificationId, userInfo.UserId, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	return &CreateOrUpdateQualificationRequestResponse{
		Request: request,
	}, nil
}

func (s *Server) getQualificationRequest(ctx context.Context, qualificationId uint64, userId int32, userInfo *userinfo.UserInfo) (*qualifications.QualificationRequest, error) {
	tQuali := tQuali.AS("qualificationshort")
	tUser := tables.Users().AS("user")
	tApprover := tUser.AS("approver")

	stmt := tQualiRequests.
		SELECT(
			tQualiRequests.CreatedAt,
			tQualiRequests.DeletedAt,
			tQualiRequests.QualificationID,
			tQuali.ID,
			tQuali.CreatedAt,
			tQuali.UpdatedAt,
			tQuali.Job,
			tQuali.Closed,
			tQuali.Abbreviation,
			tQuali.Title,
			tQuali.Description,
			tQualiRequests.UserID,
			tUser.ID,
			tUser.Job,
			tUser.JobGrade,
			tUser.Firstname,
			tUser.Lastname,
			tUser.Dateofbirth,
			tUser.PhoneNumber,
			tQualiRequests.UserComment,
			tQualiRequests.Status,
			tQualiRequests.ApprovedAt,
			tQualiRequests.ApproverComment,
			tQualiRequests.ApproverID,
			tQualiRequests.ApproverJob,
			tApprover.ID,
			tApprover.Job,
			tApprover.JobGrade,
			tApprover.Firstname,
			tApprover.Lastname,
			tApprover.Dateofbirth,
			tApprover.PhoneNumber,
		).
		FROM(tQualiRequests.
			INNER_JOIN(tQuali,
				tQuali.ID.EQ(tQualiRequests.QualificationID),
			).
			LEFT_JOIN(tUser,
				tUser.ID.EQ(tQualiRequests.UserID),
			).
			LEFT_JOIN(tApprover,
				tApprover.ID.EQ(tQualiRequests.ApproverID),
			),
		).
		GROUP_BY(tQualiRequests.QualificationID, tQualiRequests.UserID).
		ORDER_BY(tQualiRequests.CreatedAt.DESC()).
		WHERE(jet.AND(
			tQualiRequests.QualificationID.EQ(jet.Uint64(qualificationId)),
			tQualiRequests.UserID.EQ(jet.Int32(userId)),
			tQualiRequests.DeletedAt.IS_NULL(),
		)).
		LIMIT(1)

	var request qualifications.QualificationRequest
	if err := stmt.QueryContext(ctx, s.db, &request); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if request.QualificationId == 0 {
		return nil, nil
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
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	if re == nil {
		return &DeleteQualificationReqResponse{}, nil
	}

	check, err := s.access.CanUserAccessTarget(ctx, re.QualificationId, userInfo, qualifications.AccessLevel_ACCESS_LEVEL_MANAGE)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !check {
		return nil, errorsqualifications.ErrFailedQuery
	}

	if err := s.deleteQualificationRequest(ctx, s.db, re.QualificationId, re.UserId); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)

	return &DeleteQualificationReqResponse{}, nil
}

func (s *Server) deleteQualificationRequest(ctx context.Context, tx qrm.DB, qualificationId uint64, userId int32) error {
	stmt := tQualiRequests.
		UPDATE(
			tQualiRequests.DeletedAt,
		).
		SET(
			jet.CURRENT_TIMESTAMP(),
		).
		WHERE(jet.AND(
			tQualiRequests.QualificationID.EQ(jet.Uint64(qualificationId)),
			tQualiRequests.UserID.EQ(jet.Int32(userId)),
		))

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	if err := s.deleteExamUser(ctx, tx, qualificationId, userId); err != nil {
		return err
	}

	return nil
}

func (s *Server) updateRequestStatus(ctx context.Context, qualificationId uint64, userId int32, status qualifications.RequestStatus) error {
	tQualiRequests := table.FivenetQualificationsRequests
	stmt := tQualiRequests.
		INSERT(
			tQualiRequests.QualificationID,
			tQualiRequests.UserID,
			tQualiRequests.Status,
		).
		VALUES(
			qualificationId,
			userId,
			status,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tQualiRequests.Status.SET(jet.Int16(int16(status))),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return err
	}

	return nil
}
