package jobsstore

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	jobsgroups "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/groups"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

const (
	groupRuleTypeGrade         = jobsgroups.GroupRuleType_GROUP_RULE_TYPE_GRADE
	groupRuleTypeQualification = jobsgroups.GroupRuleType_GROUP_RULE_TYPE_QUALIFICATION
)

type groupRuleBuilder jobsgroups.GroupRule_builder

func (b *groupRuleBuilder) Build() *jobsgroups.GroupRule {
	if b == nil {
		return nil
	}

	return (*jobsgroups.GroupRule_builder)(b).Build()
}

func groupRuleColumns(tRules *table.FivenetJobGroupRulesTable) mysql.ProjectionList {
	return mysql.ProjectionList{
		tRules.ID.AS("group_rule_builder.id"),
		tRules.GroupID.AS("group_rule_builder.group_id"),
		tRules.RuleType.AS("group_rule_builder.type"),
		tRules.Enabled.AS("group_rule_builder.enabled"),
		tRules.CreatedByUserID.AS("group_rule_builder.created_by_user_id"),
		tRules.CreatedAt.AS("group_rule_builder.created_at"),
		tRules.UpdatedAt.AS("group_rule_builder.updated_at"),
	}
}

func nullableInt32Expression(v *int32) mysql.IntegerExpression {
	if v == nil {
		return mysql.IntExp(mysql.NULL)
	}

	return mysql.Int32(*v)
}

func int64Expressions(ids []int64) []mysql.Expression {
	exprs := make([]mysql.Expression, len(ids))
	for i := range ids {
		exprs[i] = mysql.Int64(ids[i])
	}

	return exprs
}

func finalizeRuleMemberMatches(
	matches []*GroupRuleMemberMatch,
	groupID int64,
	ruleID int64,
	label string,
) []*GroupRuleMemberMatch {
	for _, match := range matches {
		match.GroupID = groupID
		match.RuleID = ruleID
		match.Label = label
	}

	return matches
}

func nullableInt32(v *int32) any {
	if v == nil {
		return nil
	}

	return *v
}

func groupRuleType(rule *jobsgroups.GroupRule) (jobsgroups.GroupRuleType, error) {
	if rule.GetGrade() != nil {
		return groupRuleTypeGrade, nil
	}
	if rule.GetQualification() != nil {
		return groupRuleTypeQualification, nil
	}

	return 0, errors.New("job group rule has no rule body")
}

func (s *Store) ListGroupRules(
	ctx context.Context,
	db qrm.DB,
	groupID int64,
) ([]*jobsgroups.GroupRule, error) {
	tRules := table.FivenetJobGroupRules
	columns := groupRuleColumns(tRules)
	stmt := tRules.
		SELECT(columns[0], columns[1:]...).
		FROM(tRules).
		WHERE(tRules.GroupID.EQ(mysql.Int64(groupID))).
		ORDER_BY(tRules.CreatedAt.ASC(), tRules.ID.ASC())

	builders := []*groupRuleBuilder{}
	if err := stmt.QueryContext(ctx, db, &builders); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return []*jobsgroups.GroupRule{}, nil
		}
		return nil, err
	}

	rules := make([]*jobsgroups.GroupRule, 0, len(builders))
	for _, builder := range builders {
		rules = append(rules, builder.Build())
	}

	for _, rule := range rules {
		if err := s.loadGroupRuleDetails(ctx, db, rule, rule.GetType()); err != nil {
			return nil, err
		}
	}

	return rules, nil
}

func (s *Store) GetGroupRule(
	ctx context.Context,
	db qrm.DB,
	groupID int64,
	ruleID int64,
) (*jobsgroups.GroupRule, error) {
	tRules := table.FivenetJobGroupRules
	columns := groupRuleColumns(tRules)
	stmt := tRules.
		SELECT(columns[0], columns[1:]...).
		FROM(tRules).
		WHERE(mysql.AND(
			tRules.GroupID.EQ(mysql.Int64(groupID)),
			tRules.ID.EQ(mysql.Int64(ruleID)),
		)).
		LIMIT(1)

	builder := &groupRuleBuilder{}
	if err := stmt.QueryContext(ctx, db, builder); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	rule := builder.Build()
	if rule.GetId() == 0 {
		return nil, nil
	}

	if err := s.loadGroupRuleDetails(ctx, db, rule, rule.GetType()); err != nil {
		return nil, err
	}

	return rule, nil
}

