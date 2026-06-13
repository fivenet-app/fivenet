package access

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

type SubjectType int16

const (
	SubjectTypeUser          SubjectType = 1
	SubjectTypeQualification SubjectType = 2
	SubjectTypeJobGrade      SubjectType = 3
)

type AccessEffect int8

const (
	AccessEffectDeny  AccessEffect = 0
	AccessEffectAllow AccessEffect = 1
)

const (
	SubjectSpecificityUser          int32 = 500
	SubjectSpecificityQualification int32 = 300
	SubjectSpecificityJobGrade      int32 = 200
)

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

func (r *SubjectResolver) EnsureUserSubject(
	ctx context.Context,
	tx qrm.DB,
	userID int32,
) (int64, error) {
	return ensureSubject(
		ctx,
		tx,
		SubjectTypeUser,
		func() (int64, error) {
			tSubjects := table.FivenetAclSubjects.AS("user_subject")
			stmt := table.FivenetAclSubjectUsers.
				SELECT(table.FivenetAclSubjectUsers.SubjectID).
				FROM(table.FivenetAclSubjectUsers.
					INNER_JOIN(tSubjects,
						mysql.AND(
							tSubjects.ID.EQ(table.FivenetAclSubjectUsers.SubjectID),
							tSubjects.SubjectType.EQ(mysql.Int16(int16(SubjectTypeUser))),
						),
					),
				).
				WHERE(table.FivenetAclSubjectUsers.UserID.EQ(mysql.Int32(userID))).
				LIMIT(1)

			return querySubjectID(ctx, tx, stmt)
		},
		func(subjectID int64) error {
			stmt := table.FivenetAclSubjectUsers.
				INSERT(
					table.FivenetAclSubjectUsers.SubjectID,
					table.FivenetAclSubjectUsers.UserID,
				).
				VALUES(subjectID, userID)

			_, err := stmt.ExecContext(ctx, tx)
			return err
		},
	)
}

func (r *SubjectResolver) EnsureQualificationSubject(
	ctx context.Context,
	tx qrm.DB,
	qualificationID int64,
) (int64, error) {
	return ensureSubject(
		ctx,
		tx,
		SubjectTypeQualification,
		func() (int64, error) {
			tSubjects := table.FivenetAclSubjects.AS("qualification_subject")
			stmt := table.FivenetAclSubjectQualifications.
				SELECT(table.FivenetAclSubjectQualifications.SubjectID).
				FROM(table.FivenetAclSubjectQualifications.
					INNER_JOIN(tSubjects,
						mysql.AND(
							tSubjects.ID.EQ(table.FivenetAclSubjectQualifications.SubjectID),
							tSubjects.SubjectType.EQ(mysql.Int16(int16(SubjectTypeQualification))),
						),
					),
				).
				WHERE(table.FivenetAclSubjectQualifications.QualificationID.EQ(mysql.Int64(qualificationID))).
				LIMIT(1)

			return querySubjectID(ctx, tx, stmt)
		},
		func(subjectID int64) error {
			stmt := table.FivenetAclSubjectQualifications.
				INSERT(
					table.FivenetAclSubjectQualifications.SubjectID,
					table.FivenetAclSubjectQualifications.QualificationID,
				).
				VALUES(subjectID, qualificationID)

			_, err := stmt.ExecContext(ctx, tx)
			return err
		},
	)
}

