package access

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
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
	stmt := mysql.SELECT(
		table.FivenetACLSubjectUsers.SubjectID.AS("subject_id"),
		mysql.Int32(SubjectSpecificityUser).AS("specificity"),
		mysql.Int32(-1).AS("grade_specificity"),
	).
		FROM(table.FivenetACLSubjectUsers).
		WHERE(table.FivenetACLSubjectUsers.UserID.EQ(mysql.Int32(userInfo.GetUserId()))).
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
				mysql.SELECT(mysql.Int(1)).
					FROM(tSubjectUsers).
					WHERE(tSubjectUsers.SubjectID.EQ(tSubjects.ID)),
			)),
			mysql.NOT(mysql.EXISTS(
				mysql.SELECT(mysql.Int(1)).
					FROM(tSubjectQuals).
					WHERE(tSubjectQuals.SubjectID.EQ(tSubjects.ID)),
			)),
			mysql.NOT(mysql.EXISTS(
				mysql.SELECT(mysql.Int(1)).
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
				mysql.SELECT(mysql.Int(1)).
					FROM(tSubjectJobGrades).
					WHERE(tSubjectJobGrades.SubjectID.EQ(tSubjects.ID)),
			),
			mysql.NOT(mysql.EXISTS(
				mysql.SELECT(mysql.Int(1)).
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

type SubjectAccessColumns struct {
	BaseAccessColumns

	SubjectID mysql.ColumnInteger
	Effect    mysql.ColumnBool
}

type SubjectTargetTableColumns struct {
	ID         mysql.ColumnInteger
	CreatedAt  mysql.ColumnTimestamp
	UpdatedAt  mysql.ColumnTimestamp
	DeletedAt  mysql.ColumnTimestamp
	Public     mysql.ColumnBool
	CreatorID  mysql.ColumnInteger
	CreatorJob mysql.ColumnString
}

type VisibilityRuleKind int16

const (
	VisibilityRulePublic  VisibilityRuleKind = 1
	VisibilityRuleCreator VisibilityRuleKind = 2
	VisibilityRuleCustom  VisibilityRuleKind = 3
)

type VisibilityRowKind int16

const (
	VisibilityRowKindPublic  VisibilityRowKind = 1
	VisibilityRowKindCreator VisibilityRowKind = 2
	VisibilityRowKindACL     VisibilityRowKind = 3
)

type VisibilityRule struct {
	Kind   VisibilityRuleKind
	Custom func(columns *SubjectTargetTableColumns, userInfo *userinfo.UserInfo) mysql.BoolExpression
}

type VisibilityPolicy struct {
	Rules []VisibilityRule
}

func (p VisibilityPolicy) HasRules() bool {
	return len(p.Rules) > 0
}

func (p VisibilityPolicy) HasRuleKind(kind VisibilityRuleKind) bool {
	for _, rule := range p.Rules {
		if rule.Kind == kind {
			return true
		}
	}

	return false
}

func (p VisibilityPolicy) Condition(
	columns *SubjectTargetTableColumns,
	userInfo *userinfo.UserInfo,
) mysql.BoolExpression {
	if columns == nil || userInfo == nil {
		return mysql.Bool(false)
	}

	var condition mysql.BoolExpression
	hasCondition := false
	for _, rule := range p.Rules {
		var ruleCondition mysql.BoolExpression
		hasRuleCondition := false

		switch rule.Kind {
		case VisibilityRulePublic:
			if columns.Public != nil {
				ruleCondition = columns.Public.IS_TRUE()
				hasRuleCondition = true
			}
		case VisibilityRuleCreator:
			creatorCondition := mysql.Bool(false)
			hasCreatorCondition := false

			if columns.CreatorID != nil {
				creatorCondition = columns.CreatorID.EQ(mysql.Int32(userInfo.GetUserId()))
				hasCreatorCondition = true
			}
			if columns.CreatorJob != nil {
				jobCondition := columns.CreatorJob.EQ(mysql.String(userInfo.GetJob()))
				if hasCreatorCondition {
					creatorCondition = creatorCondition.AND(jobCondition)
				} else {
					creatorCondition = jobCondition
					hasCreatorCondition = true
				}
			}

			if hasCreatorCondition {
				ruleCondition = creatorCondition
				hasRuleCondition = true
			}
		case VisibilityRuleCustom:
			if rule.Custom != nil {
				ruleCondition = rule.Custom(columns, userInfo)
				hasRuleCondition = ruleCondition != nil
			}
		}

		if !hasRuleCondition || ruleCondition == nil {
			continue
		}
		if !hasCondition {
			condition = ruleCondition
			hasCondition = true
			continue
		}
		condition = condition.OR(ruleCondition)
	}

	if !hasCondition {
		return mysql.Bool(false)
	}

	return condition
}

type SubjectObjectAccessConfig struct {
	TargetTable                           mysql.Table
	TargetColumns                         *SubjectTargetTableColumns
	AccessTable                           mysql.Table
	AccessColumns                         *SubjectAccessColumns
	CalculatedVisibilityPublicTable       mysql.Table
	CalculatedVisibilityCreatorTable      mysql.Table
	CalculatedVisibilitySubjectTable      mysql.Table
	CalculatedVisibilityPublicTargetID    mysql.ColumnInteger
	CalculatedVisibilityCreatorTargetID   mysql.ColumnInteger
	CalculatedVisibilityCreatorCreatorID  mysql.ColumnInteger
	CalculatedVisibilityCreatorCreatorJob mysql.ColumnString
	CalculatedVisibilitySubjectTargetID   mysql.ColumnInteger
	CalculatedVisibilitySubjectSubjectID  mysql.ColumnInteger
	CalculatedVisibilitySubjectAccess     mysql.ColumnInteger
	CalculatedVisibilitySubjectEffect     mysql.ColumnBool
	VisibilityTable                       mysql.Table
	VisibilityCols                        *VisibilityColumns
	Visibility                            VisibilityPolicy
	CalculatedVisibilityMaps              bool
}

type SubjectObjectAccess struct {
	db                                    *sql.DB
	targetTable                           mysql.Table
	targetTableColumns                    *SubjectTargetTableColumns
	accessTable                           mysql.Table
	accessColumns                         *SubjectAccessColumns
	calculatedVisibilityPublicTable       mysql.Table
	calculatedVisibilityCreatorTable      mysql.Table
	calculatedVisibilitySubjectTable      mysql.Table
	calculatedVisibilityPublicTargetID    mysql.ColumnInteger
	calculatedVisibilityCreatorTargetID   mysql.ColumnInteger
	calculatedVisibilityCreatorCreatorID  mysql.ColumnInteger
	calculatedVisibilityCreatorCreatorJob mysql.ColumnString
	calculatedVisibilitySubjectTargetID   mysql.ColumnInteger
	calculatedVisibilitySubjectSubjectID  mysql.ColumnInteger
	calculatedVisibilitySubjectAccess     mysql.ColumnInteger
	calculatedVisibilitySubjectEffect     mysql.ColumnBool
	visibilityTable                       mysql.Table
	visibilityColumns                     *VisibilityColumns
	visibility                            VisibilityPolicy
	visibilityBackend                     VisibilityBackend
}

type canAccessIdsHelper struct {
	IDs []int64 `alias:"id"`
}

type VisibilityBackend interface {
	VisibleIDsQuery(
		userInfo *userinfo.UserInfo,
		access int32,
		targetIDs ...int64,
	) VisibilityQuery
	VisibleIDsByConditionQuery(
		userInfo *userinfo.UserInfo,
		access int32,
		condition mysql.BoolExpression,
	) VisibilityQuery
	ACLVisibleIDsByConditionQuery(
		userInfo *userinfo.UserInfo,
		access int32,
		condition mysql.BoolExpression,
	) VisibilityQuery
	VisibleIDsStatement(
		userInfo *userinfo.UserInfo,
		access int32,
		targetIDs ...int64,
	) mysql.Statement
	VisibleIDsByConditionStatement(
		userInfo *userinfo.UserInfo,
		access int32,
		condition mysql.BoolExpression,
	) mysql.Statement
	ACLVisibleIDsByConditionStatement(
		userInfo *userinfo.UserInfo,
		access int32,
		condition mysql.BoolExpression,
	) mysql.Statement
	CountVisibleByConditionStatement(
		userInfo *userinfo.UserInfo,
		access int32,
		condition mysql.BoolExpression,
	) mysql.Statement
	RefreshTargetVisibility(
		ctx context.Context,
		tx qrm.DB,
		targetID int64,
	) error
}

type VisibilityQuery struct {
	CTEs      []mysql.CommonTableExpression
	Table     mysql.SelectTable
	Statement mysql.Statement
}

func NewSubjectObjectAccess(db *sql.DB, cfg SubjectObjectAccessConfig) *SubjectObjectAccess {
	access := &SubjectObjectAccess{
		db:                                    db,
		targetTable:                           cfg.TargetTable,
		targetTableColumns:                    cfg.TargetColumns,
		accessTable:                           cfg.AccessTable,
		accessColumns:                         cfg.AccessColumns,
		calculatedVisibilityPublicTable:       cfg.CalculatedVisibilityPublicTable,
		calculatedVisibilityCreatorTable:      cfg.CalculatedVisibilityCreatorTable,
		calculatedVisibilitySubjectTable:      cfg.CalculatedVisibilitySubjectTable,
		calculatedVisibilityPublicTargetID:    cfg.CalculatedVisibilityPublicTargetID,
		calculatedVisibilityCreatorTargetID:   cfg.CalculatedVisibilityCreatorTargetID,
		calculatedVisibilityCreatorCreatorID:  cfg.CalculatedVisibilityCreatorCreatorID,
		calculatedVisibilityCreatorCreatorJob: cfg.CalculatedVisibilityCreatorCreatorJob,
		calculatedVisibilitySubjectTargetID:   cfg.CalculatedVisibilitySubjectTargetID,
		calculatedVisibilitySubjectSubjectID:  cfg.CalculatedVisibilitySubjectSubjectID,
		calculatedVisibilitySubjectAccess:     cfg.CalculatedVisibilitySubjectAccess,
		calculatedVisibilitySubjectEffect:     cfg.CalculatedVisibilitySubjectEffect,
		visibilityTable:                       cfg.VisibilityTable,
		visibilityColumns:                     cfg.VisibilityCols,
		visibility:                            cfg.Visibility,
	}
	if cfg.CalculatedVisibilityMaps {
		access.visibilityBackend = &calculatedVisibilityBackend{access: access}
	} else {
		access.visibilityBackend = &sqlVisibilityBackend{access: access}
	}
	return access
}

func NewDocumentsSubjectObjectAccess(db *sql.DB) *SubjectObjectAccess {
	return NewSubjectObjectAccess(db, SubjectObjectAccessConfig{
		TargetTable: table.FivenetDocuments,
		TargetColumns: &SubjectTargetTableColumns{
			ID:         table.FivenetDocuments.ID,
			CreatedAt:  table.FivenetDocuments.CreatedAt,
			UpdatedAt:  table.FivenetDocuments.UpdatedAt,
			DeletedAt:  table.FivenetDocuments.DeletedAt,
			Public:     table.FivenetDocuments.Public,
			CreatorID:  table.FivenetDocuments.CreatorID,
			CreatorJob: table.FivenetDocuments.CreatorJob,
		},
		AccessTable: table.FivenetDocumentsAccess,
		AccessColumns: &SubjectAccessColumns{
			BaseAccessColumns: BaseAccessColumns{
				ID:       table.FivenetDocumentsAccess.ID,
				TargetID: table.FivenetDocumentsAccess.TargetID,
				Access:   table.FivenetDocumentsAccess.Access,
			},
			SubjectID: table.FivenetDocumentsAccess.SubjectID,
			Effect:    table.FivenetDocumentsAccess.Effect,
		},
		CalculatedVisibilityPublicTable:       table.FivenetDocumentsVisibilityPublic,
		CalculatedVisibilityCreatorTable:      table.FivenetDocumentsVisibilityCreator,
		CalculatedVisibilitySubjectTable:      table.FivenetDocumentsVisibilitySubject,
		CalculatedVisibilityPublicTargetID:    table.FivenetDocumentsVisibilityPublic.TargetID,
		CalculatedVisibilityCreatorTargetID:   table.FivenetDocumentsVisibilityCreator.TargetID,
		CalculatedVisibilityCreatorCreatorID:  table.FivenetDocumentsVisibilityCreator.CreatorID,
		CalculatedVisibilityCreatorCreatorJob: table.FivenetDocumentsVisibilityCreator.CreatorJob,
		CalculatedVisibilitySubjectTargetID:   table.FivenetDocumentsVisibilitySubject.TargetID,
		CalculatedVisibilitySubjectSubjectID:  table.FivenetDocumentsVisibilitySubject.SubjectID,
		CalculatedVisibilitySubjectAccess:     table.FivenetDocumentsVisibilitySubject.Access,
		CalculatedVisibilitySubjectEffect:     table.FivenetDocumentsVisibilitySubject.Effect,
		Visibility: VisibilityPolicy{
			Rules: []VisibilityRule{
				{Kind: VisibilityRulePublic},
				{Kind: VisibilityRuleCreator},
			},
		},
		CalculatedVisibilityMaps: true,
	})
}

func NewDocumentTemplatesSubjectObjectAccess(db *sql.DB) *SubjectObjectAccess {
	return NewSubjectObjectAccess(db, SubjectObjectAccessConfig{
		TargetTable: table.FivenetDocumentsTemplates,
		TargetColumns: &SubjectTargetTableColumns{
			ID:         table.FivenetDocumentsTemplates.ID,
			CreatedAt:  table.FivenetDocumentsTemplates.CreatedAt,
			UpdatedAt:  table.FivenetDocumentsTemplates.UpdatedAt,
			DeletedAt:  table.FivenetDocumentsTemplates.DeletedAt,
			CreatorJob: table.FivenetDocumentsTemplates.CreatorJob,
		},
		AccessTable: table.FivenetDocumentsTemplatesAccess,
		AccessColumns: &SubjectAccessColumns{
			BaseAccessColumns: BaseAccessColumns{
				ID:       table.FivenetDocumentsTemplatesAccess.ID,
				TargetID: table.FivenetDocumentsTemplatesAccess.TargetID,
				Access:   table.FivenetDocumentsTemplatesAccess.Access,
			},
			SubjectID: table.FivenetDocumentsTemplatesAccess.SubjectID,
			Effect:    table.FivenetDocumentsTemplatesAccess.Effect,
		},
		CalculatedVisibilityCreatorTable:      table.FivenetDocumentsTemplatesVisibilityCreator,
		CalculatedVisibilitySubjectTable:      table.FivenetDocumentsTemplatesVisibilitySubject,
		CalculatedVisibilityCreatorTargetID:   table.FivenetDocumentsTemplatesVisibilityCreator.TargetID,
		CalculatedVisibilityCreatorCreatorID:  table.FivenetDocumentsTemplatesVisibilityCreator.CreatorID,
		CalculatedVisibilityCreatorCreatorJob: table.FivenetDocumentsTemplatesVisibilityCreator.CreatorJob,
		CalculatedVisibilitySubjectTargetID:   table.FivenetDocumentsTemplatesVisibilitySubject.TargetID,
		CalculatedVisibilitySubjectSubjectID:  table.FivenetDocumentsTemplatesVisibilitySubject.SubjectID,
		CalculatedVisibilitySubjectAccess:     table.FivenetDocumentsTemplatesVisibilitySubject.Access,
		CalculatedVisibilitySubjectEffect:     table.FivenetDocumentsTemplatesVisibilitySubject.Effect,
		Visibility: VisibilityPolicy{
			Rules: []VisibilityRule{
				{Kind: VisibilityRuleCreator},
			},
		},
		CalculatedVisibilityMaps: true,
	})
}

func NewDocumentStampsSubjectObjectAccess(db *sql.DB) *SubjectObjectAccess {
	return NewSubjectObjectAccess(db, SubjectObjectAccessConfig{
		TargetTable: table.FivenetDocumentsStamps,
		TargetColumns: &SubjectTargetTableColumns{
			ID:         table.FivenetDocumentsStamps.ID,
			CreatedAt:  table.FivenetDocumentsStamps.CreatedAt,
			UpdatedAt:  table.FivenetDocumentsStamps.UpdatedAt,
			DeletedAt:  table.FivenetDocumentsStamps.DeletedAt,
			CreatorJob: table.FivenetDocumentsStamps.Name,
		},
		AccessTable: table.FivenetDocumentsStampsAccess,
		AccessColumns: &SubjectAccessColumns{
			BaseAccessColumns: BaseAccessColumns{
				ID:       table.FivenetDocumentsStampsAccess.ID,
				TargetID: table.FivenetDocumentsStampsAccess.TargetID,
				Access:   table.FivenetDocumentsStampsAccess.Access,
			},
			SubjectID: table.FivenetDocumentsStampsAccess.SubjectID,
			Effect:    table.FivenetDocumentsStampsAccess.Effect,
		},
		CalculatedVisibilityCreatorTable:      table.FivenetDocumentsStampsVisibilityCreator,
		CalculatedVisibilitySubjectTable:      table.FivenetDocumentsStampsVisibilitySubject,
		CalculatedVisibilityCreatorTargetID:   table.FivenetDocumentsStampsVisibilityCreator.TargetID,
		CalculatedVisibilityCreatorCreatorJob: table.FivenetDocumentsStampsVisibilityCreator.CreatorJob,
		CalculatedVisibilitySubjectTargetID:   table.FivenetDocumentsStampsVisibilitySubject.TargetID,
		CalculatedVisibilitySubjectSubjectID:  table.FivenetDocumentsStampsVisibilitySubject.SubjectID,
		CalculatedVisibilitySubjectAccess:     table.FivenetDocumentsStampsVisibilitySubject.Access,
		CalculatedVisibilitySubjectEffect:     table.FivenetDocumentsStampsVisibilitySubject.Effect,
		Visibility: VisibilityPolicy{
			Rules: []VisibilityRule{{Kind: VisibilityRuleCreator}},
		},
		CalculatedVisibilityMaps: true,
	})
}

func NewCalendarSubjectObjectAccess(db *sql.DB) *SubjectObjectAccess {
	return NewSubjectObjectAccess(db, SubjectObjectAccessConfig{
		TargetTable: table.FivenetCalendar,
		TargetColumns: &SubjectTargetTableColumns{
			ID:         table.FivenetCalendar.ID,
			CreatedAt:  table.FivenetCalendar.CreatedAt,
			UpdatedAt:  table.FivenetCalendar.UpdatedAt,
			DeletedAt:  table.FivenetCalendar.DeletedAt,
			Public:     table.FivenetCalendar.Public,
			CreatorID:  table.FivenetCalendar.CreatorID,
			CreatorJob: table.FivenetCalendar.CreatorJob,
		},
		AccessTable: table.FivenetCalendarAccess,
		AccessColumns: &SubjectAccessColumns{
			BaseAccessColumns: BaseAccessColumns{
				ID:       table.FivenetCalendarAccess.ID,
				TargetID: table.FivenetCalendarAccess.TargetID,
				Access:   table.FivenetCalendarAccess.Access,
			},
			SubjectID: table.FivenetCalendarAccess.SubjectID,
			Effect:    table.FivenetCalendarAccess.Effect,
		},
		CalculatedVisibilityPublicTable:       table.FivenetCalendarVisibilityPublic,
		CalculatedVisibilityCreatorTable:      table.FivenetCalendarVisibilityCreator,
		CalculatedVisibilitySubjectTable:      table.FivenetCalendarVisibilitySubject,
		CalculatedVisibilityPublicTargetID:    table.FivenetCalendarVisibilityPublic.TargetID,
		CalculatedVisibilityCreatorTargetID:   table.FivenetCalendarVisibilityCreator.TargetID,
		CalculatedVisibilityCreatorCreatorID:  table.FivenetCalendarVisibilityCreator.CreatorID,
		CalculatedVisibilityCreatorCreatorJob: table.FivenetCalendarVisibilityCreator.CreatorJob,
		CalculatedVisibilitySubjectTargetID:   table.FivenetCalendarVisibilitySubject.TargetID,
		CalculatedVisibilitySubjectSubjectID:  table.FivenetCalendarVisibilitySubject.SubjectID,
		CalculatedVisibilitySubjectAccess:     table.FivenetCalendarVisibilitySubject.Access,
		CalculatedVisibilitySubjectEffect:     table.FivenetCalendarVisibilitySubject.Effect,
		Visibility: VisibilityPolicy{
			Rules: []VisibilityRule{
				{Kind: VisibilityRulePublic},
				{Kind: VisibilityRuleCreator},
			},
		},
		CalculatedVisibilityMaps: true,
	})
}

func NewWikiPageSubjectObjectAccess(db *sql.DB) *SubjectObjectAccess {
	return NewSubjectObjectAccess(db, SubjectObjectAccessConfig{
		TargetTable: table.FivenetWikiPages,
		TargetColumns: &SubjectTargetTableColumns{
			ID:         table.FivenetWikiPages.ID,
			CreatedAt:  table.FivenetWikiPages.CreatedAt,
			UpdatedAt:  table.FivenetWikiPages.UpdatedAt,
			DeletedAt:  table.FivenetWikiPages.DeletedAt,
			Public:     table.FivenetWikiPages.Public,
			CreatorID:  table.FivenetWikiPages.CreatorID,
			CreatorJob: table.FivenetWikiPages.Job,
		},
		AccessTable: table.FivenetWikiPagesAccess,
		AccessColumns: &SubjectAccessColumns{
			BaseAccessColumns: BaseAccessColumns{
				ID:       table.FivenetWikiPagesAccess.ID,
				TargetID: table.FivenetWikiPagesAccess.TargetID,
				Access:   table.FivenetWikiPagesAccess.Access,
			},
			SubjectID: table.FivenetWikiPagesAccess.SubjectID,
			Effect:    table.FivenetWikiPagesAccess.Effect,
		},
		CalculatedVisibilityCreatorTable:      table.FivenetWikiPagesVisibilityCreator,
		CalculatedVisibilityPublicTable:       table.FivenetWikiPagesVisibilityPublic,
		CalculatedVisibilityCreatorTargetID:   table.FivenetWikiPagesVisibilityCreator.TargetID,
		CalculatedVisibilityCreatorCreatorID:  table.FivenetWikiPagesVisibilityCreator.CreatorID,
		CalculatedVisibilityCreatorCreatorJob: table.FivenetWikiPagesVisibilityCreator.CreatorJob,
		CalculatedVisibilitySubjectTable:      table.FivenetWikiPagesVisibilitySubject,
		CalculatedVisibilityPublicTargetID:    table.FivenetWikiPagesVisibilityPublic.TargetID,
		CalculatedVisibilitySubjectTargetID:   table.FivenetWikiPagesVisibilitySubject.TargetID,
		CalculatedVisibilitySubjectSubjectID:  table.FivenetWikiPagesVisibilitySubject.SubjectID,
		CalculatedVisibilitySubjectAccess:     table.FivenetWikiPagesVisibilitySubject.Access,
		CalculatedVisibilitySubjectEffect:     table.FivenetWikiPagesVisibilitySubject.Effect,
		Visibility: VisibilityPolicy{
			Rules: []VisibilityRule{
				{Kind: VisibilityRulePublic},
				{Kind: VisibilityRuleCreator},
			},
		},
		CalculatedVisibilityMaps: true,
	})
}

func NewCitizenLabelsSubjectObjectAccess(db *sql.DB) *SubjectObjectAccess {
	return NewSubjectObjectAccess(db, SubjectObjectAccessConfig{
		TargetTable: table.FivenetUserLabelsJob,
		TargetColumns: &SubjectTargetTableColumns{
			ID:         table.FivenetUserLabelsJob.ID,
			CreatedAt:  table.FivenetUserLabelsJob.CreatedAt,
			UpdatedAt:  table.FivenetUserLabelsJob.UpdatedAt,
			DeletedAt:  table.FivenetUserLabelsJob.DeletedAt,
			CreatorJob: table.FivenetUserLabelsJob.Job,
		},
		AccessTable: table.FivenetUserLabelsJobJobAccess,
		AccessColumns: &SubjectAccessColumns{
			BaseAccessColumns: BaseAccessColumns{
				ID:       table.FivenetUserLabelsJobJobAccess.ID,
				TargetID: table.FivenetUserLabelsJobJobAccess.TargetID,
				Access:   table.FivenetUserLabelsJobJobAccess.Access,
			},
			SubjectID: table.FivenetUserLabelsJobJobAccess.SubjectID,
			Effect:    table.FivenetUserLabelsJobJobAccess.Effect,
		},
		CalculatedVisibilityCreatorTable:      table.FivenetUserLabelsJobVisibilityCreator,
		CalculatedVisibilitySubjectTable:      table.FivenetUserLabelsJobVisibilitySubject,
		CalculatedVisibilityCreatorTargetID:   table.FivenetUserLabelsJobVisibilityCreator.TargetID,
		CalculatedVisibilityCreatorCreatorJob: table.FivenetUserLabelsJobVisibilityCreator.CreatorJob,
		CalculatedVisibilitySubjectTargetID:   table.FivenetUserLabelsJobVisibilitySubject.TargetID,
		CalculatedVisibilitySubjectSubjectID:  table.FivenetUserLabelsJobVisibilitySubject.SubjectID,
		CalculatedVisibilitySubjectAccess:     table.FivenetUserLabelsJobVisibilitySubject.Access,
		CalculatedVisibilitySubjectEffect:     table.FivenetUserLabelsJobVisibilitySubject.Effect,
		Visibility: VisibilityPolicy{
			Rules: []VisibilityRule{
				{Kind: VisibilityRuleCreator},
			},
		},
		CalculatedVisibilityMaps: true,
	})
}

func NewQualificationsSubjectObjectAccess(db *sql.DB) *SubjectObjectAccess {
	return NewSubjectObjectAccess(db, SubjectObjectAccessConfig{
		TargetTable: table.FivenetQualifications,
		TargetColumns: &SubjectTargetTableColumns{
			ID:         table.FivenetQualifications.ID,
			CreatedAt:  table.FivenetQualifications.CreatedAt,
			UpdatedAt:  table.FivenetQualifications.UpdatedAt,
			DeletedAt:  table.FivenetQualifications.DeletedAt,
			Public:     table.FivenetQualifications.Public,
			CreatorID:  table.FivenetQualifications.CreatorID,
			CreatorJob: table.FivenetQualifications.CreatorJob,
		},
		AccessTable: table.FivenetQualificationsAccess,
		AccessColumns: &SubjectAccessColumns{
			BaseAccessColumns: BaseAccessColumns{
				ID:       table.FivenetQualificationsAccess.ID,
				TargetID: table.FivenetQualificationsAccess.TargetID,
				Access:   table.FivenetQualificationsAccess.Access,
			},
			SubjectID: table.FivenetQualificationsAccess.SubjectID,
			Effect:    table.FivenetQualificationsAccess.Effect,
		},
		CalculatedVisibilityPublicTable:       table.FivenetQualificationsVisibilityPublic,
		CalculatedVisibilityCreatorTable:      table.FivenetQualificationsVisibilityCreator,
		CalculatedVisibilitySubjectTable:      table.FivenetQualificationsVisibilitySubject,
		CalculatedVisibilityPublicTargetID:    table.FivenetQualificationsVisibilityPublic.TargetID,
		CalculatedVisibilityCreatorTargetID:   table.FivenetQualificationsVisibilityCreator.TargetID,
		CalculatedVisibilityCreatorCreatorID:  table.FivenetQualificationsVisibilityCreator.CreatorID,
		CalculatedVisibilityCreatorCreatorJob: table.FivenetQualificationsVisibilityCreator.CreatorJob,
		CalculatedVisibilitySubjectTargetID:   table.FivenetQualificationsVisibilitySubject.TargetID,
		CalculatedVisibilitySubjectSubjectID:  table.FivenetQualificationsVisibilitySubject.SubjectID,
		CalculatedVisibilitySubjectAccess:     table.FivenetQualificationsVisibilitySubject.Access,
		CalculatedVisibilitySubjectEffect:     table.FivenetQualificationsVisibilitySubject.Effect,
		Visibility: VisibilityPolicy{
			Rules: []VisibilityRule{
				{Kind: VisibilityRulePublic},
				{Kind: VisibilityRuleCreator},
			},
		},
		CalculatedVisibilityMaps: true,
	})
}

func NewMailerEmailsSubjectObjectAccess(db *sql.DB) *SubjectObjectAccess {
	return NewSubjectObjectAccess(db, SubjectObjectAccessConfig{
		TargetTable: table.FivenetMailerEmails,
		TargetColumns: &SubjectTargetTableColumns{
			ID:        table.FivenetMailerEmails.ID,
			CreatedAt: table.FivenetMailerEmails.CreatedAt,
			UpdatedAt: table.FivenetMailerEmails.UpdatedAt,
			DeletedAt: table.FivenetMailerEmails.DeletedAt,
			CreatorID: table.FivenetMailerEmails.UserID,
		},
		AccessTable: table.FivenetMailerEmailsAccess,
		AccessColumns: &SubjectAccessColumns{
			BaseAccessColumns: BaseAccessColumns{
				ID:       table.FivenetMailerEmailsAccess.ID,
				TargetID: table.FivenetMailerEmailsAccess.TargetID,
				Access:   table.FivenetMailerEmailsAccess.Access,
			},
			SubjectID: table.FivenetMailerEmailsAccess.SubjectID,
			Effect:    table.FivenetMailerEmailsAccess.Effect,
		},
		CalculatedVisibilityCreatorTable:      table.FivenetMailerEmailsVisibilityCreator,
		CalculatedVisibilitySubjectTable:      table.FivenetMailerEmailsVisibilitySubject,
		CalculatedVisibilityCreatorTargetID:   table.FivenetMailerEmailsVisibilityCreator.TargetID,
		CalculatedVisibilityCreatorCreatorID:  table.FivenetMailerEmailsVisibilityCreator.CreatorID,
		CalculatedVisibilityCreatorCreatorJob: table.FivenetMailerEmailsVisibilityCreator.CreatorJob,
		CalculatedVisibilitySubjectTargetID:   table.FivenetMailerEmailsVisibilitySubject.TargetID,
		CalculatedVisibilitySubjectSubjectID:  table.FivenetMailerEmailsVisibilitySubject.SubjectID,
		CalculatedVisibilitySubjectAccess:     table.FivenetMailerEmailsVisibilitySubject.Access,
		CalculatedVisibilitySubjectEffect:     table.FivenetMailerEmailsVisibilitySubject.Effect,
		Visibility: VisibilityPolicy{
			Rules: []VisibilityRule{
				{Kind: VisibilityRuleCreator},
			},
		},
		CalculatedVisibilityMaps: true,
	})
}

func NewCentrumUnitsSubjectObjectAccess(db *sql.DB) *SubjectObjectAccess {
	return NewSubjectObjectAccess(db, SubjectObjectAccessConfig{
		TargetTable: table.FivenetCentrumUnits,
		TargetColumns: &SubjectTargetTableColumns{
			ID:        table.FivenetCentrumUnits.ID,
			CreatedAt: table.FivenetCentrumUnits.CreatedAt,
			UpdatedAt: table.FivenetCentrumUnits.UpdatedAt,
			DeletedAt: table.FivenetCentrumUnits.DeletedAt,
		},
		AccessTable: table.FivenetCentrumUnitsAccess,
		AccessColumns: &SubjectAccessColumns{
			BaseAccessColumns: BaseAccessColumns{
				ID:       table.FivenetCentrumUnitsAccess.ID,
				TargetID: table.FivenetCentrumUnitsAccess.TargetID,
				Access:   table.FivenetCentrumUnitsAccess.Access,
			},
			SubjectID: table.FivenetCentrumUnitsAccess.SubjectID,
			Effect:    table.FivenetCentrumUnitsAccess.Effect,
		},
	})
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
	if userInfo != nil && userInfo.GetSuperuser() {
		return targetIDs, nil
	}
	stmt := a.backend().VisibleIDsStatement(userInfo, access, targetIDs...)

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
	return a.backend().VisibleIDsStatement(userInfo, access, targetIDs...)
}

func (a *SubjectObjectAccess) VisibleIDsQuery(
	userInfo *userinfo.UserInfo,
	access int32,
	targetIDs ...int64,
) VisibilityQuery {
	return a.backend().VisibleIDsQuery(userInfo, access, targetIDs...)
}

func (a *SubjectObjectAccess) VisibleIDsByConditionStatement(
	userInfo *userinfo.UserInfo,
	access int32,
	condition mysql.BoolExpression,
) mysql.Statement {
	return a.backend().VisibleIDsByConditionStatement(userInfo, access, condition)
}

func (a *SubjectObjectAccess) VisibleIDsByConditionQuery(
	userInfo *userinfo.UserInfo,
	access int32,
	condition mysql.BoolExpression,
) VisibilityQuery {
	return a.backend().VisibleIDsByConditionQuery(userInfo, access, condition)
}

func (a *SubjectObjectAccess) ACLVisibleIDsByConditionStatement(
	userInfo *userinfo.UserInfo,
	access int32,
	condition mysql.BoolExpression,
) mysql.Statement {
	return a.backend().ACLVisibleIDsByConditionStatement(userInfo, access, condition)
}

func (a *SubjectObjectAccess) ACLVisibleIDsByConditionQuery(
	userInfo *userinfo.UserInfo,
	access int32,
	condition mysql.BoolExpression,
) VisibilityQuery {
	return a.backend().ACLVisibleIDsByConditionQuery(userInfo, access, condition)
}

func (a *SubjectObjectAccess) CountVisibleByConditionStatement(
	userInfo *userinfo.UserInfo,
	access int32,
	condition mysql.BoolExpression,
) mysql.Statement {
	return a.backend().CountVisibleByConditionStatement(userInfo, access, condition)
}

func (a *SubjectObjectAccess) RefreshTargetVisibility(
	ctx context.Context,
	tx qrm.DB,
	targetID int64,
) error {
	return a.backend().RefreshTargetVisibility(ctx, tx, targetID)
}

func (a *SubjectObjectAccess) RefreshTargetVisibilityWithCreator(
	ctx context.Context,
	tx qrm.DB,
	targetID int64,
	creatorID int32,
	creatorJob string,
) error {
	if backend, ok := a.backend().(interface {
		RefreshTargetVisibilityWithCreator(
			ctx context.Context,
			tx qrm.DB,
			targetID int64,
			creatorID int32,
			creatorJob string,
		) error
	}); ok {
		return backend.RefreshTargetVisibilityWithCreator(ctx, tx, targetID, creatorID, creatorJob)
	}

	return a.RefreshTargetVisibility(ctx, tx, targetID)
}

func (a *SubjectObjectAccess) backend() VisibilityBackend {
	if a.visibilityBackend == nil {
		a.visibilityBackend = &sqlVisibilityBackend{access: a}
	}
	return a.visibilityBackend
}

type sqlVisibilityBackend struct {
	access *SubjectObjectAccess
}

func (b *sqlVisibilityBackend) VisibleIDsQuery(
	userInfo *userinfo.UserInfo,
	access int32,
	targetIDs ...int64,
) VisibilityQuery {
	ids := make([]mysql.Expression, 0, len(targetIDs))
	for _, targetID := range targetIDs {
		ids = append(ids, mysql.Int64(targetID))
	}

	condition := b.access.targetTableColumns.ID.IN(ids...)
	return b.VisibleIDsByConditionQuery(userInfo, access, condition)
}

func (b *sqlVisibilityBackend) VisibleIDsStatement(
	userInfo *userinfo.UserInfo,
	access int32,
	targetIDs ...int64,
) mysql.Statement {
	query := b.VisibleIDsQuery(userInfo, access, targetIDs...)
	if len(query.CTEs) == 0 {
		return query.Statement
	}
	return mysql.WITH(query.CTEs...)(query.Statement)
}

func (b *sqlVisibilityBackend) VisibleIDsByConditionStatement(
	userInfo *userinfo.UserInfo,
	access int32,
	condition mysql.BoolExpression,
) mysql.Statement {
	query := b.VisibleIDsByConditionQuery(userInfo, access, condition)
	if len(query.CTEs) == 0 {
		return query.Statement
	}
	return mysql.WITH(query.CTEs...)(query.Statement)
}

func (b *sqlVisibilityBackend) VisibleIDsByConditionQuery(
	userInfo *userinfo.UserInfo,
	access int32,
	condition mysql.BoolExpression,
) VisibilityQuery {
	if userInfo == nil {
		userInfo = &userinfo.UserInfo{}
	}
	condition = b.access.normalizeTargetCondition(condition)
	if userInfo != nil && userInfo.GetSuperuser() {
		return VisibilityQuery{
			Table:     b.access.visibleTargetIDsStatement(condition).AsTable("doc_ids"),
			Statement: b.access.visibleTargetIDsStatement(condition),
		}
	}

	implicitStmt := b.access.implicitVisibleIDsStatement(userInfo, condition)
	aclCtes, aclSelect := b.access.aclVisibleIDsComponents(userInfo, access, condition, false)
	if implicitStmt == nil {
		return VisibilityQuery{
			CTEs:      aclCtes,
			Table:     aclSelect.AsTable("doc_ids"),
			Statement: aclSelect,
		}
	}

	combined := implicitStmt.UNION(aclSelect)
	return VisibilityQuery{
		CTEs:      aclCtes,
		Table:     combined.AsTable("doc_ids"),
		Statement: combined,
	}
}

func (b *sqlVisibilityBackend) ACLVisibleIDsByConditionQuery(
	userInfo *userinfo.UserInfo,
	access int32,
	condition mysql.BoolExpression,
) VisibilityQuery {
	if userInfo == nil {
		userInfo = &userinfo.UserInfo{}
	}
	condition = b.access.normalizeTargetCondition(condition)
	if userInfo != nil && userInfo.GetSuperuser() {
		return VisibilityQuery{
			Table:     b.access.visibleTargetIDsStatement(condition).AsTable("doc_ids"),
			Statement: b.access.visibleTargetIDsStatement(condition),
		}
	}

	aclCtes, aclSelect := b.access.aclVisibleIDsComponents(userInfo, access, condition, false)
	return VisibilityQuery{
		CTEs:      aclCtes,
		Table:     aclSelect.AsTable("doc_ids"),
		Statement: aclSelect,
	}
}

func (b *sqlVisibilityBackend) ACLVisibleIDsByConditionStatement(
	userInfo *userinfo.UserInfo,
	access int32,
	condition mysql.BoolExpression,
) mysql.Statement {
	query := b.ACLVisibleIDsByConditionQuery(userInfo, access, condition)
	if len(query.CTEs) == 0 {
		return query.Statement
	}
	return mysql.WITH(query.CTEs...)(query.Statement)
}

func (b *sqlVisibilityBackend) CountVisibleByConditionStatement(
	userInfo *userinfo.UserInfo,
	access int32,
	condition mysql.BoolExpression,
) mysql.Statement {
	if userInfo == nil {
		userInfo = &userinfo.UserInfo{}
	}
	condition = b.access.normalizeTargetCondition(condition)
	if userInfo != nil && userInfo.GetSuperuser() {
		return b.access.countTargetIDsStatement(condition)
	}

	query := b.VisibleIDsByConditionQuery(userInfo, access, condition)
	visibleID := mysql.IntegerColumn("id").From(query.Table)
	countSelect := mysql.SELECT(mysql.COUNT(visibleID).AS("exact_total")).FROM(query.Table)
	if len(query.CTEs) == 0 {
		return countSelect
	}
	return mysql.WITH(
		query.CTEs...,
	)(countSelect)
}

func (b *sqlVisibilityBackend) RefreshTargetVisibility(
	ctx context.Context,
	tx qrm.DB,
	targetID int64,
) error {
	return nil
}

func (a *SubjectObjectAccess) normalizeTargetCondition(
	condition mysql.BoolExpression,
) mysql.BoolExpression {
	if a.targetTableColumns.DeletedAt != nil {
		condition = condition.AND(a.targetTableColumns.DeletedAt.IS_NULL())
	}

	return condition
}

func (a *SubjectObjectAccess) visibleTargetIDsStatement(
	condition mysql.BoolExpression,
) mysql.SelectStatement {
	return a.targetTable.
		SELECT(a.targetTableColumns.ID.AS("id")).
		FROM(a.targetTable).
		WHERE(condition).
		DISTINCT()
}

func (a *SubjectObjectAccess) countTargetIDsStatement(
	condition mysql.BoolExpression,
) mysql.Statement {
	visibleIDs := mysql.CTE("visible_ids")
	visibleID := mysql.IntegerColumn("id").From(visibleIDs)
	return mysql.WITH(
		visibleIDs.AS(a.visibleTargetIDsStatement(condition)),
	)(mysql.SELECT(mysql.COUNT(visibleID).AS("exact_total")).FROM(visibleIDs))
}

func (a *SubjectObjectAccess) implicitVisibleIDsStatement(
	userInfo *userinfo.UserInfo,
	targetCondition mysql.BoolExpression,
) mysql.SelectStatement {
	if !a.visibility.HasRules() || userInfo == nil {
		return nil
	}

	implicitCondition := a.visibility.Condition(a.targetTableColumns, userInfo)
	if implicitCondition == nil {
		return nil
	}

	return a.targetTable.
		SELECT(a.targetTableColumns.ID.AS("id")).
		FROM(a.targetTable).
		WHERE(targetCondition.AND(implicitCondition)).
		DISTINCT()
}

func (a *SubjectObjectAccess) aclVisibleIDsComponents(
	userInfo *userinfo.UserInfo,
	access int32,
	targetCondition mysql.BoolExpression,
	countOnly bool,
) ([]mysql.CommonTableExpression, mysql.SelectStatement) {
	if userInfo == nil {
		userInfo = &userinfo.UserInfo{}
	}
	actorSubjects := mysql.CTE("actor_subjects")
	matchingACL := mysql.CTE("matching_acl")
	winningSpecificity := mysql.CTE("winning_specificity")
	visibleObjects := mysql.CTE("visible_objects")

	actorSubjectID := mysql.IntegerColumn("subject_id").From(actorSubjects)
	actorSpecificity := mysql.IntegerColumn("specificity").From(actorSubjects)
	actorGradeSpecificity := mysql.IntegerColumn("grade_specificity").From(actorSubjects)

	matchingTargetID := mysql.IntegerColumn("target_id").From(matchingACL)
	matchingEffect := mysql.BoolColumn("effect").From(matchingACL)
	matchingSpecificity := mysql.IntegerColumn("specificity").From(matchingACL)
	matchingGradeSpecificity := mysql.IntegerColumn("grade_specificity").From(matchingACL)
	matchingSpecificityRank := mysql.IntegerColumn("specificity_rank").From(matchingACL)

	winningTargetID := mysql.IntegerColumn("target_id").From(winningSpecificity)
	winningSpecificityCol := mysql.IntegerColumn("specificity").From(winningSpecificity)
	winningGradeSpecificity := mysql.IntegerColumn("grade_specificity").From(winningSpecificity)

	visibleID := mysql.IntegerColumn("id").From(visibleObjects)

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

	visibleSelect = visibleSelect.GROUP_BY(a.targetTableColumns.ID)

	aclCondition := mysql.AND(
		mysql.IntExp(mysql.SUM(mysql.CASE().
			WHEN(matchingEffect.IS_FALSE()).
			THEN(mysql.Int(1)).
			ELSE(mysql.Int(0)))).EQ(mysql.Int(0)),
		mysql.IntExp(mysql.SUM(mysql.CASE().
			WHEN(matchingEffect.IS_TRUE()).
			THEN(mysql.Int(1)).
			ELSE(mysql.Int(0)))).GT(mysql.Int(0)),
	)
	visibleSelect = visibleSelect.HAVING(aclCondition)

	tSubjectUsers := table.FivenetACLSubjectUsers.AS("asu")
	tSubjectQualis := table.FivenetACLSubjectQualifications.AS("asq")
	tQualiResults := table.FivenetQualificationsResults.AS("qr")
	tSubjectJobGrade := table.FivenetACLSubjectJobGradeScopes.AS("asjg")
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
		FROM(visibleObjects)
	if countOnly {
		finalSelect = mysql.SELECT(mysql.COUNT(visibleID).AS("exact_total")).
			FROM(visibleObjects)
	}

	return []mysql.CommonTableExpression{
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
						a.accessColumns.Effect.IS_TRUE(),
						a.accessColumns.Access.GT_EQ(mysql.Int32(access)),
					),
					mysql.AND(
						a.accessColumns.Effect.IS_FALSE(),
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
	}, finalSelect
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
		VALUES(targetID, subjectID, access, effect == AccessEffectAllow)

	_, err = stmt.ExecContext(ctx, tx)
	return err
}

func subjectExists(ctx context.Context, tx qrm.DB, subjectID int64) (bool, error) {
	stmt := table.FivenetACLSubjects.
		SELECT(mysql.Int(1).AS("exists")).
		FROM(table.FivenetACLSubjects).
		WHERE(table.FivenetACLSubjects.ID.EQ(mysql.Int64(subjectID))).
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

func (a *SubjectObjectAccess) DeleteEntry(
	ctx context.Context,
	tx qrm.DB,
	targetID int64,
	id int64,
) error {
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
		SET(subjectID, access, effect == AccessEffectAllow).
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
		return errors.New("subject object access requires db")
	}
	if a.targetTable == nil || a.accessTable == nil {
		return errors.New("subject object access requires target and access tables")
	}
	if a.targetTableColumns == nil || a.targetTableColumns.ID == nil {
		return errors.New("subject object access requires target id column")
	}
	if a.accessColumns == nil || a.accessColumns.TargetID == nil ||
		a.accessColumns.SubjectID == nil || a.accessColumns.Access == nil || a.accessColumns.Effect == nil {
		return errors.New(
			"subject object access requires access target, subject, access, and effect columns",
		)
	}
	if a.visibilityTable != nil && a.visibilityColumns == nil {
		return errors.New("subject object access requires visibility columns for visibility table")
	}

	return nil
}