func (s *Store) CreateGroupRule(
	ctx context.Context,
	db qrm.DB,
	rule *jobsgroups.GroupRule,
) (*jobsgroups.GroupRule, error) {
	ruleType, err := groupRuleType(rule)
	if err != nil {
		return nil, err
	}
	rule.Type = ruleType

	tRules := table.FivenetJobGroupRules
	stmt := tRules.
		INSERT(
			tRules.GroupID,
			tRules.RuleType,
			tRules.Enabled,
			tRules.CreatedByUserID,
			tRules.CreatedAt,
		).
		VALUES(
			rule.GetGroupId(),
			int32(ruleType),
			rule.GetEnabled(),
			rule.GetCreatedByUserId(),
			mysql.CURRENT_TIMESTAMP(),
		)

	result, err := stmt.ExecContext(ctx, db)
	if err != nil {
		return nil, err
	}

	ruleID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	rule.Id = ruleID

	if err := s.insertGroupRuleDetails(ctx, db, rule); err != nil {
		return nil, err
	}

	return s.GetGroupRule(ctx, db, rule.GetGroupId(), ruleID)
}

func (s *Store) UpdateGroupRule(
	ctx context.Context,
	db qrm.DB,
	rule *jobsgroups.GroupRule,
	updatedByUserID int32,
) (*jobsgroups.GroupRule, error) {
	ruleType, err := groupRuleType(rule)
	if err != nil {
		return nil, err
	}
	rule.Type = ruleType

	tRules := table.FivenetJobGroupRules
	stmt := tRules.
		UPDATE().
		SET(
			tRules.RuleType.SET(mysql.Int32(int32(ruleType))),
			tRules.Enabled.SET(mysql.Bool(rule.GetEnabled())),
			tRules.UpdatedByUserID.SET(mysql.Int32(updatedByUserID)),
			tRules.UpdatedAt.SET(mysql.CURRENT_TIMESTAMP()),
		).
		WHERE(mysql.AND(
			tRules.GroupID.EQ(mysql.Int64(rule.GetGroupId())),
			tRules.ID.EQ(mysql.Int64(rule.GetId())),
		))

	result, err := stmt.ExecContext(ctx, db)
	if err != nil {
		return nil, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, sql.ErrNoRows
	}

	if err := s.deleteGroupRuleDetails(ctx, db, rule.GetId()); err != nil {
		return nil, err
	}
	if err := s.insertGroupRuleDetails(ctx, db, rule); err != nil {
		return nil, err
	}

	return s.GetGroupRule(ctx, db, rule.GetGroupId(), rule.GetId())
}

func (s *Store) DeleteGroupRule(ctx context.Context, db qrm.DB, groupID int64, ruleID int64) error {
	tRules := table.FivenetJobGroupRules
	stmt := tRules.
		DELETE().
		WHERE(mysql.AND(
			tRules.GroupID.EQ(mysql.Int64(groupID)),
			tRules.ID.EQ(mysql.Int64(ruleID)),
		))

	result, err := stmt.ExecContext(ctx, db)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (s *Store) loadGroupRuleDetails(
	ctx context.Context,
	db qrm.DB,
	rule *jobsgroups.GroupRule,
	ruleType jobsgroups.GroupRuleType,
) error {
	switch ruleType {
	case groupRuleTypeGrade:
		return s.loadGroupGradeRule(ctx, db, rule)
	case groupRuleTypeQualification:
		return s.loadGroupQualificationRule(ctx, db, rule)
	default:
		return fmt.Errorf("unknown job group rule type %d", ruleType)
	}
}

func (s *Store) loadGroupGradeRule(
	ctx context.Context,
	db qrm.DB,
	rule *jobsgroups.GroupRule,
) error {
	tGrades := table.FivenetJobGroupRuleGrades
	stmt := tGrades.
		SELECT(
			tGrades.GradeRuleType.AS("group_grade_rule.type"),
			tGrades.Grade.AS("group_grade_rule.grade"),
			tGrades.MinGrade.AS("group_grade_rule.min_grade"),
			tGrades.MaxGrade.AS("group_grade_rule.max_grade"),
		).
		FROM(tGrades).
		WHERE(tGrades.RuleID.EQ(mysql.Int64(rule.GetId()))).
		LIMIT(1)

	gradeRule := &jobsgroups.GroupGradeRule{}
	if err := stmt.QueryContext(ctx, db, gradeRule); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return sql.ErrNoRows
		}
		return err
	}

	rule.SetGrade(gradeRule)

	return nil
}