func (r *SubjectResolver) EnsureJobGradeSubject(
	ctx context.Context,
	tx qrm.DB,
	job string,
	minimumGrade int32,
) (int64, error) {
	return ensureSubject(
		ctx,
		tx,
		SubjectTypeJobGrade,
		func() (int64, error) {
			tSubjects := table.FivenetAclSubjects.AS("job_grade_subject")
			stmt := table.FivenetAclSubjectJobGradeScopes.
				SELECT(table.FivenetAclSubjectJobGradeScopes.SubjectID).
				FROM(table.FivenetAclSubjectJobGradeScopes.
					INNER_JOIN(tSubjects,
						mysql.AND(
							tSubjects.ID.EQ(table.FivenetAclSubjectJobGradeScopes.SubjectID),
							tSubjects.SubjectType.EQ(mysql.Int16(int16(SubjectTypeJobGrade))),
						),
					),
				).
				WHERE(mysql.AND(
					table.FivenetAclSubjectJobGradeScopes.Job.EQ(mysql.String(job)),
					table.FivenetAclSubjectJobGradeScopes.MinimumGrade.EQ(mysql.Int32(minimumGrade)),
				)).
				LIMIT(1)

			return querySubjectID(ctx, tx, stmt)
		},
		func(subjectID int64) error {
			stmt := table.FivenetAclSubjectJobGradeScopes.
				INSERT(
					table.FivenetAclSubjectJobGradeScopes.SubjectID,
					table.FivenetAclSubjectJobGradeScopes.Job,
					table.FivenetAclSubjectJobGradeScopes.MinimumGrade,
				).
				VALUES(subjectID, job, minimumGrade)

			_, err := stmt.ExecContext(ctx, tx)
			return err
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

	var dest []ActorSubject
	if err := appendActorSubjects(ctx, tx, &dest, mysql.SELECT(
		table.FivenetAclSubjectUsers.SubjectID.AS("subject_id"),
		mysql.Int32(SubjectSpecificityUser).AS("specificity"),
		mysql.Int32(-1).AS("grade_specificity"),
	).
		FROM(table.FivenetAclSubjectUsers).
		WHERE(table.FivenetAclSubjectUsers.UserID.EQ(mysql.Int32(userInfo.GetUserId())))); err != nil {
		return nil, err
	}

	tSubjectQualis := table.FivenetAclSubjectQualifications.AS("asq_resolve")
	tQualiResults := table.FivenetQualificationsResults.AS("qr_resolve")
	if err := appendActorSubjects(ctx, tx, &dest, mysql.SELECT(
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
					tQualiResults.Status.EQ(mysql.Int32(int32(qualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL))),
				),
			),
		)); err != nil {
		return nil, err
	}

	tSubjectJobGrade := table.FivenetAclSubjectJobGradeScopes.AS("asjg_resolve")
	tUserJobs := table.FivenetUserJobs.AS("uj_resolve")
	if err := appendActorSubjects(ctx, tx, &dest, mysql.SELECT(
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
		)); err != nil {
		return nil, err
	}

	if err := appendActorSubjects(ctx, tx, &dest, mysql.SELECT(
		tSubjectJobGrade.SubjectID.AS("subject_id"),
		mysql.Int32(SubjectSpecificityJobGrade).AS("specificity"),
		tSubjectJobGrade.MinimumGrade.AS("grade_specificity"),
	).
		FROM(tSubjectJobGrade).
		WHERE(mysql.AND(
			tSubjectJobGrade.Job.EQ(mysql.String(userInfo.GetJob())),
			tSubjectJobGrade.MinimumGrade.LT_EQ(mysql.Int32(userInfo.GetJobGrade())),
		))); err != nil {
		return nil, err
	}

	return dest, nil
}

func appendActorSubjects(
	ctx context.Context,
	tx qrm.DB,
	dest *[]ActorSubject,
	stmt mysql.SelectStatement,
) error {
	var out []ActorSubject
	if err := stmt.QueryContext(ctx, tx, &out); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil
		}
		return err
	}

	*dest = append(*dest, out...)
	return nil
}

func (r *SubjectResolver) CleanupOrphanSubjects(ctx context.Context, tx qrm.DB) error {
	_, err := tx.ExecContext(ctx, `
DELETE s
FROM fivenet_acl_subjects s
LEFT JOIN fivenet_acl_subject_users su ON su.subject_id = s.id
LEFT JOIN fivenet_acl_subject_qualifications sq ON sq.subject_id = s.id
LEFT JOIN fivenet_acl_subject_job_grade_scopes sj ON sj.subject_id = s.id
WHERE su.subject_id IS NULL
  AND sq.subject_id IS NULL
  AND sj.subject_id IS NULL`)
	return err
}

func (r *SubjectResolver) CleanupStaleJobGradeSubjects(ctx context.Context, tx qrm.DB) error {
	_, err := tx.ExecContext(ctx, `
DELETE s
FROM fivenet_acl_subjects s
JOIN fivenet_acl_subject_job_grade_scopes sj ON sj.subject_id = s.id
LEFT JOIN fivenet_jobs_grades jg ON jg.job_name = sj.job AND jg.grade = sj.minimum_grade
WHERE s.subject_type = ?
  AND jg.job_name IS NULL`, int16(SubjectTypeJobGrade))
	return err
}

