package qualifications

import (
	"context"
	"errors"
	"slices"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/qualifications"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/userinfo"
	permscitizens "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/citizens/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsqualifications "github.com/fivenet-app/fivenet/v2025/services/qualifications/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var (
	tQAccess = table.FivenetQualificationsAccess
	tQReqs   = table.FivenetQualificationsRequirements.AS("qualification_requirement")
)

func (s *Server) listQualificationsQuery(
	where mysql.BoolExpression,
	onlyColumns mysql.ProjectionList,
	userInfo *userinfo.UserInfo,
) mysql.SelectStatement {
	tCreator := tables.User().AS("creator")

	wheres := []mysql.BoolExpression{}
	if !userInfo.GetSuperuser() {
		accessExists := mysql.EXISTS(
			mysql.
				SELECT(mysql.Int(1)).
				FROM(tQAccess).
				WHERE(mysql.AND(
					tQAccess.TargetID.EQ(tQuali.ID),
					tQAccess.Access.IS_NOT_NULL(),
					tQAccess.Access.GT_EQ(
						mysql.Int32(int32(qualifications.AccessLevel_ACCESS_LEVEL_VIEW)),
					),
					mysql.AND(
						tQAccess.Job.EQ(mysql.String(userInfo.GetJob())),
						tQAccess.MinimumGrade.LT_EQ(mysql.Int32(userInfo.GetJobGrade())),
					),
				),
				),
		)

		wheres = append(wheres,
			mysql.AND(
				tQuali.DeletedAt.IS_NULL(),
				mysql.OR(
					tQuali.Public.IS_TRUE(),
					mysql.AND(
						tQuali.CreatorID.EQ(mysql.Int32(userInfo.GetUserId())),
						tQuali.CreatorJob.EQ(mysql.String(userInfo.GetJob())),
					),
					accessExists,
				),
			),
		)
	}

	// Select id of last result
	wheres = append(wheres, tQualiResults.ID.IS_NULL().OR(
		tQualiResults.ID.EQ(
			mysql.RawInt(
				"SELECT MAX(`qualificationresult`.`id`) FROM `fivenet_qualifications_results` AS `qualificationresult` WHERE (qualificationresult.qualification_id = qualification.id AND qualificationresult.deleted_at IS NULL AND qualificationresult.user_id = #userid)",
				mysql.RawArgs{
					"#userid": userInfo.GetUserId(),
				},
			),
		),
	))

	if where != nil {
		wheres = append(wheres, where)
	}

	var columns mysql.ProjectionList
	if onlyColumns != nil {
		columns = append(columns, onlyColumns...)
	} else {
		columns = append(columns,
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
			tCreator.PhoneNumber,
			tQuali.CreatorJob,
			tQualiResults.ID,
			tQualiResults.QualificationID,
			tQualiResults.Status,
			tQualiResults.Score,
			tQualiResults.Summary,
			tQualiResults.CreatorID,
		)

		if userInfo.GetSuperuser() {
			columns = append(columns, tQuali.DeletedAt)
		}

		// Field Permission Check
		fields, _ := s.perms.AttrStringList(userInfo, permscitizens.CitizensServicePerm, permscitizens.CitizensServiceListCitizensPerm, permscitizens.CitizensServiceListCitizensFieldsPermField)

		if fields.Contains("PhoneNumber") {
			columns = append(columns, tCreator.PhoneNumber)
		}
	}

	return tQuali.
		SELECT(columns[0], columns[1:]...).
		FROM(
			tQuali.
				LEFT_JOIN(tCreator,
					tQuali.CreatorID.EQ(tCreator.ID),
				).
				LEFT_JOIN(tQualiResults,
					tQualiResults.QualificationID.EQ(tQuali.ID).
						AND(tQualiResults.DeletedAt.IS_NULL()).
						AND(tQualiResults.UserID.EQ(mysql.Int32(userInfo.GetUserId()))),
				),
		).
		WHERE(mysql.AND(
			wheres...,
		)).
		ORDER_BY(
			tQuali.Abbreviation.ASC(),
			tQualiResults.ID.DESC(),
		)
}

func (s *Server) getQualificationQuery(
	qualificationId int64,
	where mysql.BoolExpression,
	onlyColumns mysql.ProjectionList,
	userInfo *userinfo.UserInfo,
	selectContent bool,
) mysql.SelectStatement {
	tCreator := tables.User().AS("creator")

	wheres := []mysql.BoolExpression{
		tQuali.ID.EQ(mysql.Int64(qualificationId)),
	}
	if !userInfo.GetSuperuser() {
		accessExists := mysql.EXISTS(
			mysql.
				SELECT(mysql.Int(1)).
				FROM(tQAccess).
				WHERE(mysql.AND(
					tQAccess.TargetID.EQ(tQuali.ID),
					tQAccess.Access.IS_NOT_NULL(),
					tQAccess.Access.GT_EQ(
						mysql.Int32(int32(qualifications.AccessLevel_ACCESS_LEVEL_VIEW)),
					),
					mysql.AND(
						tQAccess.Job.EQ(mysql.String(userInfo.GetJob())),
						tQAccess.MinimumGrade.LT_EQ(mysql.Int32(userInfo.GetJobGrade())),
					),
				),
				),
		)

		wheres = append(wheres,
			mysql.AND(
				tQuali.DeletedAt.IS_NULL(),
				mysql.OR(
					tQuali.Public.IS_TRUE(),
					mysql.AND(
						tQuali.CreatorID.EQ(mysql.Int32(userInfo.GetUserId())),
						tQuali.CreatorJob.EQ(mysql.String(userInfo.GetJob())),
					),
					accessExists,
				),
			),
		)
	}

	// Select id of last result
	wheres = append(wheres, tQualiResults.ID.IS_NULL().OR(
		tQualiResults.ID.EQ(
			mysql.RawInt(
				"SELECT MAX(`qualificationresult`.`id`) FROM `fivenet_qualifications_results` AS `qualificationresult` WHERE (qualificationresult.qualification_id = #qualificationId AND qualificationresult.deleted_at IS NULL AND qualificationresult.user_id = #userid)",
				mysql.RawArgs{
					"#qualificationId": qualificationId,
					"#userid":          userInfo.GetUserId(),
				},
			),
		),
	))

	if where != nil {
		wheres = append(wheres, where)
	}

	var columns mysql.ProjectionList
	if onlyColumns != nil {
		columns = append(columns, onlyColumns...)
	} else {
		columns = append(columns,
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
			tQualiResults.ID,
			tQualiResults.QualificationID,
			tQualiResults.Status,
			tQualiResults.Score,
			tQualiResults.Summary,
			tQualiResults.CreatorID,
			tQualiRequests.ApprovedAt,
			tQualiRequests.Status,
		)

		if selectContent {
			columns = append(columns, tQuali.Content)
		}

		if userInfo.GetSuperuser() {
			columns = append(columns, tQuali.DeletedAt)
		}

		// Field Permission Check
		fields, _ := s.perms.AttrStringList(userInfo, permscitizens.CitizensServicePerm, permscitizens.CitizensServiceListCitizensPerm, permscitizens.CitizensServiceListCitizensFieldsPermField)

		if fields.Contains("PhoneNumber") {
			columns = append(columns, tCreator.PhoneNumber)
		}
	}

	return tQuali.
		SELECT(
			columns[0],
			columns[1:]...,
		).
		FROM(
			tQuali.
				LEFT_JOIN(tCreator,
					tQuali.CreatorID.EQ(tCreator.ID),
				).
				LEFT_JOIN(tQualiResults,
					tQualiResults.QualificationID.EQ(tQuali.ID).
						AND(tQualiResults.DeletedAt.IS_NULL()).
						AND(tQualiResults.UserID.EQ(mysql.Int32(userInfo.GetUserId()))),
				).
				LEFT_JOIN(tQualiRequests,
					tQualiRequests.QualificationID.EQ(tQuali.ID).
						AND(tQualiRequests.DeletedAt.IS_NULL()).
						AND(tQualiRequests.UserID.EQ(mysql.Int32(userInfo.GetUserId()))).
						AND(tQualiRequests.Status.NOT_EQ(mysql.Int32(int32(qualifications.RequestStatus_REQUEST_STATUS_COMPLETED)))),
				),
		).
		WHERE(mysql.AND(
			wheres...,
		)).
		ORDER_BY(
			tQuali.CreatedAt.DESC(),
			tQuali.UpdatedAt.DESC(),
		)
}

func (s *Server) getQualificationRequirements(
	ctx context.Context,
	qualificationId int64,
) ([]*qualifications.QualificationRequirement, error) {
	tQuali := tQuali.AS("target_qualification")

	stmt := tQReqs.
		SELECT(
			tQReqs.ID,
			tQReqs.CreatedAt,
			tQReqs.TargetQualificationID,
			tQuali.ID,
			tQuali.Abbreviation,
			tQuali.Title,
		).
		FROM(tQReqs.
			INNER_JOIN(tQuali,
				tQuali.ID.EQ(tQReqs.TargetQualificationID),
			),
		).
		WHERE(tQReqs.QualificationID.EQ(mysql.Int64(qualificationId)))

	var dest []*qualifications.QualificationRequirement
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, err
	}

	return dest, nil
}

