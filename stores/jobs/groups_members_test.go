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

func TestStoreListGroupManualMembers(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	createdAt := time.Date(2026, time.July, 31, 12, 0, 0, 0, time.UTC)
	mock.ExpectQuery(`(?s)FROM fivenet_job_group_manual_members AS mm .*INNER JOIN fivenet_job_groups AS g .*INNER JOIN fivenet_user_jobs AS uj`).
		WithArgs(int64(42)).
		WillReturnRows(sqlmock.NewRows([]string{
			"group_manual_member.group_id",
			"group_manual_member.user_id",
			"group_manual_member.reason",
			"group_manual_member.created_by_user_id",
			"group_manual_member.created_at",
		}).AddRow(int64(42), int64(7), "onboarding", int64(100), createdAt))

	members, err := store.ListGroupManualMembers(
		t.Context(),
		store.db,
		GroupItemsQuery{GroupID: 42},
	)
	require.NoError(t, err)
	require.Len(t, members, 1)
	assert.Equal(t, int64(42), members[0].GetGroupId())
	assert.Equal(t, int32(7), members[0].GetUserId())
	assert.Equal(t, "onboarding", members[0].GetReason())
	assert.Equal(t, int32(100), members[0].GetCreatedByUserId())
	assert.Equal(t, createdAt, members[0].GetCreatedAt().AsTime())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreUserInJob(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	mock.ExpectQuery(`(?s)FROM fivenet_user_jobs .*INNER JOIN fivenet_user ON .*fivenet_user\.deleted_at IS NULL.*WHERE .*fivenet_user_jobs\.user_id = \?.*fivenet_user_jobs\.job = \?.*LIMIT \?.*`).
		WithArgs(sqlmock.AnyArg(), int32(7), "police", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"found"}).AddRow(1))

	ok, err := store.UserInJob(t.Context(), store.db, "police", 7)
	require.NoError(t, err)
	assert.True(t, ok)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreUserInJobSkipsDeletedUsers(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	mock.ExpectQuery(`(?s)FROM fivenet_user_jobs .*INNER JOIN fivenet_user ON .*fivenet_user\.deleted_at IS NULL.*`).
		WithArgs(sqlmock.AnyArg(), int32(7), "police", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"found"}))

	ok, err := store.UserInJob(t.Context(), store.db, "police", 7)
	require.NoError(t, err)
	assert.False(t, ok)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreRecountGroupStats(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	mock.ExpectQuery(regexp.QuoteMeta("FROM fivenet_job_groups AS `group`")).
		WithArgs(int64(42), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"type",
			"job",
			"membership_mode",
		}).AddRow(
			int32(jobsgroups.GroupType_GROUP_TYPE_MIXED),
			"police",
			int32(jobsgroups.GroupMembershipMode_GROUP_MEMBERSHIP_MODE_FLEXIBLE),
		))
	mock.ExpectQuery(`(?s)FROM fivenet_job_group_manual_members AS mm .*INNER JOIN fivenet_job_groups AS g .*INNER JOIN fivenet_user_jobs AS uj`).
		WithArgs(int64(42)).
		WillReturnRows(sqlmock.NewRows([]string{
			"group_manual_member.group_id",
			"group_manual_member.user_id",
			"group_manual_member.reason",
			"group_manual_member.created_by_user_id",
			"group_manual_member.created_at",
		}))
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
		}))
	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_job_group_member_exclusions AS me INNER JOIN fivenet_user AS u`)).
		WithArgs(int64(42)).
		WillReturnRows(sqlmock.NewRows([]string{
			"group_member_exclusion.group_id",
			"group_member_exclusion.user_id",
			"group_member_exclusion.reason_type",
			"group_member_exclusion.reason",
			"group_member_exclusion.created_by_user_id",
			"group_member_exclusion.created_at",
		}))
	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_job_group_leaders`)).
		WithArgs(int64(42)).
		WillReturnRows(sqlmock.NewRows([]string{"data_count.total"}).AddRow(int64(0)))
	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_job_group_rules`)).
		WithArgs(int64(42)).
		WillReturnRows(sqlmock.NewRows([]string{"data_count.total"}).AddRow(int64(0)))
	mock.ExpectExec(regexp.QuoteMeta("UPDATE fivenet_job_groups AS `group` SET")).
		WithArgs(int64(0), int64(0), int64(0), int64(0), int64(42)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	require.NoError(t, store.RecountGroupStats(t.Context(), store.db, 42))
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreAddGroupManualMember(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	createdAt := time.Date(2026, time.July, 31, 12, 1, 0, 0, time.UTC)
	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_job_group_manual_members`)).
		WithArgs(int64(42), int64(7), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"group_manual_member.group_id",
			"group_manual_member.user_id",
			"group_manual_member.reason",
			"group_manual_member.created_by_user_id",
			"group_manual_member.created_at",
		}))
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO fivenet_job_group_manual_members`)).
		WithArgs(int64(42), int64(7), "onboarding", int64(100)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_job_group_manual_members`)).
		WithArgs(int64(42), int64(7), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"group_manual_member.group_id",
			"group_manual_member.user_id",
			"group_manual_member.reason",
			"group_manual_member.created_by_user_id",
			"group_manual_member.created_at",
		}).AddRow(int64(42), int64(7), "onboarding", int64(100), createdAt))

	member, created, err := store.AddGroupManualMember(
		t.Context(),
		store.db,
		42,
		7,
		100,
		new("onboarding"),
	)
	require.NoError(t, err)
	require.NotNil(t, member)
	assert.True(t, created)
	assert.Equal(t, int64(42), member.GetGroupId())
	assert.Equal(t, int32(7), member.GetUserId())
	assert.Equal(t, "onboarding", member.GetReason())
	assert.Equal(t, int32(100), member.GetCreatedByUserId())
	assert.Equal(t, createdAt, member.GetCreatedAt().AsTime())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreAddGroupMemberExclusion(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	createdAt := time.Date(2026, time.July, 31, 12, 2, 0, 0, time.UTC)
	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_job_group_member_exclusions`)).
		WithArgs(int64(42), int64(7), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"group_member_exclusion.group_id",
			"group_member_exclusion.user_id",
			"group_member_exclusion.reason_type",
			"group_member_exclusion.reason",
			"group_member_exclusion.created_by_user_id",
			"group_member_exclusion.created_at",
		}))
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO fivenet_job_group_member_exclusions`)).
		WithArgs(int64(42), int64(7), int32(jobsgroups.GroupExclusionReason_GROUP_EXCLUSION_REASON_TEMPORARY), "temp reassignment", int64(100)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_job_group_member_exclusions`)).
		WithArgs(int64(42), int64(7), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"group_member_exclusion.group_id",
			"group_member_exclusion.user_id",
			"group_member_exclusion.reason_type",
			"group_member_exclusion.reason",
			"group_member_exclusion.created_by_user_id",
			"group_member_exclusion.created_at",
		}).AddRow(int64(42), int64(7), int32(jobsgroups.GroupExclusionReason_GROUP_EXCLUSION_REASON_TEMPORARY), "temp reassignment", int64(100), createdAt))

	exclusion, created, err := store.AddGroupMemberExclusion(
		t.Context(),
		store.db,
		42,
		7,
		jobsgroups.GroupExclusionReason_GROUP_EXCLUSION_REASON_TEMPORARY,
		100,
		new("temp reassignment"),
	)
	require.NoError(t, err)
	require.NotNil(t, exclusion)
	assert.True(t, created)
	assert.Equal(t, int64(42), exclusion.GetGroupId())
	assert.Equal(t, int32(7), exclusion.GetUserId())
	assert.Equal(
		t,
		jobsgroups.GroupExclusionReason_GROUP_EXCLUSION_REASON_TEMPORARY,
		exclusion.GetReasonType(),
	)
	assert.Equal(t, "temp reassignment", exclusion.GetReason())
	assert.Equal(t, int32(100), exclusion.GetCreatedByUserId())
	assert.Equal(t, createdAt, exclusion.GetCreatedAt().AsTime())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreAddGroupLeader(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	createdAt := time.Date(2026, time.July, 31, 12, 4, 0, 0, time.UTC)
	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_job_group_leaders`)).
		WithArgs(int64(42), int64(9), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"group_leader.group_id",
			"group_leader.user_id",
			"group_leader.created_by_user_id",
			"group_leader.created_at",
		}))
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO fivenet_job_group_leaders`)).
		WithArgs(int64(42), int64(9), int64(100)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_job_group_leaders`)).
		WithArgs(int64(42), int64(9), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"group_leader.group_id",
			"group_leader.user_id",
			"group_leader.created_by_user_id",
			"group_leader.created_at",
		}).AddRow(int64(42), int64(9), int64(100), createdAt))

	leader, created, err := store.AddGroupLeader(t.Context(), store.db, 42, 9, 100)
	require.NoError(t, err)
	require.NotNil(t, leader)
	assert.True(t, created)
	assert.Equal(t, int64(42), leader.GetGroupId())
	assert.Equal(t, int32(9), leader.GetUserId())
	assert.Equal(t, int32(100), leader.GetCreatedByUserId())
	assert.Equal(t, createdAt, leader.GetCreatedAt().AsTime())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreRemoveGroupLeader(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM fivenet_job_group_leaders`)).
		WithArgs(int64(42), int64(9)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	require.NoError(t, store.RemoveGroupLeader(t.Context(), store.db, 42, 9))
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreListGroupLeaders(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	createdAt := time.Date(2026, time.July, 31, 12, 3, 0, 0, time.UTC)
	mock.ExpectQuery(`(?s)FROM fivenet_job_group_leaders AS gl .*INNER JOIN fivenet_job_groups AS g .*INNER JOIN fivenet_user_jobs AS uj`).
		WithArgs(int64(42)).
		WillReturnRows(sqlmock.NewRows([]string{
			"group_leader.group_id",
			"group_leader.user_id",
			"group_leader.created_by_user_id",
			"group_leader.created_at",
		}).AddRow(int64(42), int64(9), int64(100), createdAt))

	leaders, err := store.ListGroupLeaders(t.Context(), store.db, GroupItemsQuery{GroupID: 42})
	require.NoError(t, err)
	require.Len(t, leaders, 1)
	assert.Equal(t, int64(42), leaders[0].GetGroupId())
	assert.Equal(t, int32(9), leaders[0].GetUserId())
	assert.Equal(t, int32(100), leaders[0].GetCreatedByUserId())
	assert.Equal(t, createdAt, leaders[0].GetCreatedAt().AsTime())
	require.NoError(t, mock.ExpectationsWereMet())
}