func ensureSubject(
	ctx context.Context,
	tx qrm.DB,
	subjectType SubjectType,
	find func() (int64, error),
	insertMapping func(subjectID int64) error,
) (int64, error) {
	subjectID, err := find()
	if err == nil {
		return subjectID, nil
	}
	if !errors.Is(err, qrm.ErrNoRows) {
		return 0, err
	}

	stmt := table.FivenetAclSubjects.
		INSERT(table.FivenetAclSubjects.SubjectType).
		VALUES(int16(subjectType))

	res, err := stmt.ExecContext(ctx, tx)
	if err != nil {
		return 0, err
	}

	subjectID, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}

	if err = insertMapping(subjectID); err != nil {
		if existingID, findErr := find(); findErr == nil {
			return existingID, nil
		}
		if dbutils.IsDuplicateError(err) {
			if cleanupErr := cleanupOrphanSubjectMappings(ctx, tx); cleanupErr != nil {
				return 0, cleanupErr
			}
			if retryErr := insertMapping(subjectID); retryErr == nil {
				return subjectID, nil
			}
			if existingID, findErr := find(); findErr == nil {
				return existingID, nil
			}
		}
		return 0, err
	}

	return subjectID, nil
}

func cleanupOrphanSubjectMappings(ctx context.Context, tx qrm.DB) error {
	if _, err := tx.ExecContext(ctx, `
DELETE su
FROM fivenet_acl_subject_users su
LEFT JOIN fivenet_acl_subjects s ON s.id = su.subject_id
WHERE s.id IS NULL`); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `
DELETE sq
FROM fivenet_acl_subject_qualifications sq
LEFT JOIN fivenet_acl_subjects s ON s.id = sq.subject_id
WHERE s.id IS NULL`); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `
DELETE sj
FROM fivenet_acl_subject_job_grade_scopes sj
LEFT JOIN fivenet_acl_subjects s ON s.id = sj.subject_id
WHERE s.id IS NULL`); err != nil {
		return err
	}

	return nil
}

type subjectIDDest struct {
	SubjectID int64 `alias:"subject_id"`
}

func querySubjectID(ctx context.Context, tx qrm.DB, stmt mysql.SelectStatement) (int64, error) {
	dest := subjectIDDest{}
	if err := stmt.QueryContext(ctx, tx, &dest); err != nil {
		return 0, err
	}
	return dest.SubjectID, nil
}

type SubjectAccessColumns struct {
	BaseAccessColumns
	SubjectID mysql.ColumnInteger
	Effect    mysql.ColumnInteger
}

type SubjectTargetTableColumns struct {
	ID         mysql.ColumnInteger
	DeletedAt  mysql.ColumnTimestamp
	Public     mysql.ColumnBool
	CreatorID  mysql.ColumnInteger
	CreatorJob mysql.ColumnString
}

type SubjectObjectAccess struct {
	db                 *sql.DB
	targetTable        mysql.Table
	targetTableColumns *SubjectTargetTableColumns
	accessTable        mysql.Table
	accessColumns      *SubjectAccessColumns
}

type canAccessIdsHelper struct {
	IDs []int64 `alias:"id"`
}

func NewSubjectObjectAccess(
	db *sql.DB,
	targetTable mysql.Table,
	targetColumns *SubjectTargetTableColumns,
	accessTable mysql.Table,
	accessColumns *SubjectAccessColumns,
) *SubjectObjectAccess {
	return &SubjectObjectAccess{
		db:                 db,
		targetTable:        targetTable,
		targetTableColumns: targetColumns,
		accessTable:        accessTable,
		accessColumns:      accessColumns,
	}
}

func NewDocumentsSubjectObjectAccess(db *sql.DB) *SubjectObjectAccess {
	return NewSubjectObjectAccess(
		db,
		table.FivenetDocuments,
		&SubjectTargetTableColumns{
			ID:         table.FivenetDocuments.ID,
			DeletedAt:  table.FivenetDocuments.DeletedAt,
			Public:     table.FivenetDocuments.Public,
			CreatorID:  table.FivenetDocuments.CreatorID,
			CreatorJob: table.FivenetDocuments.CreatorJob,
		},
		table.FivenetDocumentsAccess,
		&SubjectAccessColumns{
			BaseAccessColumns: BaseAccessColumns{
				ID:       table.FivenetDocumentsAccess.ID,
				TargetID: table.FivenetDocumentsAccess.TargetID,
				Access:   table.FivenetDocumentsAccess.Access,
			},
			SubjectID: table.FivenetDocumentsAccess.SubjectID,
			Effect:    table.FivenetDocumentsAccess.Effect,
		},
	)
}