func (s *Store) loadGroupQualificationRule(
	ctx context.Context,
	db qrm.DB,
	rule *jobsgroups.GroupRule,
) error {
	tQualifications := table.FivenetJobGroupRuleQualifications
	stmt := tQualifications.
		SELECT(
			tQualifications.QualificationRuleType.AS("group_qualification_rule.type"),
			tQualifications.RequireCompleted.AS("group_qualification_rule.require_completed"),
		).
		FROM(tQualifications).
		WHERE(tQualifications.RuleID.EQ(mysql.Int64(rule.GetId()))).
		LIMIT(1)

	qualificationRule := &jobsgroups.GroupQualificationRule{}
	if err := stmt.QueryContext(ctx, db, qualificationRule); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return sql.ErrNoRows
		}
		return err
	}

	tItems := table.FivenetJobGroupRuleQualificationItems
	itemStmt := tItems.
		SELECT(tItems.QualificationID.AS("qualification_id")).
		FROM(tItems).
		WHERE(tItems.RuleID.EQ(mysql.Int64(rule.GetId()))).
		ORDER_BY(tItems.QualificationID.ASC())

	qualificationIDs := []int64{}
	if err := itemStmt.QueryContext(ctx, db, &qualificationIDs); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			rule.SetQualification(qualificationRule)
			return nil
		}
		return err
	}
	qualificationRule.QualificationIds = qualificationIDs

	rule.SetQualification(qualificationRule)
	return nil
}

