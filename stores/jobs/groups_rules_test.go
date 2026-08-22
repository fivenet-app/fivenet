package jobsstore

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	jobsgroups "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/groups"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreListGroupRulesGrade(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	createdAt := time.Date(2026, time.July, 31, 12, 0, 0, 0, time.UTC)
	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_job_group_rules`)).
		WithArgs(int64(42)).
		WillReturnRows(sqlmock.NewRows([]string{
			"group_rule_builder.id",
			"group_rule_builder.group_id",
			"group_rule_builder.type",
			"group_rule_builder.enabled",
			"group_rule_builder.created_by_user_id",
			"group_rule_builder.created_at",
			"group_rule_builder.updated_at",
		}).AddRow(int64(7), int64(42), int32(groupRuleTypeGrade), true, int64(100), createdAt, nil))

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_job_group_rule_grades`)).
		WithArgs(int64(7), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"group_grade_rule.type",
			"group_grade_rule.grade",
			"group_grade_rule.min_grade",
			"group_grade_rule.max_grade",
		}).AddRow(int32(jobsgroups.GroupGradeRuleType_GROUP_GRADE_RULE_TYPE_MINIMUM), int32(3), nil, nil))

	rules, err := store.ListGroupRules(t.Context(), store.db, GroupItemsQuery{GroupID: 42})
	require.NoError(t, err)
	require.Len(t, rules, 1)
	assert.Equal(t, int64(7), rules[0].GetId())
	assert.Equal(t, int64(42), rules[0].GetGroupId())
	assert.True(t, rules[0].GetEnabled())
	require.NotNil(t, rules[0].GetGrade())
	assert.Equal(
		t,
		jobsgroups.GroupGradeRuleType_GROUP_GRADE_RULE_TYPE_MINIMUM,
		rules[0].GetGrade().GetType(),
	)
	assert.Equal(t, int32(3), rules[0].GetGrade().GetGrade())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreCreateGroupRuleQualification(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	rule := &jobsgroups.GroupRule{
		GroupId:         42,
		Enabled:         true,
		CreatedByUserId: new(int32(100)),
	}
	rule.SetQualification(&jobsgroups.GroupQualificationRule{
		Type:             jobsgroups.GroupQualificationRuleType_GROUP_QUALIFICATION_RULE_TYPE_ANY,
		QualificationIds: []int64{10, 11},
		RequireCompleted: true,
	})

	createdAt := time.Date(2026, time.July, 31, 12, 1, 0, 0, time.UTC)
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO fivenet_job_group_rules`)).
		WithArgs(int64(42), int32(groupRuleTypeQualification), true, int64(100)).
		WillReturnResult(sqlmock.NewResult(9, 1))

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO fivenet_job_group_rule_qualifications`)).
		WithArgs(int64(9), int32(jobsgroups.GroupQualificationRuleType_GROUP_QUALIFICATION_RULE_TYPE_ANY), true).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO fivenet_job_group_rule_qualification_items`)).
		WithArgs(int64(9), int64(10)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO fivenet_job_group_rule_qualification_items`)).
		WithArgs(int64(9), int64(11)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_job_group_rules`)).
		WithArgs(int64(42), int64(9), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"group_rule_builder.id",
			"group_rule_builder.group_id",
			"group_rule_builder.type",
			"group_rule_builder.enabled",
			"group_rule_builder.created_by_user_id",
			"group_rule_builder.created_at",
			"group_rule_builder.updated_at",
		}).AddRow(int64(9), int64(42), int32(groupRuleTypeQualification), true, int64(100), createdAt, nil))

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_job_group_rule_qualifications`)).
		WithArgs(int64(9), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"group_qualification_rule.type",
			"group_qualification_rule.require_completed",
		}).AddRow(int32(jobsgroups.GroupQualificationRuleType_GROUP_QUALIFICATION_RULE_TYPE_ANY), true))

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_job_group_rule_qualification_items`)).
		WithArgs(int64(9)).
		WillReturnRows(sqlmock.NewRows([]string{"qualification_id"}).AddRow(int64(10)).AddRow(int64(11)))

	created, err := store.CreateGroupRule(t.Context(), store.db, rule)
	require.NoError(t, err)
	require.NotNil(t, created)
	assert.Equal(t, int64(9), created.GetId())
	require.NotNil(t, created.GetQualification())
	assert.Equal(t, []int64{10, 11}, created.GetQualification().GetQualificationIds())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreDeleteGroupRuleNotFound(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM fivenet_job_group_rules`)).
		WithArgs(int64(42), int64(9)).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := store.DeleteGroupRule(t.Context(), store.db, 42, 9)
	require.ErrorIs(t, err, sql.ErrNoRows)
}

func TestStoreListGroupRuleMemberMatchesGrade(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	createdAt := time.Date(2026, time.July, 31, 12, 2, 0, 0, time.UTC)
	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_job_group_rules`)).
		WithArgs(int64(42)).
		WillReturnRows(sqlmock.NewRows([]string{
			"group_rule_builder.id",
			"group_rule_builder.group_id",
			"group_rule_builder.type",
			"group_rule_builder.enabled",
			"group_rule_builder.created_by_user_id",
			"group_rule_builder.created_at",
			"group_rule_builder.updated_at",
		}).AddRow(int64(7), int64(42), int32(groupRuleTypeGrade), true, int64(100), createdAt, nil))

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_job_group_rule_grades`)).
		WithArgs(int64(7), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"group_grade_rule.type",
			"group_grade_rule.grade",
			"group_grade_rule.min_grade",
			"group_grade_rule.max_grade",
		}).AddRow(int32(jobsgroups.GroupGradeRuleType_GROUP_GRADE_RULE_TYPE_MINIMUM), int32(3), nil, nil))

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_user_jobs AS uj INNER JOIN fivenet_user AS u`)).
		WithArgs("police", int32(3)).
		WillReturnRows(sqlmock.NewRows([]string{"group_rule_member_match.user_id"}).AddRow(int64(10)).AddRow(int64(11)))

	matches, err := store.ListGroupRuleMemberMatches(t.Context(), store.db, &jobsgroups.Group{
		Id:  42,
		Job: "police",
	}, "")
	require.NoError(t, err)
	require.Len(t, matches, 2)
	assert.Equal(t, int32(10), matches[0].UserID)
	assert.Equal(t, int64(7), matches[0].RuleID)
	assert.Equal(t, "Grade >= 3", matches[0].Label)
	assert.Equal(t, int32(11), matches[1].UserID)
	require.NoError(t, mock.ExpectationsWereMet())
}