func NewDocumentTemplatesSubjectObjectAccess(db *sql.DB) *SubjectObjectAccess {
	return NewSubjectObjectAccess(
		db,
		table.FivenetDocumentsTemplates,
		&SubjectTargetTableColumns{
			ID:         table.FivenetDocumentsTemplates.ID,
			DeletedAt:  table.FivenetDocumentsTemplates.DeletedAt,
			CreatorJob: table.FivenetDocumentsTemplates.CreatorJob,
		},
		table.FivenetDocumentsTemplatesAccess,
		&SubjectAccessColumns{
			BaseAccessColumns: BaseAccessColumns{
				ID:       table.FivenetDocumentsTemplatesAccess.ID,
				TargetID: table.FivenetDocumentsTemplatesAccess.TargetID,
				Access:   table.FivenetDocumentsTemplatesAccess.Access,
			},
			SubjectID: table.FivenetDocumentsTemplatesAccess.SubjectID,
			Effect:    table.FivenetDocumentsTemplatesAccess.Effect,
		},
	)
}

func NewDocumentStampsSubjectObjectAccess(db *sql.DB) *SubjectObjectAccess {
	return NewSubjectObjectAccess(
		db,
		table.FivenetDocumentsStamps,
		&SubjectTargetTableColumns{
			ID:        table.FivenetDocumentsStamps.ID,
			DeletedAt: table.FivenetDocumentsStamps.DeletedAt,
		},
		table.FivenetDocumentsStampsAccess,
		&SubjectAccessColumns{
			BaseAccessColumns: BaseAccessColumns{
				ID:       table.FivenetDocumentsStampsAccess.ID,
				TargetID: table.FivenetDocumentsStampsAccess.TargetID,
				Access:   table.FivenetDocumentsStampsAccess.Access,
			},
			SubjectID: table.FivenetDocumentsStampsAccess.SubjectID,
			Effect:    table.FivenetDocumentsStampsAccess.Effect,
		},
	)
}

func NewCalendarSubjectObjectAccess(db *sql.DB) *SubjectObjectAccess {
	return NewSubjectObjectAccess(
		db,
		table.FivenetCalendar,
		&SubjectTargetTableColumns{
			ID:         table.FivenetCalendar.ID,
			DeletedAt:  table.FivenetCalendar.DeletedAt,
			Public:     table.FivenetCalendar.Public,
			CreatorID:  table.FivenetCalendar.CreatorID,
			CreatorJob: table.FivenetCalendar.CreatorJob,
		},
		table.FivenetCalendarAccess,
		&SubjectAccessColumns{
			BaseAccessColumns: BaseAccessColumns{
				ID:       table.FivenetCalendarAccess.ID,
				TargetID: table.FivenetCalendarAccess.TargetID,
				Access:   table.FivenetCalendarAccess.Access,
			},
			SubjectID: table.FivenetCalendarAccess.SubjectID,
			Effect:    table.FivenetCalendarAccess.Effect,
		},
	)
}

func NewWikiPageSubjectObjectAccess(db *sql.DB) *SubjectObjectAccess {
	return NewSubjectObjectAccess(
		db,
		table.FivenetWikiPages,
		&SubjectTargetTableColumns{
			ID:         table.FivenetWikiPages.ID,
			DeletedAt:  table.FivenetWikiPages.DeletedAt,
			Public:     table.FivenetWikiPages.Public,
			CreatorID:  table.FivenetWikiPages.CreatorID,
			CreatorJob: table.FivenetWikiPages.Job,
		},
		table.FivenetWikiPagesAccess,
		&SubjectAccessColumns{
			BaseAccessColumns: BaseAccessColumns{
				ID:       table.FivenetWikiPagesAccess.ID,
				TargetID: table.FivenetWikiPagesAccess.TargetID,
				Access:   table.FivenetWikiPagesAccess.Access,
			},
			SubjectID: table.FivenetWikiPagesAccess.SubjectID,
			Effect:    table.FivenetWikiPagesAccess.Effect,
		},
	)
}

