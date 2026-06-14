package qualificationsstore

import (
	"context"
	"errors"
	"slices"

	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	resqualifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pbqualifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/qualifications"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) ListQualifications(
	ctx context.Context,
	req *pbqualifications.ListQualificationsRequest,
	userInfo *userinfo.UserInfo,
	where mysql.BoolExpression,
	includePhoneNumber bool,
) (*pbqualifications.ListQualificationsResponse, error) {
	tCreator := table.FivenetUser.AS("creator")

	wheres := []mysql.BoolExpression{}
	if where != nil {
		wheres = append(wheres, where)
	}

	if search := dbutils.PrepareForLikeSearch(req.GetSearch()); search != "" {
		wheres = append(wheres, mysql.OR(
			tQuali.Abbreviation.LIKE(mysql.String(search)),
			tQuali.Title.LIKE(mysql.String(search)),
		))
	}

	// Select id of last result for this user.
	wheres = append(wheres, tQualiResult.ID.IS_NULL().OR(
		tQualiResult.ID.EQ(
			mysql.RawInt(
				"SELECT MAX(`qualificationresult`.`id`) FROM `fivenet_qualifications_results` AS `qualificationresult` WHERE (qualificationresult.qualification_id = qualification.id AND qualificationresult.deleted_at IS NULL AND qualificationresult.user_id = #userid)",
				mysql.RawArgs{
					"#userid": userInfo.GetUserId(),
				},
			),
		),
	))

	countStmt := tQuali.
		SELECT(mysql.COUNT(mysql.DISTINCT(tQuali.ID)).AS("data_count.total")).
		FROM(
			tQuali.
				LEFT_JOIN(tCreator, tQuali.CreatorID.EQ(tCreator.ID)).
				LEFT_JOIN(tQualiResult, mysql.AND(
					tQualiResult.QualificationID.EQ(tQuali.ID),
					tQualiResult.DeletedAt.IS_NULL(),
					tQualiResult.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
				)),
		).
		WHERE(mysql.AND(wheres...))

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	pag, limit := req.GetPagination().GetResponseWithPageSize(count.Total, QualificationsPageSize)
	resp := &pbqualifications.ListQualificationsResponse{
		Pagination:     pag,
		Qualifications: []*resqualifications.Qualification{},
	}
	if count.Total <= 0 {
		return resp, nil
	}

	stmt := s.listQualificationsQuery(req, userInfo, wheres, includePhoneNumber, limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Qualifications); err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Store) GetQualification(
	ctx context.Context,
	qualificationId int64,
	where mysql.BoolExpression,
	userInfo *userinfo.UserInfo,
	selectContent bool,
	includePhoneNumber bool,
) (*resqualifications.Qualification, error) {
	wheres := []mysql.BoolExpression{
		tQuali.ID.EQ(mysql.Int64(qualificationId)),
	}
	if where != nil {
		wheres = append(wheres, where)
	}

	// Select id of last result for this user.
	wheres = append(wheres, tQualiResult.ID.IS_NULL().OR(
		tQualiResult.ID.EQ(
			mysql.RawInt(
				"SELECT MAX(`qualificationresult`.`id`) FROM `fivenet_qualifications_results` AS `qualificationresult` WHERE (qualificationresult.qualification_id = #qualificationId AND qualificationresult.deleted_at IS NULL AND qualificationresult.user_id = #userid)",
				mysql.RawArgs{
					"#qualificationId": qualificationId,
					"#userid":          userInfo.GetUserId(),
				},
			),
		),
	))

	stmt := s.getQualificationQuery(
		qualificationId,
		wheres,
		userInfo,
		selectContent,
		includePhoneNumber,
	)

	var quali resqualifications.Qualification
	if err := stmt.QueryContext(ctx, s.db, &quali); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}
	if quali.GetId() == 0 {
		return nil, nil
	}

	reqs, err := s.GetQualificationRequirements(ctx, qualificationId)
	if err != nil {
		return nil, err
	}
	quali.Requirements = reqs

	request, err := s.GetQualificationRequest(
		ctx,
		qualificationId,
		userInfo.GetUserId(),
		includePhoneNumber,
	)
	if err != nil {
		return nil, err
	}
	quali.Request = request

	result, err := s.GetQualificationResult(
		ctx,
		qualificationId,
		0,
		nil,
		userInfo,
		userInfo.GetUserId(),
		includePhoneNumber,
	)
	if err != nil {
		return nil, err
	}
	quali.Result = result

	return &quali, nil
}

