package qualificationsstore

import (
	"context"
	"errors"

	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	resqualifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications"
	qualificationsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications/access"
	qualificationsexam "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications/exam"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pbqualifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/qualifications"
	objectaccess "github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) ListQualificationsResults(
	ctx context.Context,
	opts ListQualificationsResultsOptions,
	userInfo *userinfo.UserInfo,
	includePhoneNumber bool,
) (*pbqualifications.ListQualificationsResultsResponse, error) {
	tQuali := table.FivenetQualifications.AS("qualificationshort")
	tUser := table.FivenetUser.AS("user")
	tCreator := tUser.AS("creator")
	userID := int32(0)
	if userInfo != nil {
		userID = userInfo.GetUserId()
	}

	condition := mysql.Bool(true)
	if !userInfo.GetSuperuser() {
		condition = condition.AND(tQualiResult.DeletedAt.IS_NULL())
	}

	countColumn := mysql.Expression(tQualiResult.QualificationID)
	if len(opts.UserIDs) > 0 {
		userIds := []mysql.Expression{}
		for _, userId := range opts.UserIDs {
			userIds = append(userIds, mysql.Int32(userId))
		}

		condition = condition.AND(mysql.AND(
			tQualiResult.UserID.IN(userIds...),
		))
	} else {
		if opts.QualificationID == 0 {
			condition = condition.AND(tUser.Job.EQ(mysql.String(userInfo.GetJob()))).
				AND(tQualiResult.UserID.EQ(mysql.Int32(userID)))
			countColumn = mysql.DISTINCT(tQualiResult.QualificationID)
		} else {
			countColumn = mysql.DISTINCT(tQualiResult.UserID)
		}
	}

	if len(opts.Status) > 0 {
		statuses := []mysql.Expression{}
		for i := range opts.Status {
			statuses = append(statuses, mysql.Int32(int32(opts.Status[i])))
		}
		condition = condition.AND(tQualiResult.Status.IN(statuses...))
	}

	orderBys := s.resultSorter.Build(opts.Sort)

	columns := mysql.ProjectionList{
		tQualiResult.ID,
		tQualiResult.CreatedAt,
		tQualiResult.QualificationID,
		tQualiResult.UserID,
		tUser.ID,
		tUser.Job,
		tUser.JobGrade,
		tUser.Firstname,
		tUser.Lastname,
		tUser.Dateofbirth,
		tQualiResult.Status,
		tQualiResult.Score,
		tQualiResult.Summary,
		tQualiResult.CreatorID,
		tCreator.ID,
		tCreator.Job,
		tCreator.JobGrade,
		tCreator.Firstname,
		tCreator.Lastname,
		tCreator.Dateofbirth,
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
	}
	if includePhoneNumber {
		columns = append(columns, tUser.PhoneNumber, tCreator.PhoneNumber)
	}

	var (
		countStmt      mysql.Statement
		ctes           []mysql.CommonTableExpression
		visibleIDs     objectaccess.VisibilityQuery
		visibleQualiID mysql.IntegerExpression
	)

	if userInfo != nil && !userInfo.GetSuperuser() {
		visibleQualificationCondition := mysql.Bool(true)
		if opts.QualificationID > 0 {
			visibleQualificationCondition = visibleQualificationCondition.AND(
				table.FivenetQualifications.ID.EQ(mysql.Int64(opts.QualificationID)),
			)
		}
		visibleIDs = s.access.VisibleIDsByConditionQuery(
			userInfo,
			int32(qualificationsaccess.AccessLevel_ACCESS_LEVEL_VIEW),
			false,
			visibleQualificationCondition,
		)
		ctes = visibleIDs.CTEs
		visibleQualiID = mysql.IntegerColumn("id").From(visibleIDs.Table)

		countStmt = tQualiResult.
			SELECT(mysql.COUNT(countColumn).AS("data_count.total")).
			FROM(
				visibleIDs.Table.
					INNER_JOIN(tQuali, tQuali.ID.EQ(visibleQualiID)).
					INNER_JOIN(tQualiResult, tQualiResult.QualificationID.EQ(tQuali.ID)).
					LEFT_JOIN(tUser, tQualiResult.UserID.EQ(tUser.ID)),
			).
			WHERE(condition)
	} else {
		countStmt = tQualiResult.
			SELECT(mysql.COUNT(countColumn).AS("data_count.total")).
			FROM(
				tQualiResult.
					INNER_JOIN(tQuali, tQuali.ID.EQ(tQualiResult.QualificationID)).
					LEFT_JOIN(tUser, tQualiResult.UserID.EQ(tUser.ID)),
			).
			WHERE(condition)
	}

	if len(ctes) > 0 {
		countStmt = mysql.WITH(ctes...)(countStmt)
	}

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	pag, limit := opts.Pagination.GetResponseWithPageSize(count.Total, QualificationsPageSize)
	resp := &pbqualifications.ListQualificationsResultsResponse{
		Pagination: pag,
		Results:    []*resqualifications.QualificationResult{},
	}
	if count.Total <= 0 {
		return resp, nil
	}

	var stmt mysql.Statement
	if userInfo != nil && !userInfo.GetSuperuser() {
		stmt = tQualiResult.
			SELECT(columns[0], columns[1:]...).
			FROM(
				visibleIDs.Table.
					INNER_JOIN(tQuali, tQuali.ID.EQ(visibleQualiID)).
					INNER_JOIN(tQualiResult, tQualiResult.QualificationID.EQ(tQuali.ID)).
					LEFT_JOIN(tUser, tQualiResult.UserID.EQ(tUser.ID)).
					LEFT_JOIN(tCreator, tQualiResult.CreatorID.EQ(tCreator.ID)),
			).
			GROUP_BY(tQualiResult.Status, tQualiResult.CreatedAt, tQualiResult.ID).
			ORDER_BY(orderBys...).
			WHERE(condition).
			OFFSET(opts.Pagination.GetOffset()).
			LIMIT(limit)
	} else {
		stmt = tQualiResult.
			SELECT(columns[0], columns[1:]...).
			FROM(tQualiResult.
				INNER_JOIN(tQuali, tQuali.ID.EQ(tQualiResult.QualificationID)).
				LEFT_JOIN(tUser, tQualiResult.UserID.EQ(tUser.ID)).
				LEFT_JOIN(tCreator, tQualiResult.CreatorID.EQ(tCreator.ID)),
			).
			GROUP_BY(tQualiResult.Status, tQualiResult.CreatedAt, tQualiResult.ID).
			ORDER_BY(orderBys...).
			WHERE(condition).
			OFFSET(opts.Pagination.GetOffset()).
			LIMIT(limit)
	}

	if len(ctes) > 0 {
		stmt = mysql.WITH(ctes...)(stmt)
	}

	if err := stmt.QueryContext(ctx, s.db, &resp.Results); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return resp, nil
}

