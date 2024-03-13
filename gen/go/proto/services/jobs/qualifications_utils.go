package jobs

import (
	"context"
	"errors"
	"slices"

	jobs "github.com/galexrt/fivenet/gen/go/proto/resources/jobs"
	permscitizenstore "github.com/galexrt/fivenet/gen/go/proto/services/citizenstore/perms"
	errorsjobs "github.com/galexrt/fivenet/gen/go/proto/services/jobs/errors"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	"github.com/galexrt/fivenet/pkg/grpc/errswrap"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var (
	tCreator    = tUser.AS("creator")
	tQJobAccess = table.FivenetJobsQualificationsJobAccess
	tQReqs      = table.FivenetJobsQualificationsRequirements.AS("qualificationrequirement")
)

func (s *Server) listQualificationsQuery(where jet.BoolExpression, onlyColumns jet.ProjectionList, userInfo *userinfo.UserInfo) jet.SelectStatement {
	wheres := []jet.BoolExpression{}
	if !userInfo.SuperUser {
		wheres = []jet.BoolExpression{
			jet.AND(
				tQuali.DeletedAt.IS_NULL(),
				jet.OR(
					tQuali.CreatorID.EQ(jet.Int32(userInfo.UserId)),
					jet.AND(
						tQJobAccess.Access.IS_NOT_NULL(),
						tQJobAccess.Access.NOT_EQ(jet.Int32(int32(jobs.AccessLevel_ACCESS_LEVEL_BLOCKED))),
					),
				),
			),
		}
	}

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
			tQuali.Content,
			tQuali.CreatorID,
			tCreator.ID,
			tCreator.Identifier,
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
					AND(tQualiResults.UserID.EQ(jet.Int32(userInfo.UserId))),
			)
	} else {
		tables = tQuali.
			LEFT_JOIN(tCreator,
				tQuali.CreatorID.EQ(tCreator.ID),
			).
			LEFT_JOIN(tQualiResults,
				tQualiResults.QualificationID.EQ(tQuali.ID).
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
			tQualiResults.Status.ASC(),
			tQuali.Weight.ASC(),
			tQuali.Abbreviation.ASC(),
		)
}

func (s *Server) getQualificationQuery(where jet.BoolExpression, onlyColumns jet.ProjectionList, userInfo *userinfo.UserInfo) jet.SelectStatement {
	var wheres []jet.BoolExpression
	if !userInfo.SuperUser {
		wheres = []jet.BoolExpression{
			jet.AND(
				tQuali.DeletedAt.IS_NULL(),
				jet.OR(
					tQuali.CreatorID.EQ(jet.Int32(userInfo.UserId)),
					jet.AND(
						tQJobAccess.Access.IS_NOT_NULL(),
						tQJobAccess.Access.NOT_EQ(jet.Int32(int32(jobs.AccessLevel_ACCESS_LEVEL_BLOCKED))),
					),
				),
			),
		}
	}

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
			tQuali.Content,
			tQuali.CreatorID,
			tCreator.ID,
			tCreator.Identifier,
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
					AND(tQualiResults.UserID.EQ(jet.Int32(userInfo.UserId))),
			)
	} else {
		tables = tQuali.
			LEFT_JOIN(tCreator,
				tQuali.CreatorID.EQ(tCreator.ID),
			).
			LEFT_JOIN(tQualiResults,
				tQualiResults.QualificationID.EQ(tQuali.ID).
					AND(tQualiResults.UserID.EQ(jet.Int32(userInfo.UserId))),
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

func (s *Server) getQualificationRequirements(ctx context.Context, qualificationId uint64) ([]*jobs.QualificationRequirement, error) {
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

	var dest []*jobs.QualificationRequirement
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, err
	}

	return dest, nil
}

func (s *Server) checkIfUserHasAccessToQuali(ctx context.Context, qualificationID uint64, userInfo *userinfo.UserInfo, access jobs.AccessLevel) (bool, error) {
	out, err := s.checkIfUserHasAccessToQualiIDs(ctx, userInfo, access, qualificationID)
	return len(out) > 0, err
}

func (s *Server) checkIfUserHasAccessToQualis(ctx context.Context, userInfo *userinfo.UserInfo, access jobs.AccessLevel, qualificationIDs ...uint64) (bool, error) {
	out, err := s.checkIfUserHasAccessToQualiIDs(ctx, userInfo, access, qualificationIDs...)
	return len(out) == len(qualificationIDs), err
}

func (s *Server) checkIfUserHasAccessToQualiIDs(ctx context.Context, userInfo *userinfo.UserInfo, access jobs.AccessLevel, qualificationIDs ...uint64) ([]uint64, error) {
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
			tQuali.CreatorID.EQ(jet.Int32(userInfo.UserId)),
			tQuali.CreatorJob.EQ(jet.String(userInfo.Job)),
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

func (s *Server) getQualification(ctx context.Context, qualificationId uint64, condition jet.BoolExpression, userInfo *userinfo.UserInfo) (*jobs.Qualification, error) {
	var quali jobs.Qualification

	stmt := s.getQualificationQuery(condition, nil, userInfo)

	if err := stmt.QueryContext(ctx, s.db, &quali); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
		}
	}

	reqs, err := s.getQualificationRequirements(ctx, qualificationId)
	if err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
		}
	}
	quali.Requirements = reqs

	if quali.Creator != nil {
		s.enricher.EnrichJobInfo(quali.Creator)
	}

	return &quali, nil
}