func (s *Store) listQualificationsQuery(
	req *pbqualifications.ListQualificationsRequest,
	userInfo *userinfo.UserInfo,
	wheres []mysql.BoolExpression,
	includePhoneNumber bool,
	limit int64,
) mysql.SelectStatement {
	tCreator := table.FivenetUser.AS("creator")

	orderBys := []mysql.OrderByClause{tQuali.Draft.ASC()}
	if req.GetSort() != nil && len(req.GetSort().GetColumns()) > 0 {
		for _, sc := range req.GetSort().GetColumns() {
			var column mysql.Column
			switch sc.GetId() {
			case "abbreviation":
				column = tQuali.Abbreviation
			case "id":
				fallthrough
			default:
				column = tQualiResult.ID
			}

			if sc.GetDesc() {
				orderBys = append(orderBys, column.DESC())
			} else {
				orderBys = append(orderBys, column.ASC())
			}
		}
	} else {
		orderBys = append(orderBys, tQualiResult.ID.DESC())
	}

	return tQuali.
		SELECT(
			tQuali.ID,
			tQuali.CreatedAt,
			tQuali.UpdatedAt,
			tQuali.Job,
			tQuali.Weight,
			tQuali.Closed,
			tQuali.Draft,
			tQuali.Public,
			tQuali.Abbreviation,
			tQuali.Title,
			tQuali.Description,
			tQuali.ExamMode,
			tQuali.ExamSettings,
			tQuali.CreatorID,
			tCreator.ID,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Dateofbirth,
			tQuali.CreatorJob,
			tQualiResult.ID,
			tQualiResult.QualificationID,
			tQualiResult.Status,
			tQualiResult.Score,
			tQualiResult.Summary,
			tQualiResult.CreatorID,
		).
		FROM(
			tQuali.
				LEFT_JOIN(tCreator, tQuali.CreatorID.EQ(tCreator.ID)).
				LEFT_JOIN(tQualiResult, mysql.AND(
					tQualiResult.QualificationID.EQ(tQuali.ID),
					tQualiResult.DeletedAt.IS_NULL(),
					tQualiResult.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
				)),
		).
		WHERE(mysql.AND(wheres...)).
		ORDER_BY(orderBys...).
		OFFSET(req.GetPagination().GetOffset()).
		LIMIT(limit)
}

func (s *Store) getQualificationQuery(
	qualificationId int64,
	wheres []mysql.BoolExpression,
	userInfo *userinfo.UserInfo,
	selectContent bool,
	includePhoneNumber bool,
) mysql.SelectStatement {
	tCreator := table.FivenetUser.AS("creator")

	columns := mysql.ProjectionList{
		tQuali.ID,
		tQuali.CreatedAt,
		tQuali.UpdatedAt,
		tQuali.Job,
		tQuali.Weight,
		tQuali.Closed,
		tQuali.Draft,
		tQuali.Public,
		tQuali.Abbreviation,
		tQuali.Title,
		tQuali.Description,
		tQuali.ExamMode,
		tQuali.ExamSettings,
		tQuali.CreatorID,
		tCreator.ID,
		tCreator.Job,
		tCreator.JobGrade,
		tCreator.Firstname,
		tCreator.Lastname,
		tCreator.Dateofbirth,
		tQuali.CreatorJob,
		tQuali.DiscordSyncEnabled,
		tQuali.DiscordSettings,
		tQuali.LabelSyncEnabled,
		tQuali.LabelSyncFormat,
		tQualiResult.ID,
		tQualiResult.QualificationID,
		tQualiResult.Status,
		tQualiResult.Score,
		tQualiResult.Summary,
		tQualiResult.CreatorID,
	}

	if selectContent {
		columns = append(columns, tQuali.Content)
	}
	if userInfo.GetSuperuser() {
		columns = append(columns, tQuali.DeletedAt)
	}
	if includePhoneNumber {
		columns = append(columns, tCreator.PhoneNumber)
	}

	return tQuali.
		SELECT(columns[0], columns[1:]...).
		FROM(
			tQuali.
				LEFT_JOIN(tCreator, tQuali.CreatorID.EQ(tCreator.ID)).
				LEFT_JOIN(tQualiResult, mysql.AND(
					tQualiResult.QualificationID.EQ(tQuali.ID),
					tQualiResult.DeletedAt.IS_NULL(),
					tQualiResult.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
				)),
		).
		WHERE(mysql.AND(wheres...)).
		ORDER_BY(tQuali.CreatedAt.DESC(), tQuali.UpdatedAt.DESC())
}