func (s *Store) insertGroupRuleDetails(
	ctx context.Context,
	db qrm.DB,
	rule *jobsgroups.GroupRule,
) error {
	if grade := rule.GetGrade(); grade != nil {
		tGrades := table.FivenetJobGroupRuleGrades
		stmt := tGrades.
			INSERT(
				tGrades.RuleID,
				tGrades.GradeRuleType,
				tGrades.Grade,
				tGrades.MinGrade,
				tGrades.MaxGrade,
			).
			VALUES(
				rule.GetId(),
				int32(grade.GetType()),
				nullableInt32(grade.Grade),
				nullableInt32(grade.MinGrade),
				nullableInt32(grade.MaxGrade),
			)

		_, err := stmt.ExecContext(ctx, db)
		return err
	}

	qualification := rule.GetQualification()
	if qualification == nil {
		return errors.New("job group qualification rule body is nil")
	}

	tQualifications := table.FivenetJobGroupRuleQualifications
	stmt := tQualifications.
		INSERT(
			tQualifications.RuleID,
			tQualifications.QualificationRuleType,
			tQualifications.RequireCompleted,
		).
		VALUES(
			rule.GetId(),
			int32(qualification.GetType()),
			qualification.GetRequireCompleted(),
		)

	_, err := stmt.ExecContext(ctx, db)
	if err != nil {
		return err
	}

	tItems := table.FivenetJobGroupRuleQualificationItems
	for _, qualificationID := range qualification.GetQualificationIds() {
		stmt := tItems.
			INSERT(tItems.RuleID, tItems.QualificationID).
			VALUES(rule.GetId(), qualificationID)

		_, err := stmt.ExecContext(ctx, db)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Store) deleteGroupRuleDetails(ctx context.Context, db qrm.DB, ruleID int64) error {
	tItems := table.FivenetJobGroupRuleQualificationItems
	if _, err := tItems.
		DELETE().
		WHERE(tItems.RuleID.EQ(mysql.Int64(ruleID))).
		ExecContext(ctx, db); err != nil {
		return err
	}

	tQualifications := table.FivenetJobGroupRuleQualifications
	if _, err := tQualifications.
		DELETE().
		WHERE(tQualifications.RuleID.EQ(mysql.Int64(ruleID))).
		ExecContext(ctx, db); err != nil {
		return err
	}

	tGrades := table.FivenetJobGroupRuleGrades
	if _, err := tGrades.
		DELETE().
		WHERE(tGrades.RuleID.EQ(mysql.Int64(ruleID))).
		ExecContext(ctx, db); err != nil {
		return err
	}

	return nil
}

func (s *Store) ListGroupRuleMemberMatches(
	ctx context.Context,
	db qrm.DB,
	group *jobsgroups.Group,
	search string,
) ([]*GroupRuleMemberMatch, error) {
	rules, err := s.ListGroupRules(ctx, db, group.GetId())
	if err != nil {
		return nil, err
	}

	searchCondition := groupMemberSearchCondition(search, table.FivenetUser.AS("u"))
	matches := []*GroupRuleMemberMatch{}
	for _, rule := range rules {
		if !rule.GetEnabled() {
			continue
		}

		var ruleMatches []*GroupRuleMemberMatch
		if grade := rule.GetGrade(); grade != nil {
			ruleMatches, err = s.listGradeRuleMemberMatches(
				ctx,
				db,
				group,
				rule.GetId(),
				grade,
				searchCondition,
			)
		} else if qualification := rule.GetQualification(); qualification != nil {
			ruleMatches, err = s.listQualificationRuleMemberMatches(
				ctx,
				db,
				group,
				rule.GetId(),
				qualification,
				searchCondition,
			)
		}
		if err != nil {
			return nil, err
		}

		matches = append(matches, ruleMatches...)
	}

	return matches, nil
}

func (s *Store) listGradeRuleMemberMatches(
	ctx context.Context,
	db qrm.DB,
	group *jobsgroups.Group,
	ruleID int64,
	rule *jobsgroups.GroupGradeRule,
	searchCondition mysql.BoolExpression,
) ([]*GroupRuleMemberMatch, error) {
	tUserJobs := table.FivenetUserJobs.AS("uj")
	tUser := table.FivenetUser.AS("u")
	condition := tUserJobs.Job.EQ(mysql.String(group.GetJob()))
	if searchCondition != nil {
		condition = condition.AND(searchCondition)
	}

	var label string
	switch rule.GetType() {
	case jobsgroups.GroupGradeRuleType_GROUP_GRADE_RULE_TYPE_MINIMUM:
		condition = condition.AND(tUserJobs.Grade.GT_EQ(nullableInt32Expression(rule.Grade)))
		label = fmt.Sprintf("Grade >= %d", rule.GetGrade())
	case jobsgroups.GroupGradeRuleType_GROUP_GRADE_RULE_TYPE_EXACT:
		condition = condition.AND(tUserJobs.Grade.EQ(nullableInt32Expression(rule.Grade)))
		label = fmt.Sprintf("Grade = %d", rule.GetGrade())
	case jobsgroups.GroupGradeRuleType_GROUP_GRADE_RULE_TYPE_RANGE:
		condition = condition.AND(
			tUserJobs.Grade.BETWEEN(
				nullableInt32Expression(rule.MinGrade),
				nullableInt32Expression(rule.MaxGrade),
			),
		)
		label = fmt.Sprintf("Grade %d-%d", rule.GetMinGrade(), rule.GetMaxGrade())
	default:
		return nil, fmt.Errorf("unsupported grade rule type %d", rule.GetType())
	}

	stmt := tUserJobs.
		SELECT(
			mysql.Raw("DISTINCT(`uj`.`user_id`) AS `group_rule_member_match.user_id`"),
		).
		FROM(tUserJobs.
			INNER_JOIN(tUser,
				tUser.ID.EQ(tUserJobs.UserID),
			),
		).
		WHERE(condition).
		ORDER_BY(tUserJobs.UserID.ASC())

	matches := []*GroupRuleMemberMatch{}
	if err := stmt.QueryContext(ctx, db, &matches); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return []*GroupRuleMemberMatch{}, nil
		}
		return nil, err
	}

	return finalizeRuleMemberMatches(matches, group.GetId(), ruleID, label), nil
}

