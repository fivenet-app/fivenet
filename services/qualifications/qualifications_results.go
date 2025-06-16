package qualifications

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common"
	database "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/notifications"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/qualifications"
	pbqualifications "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/qualifications"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsqualifications "github.com/fivenet-app/fivenet/v2025/services/qualifications/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var (
	tQualiResults = table.FivenetQualificationsResults.AS("qualification_result")

	tJobLabels = table.FivenetJobLabels
)

func (s *Server) ListQualificationsResults(ctx context.Context, req *pbqualifications.ListQualificationsResultsRequest) (*pbqualifications.ListQualificationsResultsResponse, error) {
	if req.QualificationId != nil {
		trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.qualifications.id", int64(*req.QualificationId)))
	}
	if req.UserId != nil {
		trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.qualifications.user_id", int64(*req.UserId)))
	}

	tUser := tables.User().AS("user")
	tCreator := tUser.AS("creator")

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	tQuali := tQuali.AS("qualification_short")

	condition := tQualiResults.DeletedAt.IS_NULL()

	if req.QualificationId != nil {
		check, err := s.access.CanUserAccessTarget(ctx, *req.QualificationId, userInfo, qualifications.AccessLevel_ACCESS_LEVEL_GRADE)
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
					tQAccess.Access.IS_NOT_NULL(),
					jet.OR(
						tQAccess.Access.GT_EQ(jet.Int32(int32(qualifications.AccessLevel_ACCESS_LEVEL_GRADE))),
						jet.AND(
							tQAccess.Access.GT(jet.Int32(int32(qualifications.AccessLevel_ACCESS_LEVEL_BLOCKED))),
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
		for i := range req.Status {
			statuses = append(statuses, jet.Int16(int16(req.Status[i])))
		}

		condition = condition.AND(tQualiResults.Status.IN(statuses...))
	}

	countStmt := tQualiResults.
		SELECT(
			jet.COUNT(countColumn).AS("data_count.total"),
		).
		FROM(
			tQualiResults.
				INNER_JOIN(tQuali,
					tQuali.ID.EQ(tQualiResults.QualificationID),
				).
				LEFT_JOIN(tQAccess,
					tQAccess.TargetID.EQ(tQuali.ID).
						AND(tQAccess.Job.EQ(jet.String(userInfo.Job))).
						AND(tQAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
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

	pag, limit := req.Pagination.GetResponseWithPageSize(count.Total, QualificationsPageSize)
	resp := &pbqualifications.ListQualificationsResultsResponse{
		Pagination: pag,
		Results:    []*qualifications.QualificationResult{},
	}
	if count.Total <= 0 {
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
			tQuali.Draft,
			tQuali.Public,
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
				LEFT_JOIN(tQAccess,
					tQAccess.TargetID.EQ(tQuali.ID).
						AND(tQAccess.Job.EQ(jet.String(userInfo.Job))).
						AND(tQAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
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
	for i := range resp.Results {
		if resp.Results[i].User != nil {
			jobInfoFn(resp.Results[i].User)
		}

		if resp.Results[i].Creator != nil {
			jobInfoFn(resp.Results[i].Creator)
		}
	}

	return resp, nil
}

func (s *Server) CreateOrUpdateQualificationResult(ctx context.Context, req *pbqualifications.CreateOrUpdateQualificationResultRequest) (*pbqualifications.CreateOrUpdateQualificationResultResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbqualifications.QualificationsService_ServiceDesc.ServiceName,
		Method:  "CreateOrUpdateQualificationResult",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(ctx, req.Result.QualificationId, userInfo, qualifications.AccessLevel_ACCESS_LEVEL_GRADE)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !check {
		return nil, errorsqualifications.ErrFailedQuery
	}

	resultId, err := s.createOrUpdateQualificationResult(ctx, s.db, req.Result.QualificationId, req.Result.Id, userInfo, req.Result.UserId, req.Result.Status, req.Result.Score, req.Result.Summary, req.Grading)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	result, err := s.getQualificationResult(ctx, req.Result.QualificationId, resultId, nil, userInfo, req.Result.UserId)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	return &pbqualifications.CreateOrUpdateQualificationResultResponse{
		Result: result,
	}, nil
}

func (s *Server) createOrUpdateQualificationResult(ctx context.Context, tx qrm.DB, qualificationId uint64, resultId uint64, userInfo *userinfo.UserInfo, userId int32, status qualifications.ResultStatus, score *float32, summary string, grading *qualifications.ExamGrading) (uint64, error) {
	result, err := s.getQualificationResult(ctx, qualificationId, resultId, []qualifications.ResultStatus{qualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL}, userInfo, userId)
	if err != nil {
		return 0, err
	}

	quali, err := s.getQualification(ctx, qualificationId, tQuali.ID.EQ(jet.Uint64(qualificationId)), userInfo, false)
	if err != nil {
		return 0, err
	}

	tQualiResults := table.FivenetQualificationsResults
	// There is currently no result with status successful
	if resultId <= 0 && (result == nil || (result.Status != qualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL && status != qualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL)) {
		var creatorId jet.Expression
		if userInfo.UserId <= 0 {
			creatorId = jet.NULL
		} else {
			creatorId = jet.Int32(userInfo.UserId)
		}

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
				quali.Id,
				userId,
				status,
				score,
				summary,
				creatorId,
				userInfo.Job,
			)

		res, err := stmt.ExecContext(ctx, tx)
		if err != nil {
			return 0, err
		}

		lastId, err := res.LastInsertId()
		if err != nil {
			return 0, err
		}

		resultId = uint64(lastId)
	} else {
		result, err := s.getQualificationResult(ctx, quali.Id, resultId, nil, userInfo, userId)
		if err != nil {
			return 0, err
		}

		userId = result.UserId

		stmt := tQualiResults.
			UPDATE(
				tQualiResults.QualificationID,
				tQualiResults.UserID,
				tQualiResults.Status,
				tQualiResults.Score,
				tQualiResults.Summary,
			).
			SET(
				quali.Id,
				userId,
				status,
				score,
				summary,
			).
			WHERE(jet.AND(
				tQualiResults.ID.EQ(jet.Uint64(resultId)),
				tQualiResults.DeletedAt.IS_NULL(),
			))

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return 0, err
		}
	}

	if quali.ExamMode > qualifications.QualificationExamMode_QUALIFICATION_EXAM_MODE_DISABLED && grading != nil { // Only update the exam grading info when
		// Insert/update exam grading info from tutor
		stmt := tExamResponses.
			UPDATE(
				tExamResponses.Grading,
			).
			SET(
				grading,
			).
			WHERE(jet.AND(
				tExamResponses.QualificationID.EQ(jet.Uint64(quali.Id)),
				tExamResponses.UserID.EQ(jet.Int32(userId)),
			))

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return 0, err
		}
	}

	if quali.LabelSyncEnabled {
		// Add/Remove label based on result status
		if err := s.handleColleagueLabelSync(ctx, tx, userInfo, quali, userId, status == qualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL); err != nil {
			return 0, err
		}
	}

	// If the result is successful, complete the request status
	if status == qualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL {
		if err := s.updateRequestStatus(ctx, tx, qualificationId, userId, qualifications.RequestStatus_REQUEST_STATUS_COMPLETED); err != nil {
			return 0, err
		}
	} else {
		// If failed or other status, delete the request
		if err := s.deleteQualificationRequest(ctx, tx, qualificationId, userId); err != nil {
			return 0, err
		}

		if err := s.deleteExamUser(ctx, tx, qualificationId, userId); err != nil {
			return 0, err
		}
	}

	// Only send notification when the original result had no score and wasn't in pending status
	if status != qualifications.ResultStatus_RESULT_STATUS_PENDING && (result == nil || (result.Status == qualifications.ResultStatus_RESULT_STATUS_PENDING || (result.Score == nil && score != nil))) {
		if err := s.notif.NotifyUser(ctx, &notifications.Notification{
			UserId: userId,
			Title: &common.I18NItem{
				Key: "notifications.qualifications.result_updated.title",
			},
			Content: &common.I18NItem{
				Key:        "notifications.qualifications.result_updated.content",
				Parameters: map[string]string{"abbreviation": quali.Abbreviation, "title": quali.Title},
			},
			Category: notifications.NotificationCategory_NOTIFICATION_CATEGORY_GENERAL,
			Type:     notifications.NotificationType_NOTIFICATION_TYPE_INFO,
			Data: &notifications.Data{
				Link: &notifications.Link{
					To: fmt.Sprintf("/qualifications/%d", qualificationId),
				},
			},
		}); err != nil {
			return 0, err
		}
	}

	return resultId, nil
}

func (s *Server) getQualificationResult(ctx context.Context, qualificationId uint64, resultId uint64, status []qualifications.ResultStatus, userInfo *userinfo.UserInfo, userId int32) (*qualifications.QualificationResult, error) {
	tUser := tables.User().AS("user")
	tCreator := tUser.AS("creator")

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
		for i := range status {
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

func (s *Server) DeleteQualificationResult(ctx context.Context, req *pbqualifications.DeleteQualificationResultRequest) (*pbqualifications.DeleteQualificationResultResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbqualifications.QualificationsService_ServiceDesc.ServiceName,
		Method:  "DeleteQualificationResult",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	result, err := s.getQualificationResult(ctx, 0, req.ResultId, nil, userInfo, 0)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if result == nil {
		return &pbqualifications.DeleteQualificationResultResponse{}, nil
	}

	check, err := s.access.CanUserAccessTarget(ctx, result.QualificationId, userInfo, qualifications.AccessLevel_ACCESS_LEVEL_EDIT)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !check {
		return nil, errorsqualifications.ErrFailedQuery
	}

	quali, err := s.getQualification(ctx, result.QualificationId, tQuali.ID.EQ(jet.Uint64(result.QualificationId)), userInfo, false)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	tQualiResults := table.FivenetQualificationsResults

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

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

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	if err := s.deleteExamUser(ctx, tx, result.QualificationId, result.UserId); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	if quali.LabelSyncEnabled {
		// Remove label as we are deleting the result
		if err := s.handleColleagueLabelSync(ctx, tx, userInfo, quali, result.UserId, false); err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_DELETED

	return &pbqualifications.DeleteQualificationResultResponse{}, nil
}

func (s *Server) handleColleagueLabelSync(ctx context.Context, tx qrm.DB, userInfo *userinfo.UserInfo, quali *qualifications.Qualification, targetUserId int32, addLabel bool) error {
	if quali.LabelSyncFormat == nil || *quali.LabelSyncFormat == "" {
		defaultFormat := QualificationsLabelDefaultFormat
		quali.LabelSyncFormat = &defaultFormat
	}

	labelName := strings.ReplaceAll(*quali.LabelSyncFormat, "%abbr%", quali.Abbreviation)
	labelName = strings.ReplaceAll(labelName, "%name%", quali.Title)
	labelName = strings.TrimSpace(labelName)

	// Make sure that the label isn't empty when all is screwed up
	if labelName == "" {
		labelName = fmt.Sprintf("%s: %s", quali.Abbreviation, quali.Title)
	}

	// Create label if it doesn't exist yet
	createStmt := tJobLabels.
		INSERT(
			tJobLabels.Job,
			tJobLabels.Name,
		).
		VALUES(
			userInfo.Job,
			labelName,
		)

	labelId := uint64(0)
	res, err := createStmt.ExecContext(ctx, tx)
	if err != nil {
		if !dbutils.IsDuplicateError(err) {
			return err
		}

		// Retrieve existing label
		tJobLabels := tJobLabels.AS("label")

		stmt := tJobLabels.
			SELECT(
				tJobLabels.ID.AS("id"),
			).
			FROM(tJobLabels).
			WHERE(jet.AND(
				tJobLabels.Job.EQ(jet.String(userInfo.Job)),
				tJobLabels.Name.EQ(jet.String(labelName)),
			)).
			LIMIT(1)

		dest := struct {
			ID uint64 `alias:"id"`
		}{}
		if err := stmt.QueryContext(ctx, tx, &dest); err != nil {
			return err
		}

		labelId = dest.ID
	} else {
		lastId, err := res.LastInsertId()
		if err != nil {
			return err
		}

		labelId = uint64(lastId)
	}

	tUserLabels := table.FivenetJobColleagueLabels

	// Ensure that the colleague has the label set if successful result or removed
	if addLabel {
		stmt := tUserLabels.
			INSERT(
				tUserLabels.UserID,
				tUserLabels.Job,
				tUserLabels.LabelID,
			).
			VALUES(
				targetUserId,
				userInfo.Job,
				labelId,
			)

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			if !dbutils.IsDuplicateError(err) {
				return err
			}
		}
	} else {
		stmt := tUserLabels.
			DELETE().
			WHERE(jet.AND(
				tUserLabels.UserID.EQ(jet.Int32(targetUserId)),
				tUserLabels.Job.EQ(jet.String(userInfo.Job)),
				tUserLabels.LabelID.EQ(jet.Uint64(labelId)),
			)).
			LIMIT(1)

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return err
		}
	}

	return nil
}
