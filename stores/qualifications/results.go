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
	tUser := table.FivenetUser.AS("user")
	tCreator := tUser.AS("creator")

	condition := mysql.Bool(true)
	if !userInfo.GetSuperuser() {
		condition = condition.AND(tQualiResult.DeletedAt.IS_NULL())
	}
	visibilityCondition := mysql.Bool(true)
	if !userInfo.GetSuperuser() {
		visibilityCondition = visibilityCondition.AND(
			tQuali.DeletedAt.IS_NULL(),
		)
	}
	if opts.QualificationID > 0 {
		visibilityCondition = visibilityCondition.AND(
			tQuali.ID.EQ(mysql.Int64(opts.QualificationID)),
		)
	}
	visibleGradeCondition := visibilityCondition
	visibleViewCondition := visibilityCondition
	if !userInfo.GetSuperuser() {
		visibleGradeCondition = visibleGradeCondition.AND(mysql.OR(
			tQuali.Public.IS_TRUE(),
			mysql.AND(
				tQuali.CreatorID.EQ(mysql.Int32(userInfo.GetUserId())),
				tQuali.CreatorJob.EQ(mysql.String(userInfo.GetJob())),
			),
			s.access.ACLAccessExistsCondition(
				tQuali.ID,
				userInfo,
				int32(qualificationsaccess.AccessLevel_ACCESS_LEVEL_GRADE),
			),
		))
		visibleViewCondition = visibleViewCondition.AND(mysql.OR(
			tQuali.Public.IS_TRUE(),
			mysql.AND(
				tQuali.CreatorID.EQ(mysql.Int32(userInfo.GetUserId())),
				tQuali.CreatorJob.EQ(mysql.String(userInfo.GetJob())),
			),
			s.access.ACLAccessExistsCondition(
				tQuali.ID,
				userInfo,
				int32(qualificationsaccess.AccessLevel_ACCESS_LEVEL_VIEW),
			),
		))
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
				AND(tQualiResult.UserID.EQ(mysql.Int32(userInfo.GetUserId())))
			countColumn = mysql.DISTINCT(tQualiResult.QualificationID)
		} else {
			countColumn = mysql.DISTINCT(tQualiResult.UserID)
		}
	}

	if opts.QualificationID > 0 {
		condition = condition.AND(visibleGradeCondition)
	} else {
		condition = condition.AND(mysql.OR(
			mysql.AND(
				tQualiResult.CreatorID.EQ(mysql.Int32(userInfo.GetUserId())),
				tQualiResult.CreatorJob.EQ(mysql.String(userInfo.GetJob())),
			),
			mysql.OR(
				visibleGradeCondition,
				mysql.AND(
					visibleViewCondition,
					tQualiResult.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
				),
			),
		))
	}

	if len(opts.Status) > 0 {
		statuses := []mysql.Expression{}
		for i := range opts.Status {
			statuses = append(statuses, mysql.Int32(int32(opts.Status[i])))
		}
		condition = condition.AND(tQualiResult.Status.IN(statuses...))
	}

	var countStmt mysql.Statement = tQualiResult.
		SELECT(mysql.COUNT(countColumn).AS("data_count.total")).
		FROM(
			tQualiResult.
				INNER_JOIN(tQuali, tQuali.ID.EQ(tQualiResult.QualificationID)).
				LEFT_JOIN(tUser, tQualiResult.UserID.EQ(tUser.ID)),
		).
		WHERE(condition)

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

	var stmt mysql.Statement = tQualiResult.
		SELECT(columns[0], columns[1:]...).
		FROM(
			tQualiResult.
				INNER_JOIN(tQuali, tQuali.ID.EQ(tQualiResult.QualificationID)).
				LEFT_JOIN(tUser, tQualiResult.UserID.EQ(tUser.ID)).
				LEFT_JOIN(tCreator, tQualiResult.CreatorID.EQ(tCreator.ID)),
		).
		GROUP_BY(tQualiResult.Status, tQualiResult.CreatedAt, tQualiResult.ID).
		ORDER_BY(orderBys...).
		WHERE(condition).
		OFFSET(opts.Pagination.GetOffset()).
		LIMIT(limit)

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
	stmt := tQualiResult.
		UPDATE(tQualiResult.DeletedAt).
		SET(mysql.CURRENT_TIMESTAMP()).
		WHERE(tQualiResult.ID.EQ(mysql.Int64(resultId))).
		LIMIT(1)

	_, err := stmt.ExecContext(ctx, tx)
	return err
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

	return res.LastInsertId()
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

	_, err := stmt.ExecContext(ctx, tx)
	return err
}

func (s *Store) UpdateExamResponseGrading(
	ctx context.Context,
	tx qrm.DB,
	qualificationId int64,
	userId int32,
	grading *qualificationsexam.ExamGrading,
) error {
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
