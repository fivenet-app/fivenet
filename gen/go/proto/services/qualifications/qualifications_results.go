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
	"github.com/fivenet-app/fivenet/query/fivenet/model"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var tQualiResults = table.FivenetQualificationsResults.AS("qualificationresult")

func (s *Server) ListQualificationsResults(ctx context.Context, req *ListQualificationsResultsRequest) (*ListQualificationsResultsResponse, error) {
	if req.QualificationId != nil {
		trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.qualifications.id", int64(*req.QualificationId)))
	}
	if req.UserId != nil {
		trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.qualifications.user_id", int64(*req.UserId)))
	}

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	tQuali := tQuali.AS("qualificationshort")

	condition := tQualiResults.DeletedAt.IS_NULL()

	if req.QualificationId != nil {
		check, err := s.checkIfUserHasAccessToQuali(ctx, *req.QualificationId, userInfo, qualifications.AccessLevel_ACCESS_LEVEL_GRADE)
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
		if !check {
			return nil, errorsqualifications.ErrFailedQuery
		}

		condition = condition.AND(tQualiResults.QualificationID.EQ(jet.Uint64(*req.QualificationId)))
	} else {
		condition = condition.AND(jet.AND(
			tQuali.DeletedAt.IS_NULL(),
			jet.OR(
				jet.AND(
					tQualiResults.CreatorID.EQ(jet.Int32(userInfo.UserId)),
					tQualiResults.CreatorJob.EQ(jet.String(userInfo.Job)),
				),
				jet.AND(
					tQJobAccess.Access.IS_NOT_NULL(),
					jet.OR(
						tQJobAccess.Access.GT_EQ(jet.Int32(int32(qualifications.AccessLevel_ACCESS_LEVEL_GRADE))),
						jet.AND(
							tQJobAccess.Access.GT(jet.Int32(int32(qualifications.AccessLevel_ACCESS_LEVEL_BLOCKED))),
							tQualiResults.UserID.EQ(jet.Int32(userInfo.UserId)),
						),
					),
				),
			),
		))
	}

	countColumn := jet.Expression(tQualiResults.QualificationID)
	if req.UserId != nil {
		condition = condition.AND(tUser.Job.EQ(jet.String(userInfo.Job))).AND(tQualiResults.UserID.EQ(jet.Int32(*req.UserId)))
	} else {
		if req.QualificationId == nil {
			condition = condition.AND(tUser.Job.EQ(jet.String(userInfo.Job))).AND(tQualiResults.UserID.EQ(jet.Int32(userInfo.UserId)))
			countColumn = jet.DISTINCT(tQualiResults.QualificationID)
		} else {
			countColumn = jet.DISTINCT(tQualiResults.UserID)
		}
	}

	if len(req.Status) > 0 {
		statuses := []jet.Expression{}
		for i := 0; i < len(req.Status); i++ {
			statuses = append(statuses, jet.Int16(int16(req.Status[i])))
		}

		condition = condition.AND(tQualiResults.Status.IN(statuses...))
	}

	countStmt := tQualiResults.
		SELECT(
			jet.COUNT(countColumn).AS("datacount.totalcount"),
		).
		FROM(
			tQualiResults.
				INNER_JOIN(tQuali,
					tQuali.ID.EQ(tQualiResults.QualificationID),
				).
				LEFT_JOIN(tQJobAccess,
					tQJobAccess.QualificationID.EQ(tQuali.ID).
						AND(tQJobAccess.Job.EQ(jet.String(userInfo.Job))).
						AND(tQJobAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
				).
				LEFT_JOIN(tUser,
					tQualiResults.UserID.EQ(tUser.ID),
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
	resp := &ListQualificationsResultsResponse{
		Pagination: pag,
		Results:    []*qualifications.QualificationResult{},
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
			column = tQualiResults.Status
		case "createdAt":
			fallthrough
		default:
			column = tQualiResults.CreatedAt
		}

		if req.Sort.Direction == database.AscSortDirection {
			orderBys = append(orderBys, column.ASC())
		} else {
			orderBys = append(orderBys, column.DESC())
		}
	} else {
		orderBys = append(orderBys, tQualiResults.CreatedAt.DESC())
	}

	stmt := tQualiResults.
		SELECT(
			tQualiResults.ID,
			tQualiResults.CreatedAt,
			tQualiResults.QualificationID,
			tQualiResults.UserID,
			tUser.ID,
			tUser.Job,
			tUser.JobGrade,
			tUser.Firstname,
			tUser.Lastname,
			tUser.Dateofbirth,
			tUser.PhoneNumber,
			tQualiResults.Status,
			tQualiResults.Score,
			tQualiResults.Summary,
			tQualiResults.CreatorID,
			tCreator.ID,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Dateofbirth,
			tCreator.PhoneNumber,
			tQuali.ID,
			tQuali.CreatedAt,
			tQuali.UpdatedAt,
			tQuali.Job,
			tQuali.Closed,
			tQuali.Abbreviation,
			tQuali.Title,
			tQuali.Description,
			tQuali.CreatorJob,
			tQuali.CreatorID,
		).
		FROM(
			tQualiResults.
				INNER_JOIN(tQuali,
					tQuali.ID.EQ(tQualiResults.QualificationID),
				).
				LEFT_JOIN(tUser,
					tQualiResults.UserID.EQ(tUser.ID),
				).
				LEFT_JOIN(tCreator,
					tQualiResults.CreatorID.EQ(tCreator.ID),
				).
				LEFT_JOIN(tQJobAccess,
					tQJobAccess.QualificationID.EQ(tQuali.ID).
						AND(tQJobAccess.Job.EQ(jet.String(userInfo.Job))).
						AND(tQJobAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
				),
		).
		GROUP_BY(tQualiResults.Status, tQualiResults.CreatedAt).
		ORDER_BY(orderBys...).
		WHERE(condition).
		OFFSET(req.Pagination.Offset).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Results); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := 0; i < len(resp.Results); i++ {
		if resp.Results[i].User != nil {
			jobInfoFn(resp.Results[i].User)
		}

		if resp.Results[i].Creator != nil {
			jobInfoFn(resp.Results[i].Creator)
		}
	}

	return resp, nil
}

func (s *Server) CreateOrUpdateQualificationResult(ctx context.Context, req *CreateOrUpdateQualificationResultRequest) (*CreateOrUpdateQualificationResultResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: QualificationsService_ServiceDesc.ServiceName,
		Method:  "CreateOrUpdateQualificationResult",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.checkIfUserHasAccessToQuali(ctx, req.Result.QualificationId, userInfo, qualifications.AccessLevel_ACCESS_LEVEL_GRADE)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !check {
		return nil, errorsqualifications.ErrFailedQuery
	}

	result, err := s.getQualificationResult(ctx, req.Result.QualificationId, req.Result.Id, []qualifications.ResultStatus{qualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL}, userInfo, req.Result.UserId)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	tQualiResults := table.FivenetQualificationsResults
	// There is currently no result with status successful
	if req.Result.Id <= 0 && (result == nil || (result.Status != qualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL && req.Result.Status != qualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL)) {
		stmt := tQualiResults.
			INSERT(
				tQualiResults.QualificationID,
				tQualiResults.UserID,
				tQualiResults.Status,
				tQualiResults.Score,
				tQualiResults.Summary,
				tQualiResults.CreatorID,
				tQualiResults.CreatorJob,
			).
			VALUES(
				req.Result.QualificationId,
				req.Result.UserId,
				req.Result.Status,
				req.Result.Score,
				req.Result.Summary,
				userInfo.UserId,
				userInfo.Job,
			)

		res, err := stmt.ExecContext(ctx, s.db)
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}

		lastId, err := res.LastInsertId()
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}

		req.Result.Id = uint64(lastId)

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)
	} else {
		result, err := s.getQualificationResult(ctx, req.Result.QualificationId, req.Result.Id, nil, userInfo, req.Result.UserId)
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}

		req.Result.UserId = result.UserId

		stmt := tQualiResults.
			UPDATE(
				tQualiResults.QualificationID,
				tQualiResults.UserID,
				tQualiResults.Status,
				tQualiResults.Score,
				tQualiResults.Summary,
			).
			SET(
				req.Result.QualificationId,
				req.Result.UserId,
				req.Result.Status,
				req.Result.Score,
				req.Result.Summary,
			).
			WHERE(jet.AND(
				tQualiResults.ID.EQ(jet.Uint64(req.Result.Id)),
				tQualiResults.DeletedAt.IS_NULL(),
			))

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)
	}

	quali, err := s.getQualification(ctx, req.Result.QualificationId, tQuali.ID.EQ(jet.Uint64(req.Result.QualificationId)), userInfo, false)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	// Only send notification when the original result had no score and isn't in pending status
	if req.Result.Status != qualifications.ResultStatus_RESULT_STATUS_PENDING && (result == nil || (result.Status == qualifications.ResultStatus_RESULT_STATUS_PENDING || (result.Score == nil && req.Result.Score != nil))) {
		if err := s.notif.NotifyUser(ctx, &notifications.Notification{
			UserId: req.Result.UserId,
			Title: &common.TranslateItem{
				Key: "notifications.qualifications.result_updated.title",
			},
			Content: &common.TranslateItem{
				Key:        "notifications.qualifications.result_updated.content",
				Parameters: map[string]string{"abbreviation": quali.Abbreviation, "title": quali.Title},
			},
			Category: notifications.NotificationCategory_NOTIFICATION_CATEGORY_GENERAL,
			Type:     notifications.NotificationType_NOTIFICATION_TYPE_INFO,
			Data: &notifications.Data{
				Link: &notifications.Link{
					To: fmt.Sprintf("/qualifications/%d", req.Result.QualificationId),
				},
			},
		}); err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
	}

	// If the result is successful, complete the request status
	if req.Result.Status == qualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL {
		if err := s.updateRequestStatus(ctx, req.Result.QualificationId, req.Result.UserId, qualifications.RequestStatus_REQUEST_STATUS_COMPLETED); err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
	} else {
		// If failed or other status, delete the request
		if err := s.deleteQualificationRequest(ctx, req.Result.QualificationId, req.Result.UserId); err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}

		if err := s.deleteExamUser(ctx, req.Result.QualificationId, req.Result.UserId); err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
	}

	result, err = s.getQualificationResult(ctx, quali.Id, req.Result.Id, nil, userInfo, req.Result.UserId)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	return &CreateOrUpdateQualificationResultResponse{
		Result: result,
	}, nil
}

