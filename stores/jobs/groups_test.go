package jobsstore

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	jobsgroups "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/groups"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreListGroups(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_job_groups AS job_group`)).
		WillReturnRows(sqlmock.NewRows([]string{}))

	groups, err := store.ListGroups(t.Context(), store.db, GroupsQuery{
		JobID: 7,
	})
	require.NoError(t, err)
	assert.Empty(t, groups)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreGetGroup(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_job_groups AS job_group`)).
		WithArgs(int64(7), int64(42), int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{}))

	group, err := store.GetGroup(t.Context(), store.db, GroupQuery{
		JobID:           7,
		IncludeArchived: true,
	}, 42)
	require.NoError(t, err)
	assert.Nil(t, group)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStoreCreateArchiveRestoreGroup(t *testing.T) {
	t.Parallel()

	store, mock := newTestStore(t)

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO fivenet_job_groups AS job_group`)).
		WithArgs(
			int64(7),
			"K9 Unit",
			"Certified handlers and support staff.",
			"K9",
			"file-1",
			"#123456",
			int32(jobsgroups.GroupType_GROUP_TYPE_MANUAL),
			int32(jobsgroups.GroupState_GROUP_STATE_ACTIVE),
			int32(jobsgroups.GroupMembershipMode_GROUP_MEMBERSHIP_MODE_FLEXIBLE),
			int32(5),
			int64(99),
			int64(100),
		).
		WillReturnResult(sqlmock.NewResult(42, 1))

	id, err := store.CreateGroup(t.Context(), store.db, &jobsgroups.Group{
		JobId:           7,
		Name:            "K9 Unit",
		Description:     stringPtr("Certified handlers and support staff."),
		ShortName:       stringPtr("K9"),
		LogoFileId:      stringPtr("file-1"),
		Color:           stringPtr("#123456"),
		Type:            jobsgroups.GroupType_GROUP_TYPE_MANUAL,
		State:           jobsgroups.GroupState_GROUP_STATE_ACTIVE,
		MembershipMode:  jobsgroups.GroupMembershipMode_GROUP_MEMBERSHIP_MODE_FLEXIBLE,
		SortOrder:       5,
		CreatedByUserId: 99,
		UpdatedByUserId: 100,
	})
	require.NoError(t, err)
	assert.Equal(t, int64(42), id)

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE fivenet_job_groups AS job_group SET`)).
		WithArgs(int32(jobsgroups.GroupState_GROUP_STATE_ARCHIVED), int64(100), int64(42), int64(7), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	require.NoError(t, store.ArchiveGroup(t.Context(), store.db, 7, 42, 100))

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE fivenet_job_groups AS job_group SET`)).
		WithArgs(int32(jobsgroups.GroupState_GROUP_STATE_ACTIVE), int64(100), int64(42), int64(7), int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	require.NoError(t, store.RestoreGroup(t.Context(), store.db, 7, 42, 100))
	require.NoError(t, mock.ExpectationsWereMet())
}

func stringPtr(v string) *string {
	return &v
}