func (s *Server) getQualification(
	ctx context.Context,
	qualificationId int64,
	condition mysql.BoolExpression,
	userInfo *userinfo.UserInfo,
	selectContent bool,
) (*qualifications.Qualification, error) {
	var quali qualifications.Qualification

	stmt := s.getQualificationQuery(qualificationId, condition, nil, userInfo, selectContent)

	if err := stmt.QueryContext(ctx, s.db, &quali); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
	}

	if quali.GetId() == 0 {
		return nil, nil
	}

	reqs, err := s.getQualificationRequirements(ctx, qualificationId)
	if err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
	}
	quali.Requirements = reqs

	if quali.GetCreator() != nil {
		s.enricher.EnrichJobInfoSafe(userInfo, quali.GetCreator())
	}

	request, err := s.getQualificationRequest(ctx, qualificationId, userInfo.GetUserId(), userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	quali.Request = request

	result, err := s.getQualificationResult(ctx, qualificationId, 0, nil, userInfo, 0)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	quali.Result = result

	return &quali, nil
}

func (s *Server) getQualificationShort(
	ctx context.Context,
	qualificationId int64,
	condition mysql.BoolExpression,
	userInfo *userinfo.UserInfo,
) (*qualifications.QualificationShort, error) {
	var quali qualifications.Qualification

	stmt := s.getQualificationQuery(qualificationId, condition, nil, userInfo, false)

	if err := stmt.QueryContext(ctx, s.db, &quali); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
	}

	if quali.GetId() == 0 {
		return nil, nil
	}

	if quali.GetCreator() != nil {
		s.enricher.EnrichJobInfoSafe(userInfo, quali.GetCreator())
	}

	return &qualifications.QualificationShort{
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

func (s *Server) checkRequirementsMetForQualification(
	ctx context.Context,
	qualificationId int64,
	userId int32,
) (bool, error) {
	stmt := tQReqs.
		SELECT(
			tQReqs.TargetQualificationID.AS("qualification_id"),
			tQualiResults.UserID.AS("userid"),
		).
		FROM(tQReqs.
			LEFT_JOIN(tQualiResults,
				tQualiResults.QualificationID.EQ(tQReqs.TargetQualificationID).
					AND(tQualiResults.DeletedAt.IS_NULL()).
					AND(tQualiResults.UserID.EQ(mysql.Int32(userId))).
					AND(tQualiResults.Status.EQ(mysql.Int32(int32(qualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL)))),
			),
		).
		WHERE(
			tQReqs.QualificationID.EQ(mysql.Int64(qualificationId)),
		)

	var dest []*struct {
		QualificationID int64
		UserID          int32
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return false, err
		}
	}

	// No requirements on qualification
	if len(dest) == 0 {
		return true, nil
	}

	// Remove all requirements which the user has fulfilled
	dest = slices.DeleteFunc(dest, func(s *struct {
		QualificationID int64
		UserID          int32
	},
	) bool {
		return s.UserID > 0
	})

	return len(dest) == 0, nil
}

