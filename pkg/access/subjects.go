package access

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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

type BaseAccessColumns struct {
	ID       mysql.ColumnInteger
	TargetID mysql.ColumnInteger
	Access   mysql.ColumnInteger
}

type VisibilityColumns struct {
	BaseAccessColumns

	RuleKind   mysql.ColumnInteger
	SubjectID  mysql.ColumnInteger
	CreatorID  mysql.ColumnInteger
	CreatorJob mysql.ColumnString
	Effect     mysql.ColumnBool
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
		includeDeleted bool,
		targetIDs ...int64,
	) VisibilityQuery
	VisibleIDsByConditionQuery(
		userInfo *userinfo.UserInfo,
		access int32,
		includeDeleted bool,
		condition mysql.BoolExpression,
	) VisibilityQuery
	ACLVisibleIDsByConditionQuery(
		userInfo *userinfo.UserInfo,
		access int32,
		includeDeleted bool,
		condition mysql.BoolExpression,
	) VisibilityQuery
	VisibleIDsStatement(
		userInfo *userinfo.UserInfo,
		access int32,
		includeDeleted bool,
		targetIDs ...int64,
	) mysql.Statement
	VisibleIDsByConditionStatement(
		userInfo *userinfo.UserInfo,
		access int32,
		includeDeleted bool,
		condition mysql.BoolExpression,
	) mysql.Statement
	ACLVisibleIDsByConditionStatement(
		userInfo *userinfo.UserInfo,
		access int32,
		includeDeleted bool,
		condition mysql.BoolExpression,
	) mysql.Statement
	CountVisibleByConditionStatement(
		userInfo *userinfo.UserInfo,
		access int32,
		includeDeleted bool,
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

	if err := access.Validate(); err != nil {
		panic(err)
	}

	return access
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
	stmt := a.backend().VisibleIDsStatement(userInfo, access, false, targetIDs...)

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
	includeDeleted bool,
	targetIDs ...int64,
) mysql.Statement {
	return a.backend().VisibleIDsStatement(userInfo, access, includeDeleted, targetIDs...)
}

func (a *SubjectObjectAccess) VisibleIDsQuery(
	userInfo *userinfo.UserInfo,
	access int32,
	includeDeleted bool,
	targetIDs ...int64,
) VisibilityQuery {
	return a.backend().VisibleIDsQuery(userInfo, access, includeDeleted, targetIDs...)
}

func (a *SubjectObjectAccess) VisibleIDsByConditionStatement(
	userInfo *userinfo.UserInfo,
	access int32,
	includeDeleted bool,
	condition mysql.BoolExpression,
) mysql.Statement {
	return a.backend().VisibleIDsByConditionStatement(userInfo, access, includeDeleted, condition)
}

func (a *SubjectObjectAccess) VisibleIDsByConditionQuery(
	userInfo *userinfo.UserInfo,
	access int32,
	includeDeleted bool,
	condition mysql.BoolExpression,
) VisibilityQuery {
	return a.backend().VisibleIDsByConditionQuery(userInfo, access, includeDeleted, condition)
}

func (a *SubjectObjectAccess) ACLVisibleIDsByConditionStatement(
	userInfo *userinfo.UserInfo,
	access int32,
	includeDeleted bool,
	condition mysql.BoolExpression,
) mysql.Statement {
	return a.backend().
		ACLVisibleIDsByConditionStatement(userInfo, access, includeDeleted, condition)
}

func (a *SubjectObjectAccess) ACLVisibleIDsByConditionQuery(
	userInfo *userinfo.UserInfo,
	access int32,
	includeDeleted bool,
	condition mysql.BoolExpression,
) VisibilityQuery {
	return a.backend().ACLVisibleIDsByConditionQuery(userInfo, access, includeDeleted, condition)
}

func (a *SubjectObjectAccess) CountVisibleByConditionStatement(
	userInfo *userinfo.UserInfo,
	access int32,
	includeDeleted bool,
	condition mysql.BoolExpression,
) mysql.Statement {
	return a.backend().CountVisibleByConditionStatement(userInfo, access, includeDeleted, condition)
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

func (a *SubjectObjectAccess) normalizeTargetCondition(
	condition mysql.BoolExpression,
	includeDeleted bool,
) mysql.BoolExpression {
	if !includeDeleted && a.targetTableColumns.DeletedAt != nil {
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
	)(mysql.
		SELECT(mysql.COUNT(visibleID).AS("exact_total")).FROM(visibleIDs))
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

	visibleSelect := mysql.
		SELECT(a.targetTableColumns.ID.AS("id")).
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

	visibleSelect = visibleSelect.
		GROUP_BY(a.targetTableColumns.ID)

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
	tQualiSuccess := table.FivenetQualificationsResultSuccessMap.AS("qrsm")
	tSubjectJobGrade := table.FivenetACLSubjectJobGradeScopes.AS("asjg")
	tUserJobs := table.FivenetUserJobs.AS("uj")

	actorSubjectsUnion := mysql.
		SELECT(
			tSubjectUsers.SubjectID.AS("subject_id"),
			mysql.Int32(SubjectSpecificityUser).AS("specificity"),
			mysql.Int32(-1).AS("grade_specificity"),
		).
		FROM(tSubjectUsers).
		WHERE(tSubjectUsers.UserID.EQ(mysql.Int32(userInfo.GetUserId()))).
		UNION_ALL(
			mysql.
				SELECT(
					tSubjectQualis.SubjectID.AS("subject_id"),
					mysql.Int32(SubjectSpecificityQualification).AS("specificity"),
					mysql.Int32(-1).AS("grade_specificity"),
				).
				FROM(tSubjectQualis.
					INNER_JOIN(tQualiSuccess,
						mysql.AND(
							tQualiSuccess.QualificationID.EQ(tSubjectQualis.QualificationID),
							tQualiSuccess.UserID.EQ(mysql.Int32(userInfo.GetUserId())),
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
							tUserJobs.IsPrimary.IS_TRUE(),
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

	finalSelect := mysql.
		SELECT(visibleID.AS("id")).
		FROM(visibleObjects)
	if countOnly {
		finalSelect = mysql.
			SELECT(mysql.COUNT(visibleID).AS("exact_total")).
			FROM(visibleObjects)
	}

	return []mysql.CommonTableExpression{
		actorSubjects.AS(actorSubjectsUnion),
		matchingACL.AS(
			mysql.
				SELECT(
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
						a.accessColumns.Access.GT_EQ(mysql.Int32(access)),
					),
				)),
		),
		winningSpecificity.AS(
			mysql.
				SELECT(
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
	if a.calculatedVisibilitySubjectTable != nil &&
		(a.calculatedVisibilitySubjectTargetID == nil ||
			a.calculatedVisibilitySubjectSubjectID == nil ||
			a.calculatedVisibilitySubjectAccess == nil ||
			a.calculatedVisibilitySubjectEffect == nil) {
		return errors.New(
			"subject object access requires calculated visibility subject target, subject, access, and effect columns",
		)
	}

	return nil
}
