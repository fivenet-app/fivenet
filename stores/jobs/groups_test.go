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

func TestStoreListGroups(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_job_groups AS job_group LEFT JOIN fivenet_files AS logo_file`)).
		WithArgs("police", int32(jobsgroups.GroupState_GROUP_STATE_ACTIVE), int64(0), int64(0)).
		WillReturnRows(sqlmock.NewRows([]string{}))

	groups, err := store.ListGroups(t.Context(), store.db, GroupsQuery{
		Job: "police",
	})
	require.NoError(t, err)
	assert.Empty(t, groups)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreGetGroup(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_job_groups AS job_group LEFT JOIN fivenet_files AS logo_file`)).
		WithArgs("police", int64(42), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{}))

	group, err := store.GetGroup(t.Context(), store.db, GroupQuery{
		Job:             "police",
		IncludeArchived: true,
	}, 42)
	require.NoError(t, err)
	assert.Nil(t, group)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreGetGroupWithLogo(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)
	now := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_job_groups AS job_group LEFT JOIN fivenet_files AS logo_file`)).
		WithArgs("police", int64(42), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"group.id",
			"group.job",
			"group.name",
			"group.description",
			"group.short_name",
			"group.logo_file_id",
			"group.color",
			"group.type",
			"group.state",
			"group.membership_mode",
			"group.sort_rank",
			"group.members_count",
			"group.leaders_count",
			"group.rules_count",
			"group.exclusions_count",
			"group.created_by_user_id",
			"group.updated_by_user_id",
			"group.created_at",
			"group.updated_at",
			"group.deleted_at",
			"logo_file.id",
			"logo_file.file_path",
			"logo_file.byte_size",
			"logo_file.content_type",
			"logo_file.created_at",
		}).AddRow(
			int64(42),
			"police",
			"K9 Unit",
			"Certified handlers and support staff.",
			"K9",
			int64(5),
			"#123456",
			int32(jobsgroups.GroupType_GROUP_TYPE_MANUAL),
			int32(jobsgroups.GroupState_GROUP_STATE_ACTIVE),
			int32(jobsgroups.GroupMembershipMode_GROUP_MEMBERSHIP_MODE_FLEXIBLE),
			"0|zzzzzz:",
			int32(7),
			int32(2),
			int32(1),
			int32(0),
			int64(99),
			int64(100),
			now,
			now,
			nil,
			int64(5),
			"jobgrouplogos/police/42.png",
			int64(1024),
			"image/png",
			now,
		))

	group, err := store.GetGroup(t.Context(), store.db, GroupQuery{
		Job:             "police",
		IncludeArchived: true,
	}, 42)
	require.NoError(t, err)
	require.NotNil(t, group)
	assert.Equal(t, int64(42), group.GetId())
	assert.Equal(t, "0|zzzzzz:", group.GetSortRank())
	assert.Equal(t, int64(5), group.GetLogoFileId())
	require.NotNil(t, group.GetLogoFile())
	assert.Equal(t, int64(5), group.GetLogoFile().GetId())
	assert.Equal(t, "jobgrouplogos/police/42.png", group.GetLogoFile().GetFilePath())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreCreateArchiveRestoreGroup(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO fivenet_job_groups`)).
		WithArgs(
			"police",
			"K9 Unit",
			"Certified handlers and support staff.",
			"K9",
			int64(5),
			"#123456",
			int32(jobsgroups.GroupType_GROUP_TYPE_MANUAL),
			int32(jobsgroups.GroupState_GROUP_STATE_ACTIVE),
			int32(jobsgroups.GroupMembershipMode_GROUP_MEMBERSHIP_MODE_FLEXIBLE),
			"0|zzzzzz:",
			int64(99),
			int64(100),
		).
		WillReturnResult(sqlmock.NewResult(42, 1))

	id, err := store.CreateGroup(t.Context(), store.db, &jobsgroups.Group{
		Job:             "police",
		Name:            "K9 Unit",
		Description:     stringPtr("Certified handlers and support staff."),
		ShortName:       stringPtr("K9"),
		LogoFileId:      int64Ptr(5),
		Color:           stringPtr("#123456"),
		Type:            jobsgroups.GroupType_GROUP_TYPE_MANUAL,
		State:           jobsgroups.GroupState_GROUP_STATE_ACTIVE,
		MembershipMode:  jobsgroups.GroupMembershipMode_GROUP_MEMBERSHIP_MODE_FLEXIBLE,
		SortRank:        "0|zzzzzz:",
		CreatedByUserId: new(int32(99)),
		UpdatedByUserId: new(int32(100)),
	})
	require.NoError(t, err)
	assert.Equal(t, int64(42), id)

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE fivenet_job_groups SET`)).
		WithArgs(int32(jobsgroups.GroupState_GROUP_STATE_ARCHIVED), int64(100), int64(42), "police", int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	require.NoError(t, store.ArchiveGroup(t.Context(), store.db, "police", 42, 100))

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE fivenet_job_groups SET`)).
		WithArgs(int32(jobsgroups.GroupState_GROUP_STATE_ACTIVE), int64(100), int64(42), "police", int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	require.NoError(t, store.RestoreGroup(t.Context(), store.db, "police", 42, 100))
	require.NoError(t, mock.ExpectationsWereMet())
}

func stringPtr(v string) *string {
	return &v
}

func int64Ptr(v int64) *int64 {
	return &v
}