func (s *Store) GetQualificationResult(
	ctx context.Context,
	qualificationId int64,
	resultId int64,
	status []resqualifications.ResultStatus,
	userInfo *userinfo.UserInfo,
	userId int32,
	includePhoneNumber bool,
) (*resqualifications.QualificationResult, error) {
	tUser := table.FivenetUser.AS("user")
	tCreator := tUser.AS("creator")

	condition := mysql.Bool(true)
	if userInfo == nil || !userInfo.GetSuperuser() {
		condition = condition.AND(tQualiResult.DeletedAt.IS_NULL())
	}
	if resultId > 0 {
		condition = condition.AND(tQualiResult.ID.EQ(mysql.Int64(resultId)))
	} else if userId > 0 {
		condition = condition.AND(tQualiResult.UserID.EQ(mysql.Int32(userId)))
	} else {
		condition = condition.AND(tQualiResult.UserID.EQ(mysql.Int32(userInfo.GetUserId())))
	}
	if qualificationId > 0 {
		condition = condition.AND(tQualiResult.QualificationID.EQ(mysql.Int64(qualificationId)))
	}
	if len(status) > 0 {
		statusConds := make([]mysql.Expression, len(status))
		for i := range status {
			statusConds[i] = mysql.Int32(int32(status[i]))
		}
		condition = condition.AND(tQualiResult.Status.IN(statusConds...))
	}

	columns := mysql.ProjectionList{
		tQualiResult.ID,
		tQualiResult.CreatedAt,
		tQualiResult.DeletedAt,
		tQualiResult.QualificationID,
		tQualiResult.UserID,
		tUser.ID,
		tUser.Job,
		tUser.JobGrade,
		tUser.Firstname,
		tUser.Lastname,
		tUser.Dateofbirth,
		tQualiResult.Status,
		tQualiResult.Score,
		tQualiResult.Summary,
		tQualiResult.CreatorID,
		tQualiResult.CreatorJob,
		tCreator.ID,
		tCreator.Job,
		tCreator.JobGrade,
		tCreator.Firstname,
		tCreator.Lastname,
		tCreator.Dateofbirth,
	}
	if includePhoneNumber {
		columns = append(columns, tUser.PhoneNumber, tCreator.PhoneNumber)
	}

	stmt := tQualiResult.
		SELECT(columns[0], columns[1:]...).
		FROM(tQualiResult.
			LEFT_JOIN(tUser, tUser.ID.EQ(tQualiResult.UserID)).
			LEFT_JOIN(tCreator, tCreator.ID.EQ(tQualiResult.CreatorID)),
		).
		GROUP_BY(tQualiResult.ID).
		ORDER_BY(tQualiResult.ID.DESC()).
		WHERE(condition).
		LIMIT(1)

	var result resqualifications.QualificationResult
	if err := stmt.QueryContext(ctx, s.db, &result); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}
	if result.GetId() == 0 {
		return nil, nil
	}

	return &result, nil
}

