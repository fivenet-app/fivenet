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
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/userinfo"
	pbqualifications "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/qualifications"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsqualifications "github.com/fivenet-app/fivenet/v2025/services/qualifications/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	logging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

var (
	tQualiResults = table.FivenetQualificationsResults.AS("qualification_result")

	tJobLabels = table.FivenetJobLabels
)

func (s *Server) ListQualificationsResults(
	ctx context.Context,
	req *pbqualifications.ListQualificationsResultsRequest,
) (*pbqualifications.ListQualificationsResultsResponse, error) {
	if req.QualificationId != nil {
		logging.InjectFields(
			ctx,
			logging.Fields{"fivenet.qualifications.id", req.GetQualificationId()},
		)
	}
	if req.UserId != nil {
		logging.InjectFields(ctx, logging.Fields{"fivenet.qualifications.user_id", req.GetUserId()})
	}

	tUser := tables.User().AS("user")
	tCreator := tUser.AS("creator")

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	tQuali := tQuali.AS("qualification_short")

	condition := tQualiResults.DeletedAt.IS_NULL()

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
			tQualiResults.QualificationID.EQ(mysql.Int64(req.GetQualificationId())),
		)
	} else {
		accessExists := mysql.EXISTS(
			mysql.SELECT(mysql.Int(1)).
				FROM(tQAccess).
				WHERE(mysql.AND(
					tQAccess.TargetID.EQ(tQualiResults.QualificationID),
					tQAccess.Job.EQ(mysql.String(userInfo.GetJob())),
					tQAccess.MinimumGrade.LT_EQ(mysql.Int32(userInfo.GetJobGrade())),

					mysql.OR(
						tQAccess.Access.GT_EQ(mysql.Int32(int32(qualifications.AccessLevel_ACCESS_LEVEL_GRADE))),
						mysql.AND(
							tQAccess.Access.GT_EQ(mysql.Int32(int32(qualifications.AccessLevel_ACCESS_LEVEL_VIEW))),
							tQualiResults.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
						),
					),
				)),
		)

		condition = condition.AND(mysql.AND(
			tQuali.DeletedAt.IS_NULL(),
			mysql.OR(
				mysql.AND(
					tQualiResults.CreatorID.EQ(mysql.Int32(userInfo.GetUserId())),
					tQualiResults.CreatorJob.EQ(mysql.String(userInfo.GetJob())),
				),
				accessExists,
			),
		))
	}

	countColumn := mysql.Expression(tQualiResults.QualificationID)
	if req.UserId != nil {
		condition = condition.AND(mysql.AND(
			tUser.Job.EQ(mysql.String(userInfo.GetJob())),
			tQualiResults.UserID.EQ(mysql.Int32(req.GetUserId())),
		))
	} else {
		if req.QualificationId == nil {
			condition = condition.AND(tUser.Job.EQ(mysql.String(userInfo.GetJob()))).AND(tQualiResults.UserID.EQ(mysql.Int32(userInfo.GetUserId())))
			countColumn = mysql.DISTINCT(tQualiResults.QualificationID)
		} else {
			countColumn = mysql.DISTINCT(tQualiResults.UserID)
		}
	}

	if len(req.GetStatus()) > 0 {
		statuses := []mysql.Expression{}
		for i := range req.GetStatus() {
			statuses = append(statuses, mysql.Int32(int32(req.GetStatus()[i])))
		}

		condition = condition.AND(tQualiResults.Status.IN(statuses...))
	}

	countStmt := tQualiResults.
		SELECT(
			mysql.COUNT(countColumn).AS("data_count.total"),
		).
		FROM(
			tQualiResults.
				INNER_JOIN(tQuali,
					tQuali.ID.EQ(tQualiResults.QualificationID),
				).
				LEFT_JOIN(tQAccess,
					mysql.AND(
						tQAccess.TargetID.EQ(tQuali.ID),
						tQAccess.Job.EQ(mysql.String(userInfo.GetJob())),
						tQAccess.MinimumGrade.LT_EQ(mysql.Int32(userInfo.GetJobGrade())),
					),
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

	pag, limit := req.GetPagination().GetResponseWithPageSize(count.Total, QualificationsPageSize)
	resp := &pbqualifications.ListQualificationsResultsResponse{
		Pagination: pag,
		Results:    []*qualifications.QualificationResult{},
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
			column = tQualiResults.Status
		case "createdAt":
			fallthrough
		default:
			column = tQualiResults.CreatedAt
		}

		if req.GetSort().GetDirection() == database.AscSortDirection {
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
					mysql.AND(
						tQAccess.TargetID.EQ(tQuali.ID),
						tQAccess.Job.EQ(mysql.String(userInfo.GetJob())),
						tQAccess.MinimumGrade.LT_EQ(mysql.Int32(userInfo.GetJobGrade())),
					),
				),
		).
		GROUP_BY(tQualiResults.Status, tQualiResults.CreatedAt, tQualiResults.ID).
		ORDER_BY(orderBys...).
		WHERE(condition).
		OFFSET(req.GetPagination().GetOffset()).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Results); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.GetResults() {
		if resp.GetResults()[i].GetUser() != nil {
			jobInfoFn(resp.GetResults()[i].GetUser())
		}

		if resp.GetResults()[i].GetCreator() != nil {
			jobInfoFn(resp.GetResults()[i].GetCreator())
		}
	}

	return resp, nil
}