func (s *Server) getQualificationResult(ctx context.Context, qualificationId uint64, resultId uint64, status []qualifications.ResultStatus, userInfo *userinfo.UserInfo, userId int32) (*qualifications.QualificationResult, error) {
	tUser := tUser.AS("user")

	condition := tQualiResults.DeletedAt.IS_NULL()

	if resultId > 0 {
		condition = condition.AND(tQualiResults.ID.EQ(jet.Uint64(resultId)))
	} else if userId > 0 {
		condition = condition.AND(tQualiResults.UserID.EQ(jet.Int32(userId)))
	} else {
		condition = condition.AND(tQualiResults.UserID.EQ(jet.Int32(userInfo.UserId)))
	}
	if qualificationId > 0 {
		condition = condition.AND(tQualiResults.QualificationID.EQ(jet.Uint64(qualificationId)))
	}

	if len(status) > 0 {
		statusConds := make([]jet.Expression, len(status))
		for i := 0; i < len(status); i++ {
			statusConds[i] = jet.Int16(int16(status[i]))
		}

		condition = condition.AND(tQualiResults.Status.IN(statusConds...))
	}

	stmt := tQualiResults.
		SELECT(
			tQualiResults.ID,
			tQualiResults.CreatedAt,
			tQualiResults.DeletedAt,
			tQualiResults.QualificationID,
			tQualiResults.UserID,
			tUser.ID,
			tUser.Job,
			tUser.JobGrade,
			tUser.Firstname,
			tUser.Lastname,
			tUser.Dateofbirth,
			tUser.PhoneNumber,
			tQualiResults.Status,
			tQualiResults.Score,
			tQualiResults.Summary,
			tQualiResults.CreatorID,
			tQualiResults.CreatorJob,
			tCreator.ID,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Dateofbirth,
			tCreator.PhoneNumber,
		).
		FROM(tQualiResults.
			LEFT_JOIN(tUser,
				tUser.ID.EQ(tQualiResults.UserID),
			).
			LEFT_JOIN(tCreator,
				tCreator.ID.EQ(tQualiResults.CreatorID),
			),
		).
		GROUP_BY(tQualiResults.ID).
		ORDER_BY(tQualiResults.ID.DESC()).
		WHERE(condition).
		LIMIT(1)

	var result qualifications.QualificationResult
	if err := stmt.QueryContext(ctx, s.db, &result); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if result.Id == 0 {
		return nil, nil
	}

	if result.User != nil {
		s.enricher.EnrichJobInfoSafe(userInfo, result.User)
	}

	if result.Creator != nil {
		s.enricher.EnrichJobInfoSafe(userInfo, result.Creator)
	}

	return &result, nil
}

func (s *Server) DeleteQualificationResult(ctx context.Context, req *DeleteQualificationResultRequest) (*DeleteQualificationResultResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: QualificationsService_ServiceDesc.ServiceName,
		Method:  "DeleteQualificationResult",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	result, err := s.getQualificationResult(ctx, 0, req.ResultId, nil, userInfo, 0)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	if result == nil {
		return &DeleteQualificationResultResponse{}, nil
	}

	check, err := s.checkIfUserHasAccessToQuali(ctx, result.QualificationId, userInfo, qualifications.AccessLevel_ACCESS_LEVEL_MANAGE)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !check {
		return nil, errorsqualifications.ErrFailedQuery
	}

	stmt := tQualiResults.
		UPDATE(
			tQualiResults.DeletedAt,
		).
		SET(
			jet.CURRENT_TIMESTAMP(),
		).
		WHERE(jet.AND(
			tQualiResults.ID.EQ(jet.Uint64(result.Id)),
			tQualiResults.ID.EQ(jet.Uint64(req.ResultId)),
		))

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	if err := s.deleteExamUser(ctx, result.QualificationId, result.UserId); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)

	return &DeleteQualificationResultResponse{}, nil
}