func NewCitizenLabelsSubjectObjectAccess(db *sql.DB) *SubjectObjectAccess {
	return NewSubjectObjectAccess(
		db,
		table.FivenetUserLabelsJob,
		&SubjectTargetTableColumns{
			ID:         table.FivenetUserLabelsJob.ID,
			DeletedAt:  table.FivenetUserLabelsJob.DeletedAt,
			CreatorJob: table.FivenetUserLabelsJob.Job,
		},
		table.FivenetUserLabelsJobJobAccess,
		&SubjectAccessColumns{
			BaseAccessColumns: BaseAccessColumns{
				ID:       table.FivenetUserLabelsJobJobAccess.ID,
				TargetID: table.FivenetUserLabelsJobJobAccess.TargetID,
				Access:   table.FivenetUserLabelsJobJobAccess.Access,
			},
			SubjectID: table.FivenetUserLabelsJobJobAccess.SubjectID,
			Effect:    table.FivenetUserLabelsJobJobAccess.Effect,
		},
	)
}

func NewQualificationsSubjectObjectAccess(db *sql.DB) *SubjectObjectAccess {
	return NewSubjectObjectAccess(
		db,
		table.FivenetQualifications,
		&SubjectTargetTableColumns{
			ID:         table.FivenetQualifications.ID,
			DeletedAt:  table.FivenetQualifications.DeletedAt,
			Public:     table.FivenetQualifications.Public,
			CreatorID:  table.FivenetQualifications.CreatorID,
			CreatorJob: table.FivenetQualifications.CreatorJob,
		},
		table.FivenetQualificationsAccess,
		&SubjectAccessColumns{
			BaseAccessColumns: BaseAccessColumns{
				ID:       table.FivenetQualificationsAccess.ID,
				TargetID: table.FivenetQualificationsAccess.TargetID,
				Access:   table.FivenetQualificationsAccess.Access,
			},
			SubjectID: table.FivenetQualificationsAccess.SubjectID,
			Effect:    table.FivenetQualificationsAccess.Effect,
		},
	)
}

func NewMailerEmailsSubjectObjectAccess(db *sql.DB) *SubjectObjectAccess {
	return NewSubjectObjectAccess(
		db,
		table.FivenetMailerEmails,
		&SubjectTargetTableColumns{
			ID:        table.FivenetMailerEmails.ID,
			DeletedAt: table.FivenetMailerEmails.DeletedAt,
			CreatorID: table.FivenetMailerEmails.UserID,
		},
		table.FivenetMailerEmailsAccess,
		&SubjectAccessColumns{
			BaseAccessColumns: BaseAccessColumns{
				ID:       table.FivenetMailerEmailsAccess.ID,
				TargetID: table.FivenetMailerEmailsAccess.TargetID,
				Access:   table.FivenetMailerEmailsAccess.Access,
			},
			SubjectID: table.FivenetMailerEmailsAccess.SubjectID,
			Effect:    table.FivenetMailerEmailsAccess.Effect,
		},
	)
}

func NewCentrumUnitsSubjectObjectAccess(db *sql.DB) *SubjectObjectAccess {
	return NewSubjectObjectAccess(
		db,
		table.FivenetCentrumUnits,
		&SubjectTargetTableColumns{
			ID:        table.FivenetCentrumUnits.ID,
			DeletedAt: table.FivenetCentrumUnits.DeletedAt,
		},
		table.FivenetCentrumUnitsAccess,
		&SubjectAccessColumns{
			BaseAccessColumns: BaseAccessColumns{
				ID:       table.FivenetCentrumUnitsAccess.ID,
				TargetID: table.FivenetCentrumUnitsAccess.TargetID,
				Access:   table.FivenetCentrumUnitsAccess.Access,
			},
			SubjectID: table.FivenetCentrumUnitsAccess.SubjectID,
			Effect:    table.FivenetCentrumUnitsAccess.Effect,
		},
	)
}

func (a *SubjectObjectAccess) CanUserAccessTarget(
	ctx context.Context,
	targetID int64,
	userInfo *userinfo.UserInfo,
	access int32,
) (bool, error) {
	out, err := a.CanUserAccessTargetIDs(ctx, userInfo, access, targetID)
	return len(out) > 0, err
}

func (a *SubjectObjectAccess) CanUserAccessTargetIDs(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
	access int32,
	targetIDs ...int64,
) ([]int64, error) {
	if len(targetIDs) == 0 {
		return nil, nil
	}
	if userInfo.GetSuperuser() {
		return targetIDs, nil
	}

	stmt := a.VisibleIDsStatement(userInfo, access, targetIDs...)

	dest := &canAccessIdsHelper{}
	if err := stmt.QueryContext(ctx, a.db, &dest.IDs); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest.IDs, nil
}

