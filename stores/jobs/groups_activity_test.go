package jobsstore

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	jobsgroups "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/groups"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreCreateGroupActivity(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	reason := "initial setup"
	targetUserID := int32(55)
	ruleID := int64(77)
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO fivenet_job_group_activity`)).
		WithArgs(
			"police",
			int64(42),
			int32(jobsgroups.GroupActivityType_GROUP_ACTIVITY_TYPE_RULE_ADDED),
			int64(100),
			targetUserID,
			ruleID,
			reason,
			nil,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := store.CreateGroupActivity(t.Context(), store.db, &jobsgroups.GroupActivity{
		Job:          "police",
		GroupId:      42,
		Type:         jobsgroups.GroupActivityType_GROUP_ACTIVITY_TYPE_RULE_ADDED,
		ActorUserId:  new(int32(100)),
		TargetUserId: &targetUserID,
		RuleId:       &ruleID,
		Reason:       &reason,
	})
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreCountGroupActivity(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_job_group_activity`)).
		WithArgs("police").
		WillReturnRows(sqlmock.NewRows([]string{"data_count.total"}).AddRow(int64(3)))

	count, err := store.CountGroupActivity(t.Context(), store.db, ListQuery{Job: "police"})
	require.NoError(t, err)
	assert.Equal(t, int64(3), count)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreListGroupActivity(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	createdAt := time.Date(2026, time.August, 1, 9, 30, 0, 0, time.UTC)
	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_job_group_activity`)).
		WithArgs("police", int64(20), int64(5)).
		WillReturnRows(sqlmock.NewRows([]string{
			"group_activity.id",
			"group_activity.job",
			"group_activity.group_id",
			"group_activity.type",
			"group_activity.actor_user_id",
			"group_activity.target_user_id",
			"group_activity.rule_id",
			"group_activity.reason",
			"group_activity.data",
			"group_activity.created_at",
		}).AddRow(
			int64(9),
			"police",
			int64(42),
			int32(jobsgroups.GroupActivityType_GROUP_ACTIVITY_TYPE_MEMBER_ADDED),
			int64(100),
			int64(55),
			nil,
			"manual add",
			nil,
			createdAt,
		))

	activity, err := store.ListGroupActivity(t.Context(), store.db, ListQuery{
		Job:    "police",
		Offset: 5,
		Limit:  20,
	})
	require.NoError(t, err)
	require.Len(t, activity, 1)
	assert.Equal(t, int64(9), activity[0].GetId())
	assert.Equal(t, int64(42), activity[0].GetGroupId())
	assert.Equal(t, int32(100), activity[0].GetActorUserId())
	assert.Equal(t, int32(55), activity[0].GetTargetUserId())
	assert.Equal(t, "manual add", activity[0].GetReason())
	assert.Equal(
		t,
		jobsgroups.GroupActivityType_GROUP_ACTIVITY_TYPE_MEMBER_ADDED,
		activity[0].GetType(),
	)
	require.NoError(t, mock.ExpectationsWereMet())
}