func (s *Store) DeleteQualificationResult(ctx context.Context, tx qrm.DB, resultId int64) error {
	tQualiResult := table.FivenetQualificationsResults
	stmt := tQualiResult.
		UPDATE(tQualiResult.DeletedAt).
		SET(mysql.CURRENT_TIMESTAMP()).
		WHERE(tQualiResult.ID.EQ(mysql.Int64(resultId))).
		LIMIT(1)

	res, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return nil
	}

	return s.deleteQualificationResultSuccessMapByResultID(ctx, tx, resultId)
}

func (s *Store) CreateQualificationResult(
	ctx context.Context,
	tx qrm.DB,
	qualificationId int64,
	userId int32,
	status resqualifications.ResultStatus,
	score *float32,
	summary string,
	creator *userinfo.UserInfo,
) (int64, error) {
	var creatorId mysql.Expression
	if creator.GetUserId() <= 0 {
		creatorId = mysql.NULL
	} else {
		creatorId = mysql.Int32(creator.GetUserId())
	}

	tQualiResult := table.FivenetQualificationsResults
	stmt := tQualiResult.
		INSERT(
			tQualiResult.QualificationID,
			tQualiResult.UserID,
			tQualiResult.Status,
			tQualiResult.Score,
			tQualiResult.Summary,
			tQualiResult.CreatorID,
			tQualiResult.CreatorJob,
		).
		VALUES(
			qualificationId,
			userId,
			status,
			score,
			summary,
			creatorId,
			creator.GetJob(),
		)

	res, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return 0, err
	}

	lastInsertId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	if status == resqualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL {
		if err := s.syncQualificationResultSuccessMap(
			ctx,
			tx,
			lastInsertId,
			qualificationId,
			userId,
			true,
		); err != nil {
			return 0, err
		}
	}

	return lastInsertId, nil
}

func (s *Store) UpdateQualificationResult(
	ctx context.Context,
	tx qrm.DB,
	qualificationId int64,
	resultId int64,
	userId int32,
	status resqualifications.ResultStatus,
	score *float32,
	summary string,
) error {
	tQualiResult := table.FivenetQualificationsResults
	stmt := tQualiResult.
		UPDATE(
			tQualiResult.QualificationID,
			tQualiResult.UserID,
			tQualiResult.Status,
			tQualiResult.Score,
			tQualiResult.Summary,
		).
		SET(
			qualificationId,
			userId,
			status,
			score,
			summary,
		).
		WHERE(mysql.AND(
			tQualiResult.ID.EQ(mysql.Int64(resultId)),
			tQualiResult.DeletedAt.IS_NULL(),
		)).
		LIMIT(1)

	res, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return nil
	}

	return s.syncQualificationResultSuccessMap(
		ctx,
		tx,
		resultId,
		qualificationId,
		userId,
		status == resqualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL,
	)
}

func (s *Store) UpdateExamResponseGrading(
	ctx context.Context,
	tx qrm.DB,
	qualificationId int64,
	userId int32,
	grading *qualificationsexam.ExamGrading,
) error {
	tExamResponse := table.FivenetQualificationsExamResponses
	stmt := tExamResponse.
		UPDATE(tExamResponse.Grading).
		SET(grading).
		WHERE(mysql.AND(
			tExamResponse.QualificationID.EQ(mysql.Int64(qualificationId)),
			tExamResponse.UserID.EQ(mysql.Int32(userId)),
		)).
		LIMIT(1)

	_, err := stmt.ExecContext(ctx, tx)
	return err
}