func (a *SubjectObjectAccess) VisibleIDsStatement(
	userInfo *userinfo.UserInfo,
	access int32,
	targetIDs ...int64,
) mysql.Statement {
	ids := make([]mysql.Expression, 0, len(targetIDs))
	for _, targetID := range targetIDs {
		ids = append(ids, mysql.Int64(targetID))
	}

	condition := a.targetTableColumns.ID.IN(ids...)
	if a.targetTableColumns.DeletedAt != nil {
		condition = condition.AND(a.targetTableColumns.DeletedAt.IS_NULL())
	}

	return a.visibleObjectsStatement(userInfo, access, condition, false, true)
}

func (a *SubjectObjectAccess) VisibleIDsByConditionStatement(
	userInfo *userinfo.UserInfo,
	access int32,
	condition mysql.BoolExpression,
) mysql.Statement {
	if a.targetTableColumns.DeletedAt != nil {
		condition = condition.AND(a.targetTableColumns.DeletedAt.IS_NULL())
	}

	return a.visibleObjectsStatement(userInfo, access, condition, false, true)
}

func (a *SubjectObjectAccess) ACLVisibleIDsByConditionStatement(
	userInfo *userinfo.UserInfo,
	access int32,
	condition mysql.BoolExpression,
) mysql.Statement {
	if a.targetTableColumns.DeletedAt != nil {
		condition = condition.AND(a.targetTableColumns.DeletedAt.IS_NULL())
	}

	return a.visibleObjectsStatement(userInfo, access, condition, false, false)
}

func (a *SubjectObjectAccess) CountVisibleByConditionStatement(
	userInfo *userinfo.UserInfo,
	access int32,
	condition mysql.BoolExpression,
) mysql.Statement {
	if a.targetTableColumns.DeletedAt != nil {
		condition = condition.AND(a.targetTableColumns.DeletedAt.IS_NULL())
	}

	return a.visibleObjectsStatement(userInfo, access, condition, true, true)
}

