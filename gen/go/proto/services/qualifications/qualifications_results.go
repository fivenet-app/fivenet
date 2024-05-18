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

var (
	tQualiResults = table.FivenetQualificationsResults.AS("qualificationresult")
)

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

	if req.UserId != nil {
		condition = condition.AND(tUser.Job.EQ(jet.String(userInfo.Job))).AND(tQualiResults.UserID.EQ(jet.Int32(*req.UserId)))
	} else if req.QualificationId == nil {
		condition = condition.AND(tUser.Job.EQ(jet.String(userInfo.Job))).AND(tQualiResults.UserID.EQ(jet.Int32(userInfo.UserId)))
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
			jet.COUNT(tQualiResults.ID).AS("datacount.totalcount"),
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

	stmt := tQualiResults.
		SELECT(
			tQualiResults.ID,
			tQualiResults.CreatedAt,
			tQualiResults.QualificationID,
			tQualiResults.UserID,
			tUser.ID,
			tUser.Identifier,
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
			tCreator.Identifier,
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
		GROUP_BY(tQualiResults.ID).
		ORDER_BY(tQualiResults.CreatedAt.DESC()).
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

	tQualiResults := table.FivenetQualificationsResults
	if req.Result.Id <= 0 {
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
		result, err := s.getQualificationResult(ctx, req.Result.QualificationId, req.Result.Id, userInfo)
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

	quali, err := s.getQualification(ctx, req.Result.QualificationId, nil, userInfo, false)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	result, err := s.getQualificationResult(ctx, quali.Id, req.Result.Id, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	// Only send notification when there is currently no score
	if result.Score == nil {
		if err := s.notif.NotifyUser(ctx, &notifications.Notification{
			UserId: result.UserId,
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
					To: fmt.Sprintf("/qualifications/%d", result.QualificationId),
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
	} else if req.Result.Status == qualifications.ResultStatus_RESULT_STATUS_FAILED {
		// If failed status, delete the request
		if err := s.deleteQualificationRequest(ctx, req.Result.QualificationId, req.Result.UserId); err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
	}

	return &CreateOrUpdateQualificationResultResponse{
		Result: result,
	}, nil
}

func (s *Server) getQualificationResult(ctx context.Context, qualificationId uint64, resultId uint64, userInfo *userinfo.UserInfo) (*qualifications.QualificationResult, error) {
	tUser := tUser.AS("user")

	condition := tQualiResults.DeletedAt.IS_NULL()

	if resultId > 0 {
		condition = condition.AND(tQualiResults.ID.EQ(jet.Uint64(resultId)))
	} else {
		condition = condition.AND(tQualiResults.UserID.EQ(jet.Int32(userInfo.UserId)))
	}
	if qualificationId > 0 {
		condition = condition.AND(tQualiResults.QualificationID.EQ(jet.Uint64(qualificationId)))
	}

	stmt := tQualiResults.
		SELECT(
			tQualiResults.ID,
			tQualiResults.CreatedAt,
			tQualiResults.DeletedAt,
			tQualiResults.QualificationID,
			tQualiResults.UserID,
			tUser.ID,
			tUser.Identifier,
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
			tCreator.Identifier,
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

	result, err := s.getQualificationResult(ctx, 0, req.ResultId, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
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
		WHERE(
			tQualiResults.ID.EQ(jet.Uint64(result.Id)),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)

	return &DeleteQualificationResultResponse{}, nil
}
