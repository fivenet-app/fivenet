package qualifications

import (
	"context"
	"fmt"
	"strings"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/notifications"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications"
	qualificationsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications/access"
	qualificationsexam "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications/exam"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	permscitizens "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/citizens/perms"
	pbqualifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/qualifications"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorsqualifications "github.com/fivenet-app/fivenet/v2026/services/qualifications/errors"
	qualificationsstore "github.com/fivenet-app/fivenet/v2026/stores/qualifications"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	logging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

var tJobLabels = table.FivenetJobLabels

func (s *Server) ListQualificationsResults(
	ctx context.Context,
	req *pbqualifications.ListQualificationsResultsRequest,
) (*pbqualifications.ListQualificationsResultsResponse, error) {
	if req.GetQualificationId() > 0 {
		logging.InjectFields(
			ctx,
			logging.Fields{"fivenet.qualifications.id", req.GetQualificationId()},
		)
	}

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	includePhoneNumber := false
	if fields, err := permscitizens.CitizensService.ListCitizens.FieldsTyped.Get(
		s.perms,
		userInfo,
	); err == nil {
		includePhoneNumber = fields.Contains(
			permscitizens.CitizensServiceListCitizensFieldsPermValuePhoneNumber,
		)
	}

	resp, err := s.store.ListQualificationsResults(
		ctx,
		qualificationsstore.ListQualificationsResultsOptions{
			Pagination:      req.GetPagination(),
			Sort:            req.GetSort(),
			QualificationID: req.GetQualificationId(),
			Status:          req.GetStatus(),
			UserIDs:         req.GetUserIds(),
		},
		userInfo,
		includePhoneNumber,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
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

	check, err := s.access.CanUserAccessTarget(
		ctx,
		req.GetResult().GetQualificationId(),
		userInfo,
		int32(qualificationsaccess.AccessLevel_ACCESS_LEVEL_GRADE),
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
	grading *qualificationsexam.ExamGrading,
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

	quali, err := s.store.GetQualification(
		ctx,
		qualificationId,
		userInfo,
		false,
		false,
	)
	if err != nil {
		return 0, err
	}

	// There is currently no result with status successful
	if resultId <= 0 &&
		(currentResult == nil || (currentResult.GetStatus() != qualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL && status != qualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL)) {
		lastId, err := s.store.CreateQualificationResult(
			ctx,
			tx,
			quali.GetId(),
			userId,
			status,
			score,
			summary,
			userInfo,
		)
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
		if err := s.store.UpdateQualificationResult(
			ctx,
			tx,
			quali.GetId(),
			resultId,
			userId,
			status,
			score,
			summary,
		); err != nil {
			return 0, err
		}
	}

	if quali.GetExamMode() > qualificationsexam.QualificationExamMode_QUALIFICATION_EXAM_MODE_DISABLED &&
		grading != nil {
		if err := s.store.UpdateExamResponseGrading(
			ctx,
			tx,
			quali.GetId(),
			userId,
			grading,
		); err != nil {
			return 0, err
		}
	}

	if quali.GetLabelSyncEnabled() {
		// Add/Remove label based on result status
		if err := s.handleColleagueLabelSync(
			ctx,
			tx,
			userInfo,
			quali,
			userId,
			status == qualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL,
		); err != nil {
			return 0, err
		}
	}

	// If the result is successful, complete the request status
	if status == qualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL {
		if err := s.store.UpdateRequestStatus(
			ctx,
			tx,
			qualificationId,
			userId,
			qualifications.RequestStatus_REQUEST_STATUS_COMPLETED,
		); err != nil {
			return 0, err
		}
	} else {
		// If failed or other status, delete the request + exam user
		if err := s.deleteQualificationRequest(ctx, tx, qualificationId, userId); err != nil {
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
				Key: "notifications.qualifications.result_updated.content",
				Parameters: map[string]string{
					"abbreviation": quali.GetAbbreviation(),
					"title":        quali.GetTitle(),
				},
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
	result, err := s.store.GetQualificationResult(
		ctx,
		qualificationId,
		resultId,
		status,
		userInfo,
		userId,
		false,
	)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}

	if result.GetUser() != nil {
		s.enricher.EnrichJobInfoSafe(userInfo, result.GetUser())
	}

	if result.GetCreator() != nil {
		s.enricher.EnrichJobInfoSafe(userInfo, result.GetCreator())
	}

	return result, nil
}

func (s *Server) DeleteQualificationResult(
	ctx context.Context,
	req *pbqualifications.DeleteQualificationResultRequest,
) (*pbqualifications.DeleteQualificationResultResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

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
		int32(qualificationsaccess.AccessLevel_ACCESS_LEVEL_EDIT),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	if !check {
		return nil, errorsqualifications.ErrFailedQuery
	}

	quali, err := s.store.GetQualification(
		ctx,
		result.GetQualificationId(),
		userInfo,
		false,
		false,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	if err := s.store.DeleteQualificationResult(ctx, tx, result.GetId()); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	if err := s.store.DeleteExamUser(
		ctx,
		tx,
		result.GetQualificationId(),
		result.GetUserId(),
	); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	if quali.GetLabelSyncEnabled() {
		// Remove label as we are deleting the result
		if err := s.handleColleagueLabelSync(
			ctx,
			tx,
			userInfo,
			quali,
			result.GetUserId(),
			false,
		); err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_DELETED)

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
			tJobLabels.Color,
		).
		VALUES(
			userInfo.GetJob(),
			labelName,
			"#5c7aff", // Default color if not set
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
