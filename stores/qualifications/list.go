package qualificationsstore

import (
	"context"
	"errors"
	"slices"

	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	resqualifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications"
	qualificationsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications/access"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pbqualifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/qualifications"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) ListQualifications(
	ctx context.Context,
	opts ListQualificationsOptions,
	userInfo *userinfo.UserInfo,
	includePhoneNumber bool,
) (*pbqualifications.ListQualificationsResponse, error) {
	tCreator := table.FivenetUser.AS("creator")
	searchCondition := mysql.Bool(true)

	if search := dbutils.PrepareForLikeSearch(opts.Search); search != "" {
		searchCondition = searchCondition.AND(mysql.OR(
			table.FivenetQualifications.Abbreviation.LIKE(mysql.String(search)),
			table.FivenetQualifications.Title.LIKE(mysql.String(search)),
		))
	}

	includeDeleted := userInfo != nil && userInfo.GetSuperuser()
	userID := int32(0)
	if userInfo != nil {
		userID = userInfo.GetUserId()
	}
	visibleIDs := s.access.VisibleIDsByConditionQuery(
		userInfo,
		int32(qualificationsaccess.AccessLevel_ACCESS_LEVEL_VIEW),
		includeDeleted,
		searchCondition,
	)
	ctes := visibleIDs.CTEs
	visibleQualiID := mysql.IntegerColumn("id").From(visibleIDs.Table)

	tQualiResultLatest := table.FivenetQualificationsResults.AS("qualificationresult_latest")

	// Select id of last result for this user.
	lastResultID := tQualiResultLatest.
		SELECT(mysql.MAX(tQualiResultLatest.ID)).
		FROM(tQualiResultLatest).
		WHERE(mysql.AND(
			tQualiResultLatest.QualificationID.EQ(tQuali.ID),
			tQualiResultLatest.DeletedAt.IS_NULL(),
			tQualiResultLatest.UserID.EQ(mysql.Int32(userID)),
		)).
		LIMIT(1)
	lastResultFilter := tQualiResult.ID.IS_NULL().OR(tQualiResult.ID.EQ(mysql.IntExp(lastResultID)))
	visibleQuery := tQuali.
		SELECT(mysql.COUNT(mysql.DISTINCT(tQuali.ID)).AS("data_count.total")).
		FROM(
			visibleIDs.Table.
				INNER_JOIN(tQuali,
					tQuali.ID.EQ(visibleQualiID),
				).
				LEFT_JOIN(tCreator, tQuali.CreatorID.EQ(tCreator.ID)).
				LEFT_JOIN(tQualiResult, mysql.AND(
					tQualiResult.QualificationID.EQ(tQuali.ID),
					tQualiResult.DeletedAt.IS_NULL(),
					tQualiResult.UserID.EQ(mysql.Int32(userID)),
				)),
		).
		WHERE(lastResultFilter)

	var countStmt mysql.Statement = visibleQuery
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
	resp := &pbqualifications.ListQualificationsResponse{
		Pagination:     pag,
		Qualifications: []*resqualifications.Qualification{},
	}
	if count.Total <= 0 {
		return resp, nil
	}

	stmt := s.listQualificationsQuery(
		opts,
		userInfo,
		visibleIDs.Table,
		ctes,
		lastResultFilter,
		includePhoneNumber,
		limit,
	)

	if err := stmt.QueryContext(ctx, s.db, &resp.Qualifications); err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Store) GetQualification(
	ctx context.Context,
	qualificationId int64,
	userInfo *userinfo.UserInfo,
	selectContent bool,
	includePhoneNumber bool,
) (*resqualifications.Qualification, error) {
	wheres := []mysql.BoolExpression{
		tQuali.ID.EQ(mysql.Int64(qualificationId)),
	}

	tQualiResultLatest := table.FivenetQualificationsResults.AS("qualificationresult_latest")

	// Select id of last result for this user.
	lastResultID := tQualiResultLatest.
		SELECT(mysql.MAX(tQualiResultLatest.ID)).
		FROM(tQualiResultLatest).
		WHERE(mysql.AND(
			tQualiResultLatest.QualificationID.EQ(tQuali.ID),
			tQualiResultLatest.DeletedAt.IS_NULL(),
			tQualiResultLatest.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
		)).
		LIMIT(1)
	wheres = append(
		wheres,
		tQualiResult.ID.IS_NULL().OR(tQualiResult.ID.EQ(mysql.IntExp(lastResultID))),
	)

	stmt := s.getQualificationQuery(
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
		userInfo,
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
	opts ListQualificationsOptions,
	userInfo *userinfo.UserInfo,
	visibleIDs mysql.SelectTable,
	ctes []mysql.CommonTableExpression,
	lastResultFilter mysql.BoolExpression,
	includePhoneNumber bool,
	limit int64,
) mysql.Statement {
	tCreator := table.FivenetUser.AS("creator")
	visibleQualiID := mysql.IntegerColumn("id").From(visibleIDs)
	userID := int32(0)
	if userInfo != nil {
		userID = userInfo.GetUserId()
	}

	orderBys := append(
		[]mysql.OrderByClause{tQuali.Draft.ASC()},
		s.qualificationSorter.Build(opts.Sort)...)

	columns := []mysql.Projection{
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
	}
	if includePhoneNumber {
		columns = append(columns, tCreator.PhoneNumber)
	}

	var stmt mysql.Statement = tQuali.
		SELECT(
			columns[0],
			columns[1:]...,
		).
		FROM(
			visibleIDs.
				INNER_JOIN(tQuali,
					tQuali.ID.EQ(visibleQualiID),
				).
				LEFT_JOIN(tCreator, tQuali.CreatorID.EQ(tCreator.ID)).
				LEFT_JOIN(tQualiResult, mysql.AND(
					tQualiResult.QualificationID.EQ(tQuali.ID),
					tQualiResult.DeletedAt.IS_NULL(),
					tQualiResult.UserID.EQ(mysql.Int32(userID)),
				)),
		).
		WHERE(lastResultFilter).
		ORDER_BY(orderBys...).
		OFFSET(opts.Pagination.GetOffset()).
		LIMIT(limit)

	if len(ctes) > 0 {
		stmt = mysql.WITH(ctes...)(stmt)
	}

	return stmt
}

func (s *Store) getQualificationQuery(
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
	columns = append(columns, tQuali.DeletedAt)
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
	userInfo *userinfo.UserInfo,
	includePhoneNumber bool,
) (*resqualifications.QualificationShort, error) {
	quali, err := s.GetQualification(
		ctx,
		qualificationId,
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
			tQualiResultSuccess.UserID.AS("userid"),
		).
		FROM(tQualiReqs.LEFT_JOIN(tQualiResultSuccess, mysql.AND(
			tQualiResultSuccess.QualificationID.EQ(tQualiReqs.TargetQualificationID),
			tQualiResultSuccess.UserID.EQ(mysql.Int32(userId)),
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
