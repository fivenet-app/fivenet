package qualifications

import (
	"context"
	"errors"
	"slices"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/qualifications"
	permscitizenstore "github.com/fivenet-app/fivenet/gen/go/proto/services/citizenstore/perms"
	errorsqualifications "github.com/fivenet-app/fivenet/gen/go/proto/services/qualifications/errors"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/pkg/perms"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var (
	tUser       = table.Users.AS("user")
	tCreator    = table.Users.AS("creator")
	tQJobAccess = table.FivenetQualificationsJobAccess
	tQReqs      = table.FivenetQualificationsRequirements.AS("qualificationrequirement")
)

func (s *Server) listQualificationsQuery(where jet.BoolExpression, onlyColumns jet.ProjectionList, userInfo *userinfo.UserInfo) jet.SelectStatement {
	wheres := []jet.BoolExpression{}
	if !userInfo.SuperUser {
		wheres = append(wheres,
			jet.AND(
				tQuali.DeletedAt.IS_NULL(),
				jet.OR(
					jet.AND(
						tQuali.CreatorID.EQ(jet.Int32(userInfo.UserId)),
						tQuali.CreatorJob.EQ(jet.String(userInfo.Job)),
					),
					jet.AND(
						tQJobAccess.Access.IS_NOT_NULL(),
						tQJobAccess.Access.GT(jet.Int32(int32(qualifications.AccessLevel_ACCESS_LEVEL_BLOCKED))),
					),
				),
			),
		)
	}

	// Select id of last result
	wheres = append(wheres, tQualiResults.ID.IS_NULL().OR(
		tQualiResults.ID.EQ(
			jet.RawInt("SELECT MAX(`qualificationresult`.`id`) FROM `fivenet_qualifications_results` AS `qualificationresult` WHERE (qualificationresult.qualification_id = qualification.id AND qualificationresult.deleted_at IS NULL AND qualificationresult.user_id = #userid)",
				jet.RawArgs{
					"#userid": userInfo.UserId,
				},
			),
		),
	))

	if where != nil {
		wheres = append(wheres, where)
	}

	var q jet.SelectStatement
	if onlyColumns != nil {
		q = tQuali.
			SELECT(
				onlyColumns,
			)
	} else {
		columns := jet.ProjectionList{
			tQuali.ID,
			tQuali.CreatedAt,
			tQuali.UpdatedAt,
			tQuali.Job,
			tQuali.Weight,
			tQuali.Closed,
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
		}

		if userInfo.SuperUser {
			columns = append(columns, tQuali.DeletedAt)
		}

		// Field Permission Check
		fieldsAttr, _ := s.ps.Attr(userInfo, permscitizenstore.CitizenStoreServicePerm, permscitizenstore.CitizenStoreServiceListCitizensPerm, permscitizenstore.CitizenStoreServiceListCitizensFieldsPermField)
		var fields perms.StringList
		if fieldsAttr != nil {
			fields = fieldsAttr.([]string)
		}

		if slices.Contains(fields, "PhoneNumber") {
			columns = append(columns, tCreator.PhoneNumber)
		}

		q = tQuali.SELECT(columns[0], columns[1:])
	}

	var tables jet.ReadableTable
	if !userInfo.SuperUser {
		tables = tQuali.
			LEFT_JOIN(tQJobAccess,
				tQJobAccess.QualificationID.EQ(tQuali.ID).
					AND(tQJobAccess.Job.EQ(jet.String(userInfo.Job))).
					AND(tQJobAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
			).
			LEFT_JOIN(tCreator,
				tQuali.CreatorID.EQ(tCreator.ID),
			).
			LEFT_JOIN(tQualiResults,
				tQualiResults.QualificationID.EQ(tQuali.ID).
					AND(tQualiResults.DeletedAt.IS_NULL()).
					AND(tQualiResults.UserID.EQ(jet.Int32(userInfo.UserId))),
			)
	} else {
		tables = tQuali.
			LEFT_JOIN(tCreator,
				tQuali.CreatorID.EQ(tCreator.ID),
			).
			LEFT_JOIN(tQualiResults,
				tQualiResults.QualificationID.EQ(tQuali.ID).
					AND(tQualiResults.DeletedAt.IS_NULL()).
					AND(tQualiResults.UserID.EQ(jet.Int32(userInfo.UserId))),
			)
	}

	return q.
		FROM(tables).
		WHERE(
			jet.AND(
				wheres...,
			),
		).
		ORDER_BY(
			tQuali.Weight.ASC(),
			tQuali.Abbreviation.ASC(),
			tQualiResults.ID.DESC(),
		)
}

func (s *Server) getQualificationQuery(qualificationId uint64, where jet.BoolExpression, onlyColumns jet.ProjectionList, userInfo *userinfo.UserInfo, selectContent bool) jet.SelectStatement {
	wheres := []jet.BoolExpression{jet.Bool(true)}
	if !userInfo.SuperUser {
		wheres = append(wheres,
			jet.AND(
				tQuali.DeletedAt.IS_NULL(),
				jet.OR(
					jet.AND(
						tQuali.CreatorID.EQ(jet.Int32(userInfo.UserId)),
						tQuali.CreatorJob.EQ(jet.String(userInfo.Job)),
					),
					jet.AND(
						tQJobAccess.Access.IS_NOT_NULL(),
						tQJobAccess.Access.GT(jet.Int32(int32(qualifications.AccessLevel_ACCESS_LEVEL_BLOCKED))),
					),
				),
			),
		)
	}

	// Select id of last result
	wheres = append(wheres, tQualiResults.ID.IS_NULL().OR(
		tQualiResults.ID.EQ(
			jet.RawInt("SELECT MAX(`qualificationresult`.`id`) FROM `fivenet_qualifications_results` AS `qualificationresult` WHERE (qualificationresult.qualification_id = #qualificationId AND qualificationresult.deleted_at IS NULL AND qualificationresult.user_id = #userid)",
				jet.RawArgs{
					"#qualificationId": qualificationId,
					"#userid":          userInfo.UserId,
				},
			),
		),
	))

	if where != nil {
		wheres = append(wheres, where)
	}

	var q jet.SelectStatement
	if onlyColumns != nil {
		q = tQuali.
			SELECT(
				onlyColumns,
			)
	} else {
		columns := jet.ProjectionList{
			tQuali.ID,
			tQuali.CreatedAt,
			tQuali.UpdatedAt,
			tQuali.Job,
			tQuali.Weight,
			tQuali.Closed,
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
			tQualiResults.ID,
			tQualiResults.QualificationID,
			tQualiResults.Status,
			tQualiResults.Score,
			tQualiResults.Summary,
			tQualiRequests.ApprovedAt,
			tQualiRequests.Status,
		}

		if selectContent {
			columns = append(columns, tQuali.Content)
		}

		if userInfo.SuperUser {
			columns = append(columns, tQuali.DeletedAt)
		}

		// Field Permission Check
		fieldsAttr, _ := s.ps.Attr(userInfo, permscitizenstore.CitizenStoreServicePerm, permscitizenstore.CitizenStoreServiceListCitizensPerm, permscitizenstore.CitizenStoreServiceListCitizensFieldsPermField)
		var fields perms.StringList
		if fieldsAttr != nil {
			fields = fieldsAttr.([]string)
		}

		if slices.Contains(fields, "PhoneNumber") {
			columns = append(columns, tCreator.PhoneNumber)
		}

		q = tQuali.SELECT(columns[0], columns[1:])
	}

	var tables jet.ReadableTable
	if !userInfo.SuperUser {
		tables = tQuali.
			LEFT_JOIN(tQJobAccess,
				tQJobAccess.QualificationID.EQ(tQuali.ID).
					AND(tQJobAccess.Job.EQ(jet.String(userInfo.Job))).
					AND(tQJobAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
			).
			LEFT_JOIN(tCreator,
				tQuali.CreatorID.EQ(tCreator.ID),
			).
			LEFT_JOIN(tQualiResults,
				tQualiResults.QualificationID.EQ(tQuali.ID).
					AND(tQualiResults.DeletedAt.IS_NULL()).
					AND(tQualiResults.UserID.EQ(jet.Int32(userInfo.UserId))),
			).
			LEFT_JOIN(tQualiRequests,
				tQualiRequests.QualificationID.EQ(tQuali.ID).
					AND(tQualiRequests.DeletedAt.IS_NULL()).
					AND(tQualiRequests.UserID.EQ(jet.Int32(userInfo.UserId))).
					AND(tQualiRequests.Status.NOT_EQ(jet.Int16(int16(qualifications.RequestStatus_REQUEST_STATUS_COMPLETED)))),
			)
	} else {
		tables = tQuali.
			LEFT_JOIN(tCreator,
				tQuali.CreatorID.EQ(tCreator.ID),
			).
			LEFT_JOIN(tQualiResults,
				tQualiResults.QualificationID.EQ(tQuali.ID).
					AND(tQualiResults.DeletedAt.IS_NULL()).
					AND(tQualiResults.UserID.EQ(jet.Int32(userInfo.UserId))),
			).
			LEFT_JOIN(tQualiRequests,
				tQualiRequests.QualificationID.EQ(tQuali.ID).
					AND(tQualiRequests.DeletedAt.IS_NULL()).
					AND(tQualiRequests.UserID.EQ(jet.Int32(userInfo.UserId))).
					AND(tQualiRequests.Status.NOT_EQ(jet.Int16(int16(qualifications.RequestStatus_REQUEST_STATUS_COMPLETED)))),
			)
	}

	return q.
		FROM(tables).
		WHERE(jet.AND(
			wheres...,
		)).
		ORDER_BY(
			tQuali.CreatedAt.DESC(),
			tQuali.UpdatedAt.DESC(),
		)
}

func (s *Server) getQualificationRequirements(ctx context.Context, qualificationId uint64) ([]*qualifications.QualificationRequirement, error) {
	tQuali := tQuali.AS("targetqualification")

	stmt := tQReqs.
		SELECT(
			tQReqs.ID,
			tQuali.ID,
			tQuali.Abbreviation,
			tQuali.Title,
			tQReqs.TargetQualificationID,
			tQualiResults.ID,
			tQualiResults.QualificationID,
			tQualiResults.Status,
			tQualiResults.Score,
			tQualiResults.Summary,
		).
		FROM(tQReqs.
			INNER_JOIN(tQuali,
				tQuali.ID.EQ(tQReqs.TargetQualificationID),
			).
			LEFT_JOIN(tQualiResults,
				tQualiResults.QualificationID.EQ(tQReqs.TargetQualificationID),
			),
		).
		WHERE(tQReqs.QualificationID.EQ(jet.Uint64(qualificationId)))

	var dest []*qualifications.QualificationRequirement
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, err
	}

	return dest, nil
}

func (s *Server) checkIfUserHasAccessToQuali(ctx context.Context, qualificationID uint64, userInfo *userinfo.UserInfo, access qualifications.AccessLevel) (bool, error) {
	out, err := s.checkIfUserHasAccessToQualiIDs(ctx, userInfo, access, qualificationID)
	return len(out) > 0, err
}

func (s *Server) checkIfUserHasAccessToQualiIDs(ctx context.Context, userInfo *userinfo.UserInfo, access qualifications.AccessLevel, qualificationIDs ...uint64) ([]uint64, error) {
	if len(qualificationIDs) == 0 {
		return qualificationIDs, nil
	}

	// Allow superusers access to any docs
	if userInfo.SuperUser {
		return qualificationIDs, nil
	}

	ids := make([]jet.Expression, len(qualificationIDs))
	for i := 0; i < len(qualificationIDs); i++ {
		ids[i] = jet.Uint64(qualificationIDs[i])
	}

	condition := jet.AND(
		tQuali.ID.IN(ids...),
		tQuali.DeletedAt.IS_NULL(),
		jet.OR(
			jet.AND(
				tQuali.CreatorID.EQ(jet.Int32(userInfo.UserId)),
				tQuali.CreatorJob.EQ(jet.String(userInfo.Job)),
			),
			jet.AND(
				tQJobAccess.Access.IS_NOT_NULL(),
				tQJobAccess.Access.GT_EQ(jet.Int32(int32(access))),
			),
		),
	)

	stmt := tQuali.
		SELECT(
			tQuali.ID,
		).
		FROM(
			tQuali.
				LEFT_JOIN(tQJobAccess,
					tQJobAccess.QualificationID.EQ(tQuali.ID).
						AND(tQJobAccess.Job.EQ(jet.String(userInfo.Job))).
						AND(tQJobAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
				),
		).
		WHERE(condition).
		GROUP_BY(tQuali.ID).
		ORDER_BY(tQuali.ID.DESC(), tQJobAccess.MinimumGrade)

	var dest struct {
		IDs []uint64 `alias:"qualification.id"`
	}
	if err := stmt.QueryContext(ctx, s.db, &dest.IDs); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest.IDs, nil
}

func (s *Server) getQualification(ctx context.Context, qualificationId uint64, condition jet.BoolExpression, userInfo *userinfo.UserInfo, selectContent bool) (*qualifications.Qualification, error) {
	var quali qualifications.Qualification

	stmt := s.getQualificationQuery(qualificationId, condition, nil, userInfo, selectContent)

	if err := stmt.QueryContext(ctx, s.db, &quali); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
	}

	if quali.Id == 0 {
		return nil, nil
	}

	reqs, err := s.getQualificationRequirements(ctx, qualificationId)
	if err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
	}
	quali.Requirements = reqs

	if quali.Creator != nil {
		s.enricher.EnrichJobInfoSafe(userInfo, quali.Creator)
	}

	request, err := s.getQualificationRequest(ctx, qualificationId, userInfo.UserId, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	quali.Request = request

	result, err := s.getQualificationResult(ctx, qualificationId, 0, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}
	quali.Result = result

	return &quali, nil
}

func (s *Server) getQualificationShort(ctx context.Context, qualificationId uint64, condition jet.BoolExpression, userInfo *userinfo.UserInfo) (*qualifications.QualificationShort, error) {
	var quali qualifications.Qualification

	stmt := s.getQualificationQuery(qualificationId, condition, nil, userInfo, false)

	if err := stmt.QueryContext(ctx, s.db, &quali); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
	}

	if quali.Id == 0 {
		return nil, nil
	}

	if quali.Creator != nil {
		s.enricher.EnrichJobInfoSafe(userInfo, quali.Creator)
	}

	return &qualifications.QualificationShort{
		Id:           quali.Id,
		CreatedAt:    quali.CreatedAt,
		UpdatedAt:    quali.UpdatedAt,
		DeletedAt:    quali.DeletedAt,
		Job:          quali.Job,
		Weight:       quali.Weight,
		Closed:       quali.Closed,
		Abbreviation: quali.Abbreviation,
		Title:        quali.Title,
		Description:  quali.Description,
		CreatorId:    quali.CreatorId,
		Creator:      quali.Creator,
		Requirements: quali.Requirements,
		Result:       quali.Result,
	}, nil
}