func (s *Store) GetQualificationShort(
	ctx context.Context,
	qualificationId int64,
	where mysql.BoolExpression,
	userInfo *userinfo.UserInfo,
	includePhoneNumber bool,
) (*resqualifications.QualificationShort, error) {
	quali, err := s.GetQualification(
		ctx,
		qualificationId,
		where,
		userInfo,
		false,
		includePhoneNumber,
	)
	if err != nil {
		return nil, err
	}
	if quali == nil {
		return nil, nil
	}

	return &resqualifications.QualificationShort{
		Id:           quali.GetId(),
		CreatedAt:    quali.GetCreatedAt(),
		UpdatedAt:    quali.GetUpdatedAt(),
		DeletedAt:    quali.GetDeletedAt(),
		Job:          quali.GetJob(),
		Weight:       quali.GetWeight(),
		Closed:       quali.GetClosed(),
		Draft:        quali.GetDraft(),
		Public:       quali.GetPublic(),
		Abbreviation: quali.GetAbbreviation(),
		Title:        quali.GetTitle(),
		Description:  quali.Description,
		CreatorId:    quali.CreatorId,
		Creator:      quali.GetCreator(),
		ExamMode:     quali.GetExamMode(),
		ExamSettings: quali.GetExamSettings(),
		Requirements: quali.GetRequirements(),
		Result:       quali.GetResult(),
	}, nil
}

func (s *Store) GetQualificationRequirements(
	ctx context.Context,
	qualificationId int64,
) ([]*resqualifications.QualificationRequirement, error) {
	tQuali := tQuali.AS("target_qualification")

	stmt := tQualiReqs.
		SELECT(
			tQualiReqs.ID,
			tQualiReqs.CreatedAt,
			tQualiReqs.TargetQualificationID,
			tQuali.ID,
			tQuali.Abbreviation,
			tQuali.Title,
		).
		FROM(tQualiReqs.INNER_JOIN(tQuali, tQuali.ID.EQ(tQualiReqs.TargetQualificationID))).
		WHERE(tQualiReqs.QualificationID.EQ(mysql.Int64(qualificationId)))

	var dest []*resqualifications.QualificationRequirement
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, err
	}

	return dest, nil
}

func (s *Store) CheckRequirementsMetForQualification(
	ctx context.Context,
	qualificationId int64,
	userId int32,
) (bool, error) {
	stmt := tQualiReqs.
		SELECT(
			tQualiReqs.TargetQualificationID.AS("qualification_id"),
			tQualiResult.UserID.AS("userid"),
		).
		FROM(tQualiReqs.LEFT_JOIN(tQualiResult, mysql.AND(
			tQualiResult.QualificationID.EQ(tQualiReqs.TargetQualificationID),
			tQualiResult.DeletedAt.IS_NULL(),
			tQualiResult.UserID.EQ(mysql.Int32(userId)),
			tQualiResult.Status.EQ(
				mysql.Int32(int32(resqualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL)),
			),
		))).
		WHERE(tQualiReqs.QualificationID.EQ(mysql.Int64(qualificationId)))

	var dest []*struct {
		QualificationID int64
		UserID          int32
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return false, err
		}
	}
	if len(dest) == 0 {
		return true, nil
	}

	dest = slices.DeleteFunc(dest, func(s *struct {
		QualificationID int64
		UserID          int32
	},
	) bool {
		return s.UserID > 0
	})

	return len(dest) == 0, nil
}