func (s *Server) handleQualificationRequirementsChanges(
	ctx context.Context,
	tx qrm.DB,
	qualificationId int64,
	reqs []*qualifications.QualificationRequirement,
) error {
	current, err := s.getQualificationRequirements(ctx, qualificationId)
	if err != nil {
		return err
	}

	toCreate, toDelete := s.compareQualificationRequirements(current, reqs)

	tQReqs := table.FivenetQualificationsRequirements

	for _, req := range toDelete {
		stmt := tQReqs.
			DELETE().
			WHERE(mysql.AND(
				tQReqs.ID.EQ(mysql.Int64(req.GetId())),
			)).
			LIMIT(1)

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return err
		}
	}

	for _, req := range toCreate {
		stmt := tQReqs.
			INSERT(
				tQReqs.QualificationID,
				tQReqs.TargetQualificationID,
			).
			VALUES(
				qualificationId,
				req.GetTargetQualificationId(),
			)

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) compareQualificationRequirements(
	current, in []*qualifications.QualificationRequirement,
) ([]*qualifications.QualificationRequirement, []*qualifications.QualificationRequirement) {
	toCreate := []*qualifications.QualificationRequirement{}
	toDelete := []*qualifications.QualificationRequirement{}

	if current == nil {
		return in, toDelete
	}

	if len(current) == 0 {
		if len(in) == 0 {
			toDelete = current
		} else {
			toCreate = in
		}
	} else {
		foundTracker := []int{}
		for _, cq := range current {
			var found *qualifications.QualificationRequirement
			var foundIdx int
			for i, qj := range in {
				if cq.GetTargetQualificationId() != qj.GetTargetQualificationId() {
					continue
				}
				found = qj
				foundIdx = i
				break
			}
			// No match in incoming requirement, needs to be deleted
			if found == nil {
				toDelete = append(toDelete, cq)
				continue
			}

			foundTracker = append(foundTracker, foundIdx)
		}

		for i, uj := range in {
			idx := slices.Index(foundTracker, i)
			if idx == -1 {
				toCreate = append(toCreate, uj)
			}
		}
	}

	return toCreate, toDelete
}