func (s *Server) checkRequirementsMetForQualification(ctx context.Context, qualificationId uint64, userId int32) (bool, error) {
	stmt := tQReqs.
		SELECT(
			tQReqs.TargetQualificationID.AS("qualiid"),
			tQualiResults.UserID.AS("userid"),
		).
		FROM(tQReqs.
			LEFT_JOIN(tQualiResults,
				tQualiResults.QualificationID.EQ(tQReqs.TargetQualificationID).
					AND(tQualiResults.DeletedAt.IS_NULL()).
					AND(tQualiResults.UserID.EQ(jet.Int32(userId))).
					AND(tQualiResults.Status.EQ(jet.Int16(int16(qualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL)))),
			),
		).
		WHERE(
			tQReqs.QualificationID.EQ(jet.Uint64(qualificationId)),
		)

	var dest []*struct {
		QualiID uint64
		UserID  int32
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

	// Remove all requirements which the user has fullfilled
	dest = slices.DeleteFunc(dest, func(s *struct {
		QualiID uint64
		UserID  int32
	}) bool {
		return s.UserID > 0
	})

	return len(dest) == 0, nil
}

func (s *Server) handleQualificationRequirementsChanges(ctx context.Context, tx qrm.DB, qualificationId uint64, reqs []*qualifications.QualificationRequirement) error {
	current, err := s.getQualificationRequirements(ctx, qualificationId)
	if err != nil {
		return err
	}

	toCreate, toDelete := s.compareQualificationRequirements(current, reqs)

	tQReqs := table.FivenetQualificationsRequirements

	for _, req := range toDelete {
		stmt := tQReqs.
			DELETE().
			WHERE(jet.AND(
				tQReqs.ID.EQ(jet.Uint64(req.Id)),
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
				req.TargetQualificationId,
			)

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) compareQualificationRequirements(current, in []*qualifications.QualificationRequirement) (toCreate []*qualifications.QualificationRequirement, toDelete []*qualifications.QualificationRequirement) {
	if current == nil {
		return in, toDelete
	}

	if len(current) == 0 {
		toCreate = in
	} else {
		foundTracker := []int{}
		for _, cq := range current {
			var found *qualifications.QualificationRequirement
			var foundIdx int
			for i, qj := range in {
				if cq.TargetQualificationId != qj.TargetQualificationId {
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

	return
}