func (a *SubjectObjectAccess) visibleObjectsStatement(
	userInfo *userinfo.UserInfo,
	access int32,
	targetCondition mysql.BoolExpression,
	countOnly bool,
	includeImplicitAccess bool,
) mysql.Statement {
	actorSubjects := mysql.CTE("actor_subjects")
	matchingACL := mysql.CTE("matching_acl")
	winningSpecificity := mysql.CTE("winning_specificity")
	visibleObjects := mysql.CTE("visible_objects")

	actorSubjectID := mysql.IntegerColumn("subject_id").From(actorSubjects)
	actorSpecificity := mysql.IntegerColumn("specificity").From(actorSubjects)
	actorGradeSpecificity := mysql.IntegerColumn("grade_specificity").From(actorSubjects)

	matchingTargetID := mysql.IntegerColumn("target_id").From(matchingACL)
	matchingEffect := mysql.IntegerColumn("effect").From(matchingACL)
	matchingSpecificity := mysql.IntegerColumn("specificity").From(matchingACL)
	matchingGradeSpecificity := mysql.IntegerColumn("grade_specificity").From(matchingACL)
	matchingSpecificityRank := mysql.IntegerColumn("specificity_rank").From(matchingACL)

	winningTargetID := mysql.IntegerColumn("target_id").From(winningSpecificity)
	winningSpecificityCol := mysql.IntegerColumn("specificity").From(winningSpecificity)
	winningGradeSpecificity := mysql.IntegerColumn("grade_specificity").From(winningSpecificity)

	visibleID := mysql.IntegerColumn("id").From(visibleObjects)

	publicCondition := mysql.Bool(false)
	if a.targetTableColumns.Public != nil {
		publicCondition = a.targetTableColumns.Public.IS_TRUE()
	}
	creatorCondition := mysql.Bool(false)
	if a.targetTableColumns.CreatorID != nil {
		creatorCondition = a.targetTableColumns.CreatorID.EQ(mysql.Int32(userInfo.GetUserId()))
		if a.targetTableColumns.CreatorJob != nil {
			creatorCondition = creatorCondition.AND(
				a.targetTableColumns.CreatorJob.EQ(mysql.String(userInfo.GetJob())),
			)
		}
	} else if a.targetTableColumns.CreatorJob != nil {
		creatorCondition = a.targetTableColumns.CreatorJob.EQ(mysql.String(userInfo.GetJob()))
	}

	visibleSelect := mysql.SELECT(a.targetTableColumns.ID.AS("id")).
		FROM(a.targetTable.
			LEFT_JOIN(winningSpecificity,
				winningTargetID.EQ(a.targetTableColumns.ID),
			).
			LEFT_JOIN(matchingACL,
				mysql.AND(
					matchingTargetID.EQ(winningTargetID),
					matchingSpecificity.EQ(winningSpecificityCol),
					matchingGradeSpecificity.EQ(winningGradeSpecificity),
				),
			),
		).
		WHERE(targetCondition)

	groupBys := []mysql.GroupByClause{a.targetTableColumns.ID}
	if a.targetTableColumns.Public != nil {
		groupBys = append(groupBys, a.targetTableColumns.Public)
	}
	if a.targetTableColumns.CreatorID != nil {
		groupBys = append(groupBys, a.targetTableColumns.CreatorID)
	}
	if a.targetTableColumns.CreatorJob != nil {
		groupBys = append(groupBys, a.targetTableColumns.CreatorJob)
	}
	visibleSelect = visibleSelect.GROUP_BY(groupBys...)

	aclCondition := mysql.AND(
		mysql.IntExp(mysql.SUM(mysql.CASE().
			WHEN(matchingEffect.EQ(mysql.Int8(int8(AccessEffectDeny)))).
			THEN(mysql.Int(1)).
			ELSE(mysql.Int(0)))).EQ(mysql.Int(0)),
		mysql.IntExp(mysql.SUM(mysql.CASE().
			WHEN(matchingEffect.EQ(mysql.Int8(int8(AccessEffectAllow)))).
			THEN(mysql.Int(1)).
			ELSE(mysql.Int(0)))).GT(mysql.Int(0)),
	)
	if includeImplicitAccess {
		visibleSelect = visibleSelect.HAVING(mysql.OR(publicCondition, creatorCondition, aclCondition))
	} else {
		visibleSelect = visibleSelect.HAVING(aclCondition)
	}

	tSubjectUsers := table.FivenetAclSubjectUsers.AS("asu")
	tSubjectQualis := table.FivenetAclSubjectQualifications.AS("asq")
	tQualiResults := table.FivenetQualificationsResults.AS("qr")
	tSubjectJobGrade := table.FivenetAclSubjectJobGradeScopes.AS("asjg")
	tUserJobs := table.FivenetUserJobs.AS("uj")

	actorSubjectsUnion := mysql.SELECT(
		tSubjectUsers.SubjectID.AS("subject_id"),
		mysql.Int32(SubjectSpecificityUser).AS("specificity"),
		mysql.Int32(-1).AS("grade_specificity"),
	).
		FROM(tSubjectUsers).
		WHERE(tSubjectUsers.UserID.EQ(mysql.Int32(userInfo.GetUserId()))).
		UNION_ALL(
			mysql.SELECT(
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
							tQualiResults.Status.EQ(mysql.Int32(int32(qualifications.ResultStatus_RESULT_STATUS_SUCCESSFUL))),
						),
					),
				),
		).
		UNION_ALL(
			mysql.SELECT(
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
			mysql.SELECT(
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

	finalSelect := mysql.SELECT(visibleID.AS("id")).
		FROM(visibleObjects).
		ORDER_BY(visibleID.DESC())
	if countOnly {
		finalSelect = mysql.SELECT(mysql.COUNT(visibleID).AS("exact_total")).
			FROM(visibleObjects)
	}

	return mysql.WITH(
		actorSubjects.AS(actorSubjectsUnion),
		matchingACL.AS(
			mysql.SELECT(
				a.accessColumns.TargetID.AS("target_id"),
				a.accessColumns.Effect.AS("effect"),
				actorSpecificity.AS("specificity"),
				mysql.COALESCE(actorGradeSpecificity, mysql.Int32(-1)).AS("grade_specificity"),
				mysql.ROW_NUMBER().OVER(
					mysql.PARTITION_BY(a.accessColumns.TargetID).
						ORDER_BY(
							actorSpecificity.DESC(),
							mysql.COALESCE(actorGradeSpecificity, mysql.Int32(-1)).DESC(),
						),
				).AS("specificity_rank"),
			).
				FROM(a.accessTable.
					INNER_JOIN(actorSubjects,
						actorSubjectID.EQ(a.accessColumns.SubjectID),
					),
				).
				WHERE(mysql.OR(
					mysql.AND(
						a.accessColumns.Effect.EQ(mysql.Int8(int8(AccessEffectAllow))),
						a.accessColumns.Access.GT_EQ(mysql.Int32(access)),
					),
					mysql.AND(
						a.accessColumns.Effect.EQ(mysql.Int8(int8(AccessEffectDeny))),
						a.accessColumns.Access.EQ(mysql.Int32(access)),
					),
				)),
		),
		winningSpecificity.AS(
			mysql.SELECT(
				matchingTargetID.AS("target_id"),
				matchingSpecificity.AS("specificity"),
				matchingGradeSpecificity.AS("grade_specificity"),
			).
				FROM(matchingACL).
				WHERE(matchingSpecificityRank.EQ(mysql.Int(1))),
		),
		visibleObjects.AS(visibleSelect),
	)(finalSelect)
}

func (a *SubjectObjectAccess) CreateEntry(
	ctx context.Context,
	tx qrm.DB,
	targetID int64,
	subjectID int64,
	access int32,
	effect AccessEffect,
) error {
	exists, err := subjectExists(ctx, tx, subjectID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("acl subject %d does not exist", subjectID)
	}

	stmt := a.accessTable.
		INSERT(
			a.accessColumns.TargetID,
			a.accessColumns.SubjectID,
			a.accessColumns.Access,
			a.accessColumns.Effect,
		).
		VALUES(targetID, subjectID, access, int8(effect))

	_, err = stmt.ExecContext(ctx, tx)
	return err
}

func subjectExists(ctx context.Context, tx qrm.DB, subjectID int64) (bool, error) {
	stmt := table.FivenetAclSubjects.
		SELECT(mysql.Int(1).AS("exists")).
		FROM(table.FivenetAclSubjects).
		WHERE(table.FivenetAclSubjects.ID.EQ(mysql.Int64(subjectID))).
		LIMIT(1)

	var dest struct {
		Exists int32 `alias:"exists"`
	}
	if err := stmt.QueryContext(ctx, tx, &dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return dest.Exists > 0, nil
}

func (a *SubjectObjectAccess) ClearTarget(ctx context.Context, tx qrm.DB, targetID int64) error {
	stmt := a.accessTable.
		DELETE().
		WHERE(a.accessColumns.TargetID.EQ(mysql.Int64(targetID))).
		LIMIT(1000)

	_, err := stmt.ExecContext(ctx, tx)
	return err
}

func (a *SubjectObjectAccess) DeleteEntry(ctx context.Context, tx qrm.DB, targetID int64, id int64) error {
	stmt := a.accessTable.
		DELETE().
		WHERE(mysql.AND(
			a.accessColumns.ID.EQ(mysql.Int64(id)),
			a.accessColumns.TargetID.EQ(mysql.Int64(targetID)),
		)).
		LIMIT(1)

	_, err := stmt.ExecContext(ctx, tx)
	return err
}

func (a *SubjectObjectAccess) UpdateEntry(
	ctx context.Context,
	tx qrm.DB,
	targetID int64,
	id int64,
	subjectID int64,
	access int32,
	effect AccessEffect,
) error {
	stmt := a.accessTable.
		UPDATE(
			a.accessColumns.SubjectID,
			a.accessColumns.Access,
			a.accessColumns.Effect,
		).
		SET(subjectID, access, int8(effect)).
		WHERE(mysql.AND(
			a.accessColumns.ID.EQ(mysql.Int64(id)),
			a.accessColumns.TargetID.EQ(mysql.Int64(targetID)),
		)).
		LIMIT(1)

	_, err := stmt.ExecContext(ctx, tx)
	return err
}

func (a *SubjectObjectAccess) Validate() error {
	if a.db == nil {
		return fmt.Errorf("subject object access requires db")
	}
	if a.targetTable == nil || a.accessTable == nil {
		return fmt.Errorf("subject object access requires target and access tables")
	}
	if a.targetTableColumns == nil || a.targetTableColumns.ID == nil {
		return fmt.Errorf("subject object access requires target id column")
	}
	if a.accessColumns == nil || a.accessColumns.TargetID == nil ||
		a.accessColumns.SubjectID == nil || a.accessColumns.Access == nil || a.accessColumns.Effect == nil {
		return fmt.Errorf("subject object access requires access target, subject, access, and effect columns")
	}

	return nil
}
