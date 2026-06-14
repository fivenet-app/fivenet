package qualifications

import (
	"context"
	"errors"

	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	resqualifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications"
	qualificationsexam "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications/exam"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pbqualifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/qualifications"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) ListQualificationsResults(
	ctx context.Context,
	req *pbqualifications.ListQualificationsResultsRequest,
	userInfo *userinfo.UserInfo,
	where mysql.BoolExpression,
	includePhoneNumber bool,
) (*pbqualifications.ListQualificationsResultsResponse, error) {
	tUser := table.FivenetUser.AS("user")
	tCreator := tUser.AS("creator")
	tQuali := tQuali.AS("qualification_short")

	condition := tQualiResult.DeletedAt.IS_NULL()
	if where != nil {
		condition = condition.AND(where)
	}

	countColumn := mysql.Expression(tQualiResult.QualificationID)
	if req.GetUserId() != 0 {
		condition = condition.AND(mysql.AND(
			tUser.Job.EQ(mysql.String(userInfo.GetJob())),
			tQualiResult.UserID.EQ(mysql.Int32(req.GetUserId())),
		))
	} else {
		if req.QualificationId == nil {
			condition = condition.AND(tUser.Job.EQ(mysql.String(userInfo.GetJob()))).
				AND(tQualiResult.UserID.EQ(mysql.Int32(userInfo.GetUserId())))
			countColumn = mysql.DISTINCT(tQualiResult.QualificationID)
		} else {
			countColumn = mysql.DISTINCT(tQualiResult.UserID)
		}
	}

	if len(req.GetStatus()) > 0 {
		statuses := []mysql.Expression{}
		for i := range req.GetStatus() {
			statuses = append(statuses, mysql.Int32(int32(req.GetStatus()[i])))
		}
		condition = condition.AND(tQualiResult.Status.IN(statuses...))
	}

	countStmt := tQualiResult.
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

	pag, limit := req.GetPagination().GetResponseWithPageSize(count.Total, QualificationsPageSize)
	resp := &pbqualifications.ListQualificationsResultsResponse{
		Pagination: pag,
		Results:    []*resqualifications.QualificationResult{},
	}
	if count.Total <= 0 {
		return resp, nil
	}

	orderBys := []mysql.OrderByClause{}
	if req.GetSort() != nil && len(req.GetSort().GetColumns()) > 0 {
		for _, sc := range req.GetSort().GetColumns() {
			var column mysql.Column
			switch sc.GetId() {
			case "status":
				column = tQualiResult.Status
			case "createdAt":
				fallthrough
			default:
				column = tQualiResult.CreatedAt
			}

			if sc.GetDesc() {
				orderBys = append(orderBys, column.DESC())
			} else {
				orderBys = append(orderBys, column.ASC())
			}
		}
	} else {
		orderBys = append(orderBys, tQualiResult.CreatedAt.DESC())
	}

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

	stmt := tQualiResult.
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
		OFFSET(req.GetPagination().GetOffset()).
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

	condition := tQualiResult.DeletedAt.IS_NULL()
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
