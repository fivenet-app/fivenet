package qualifications

import (
	"context"
	"errors"
	"fmt"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common"
	database "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/notifications"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/qualifications"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/userinfo"
	pbqualifications "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/qualifications"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsqualifications "github.com/fivenet-app/fivenet/v2025/services/qualifications/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	logging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

var tQualiRequests = table.FivenetQualificationsRequests.AS("qualification_request")

func (s *Server) ListQualificationRequests(
	ctx context.Context,
	req *pbqualifications.ListQualificationRequestsRequest,
) (*pbqualifications.ListQualificationRequestsResponse, error) {
	if req.QualificationId != nil {
		logging.InjectFields(
			ctx,
			logging.Fields{"fivenet.qualifications.id", req.GetQualificationId()},
		)
	}
	if req.UserId != nil {
		logging.InjectFields(ctx, logging.Fields{"fivenet.qualifications.user_id", req.GetUserId()})
	}

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	tQuali := tQuali.AS("qualification_short")
	tUser := tables.User().AS("user")
	tApprover := tUser.AS("approver")

	condition := tQualiRequests.DeletedAt.IS_NULL().
		AND(tQualiRequests.Status.NOT_EQ(mysql.Int32(int32(qualifications.RequestStatus_REQUEST_STATUS_COMPLETED))))

	if req.QualificationId != nil {
		check, err := s.access.CanUserAccessTarget(
			ctx,
			req.GetQualificationId(),
			userInfo,
			qualifications.AccessLevel_ACCESS_LEVEL_GRADE,
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
		if !check {
			return nil, errorsqualifications.ErrFailedQuery
		}

		condition = condition.AND(
			tQualiRequests.QualificationID.EQ(mysql.Int64(req.GetQualificationId())),
		)
	} else {
		accessExists := mysql.EXISTS(
			mysql.SELECT(mysql.Int(1)).
				FROM(tQAccess).
				WHERE(mysql.AND(
					tQAccess.TargetID.EQ(tQualiRequests.QualificationID),
					tQAccess.Job.EQ(mysql.String(userInfo.GetJob())),
					tQAccess.MinimumGrade.LT_EQ(mysql.Int32(userInfo.GetJobGrade())),

					mysql.OR(
						tQAccess.Access.GT_EQ(mysql.Int32(int32(qualifications.AccessLevel_ACCESS_LEVEL_GRADE))),
						mysql.OR(
							tQAccess.Access.GT_EQ(mysql.Int32(int32(qualifications.AccessLevel_ACCESS_LEVEL_VIEW))).
								AND(tQualiRequests.UserID.EQ(mysql.Int32(userInfo.GetUserId()))),
						),
					),
				),
				),
		)

		condition = condition.AND(mysql.AND(
			tQuali.DeletedAt.IS_NULL(),
			mysql.OR(
				tQualiRequests.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
				accessExists,
			),
		))
	}

	countColumn := mysql.Expression(tQualiRequests.QualificationID)
	if req.UserId != nil {
		condition = condition.AND(tUser.Job.EQ(mysql.String(userInfo.GetJob()))).
			AND(tQualiRequests.UserID.EQ(mysql.Int32(req.GetUserId())))
	} else {
		if req.QualificationId == nil {
			condition = condition.AND(tUser.Job.EQ(mysql.String(userInfo.GetJob()))).AND(tQualiRequests.UserID.EQ(mysql.Int32(userInfo.GetUserId())))
			countColumn = mysql.DISTINCT(tQualiRequests.QualificationID)
		} else {
			countColumn = mysql.DISTINCT(tQualiRequests.UserID)
		}
	}

	if len(req.GetStatus()) > 0 {
		statuses := []mysql.Expression{}
		for i := range req.GetStatus() {
			statuses = append(statuses, mysql.Int32(int32(req.GetStatus()[i])))
		}

		condition = condition.AND(tQualiRequests.Status.IN(statuses...))
	} else {
		condition = condition.AND(tQualiRequests.Status.NOT_EQ(mysql.Int32(int32(qualifications.RequestStatus_REQUEST_STATUS_COMPLETED))))
	}

	countStmt := tQualiRequests.
		SELECT(
			mysql.COUNT(countColumn).AS("data_count.total"),
		).
		FROM(
			tQualiRequests.
				INNER_JOIN(tQuali,
					tQuali.ID.EQ(tQualiRequests.QualificationID),
				).
				LEFT_JOIN(tQAccess,
					tQAccess.TargetID.EQ(tQuali.ID).
						AND(tQAccess.Job.EQ(mysql.String(userInfo.GetJob()))).
						AND(tQAccess.MinimumGrade.LT_EQ(mysql.Int32(userInfo.GetJobGrade()))),
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

	pag, limit := req.GetPagination().GetResponseWithPageSize(count.Total, QualificationsPageSize)
	resp := &pbqualifications.ListQualificationRequestsResponse{
		Pagination: pag,
		Requests:   []*qualifications.QualificationRequest{},
	}
	if count.Total <= 0 {
		return resp, nil
	}

	// Convert proto sort to db sorting
	orderBys := []mysql.OrderByClause{}
	if req.GetSort() != nil {
		var column mysql.Column
		switch req.GetSort().GetColumn() {
		case "status":
			column = tQualiRequests.Status
		case "approvedAt":
			column = tQualiRequests.ApprovedAt
		case "createdAt":
			fallthrough
		default:
			column = tQualiRequests.CreatedAt
		}

		if req.GetSort().GetDirection() == database.AscSortDirection {
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
			tQuali.Draft,
			tQuali.Public,
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
				LEFT_JOIN(tQAccess,
					tQAccess.TargetID.EQ(tQuali.ID).
						AND(tQAccess.Job.EQ(mysql.String(userInfo.GetJob()))).
						AND(tQAccess.MinimumGrade.LT_EQ(mysql.Int32(userInfo.GetJobGrade()))),
				),
		).
		GROUP_BY(tQualiRequests.QualificationID, tQualiRequests.UserID).
		ORDER_BY(orderBys...).
		WHERE(condition).
		OFFSET(req.GetPagination().GetOffset()).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Requests); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.GetRequests() {
		if resp.GetRequests()[i].GetUser() != nil {
			jobInfoFn(resp.GetRequests()[i].GetUser())
		}

		if resp.GetRequests()[i].GetApprover() != nil {
			jobInfoFn(resp.GetRequests()[i].GetApprover())
		}
	}

	return resp, nil
}

func (s *Server) CreateOrUpdateQualificationRequest(
	ctx context.Context,
	req *pbqualifications.CreateOrUpdateQualificationRequestRequest,
) (*pbqualifications.CreateOrUpdateQualificationRequestResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbqualifications.QualificationsService_ServiceDesc.ServiceName,
		Method:  "CreateOrUpdateQualificationRequest",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	canGrade, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetRequest().GetQualificationId(),
		userInfo,
		qualifications.AccessLevel_ACCESS_LEVEL_GRADE,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	quali, err := s.getQualification(
		ctx,
		req.GetRequest().GetQualificationId(),
		tQuali.ID.EQ(mysql.Int64(req.GetRequest().GetQualificationId())),
		userInfo,
		false,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	// If the qualification is closed and user is not a grade tutor
	if !canGrade && quali.GetClosed() {
		return nil, errorsqualifications.ErrQualificationClosed
	}

	// If user can grade a qualification, they are treated as an "approver" of requests
	if canGrade && req.GetRequest().GetUserId() > 0 {
		stmt := tQualiRequests.
			UPDATE(
				tQualiRequests.Status,
				tQualiRequests.ApprovedAt,
				tQualiRequests.ApproverComment,
				tQualiRequests.ApproverID,
				tQualiRequests.ApproverJob,
			).
			SET(
				req.GetRequest().GetStatus(),
				mysql.CURRENT_TIMESTAMP(),
				req.GetRequest().ApproverComment,
				userInfo.GetUserId(),
				userInfo.GetJob(),
			).
			WHERE(mysql.AND(
				tQualiRequests.QualificationID.EQ(
					mysql.Int64(req.GetRequest().GetQualificationId()),
				),
				tQualiRequests.UserID.EQ(mysql.Int32(req.GetRequest().GetUserId())),
			))

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}

		request, err := s.getQualificationRequest(
			ctx,
			req.GetRequest().GetQualificationId(),
			req.GetRequest().GetUserId(),
			userInfo,
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}

		// Only send notification when it wasn't already in the same status
		if request == nil || request.Status == nil ||
			request.GetStatus().Enum() != req.GetRequest().GetStatus().Enum() {
			if err := s.notif.NotifyUser(ctx, &notifications.Notification{
				UserId: request.GetUserId(),
				Title: &common.I18NItem{
					Key: "notifications.qualifications.request_updated.title",
				},
				Content: &common.I18NItem{
					Key:        "notifications.qualifications.request_updated.content",
					Parameters: map[string]string{"abbreviation": quali.GetAbbreviation(), "title": quali.GetTitle()},
				},
				Category: notifications.NotificationCategory_NOTIFICATION_CATEGORY_GENERAL,
				Type:     notifications.NotificationType_NOTIFICATION_TYPE_INFO,
				Data: &notifications.Data{
					Link: &notifications.Link{
						To: fmt.Sprintf("/qualifications/%d", request.GetQualificationId()),
					},
				},
			}); err != nil {
				return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
			}
		}

		auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED
	} else {
		canRequest, err := s.access.CanUserAccessTarget(ctx, req.GetRequest().GetQualificationId(), userInfo, qualifications.AccessLevel_ACCESS_LEVEL_REQUEST)
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
		if !canRequest {
			return nil, errorsqualifications.ErrFailedQuery
		}

		// Make sure the requirements of the qualification are fullfiled by the user, ErrRequirementsMissing
		reqsFullfilled, err := s.checkRequirementsMetForQualification(ctx, req.GetRequest().GetQualificationId(), userInfo.GetUserId())
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
		if !reqsFullfilled {
			return nil, errorsqualifications.ErrRequirementsMissing
		}

		request, err := s.getQualificationRequest(ctx, req.GetRequest().GetQualificationId(), userInfo.GetUserId(), userInfo)
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}

		if request != nil &&
			(request.Status == nil || (request.GetStatus() != qualifications.RequestStatus_REQUEST_STATUS_PENDING &&
				request.GetStatus() != qualifications.RequestStatus_REQUEST_STATUS_COMPLETED)) {
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
				req.GetRequest().GetQualificationId(),
				userInfo.GetUserId(),
				req.GetRequest().UserComment,
				qualifications.RequestStatus_REQUEST_STATUS_PENDING,
			).
			ON_DUPLICATE_KEY_UPDATE(
				tQualiRequests.DeletedAt.SET(mysql.TimestampExp(mysql.NULL)),
				tQualiRequests.UserComment.SET(mysql.StringExp(mysql.Raw("VALUES(`user_comment`)"))),
				tQualiRequests.Status.SET(mysql.Int32(int32(qualifications.RequestStatus_REQUEST_STATUS_PENDING))),
				tQualiRequests.ApprovedAt.SET(mysql.DateTimeExp(mysql.NULL)),
				tQualiRequests.ApproverComment.SET(mysql.StringExp(mysql.NULL)),
				tQualiRequests.ApproverID.SET(mysql.IntExp(mysql.NULL)),
				tQualiRequests.ApproverJob.SET(mysql.StringExp(mysql.NULL)),
			)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}

		auditEntry.State = audit.EventType_EVENT_TYPE_CREATED
	}

	request, err := s.getQualificationRequest(
		ctx,
		req.GetRequest().GetQualificationId(),
		userInfo.GetUserId(),
		userInfo,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	return &pbqualifications.CreateOrUpdateQualificationRequestResponse{
		Request: request,
	}, nil
}

func (s *Server) getQualificationRequest(
	ctx context.Context,
	qualificationId int64,
	userId int32,
	userInfo *userinfo.UserInfo,
) (*qualifications.QualificationRequest, error) {
	tQuali := tQuali.AS("qualificationshort")
	tUser := tables.User().AS("user")
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
			tQuali.Draft,
			tQuali.Public,
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
		WHERE(mysql.AND(
			tQualiRequests.QualificationID.EQ(mysql.Int64(qualificationId)),
			tQualiRequests.UserID.EQ(mysql.Int32(userId)),
			tQualiRequests.DeletedAt.IS_NULL(),
		)).
		LIMIT(1)

	var request qualifications.QualificationRequest
	if err := stmt.QueryContext(ctx, s.db, &request); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if request.GetQualificationId() == 0 {
		return nil, nil
	}

	if request.GetUser() != nil {
		s.enricher.EnrichJobInfoSafe(userInfo, request.GetUser())
	}

	if request.GetApprover() != nil {
		s.enricher.EnrichJobInfoSafe(userInfo, request.GetApprover())
	}

	return &request, nil
}

func (s *Server) DeleteQualificationReq(
	ctx context.Context,
	req *pbqualifications.DeleteQualificationReqRequest,
) (*pbqualifications.DeleteQualificationReqResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbqualifications.QualificationsService_ServiceDesc.ServiceName,
		Method:  "DeleteQualificationReq",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	re, err := s.getQualificationRequest(ctx, req.GetQualificationId(), req.GetUserId(), userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	if re == nil {
		return &pbqualifications.DeleteQualificationReqResponse{}, nil
	}

	check, err := s.access.CanUserAccessTarget(
		ctx,
		re.GetQualificationId(),
		userInfo,
		qualifications.AccessLevel_ACCESS_LEVEL_EDIT,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !check {
		return nil, errorsqualifications.ErrFailedQuery
	}

	if err := s.deleteQualificationRequest(ctx, s.db, re.GetQualificationId(), re.GetUserId()); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_DELETED

	return &pbqualifications.DeleteQualificationReqResponse{}, nil
}

func (s *Server) deleteQualificationRequest(
	ctx context.Context,
	tx qrm.DB,
	qualificationId int64,
	userId int32,
) error {
	stmt := tQualiRequests.
		UPDATE(
			tQualiRequests.DeletedAt,
		).
		SET(
			mysql.CURRENT_TIMESTAMP(),
		).
		WHERE(mysql.AND(
			tQualiRequests.QualificationID.EQ(mysql.Int64(qualificationId)),
			tQualiRequests.UserID.EQ(mysql.Int32(userId)),
		))

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	if err := s.deleteExamUser(ctx, tx, qualificationId, userId); err != nil {
		return err
	}

	return nil
}

func (s *Server) updateRequestStatus(
	ctx context.Context,
	tx qrm.DB,
	qualificationId int64,
	userId int32,
	status qualifications.RequestStatus,
) error {
	tQualiRequests := table.FivenetQualificationsRequests
	stmt := tQualiRequests.
		INSERT(
			tQualiResults.DeletedAt,
			tQualiRequests.QualificationID,
			tQualiRequests.UserID,
			tQualiRequests.Status,
		).
		VALUES(
			mysql.NULL,
			qualificationId,
			userId,
			status,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tQualiRequests.Status.SET(mysql.Int32(int32(status))),
		)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	return nil
}
