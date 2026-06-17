package access

import (
	"context"
	"database/sql"
	"errors"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

const subjectCleanupDeleteLimit int64 = 500

type ActorSubject struct {
	SubjectID        int64 `alias:"subject_id"`
	Specificity      int32 `alias:"specificity"`
	GradeSpecificity int32 `alias:"grade_specificity"`
}

type SubjectResolver struct {
	db *sql.DB
}

func NewSubjectResolver(db *sql.DB) *SubjectResolver {
	return &SubjectResolver{db: db}
}

func (r *SubjectResolver) ensureSubject(
	ctx context.Context,
	tx qrm.DB,
	subjectType SubjectType,
	find func() (int64, error),
	upsert func(subjectID int64) (int64, error),
) (int64, error) {
	if subjectID, err := find(); err == nil {
		return subjectID, nil
	} else if !errors.Is(err, qrm.ErrNoRows) {
		return 0, err
	}

	subjectID, err := createSubject(ctx, tx, subjectType)
	if err != nil {
		return 0, err
	}

	mappedID, err := upsert(subjectID)
	if err != nil {
		return 0, err
	}
	if mappedID != 0 && mappedID != subjectID {
		_ = deleteSubject(ctx, tx, subjectID)
		return mappedID, nil
	}

	return subjectID, nil
}

func lookupSubjectIDFromMapping(
	ctx context.Context,
	tx qrm.DB,
	subjectAlias string,
	subjectType SubjectType,
	mappingTable mysql.ReadableTable,
	subjectID mysql.IntegerExpression,
	where mysql.BoolExpression,
) (int64, error) {
	tSubjects := table.FivenetACLSubjects.AS(subjectAlias)
	stmt := mappingTable.
		SELECT(subjectID).
		FROM(mappingTable.INNER_JOIN(tSubjects,
			mysql.AND(
				tSubjects.ID.EQ(subjectID),
				tSubjects.SubjectType.EQ(mysql.Int16(int16(subjectType))),
			),
		)).
		WHERE(where).
		LIMIT(1)

	return querySubjectID(ctx, tx, stmt)
}

func (r *SubjectResolver) EnsureUserSubject(
	ctx context.Context,
	tx qrm.DB,
	userID int32,
) (int64, error) {
	return r.ensureSubject(ctx, tx, SubjectTypeUser, func() (int64, error) {
		return lookupSubjectIDFromMapping(
			ctx,
			tx,
			"user_subject",
			SubjectTypeUser,
			table.FivenetACLSubjectUsers,
			table.FivenetACLSubjectUsers.SubjectID,
			table.FivenetACLSubjectUsers.UserID.EQ(mysql.Int32(userID)),
		)
	}, func(subjectID int64) (int64, error) {
		return upsertUserSubject(ctx, tx, subjectID, userID)
	})
}

func (r *SubjectResolver) EnsureQualificationSubject(
	ctx context.Context,
	tx qrm.DB,
	qualificationID int64,
) (int64, error) {
	return r.ensureSubject(
		ctx,
		tx,
		SubjectTypeQualification,
		func() (int64, error) {
			return lookupSubjectIDFromMapping(
				ctx,
				tx,
				"qualification_subject",
				SubjectTypeQualification,
				table.FivenetACLSubjectQualifications,
				table.FivenetACLSubjectQualifications.SubjectID,
				table.FivenetACLSubjectQualifications.QualificationID.EQ(
					mysql.Int64(qualificationID),
				),
			)
		},
		func(subjectID int64) (int64, error) {
			return upsertQualificationSubject(ctx, tx, subjectID, qualificationID)
		},
	)
}

func (r *SubjectResolver) EnsureJobGradeSubject(
	ctx context.Context,
	tx qrm.DB,
	job string,
	minimumGrade int32,
) (int64, error) {
	return r.ensureSubject(
		ctx,
		tx,
		SubjectTypeJobGrade,
		func() (int64, error) {
			return lookupSubjectIDFromMapping(
				ctx,
				tx,
				"job_grade_subject",
				SubjectTypeJobGrade,
				table.FivenetACLSubjectJobGradeScopes,
				table.FivenetACLSubjectJobGradeScopes.SubjectID,
				mysql.AND(
					table.FivenetACLSubjectJobGradeScopes.Job.EQ(mysql.String(job)),
					table.FivenetACLSubjectJobGradeScopes.MinimumGrade.EQ(
						mysql.Int32(minimumGrade),
					),
				),
			)
		},
		func(subjectID int64) (int64, error) {
			return upsertJobGradeSubject(ctx, tx, subjectID, job, minimumGrade)
		},
	)
}

func (r *SubjectResolver) ResolveActorSubjects(
	ctx context.Context,
	tx qrm.DB,
	userInfo *userinfo.UserInfo,
) ([]ActorSubject, error) {
	if userInfo == nil || userInfo.GetUserId() <= 0 {
		return nil, nil
	}

	tSubjectQualis := table.FivenetACLSubjectQualifications.AS("asq_resolve")
	tQualiResults := table.FivenetQualificationsResults.AS("qr_resolve")
	tSubjectJobGrade := table.FivenetACLSubjectJobGradeScopes.AS("asjg_resolve")
	tUserJobs := table.FivenetUserJobs.AS("uj_resolve")
	stmt := mysql.
		SELECT(
			table.FivenetACLSubjectUsers.SubjectID.AS("subject_id"),
			mysql.Int32(SubjectSpecificityUser).AS("specificity"),
			mysql.Int32(-1).AS("grade_specificity"),
		).
		FROM(table.FivenetACLSubjectUsers).
		WHERE(table.FivenetACLSubjectUsers.UserID.EQ(mysql.Int32(userInfo.GetUserId()))).
		UNION_ALL(
			mysql.
				SELECT(
					tSubjectQualis.SubjectID.AS("subject_id"),
					mysql.Int32(SubjectSpecificityQualification).AS("specificity"),
					mysql.Int32(-1).AS("grade_specificity"),
				).
				FROM(tSubjectQualis.
					INNER_JOIN(tQualiResults,
						mysql.AND(
							tQualiResults.QualificationID.EQ(tSubjectQualis.QualificationID),
							tQualiResults.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
							tQualiResults.DeletedAt.IS_NULL(),
							tQualiResults.Status.EQ(
								mysql.Int32(
									int32(qualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL),
								),
							),
						),
					),
				),
		).
		UNION_ALL(
			mysql.
				SELECT(
					tSubjectJobGrade.SubjectID.AS("subject_id"),
					mysql.Int32(SubjectSpecificityJobGrade).AS("specificity"),
					tSubjectJobGrade.MinimumGrade.AS("grade_specificity"),
				).
				FROM(tSubjectJobGrade.
					INNER_JOIN(tUserJobs,
						mysql.AND(
							tUserJobs.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
							tUserJobs.Job.EQ(tSubjectJobGrade.Job),
							tUserJobs.Grade.GT_EQ(tSubjectJobGrade.MinimumGrade),
						),
					),
				),
		).
		UNION_ALL(
			mysql.
				SELECT(
					tSubjectJobGrade.SubjectID.AS("subject_id"),
					mysql.Int32(SubjectSpecificityJobGrade).AS("specificity"),
					tSubjectJobGrade.MinimumGrade.AS("grade_specificity"),
				).
				FROM(tSubjectJobGrade).
				WHERE(mysql.AND(
					tSubjectJobGrade.Job.EQ(mysql.String(userInfo.GetJob())),
					tSubjectJobGrade.MinimumGrade.LT_EQ(mysql.Int32(userInfo.GetJobGrade())),
				)),
		)

	var dest []ActorSubject
	if err := stmt.QueryContext(ctx, tx, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest, nil
}

func (r *SubjectResolver) CleanupOrphanSubjects(ctx context.Context, tx qrm.DB) error {
	_, err := r.cleanupOrphanSubjectsStmt().ExecContext(ctx, tx)
	return err
}

func (r *SubjectResolver) cleanupOrphanSubjectsStmt() mysql.DeleteStatement {
	tSubjects := table.FivenetACLSubjects.AS("orphan_subject")
	tSubjectUsers := table.FivenetACLSubjectUsers.AS("orphan_subject_user")
	tSubjectQuals := table.FivenetACLSubjectQualifications.AS("orphan_subject_qual")
	tSubjectJobGrades := table.FivenetACLSubjectJobGradeScopes.AS("orphan_subject_job_grade")

	return tSubjects.
		DELETE().
		WHERE(mysql.AND(
			mysql.NOT(mysql.EXISTS(
				mysql.
					SELECT(mysql.Int(1)).
					FROM(tSubjectUsers).
					WHERE(tSubjectUsers.SubjectID.EQ(tSubjects.ID)),
			)),
			mysql.NOT(mysql.EXISTS(
				mysql.
					SELECT(mysql.Int(1)).
					FROM(tSubjectQuals).
					WHERE(tSubjectQuals.SubjectID.EQ(tSubjects.ID)),
			)),
			mysql.NOT(mysql.EXISTS(
				mysql.
					SELECT(mysql.Int(1)).
					FROM(tSubjectJobGrades).
					WHERE(tSubjectJobGrades.SubjectID.EQ(tSubjects.ID)),
			)),
		)).
		LIMIT(subjectCleanupDeleteLimit)
}

func (r *SubjectResolver) CleanupStaleJobGradeSubjects(ctx context.Context, tx qrm.DB) error {
	_, err := r.cleanupStaleJobGradeSubjectsStmt().ExecContext(ctx, tx)
	return err
}

func (r *SubjectResolver) cleanupStaleJobGradeSubjectsStmt() mysql.DeleteStatement {
	tSubjects := table.FivenetACLSubjects.AS("stale_job_grade_subject")
	tSubjectJobGrades := table.FivenetACLSubjectJobGradeScopes.AS("stale_job_grade_scope")
	tJobGrades := table.FivenetJobsGrades.AS("stale_job_grade")

	return tSubjects.
		DELETE().
		WHERE(mysql.AND(
			tSubjects.SubjectType.EQ(mysql.Int16(int16(SubjectTypeJobGrade))),
			mysql.EXISTS(
				mysql.
					SELECT(mysql.Int(1)).
					FROM(tSubjectJobGrades).
					WHERE(tSubjectJobGrades.SubjectID.EQ(tSubjects.ID)),
			),
			mysql.NOT(mysql.EXISTS(
				mysql.
					SELECT(mysql.Int(1)).
					FROM(tSubjectJobGrades.
						INNER_JOIN(tJobGrades,
							mysql.AND(
								tJobGrades.JobName.EQ(tSubjectJobGrades.Job),
								tJobGrades.Grade.EQ(tSubjectJobGrades.MinimumGrade),
							),
						),
					).
					WHERE(tSubjectJobGrades.SubjectID.EQ(tSubjects.ID)),
			)),
		)).
		LIMIT(subjectCleanupDeleteLimit)
}

func createSubject(ctx context.Context, tx qrm.DB, subjectType SubjectType) (int64, error) {
	stmt := table.FivenetACLSubjects.
		INSERT(table.FivenetACLSubjects.SubjectType).
		VALUES(int16(subjectType))

	res, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func deleteSubject(ctx context.Context, tx qrm.DB, subjectID int64) error {
	_, err := table.FivenetACLSubjects.
		DELETE().
		WHERE(table.FivenetACLSubjects.ID.EQ(mysql.Int64(subjectID))).
		LIMIT(1).
		ExecContext(ctx, tx)
	return err
}

func upsertUserSubject(
	ctx context.Context,
	tx qrm.DB,
	subjectID int64,
	userID int32,
) (int64, error) {
	stmt := table.FivenetACLSubjectUsers.
		INSERT(
			table.FivenetACLSubjectUsers.SubjectID,
			table.FivenetACLSubjectUsers.UserID,
		).
		VALUES(
			subjectID,
			userID,
		).
		ON_DUPLICATE_KEY_UPDATE(
			table.FivenetACLSubjectUsers.SubjectID.SET(
				mysql.RawInt("LAST_INSERT_ID(`subject_id`)"),
			),
		)

	return upsertSubjectMapping(ctx, tx, stmt)
}

func upsertQualificationSubject(
	ctx context.Context,
	tx qrm.DB,
	subjectID int64,
	qualificationID int64,
) (int64, error) {
	stmt := table.FivenetACLSubjectQualifications.
		INSERT(
			table.FivenetACLSubjectQualifications.SubjectID,
			table.FivenetACLSubjectQualifications.QualificationID,
		).
		VALUES(
			subjectID,
			qualificationID,
		).
		ON_DUPLICATE_KEY_UPDATE(
			table.FivenetACLSubjectQualifications.SubjectID.SET(
				mysql.RawInt("LAST_INSERT_ID(`subject_id`)"),
			),
		)

	return upsertSubjectMapping(ctx, tx, stmt)
}

func upsertJobGradeSubject(
	ctx context.Context,
	tx qrm.DB,
	subjectID int64,
	job string,
	minimumGrade int32,
) (int64, error) {
	stmt := table.FivenetACLSubjectJobGradeScopes.
		INSERT(
			table.FivenetACLSubjectJobGradeScopes.SubjectID,
			table.FivenetACLSubjectJobGradeScopes.Job,
			table.FivenetACLSubjectJobGradeScopes.MinimumGrade,
		).
		VALUES(
			subjectID,
			job,
			minimumGrade,
		).
		ON_DUPLICATE_KEY_UPDATE(
			table.FivenetACLSubjectJobGradeScopes.SubjectID.SET(
				mysql.RawInt("LAST_INSERT_ID(`subject_id`)"),
			),
		)

	return upsertSubjectMapping(ctx, tx, stmt)
}

func upsertSubjectMapping(
	ctx context.Context,
	tx qrm.DB,
	stmt mysql.InsertStatement,
) (int64, error) {
	res, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

type subjectIDDest struct {
	SubjectID int64 `alias:"subject_id"`
}

func querySubjectID(ctx context.Context, tx qrm.DB, stmt mysql.SelectStatement) (int64, error) {
	dest := subjectIDDest{}
	if err := stmt.QueryContext(ctx, tx, &dest); err != nil {
		return 0, err
	}
	if dest.SubjectID == 0 {
		return 0, qrm.ErrNoRows
	}
	return dest.SubjectID, nil
}