func (s *Server) CreateOrUpdateQualificationResult(
	ctx context.Context,
	req *pbqualifications.CreateOrUpdateQualificationResultRequest,
) (*pbqualifications.CreateOrUpdateQualificationResultResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbqualifications.QualificationsService_ServiceDesc.ServiceName,
		Method:  "CreateOrUpdateQualificationResult",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetResult().GetQualificationId(),
		userInfo,
		qualifications.AccessLevel_ACCESS_LEVEL_GRADE,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !check {
		return nil, errorsqualifications.ErrFailedQuery
	}

	resultId, err := s.createOrUpdateQualificationResult(
		ctx,
		s.db,
		req.GetResult().GetQualificationId(),
		req.GetResult().GetId(),
		userInfo,
		req.GetResult().GetUserId(),
		req.GetResult().GetStatus(),
		//nolint:protogetter // The value is needed as a pointer
		req.GetResult().Score,
		req.GetResult().GetSummary(),
		req.GetGrading(),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	result, err := s.getQualificationResult(
		ctx,
		req.GetResult().GetQualificationId(),
		resultId,
		nil,
		userInfo,
		req.GetResult().GetUserId(),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	return &pbqualifications.CreateOrUpdateQualificationResultResponse{
		Result: result,
	}, nil
}

func (s *Server) createOrUpdateQualificationResult(
	ctx context.Context,
	tx qrm.DB,
	qualificationId int64,
	resultId int64,
	userInfo *userinfo.UserInfo,
	userId int32,
	status qualifications.ResultStatus,
	score *float32,
	summary string,
	grading *qualifications.ExamGrading,
) (int64, error) {
	currentResult, err := s.getQualificationResult(
		ctx,
		qualificationId,
		resultId,
		[]qualifications.ResultStatus{qualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL},
		userInfo,
		userId,
	)
	if err != nil {
		return 0, err
	}

	quali, err := s.getQualification(
		ctx,
		qualificationId,
		tQuali.ID.EQ(mysql.Int64(qualificationId)),
		userInfo,
		false,
	)
	if err != nil {
		return 0, err
	}

	tQualiResults := table.FivenetQualificationsResults
	// There is currently no result with status successful
	if resultId <= 0 &&
		(currentResult == nil || (currentResult.GetStatus() != qualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL && status != qualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL)) {
		var creatorId mysql.Expression
		if userInfo.GetUserId() <= 0 {
			creatorId = mysql.NULL
		} else {
			creatorId = mysql.Int32(userInfo.GetUserId())
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
				quali.GetId(),
				userId,
				status,
				score,
				summary,
				creatorId,
				userInfo.GetJob(),
			)

		res, err := stmt.ExecContext(ctx, tx)
		if err != nil {
			return 0, err
		}

		lastId, err := res.LastInsertId()
		if err != nil {
			return 0, err
		}

		resultId = lastId
	} else {
		result, err := s.getQualificationResult(ctx, quali.GetId(), resultId, nil, userInfo, userId)
		if err != nil {
			return 0, err
		}

		userId = result.GetUserId()

		stmt := tQualiResults.
			UPDATE(
				tQualiResults.QualificationID,
				tQualiResults.UserID,
				tQualiResults.Status,
				tQualiResults.Score,
				tQualiResults.Summary,
			).
			SET(
				quali.GetId(),
				userId,
				status,
				score,
				summary,
			).
			WHERE(mysql.AND(
				tQualiResults.ID.EQ(mysql.Int64(resultId)),
				tQualiResults.DeletedAt.IS_NULL(),
			))

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return 0, err
		}
	}

	if quali.GetExamMode() > qualifications.QualificationExamMode_QUALIFICATION_EXAM_MODE_DISABLED &&
		grading != nil { // Only update the exam grading info when
		// Insert/update exam grading info from tutor
		stmt := tExamResponses.
			UPDATE(
				tExamResponses.Grading,
			).
			SET(
				grading,
			).
			WHERE(mysql.AND(
				tExamResponses.QualificationID.EQ(mysql.Int64(quali.GetId())),
				tExamResponses.UserID.EQ(mysql.Int32(userId)),
			))

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return 0, err
		}
	}

	if quali.GetLabelSyncEnabled() {
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
	if status != qualifications.ResultStatus_RESULT_STATUS_PENDING &&
		(currentResult == nil || (currentResult.GetStatus() == qualifications.ResultStatus_RESULT_STATUS_PENDING || (currentResult.Score == nil && score != nil))) {
		if err := s.notif.NotifyUser(ctx, &notifications.Notification{
			UserId: userId,
			Title: &common.I18NItem{
				Key: "notifications.qualifications.result_updated.title",
			},
			Content: &common.I18NItem{
				Key:        "notifications.qualifications.result_updated.content",
				Parameters: map[string]string{"abbreviation": quali.GetAbbreviation(), "title": quali.GetTitle()},
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

func (s *Server) getQualificationResult(
	ctx context.Context,
	qualificationId int64,
	resultId int64,
	status []qualifications.ResultStatus,
	userInfo *userinfo.UserInfo,
	userId int32,
) (*qualifications.QualificationResult, error) {
	tUser := tables.User().AS("user")
	tCreator := tUser.AS("creator")

	condition := tQualiResults.DeletedAt.IS_NULL()

	if resultId > 0 {
		condition = condition.AND(tQualiResults.ID.EQ(mysql.Int64(resultId)))
	} else if userId > 0 {
		condition = condition.AND(tQualiResults.UserID.EQ(mysql.Int32(userId)))
	} else {
		condition = condition.AND(tQualiResults.UserID.EQ(mysql.Int32(userInfo.GetUserId())))
	}
	if qualificationId > 0 {
		condition = condition.AND(tQualiResults.QualificationID.EQ(mysql.Int64(qualificationId)))
	}

	if len(status) > 0 {
		statusConds := make([]mysql.Expression, len(status))
		for i := range status {
			statusConds[i] = mysql.Int32(int32(status[i]))
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

	if result.GetId() == 0 {
		return nil, nil
	}

	if result.GetUser() != nil {
		s.enricher.EnrichJobInfoSafe(userInfo, result.GetUser())
	}

	if result.GetCreator() != nil {
		s.enricher.EnrichJobInfoSafe(userInfo, result.GetCreator())
	}

	return &result, nil
}

func (s *Server) DeleteQualificationResult(
	ctx context.Context,
	req *pbqualifications.DeleteQualificationResultRequest,
) (*pbqualifications.DeleteQualificationResultResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbqualifications.QualificationsService_ServiceDesc.ServiceName,
		Method:  "DeleteQualificationResult",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	result, err := s.getQualificationResult(ctx, 0, req.GetResultId(), nil, userInfo, 0)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if result == nil {
		return &pbqualifications.DeleteQualificationResultResponse{}, nil
	}

	check, err := s.access.CanUserAccessTarget(
		ctx,
		result.GetQualificationId(),
		userInfo,
		qualifications.AccessLevel_ACCESS_LEVEL_EDIT,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !check {
		return nil, errorsqualifications.ErrFailedQuery
	}

	quali, err := s.getQualification(
		ctx,
		result.GetQualificationId(),
		tQuali.ID.EQ(mysql.Int64(result.GetQualificationId())),
		userInfo,
		false,
	)
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
			mysql.CURRENT_TIMESTAMP(),
		).
		WHERE(mysql.AND(
			tQualiResults.ID.EQ(mysql.Int64(result.GetId())),
			tQualiResults.ID.EQ(mysql.Int64(req.GetResultId())),
		))

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	if err := s.deleteExamUser(ctx, tx, result.GetQualificationId(), result.GetUserId()); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	if quali.GetLabelSyncEnabled() {
		// Remove label as we are deleting the result
		if err := s.handleColleagueLabelSync(ctx, tx, userInfo, quali, result.GetUserId(), false); err != nil {
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

func (s *Server) handleColleagueLabelSync(
	ctx context.Context,
	tx qrm.DB,
	userInfo *userinfo.UserInfo,
	quali *qualifications.Qualification,
	targetUserId int32,
	addLabel bool,
) error {
	if quali.LabelSyncFormat == nil || quali.GetLabelSyncFormat() == "" {
		defaultFormat := QualificationsLabelDefaultFormat
		quali.LabelSyncFormat = &defaultFormat
	}

	labelName := strings.ReplaceAll(quali.GetLabelSyncFormat(), "%abbr%", quali.GetAbbreviation())
	labelName = strings.ReplaceAll(labelName, "%name%", quali.GetTitle())
	labelName = strings.TrimSpace(labelName)

	// Make sure that the label isn't empty when all is screwed up
	if labelName == "" {
		labelName = fmt.Sprintf("%s: %s", quali.GetAbbreviation(), quali.GetTitle())
	}

	// Create label if it doesn't exist yet
	createStmt := tJobLabels.
		INSERT(
			tJobLabels.Job,
			tJobLabels.Name,
		).
		VALUES(
			userInfo.GetJob(),
			labelName,
		)

	var labelId int64
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
			WHERE(mysql.AND(
				tJobLabels.Job.EQ(mysql.String(userInfo.GetJob())),
				tJobLabels.Name.EQ(mysql.String(labelName)),
			)).
			LIMIT(1)

		dest := struct {
			ID int64 `alias:"id"`
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

		labelId = lastId
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
				userInfo.GetJob(),
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
			WHERE(mysql.AND(
				tUserLabels.UserID.EQ(mysql.Int32(targetUserId)),
				tUserLabels.Job.EQ(mysql.String(userInfo.GetJob())),
				tUserLabels.LabelID.EQ(mysql.Int64(labelId)),
			)).
			LIMIT(1)

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return err
		}
	}

	return nil
}