func (s *Store) listQualificationRuleMemberMatches(
	ctx context.Context,
	db qrm.DB,
	group *jobsgroups.Group,
	ruleID int64,
	rule *jobsgroups.GroupQualificationRule,
	searchCondition mysql.BoolExpression,
) ([]*GroupRuleMemberMatch, error) {
	if len(rule.GetQualificationIds()) == 0 {
		return nil, nil
	}

	tUserJobs := table.FivenetUserJobs.AS("uj")
	tUser := table.FivenetUser.AS("u")
	tItems := table.FivenetJobGroupRuleQualificationItems.AS("qi")
	condition := mysql.AND(
		tUserJobs.Job.EQ(mysql.String(group.GetJob())),
		tItems.RuleID.EQ(mysql.Int64(ruleID)),
		tItems.QualificationID.IN(int64Expressions(rule.GetQualificationIds())...),
	)
	if searchCondition != nil {
		condition = condition.AND(searchCondition)
	}
	var from mysql.ReadableTable
	if rule.GetRequireCompleted() {
		tSuccess := table.FivenetQualificationsResultSuccessMap.AS("qr")
		from = tUserJobs.
			INNER_JOIN(tUser, tUser.ID.EQ(tUserJobs.UserID)).
			INNER_JOIN(tItems,
				mysql.AND(
					tItems.RuleID.EQ(mysql.Int64(ruleID)),
					tItems.QualificationID.IN(int64Expressions(rule.GetQualificationIds())...),
				),
			).
			INNER_JOIN(tSuccess,
				mysql.AND(
					tSuccess.UserID.EQ(tUserJobs.UserID),
					tSuccess.QualificationID.EQ(tItems.QualificationID),
				),
			)
	} else {
		tResults := table.FivenetQualificationsResults.AS("qr")
		from = tUserJobs.
			INNER_JOIN(tUser, tUser.ID.EQ(tUserJobs.UserID)).
			INNER_JOIN(tItems, mysql.AND(
				tItems.RuleID.EQ(mysql.Int64(ruleID)),
				tItems.QualificationID.IN(int64Expressions(rule.GetQualificationIds())...),
			)).
			INNER_JOIN(tResults, mysql.AND(
				tResults.UserID.EQ(tUserJobs.UserID),
				tResults.QualificationID.EQ(tItems.QualificationID),
				tResults.DeletedAt.IS_NULL(),
			))
	}

	stmt := tUserJobs.
		SELECT(tUserJobs.UserID.AS("group_rule_member_match.user_id")).
		FROM(from).
		WHERE(condition).
		GROUP_BY(tUserJobs.UserID).
		ORDER_BY(tUserJobs.UserID.ASC())

	matchCount := mysql.COUNT(mysql.DISTINCT(tItems.QualificationID))
	label := "Matches any qualification"
	switch rule.GetType() {
	case jobsgroups.GroupQualificationRuleType_GROUP_QUALIFICATION_RULE_TYPE_ANY:
		stmt = stmt.HAVING(matchCount.GT(mysql.Int(0)))
	case jobsgroups.GroupQualificationRuleType_GROUP_QUALIFICATION_RULE_TYPE_ALL:
		stmt = stmt.HAVING(
			matchCount.EQ(mysql.Int(int64(len(rule.GetQualificationIds())))),
		)
		label = "Matches all qualifications"
	default:
		return nil, fmt.Errorf("unsupported qualification rule type %d", rule.GetType())
	}
	if rule.GetRequireCompleted() {
		label += " (completed)"
	}

	matches := []*GroupRuleMemberMatch{}
	if err := stmt.QueryContext(ctx, db, &matches); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return []*GroupRuleMemberMatch{}, nil
		}
		return nil, err
	}

	return finalizeRuleMemberMatches(matches, group.GetId(), ruleID, label), nil
}
