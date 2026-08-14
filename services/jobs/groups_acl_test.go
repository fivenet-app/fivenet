package jobs

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	resourcesaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/access"
	jobscolleagues "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/colleagues"
	jobsgroups "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/groups"
	pbuserinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pbjobs "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/jobs"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	grpcauth "github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/mstlystcdata"
	colleaguehydrator "github.com/fivenet-app/fivenet/v2026/services/jobs/colleagues"
	errorsjobs "github.com/fivenet-app/fivenet/v2026/services/jobs/errors"
	jobsstore "github.com/fivenet-app/fivenet/v2026/stores/jobs"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/stretchr/testify/require"
)

type stubGroupAccess struct {
	allowed          bool
	replaceCallCount int
	replacedTargetID int64
	replacedAccess   *resourcesaccess.Access
}

func (s *stubGroupAccess) CanUserAccessTarget(
	_ context.Context,
	_ int64,
	_ *pbuserinfo.UserInfo,
	_ int32,
) (bool, error) {
	return s.allowed, nil
}

func (s *stubGroupAccess) CanUserAccessTargetIncludingDeleted(
	_ context.Context,
	_ int64,
	_ *pbuserinfo.UserInfo,
	_ int32,
) (bool, error) {
	return s.allowed, nil
}

func (s *stubGroupAccess) CanUserAccessTargetIDs(
	_ context.Context,
	_ *pbuserinfo.UserInfo,
	_ int32,
	targetIDs ...int64,
) ([]int64, error) {
	if s.allowed {
		return targetIDs, nil
	}
	return nil, nil
}

func (s *stubGroupAccess) ListTargetAccess(
	_ context.Context,
	_ qrm.DB,
	_ int64,
	_ access.SubjectAccessOptions,
) (*resourcesaccess.Access, error) {
	return &resourcesaccess.Access{}, nil
}

func (s *stubGroupAccess) ReplaceTargetAccess(
	_ context.Context,
	_ qrm.DB,
	_ *access.SubjectResolver,
	targetID int64,
	access *resourcesaccess.Access,
	_ access.SubjectAccessOptions,
) (*access.SubjectAccessChanges, error) {
	s.replaceCallCount++
	s.replacedTargetID = targetID
	s.replacedAccess = access

	return nil, nil
}

func (s *stubGroupAccess) VisibleIDsByConditionQuery(
	_ *pbuserinfo.UserInfo,
	_ int32,
	_ bool,
	_ mysql.BoolExpression,
) access.VisibilityQuery {
	return access.VisibilityQuery{}
}

type noopColleagueHydrator struct{}

func (noopColleagueHydrator) ListByUserID(
	_ context.Context,
	_ qrm.DB,
	_ *pbuserinfo.UserInfo,
	_ string,
	_ []int32,
	_ bool,
) ([]*jobscolleagues.Colleague, error) {
	return nil, nil
}

func (noopColleagueHydrator) HydrateByUserID(
	_ context.Context,
	_ qrm.DB,
	_ *pbuserinfo.UserInfo,
	_ string,
	_ []int32,
	_ bool,
) (map[int32]*jobscolleagues.Colleague, error) {
	return map[int32]*jobscolleagues.Colleague{}, nil
}

func (noopColleagueHydrator) HydrateTargets(
	_ context.Context,
	_ qrm.DB,
	_ *pbuserinfo.UserInfo,
	_ string,
	_ []colleaguehydrator.Target,
	_ bool,
) error {
	return nil
}

func newJobsGroupACLTestServer(
	t *testing.T,
	groupAccess access.JobGroupsAccess,
) (*Server, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	if groupAccess == nil {
		groupAccess = access.NewJobGroupsSubjectObjectAccess(db)
	}
	storeAccess := access.NewJobGroupsSubjectObjectAccess(db)

	return &Server{
		db:                  db,
		store:               jobsstore.New(db, &config.CustomDB{}, storeAccess).Store,
		groupAccess:         groupAccess,
		groupAccessResolver: access.NewSubjectResolver(db),
		colleagueHydrator:   noopColleagueHydrator{},
	}, mock
}

func expectGroupGetRows(now time.Time, id int64, name string, updatedByUserID int64) *sqlmock.Rows {
	return sqlmock.NewRows([]string{
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
		id,
		"police",
		name,
		"Certified handlers and support staff.",
		"K9",
		nil,
		"#123456",
		int32(1),
		int32(1),
		int32(1),
		"0|zzzzzz:",
		int32(0),
		int32(0),
		int32(0),
		int32(0),
		int64(7),
		updatedByUserID,
		now,
		now,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	)
}

func expectGroupCreateCounts(mock sqlmock.Sqlmock, groupID int64) {
	mock.ExpectQuery(regexp.QuoteMeta("FROM fivenet_job_groups AS `group`")).
		WithArgs(groupID, int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"type", "job", "membership_mode"}).
			AddRow(int32(jobsgroups.GroupType_GROUP_TYPE_MANUAL), "police", int32(1)))

	mock.ExpectQuery(`(?s)FROM fivenet_job_group_manual_members AS mm .*INNER JOIN fivenet_job_groups AS g .*INNER JOIN fivenet_user_jobs AS uj`).
		WithArgs(groupID).
		WillReturnRows(sqlmock.NewRows([]string{
			"group_manual_member.group_id",
			"group_manual_member.user_id",
			"group_manual_member.reason",
			"group_manual_member.created_by_user_id",
			"group_manual_member.created_at",
		}))
	mock.ExpectQuery(regexp.QuoteMeta(`FROM fivenet_job_group_leaders`)).
		WithArgs(groupID).
		WillReturnRows(sqlmock.NewRows([]string{"data_count.total"}).AddRow(int64(0)))
	mock.ExpectExec(regexp.QuoteMeta("UPDATE fivenet_job_groups AS `group` SET")).
		WithArgs(int64(0), int64(0), int64(0), int64(0), groupID).
		WillReturnResult(sqlmock.NewResult(0, 1))
}

func TestCreateGroupWithoutAccessStillWorks(t *testing.T) {
	t.Parallel()

	accessStub := &stubGroupAccess{allowed: true}
	srv, mock := newJobsGroupACLTestServer(t, accessStub)
	srv.enricher = mstlystcdata.NewDummyUserAwareEnricher()
	now := time.Date(2026, time.August, 12, 13, 0, 0, 0, time.UTC)

	ctx := grpcauth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		UserId: 7,
		Job:    "police",
	})

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO fivenet_job_groups`)).
		WithArgs(
			"police",
			"K9 Unit",
			"Certified handlers and support staff.",
			"K9",
			"#123456",
			int32(1),
			int32(1),
			int32(1),
			"0|zzzzzz:",
			int64(7),
			int64(7),
		).
		WillReturnResult(sqlmock.NewResult(42, 1))
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO fivenet_job_group_activity`)).
		WithArgs("police", int64(42), int32(1), int32(7), nil).
		WillReturnResult(sqlmock.NewResult(0, 1))
	expectGroupCreateCounts(mock, 42)
	mock.ExpectQuery(regexp.QuoteMeta("FROM fivenet_job_groups AS `group` LEFT JOIN fivenet_files AS logo_file")).
		WithArgs("police", int64(42), int64(1)).
		WillReturnRows(expectGroupGetRows(now, 42, "K9 Unit", 7))
	mock.ExpectCommit()

	resp, err := srv.CreateGroup(ctx, &pbjobs.CreateGroupRequest{
		Name:        "K9 Unit",
		Description: new("Certified handlers and support staff."),
		ShortName:   new("K9"),
		Color:       new("#123456"),
		SortRank:    new("0|zzzzzz:"),
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, int64(42), resp.GetGroup().GetId())
	require.Equal(t, 1, accessStub.replaceCallCount)
	require.Equal(t, int64(42), accessStub.replacedTargetID)
	require.Len(t, accessStub.replacedAccess.GetJobs(), 1)
	require.Equal(t, "police", accessStub.replacedAccess.GetJobs()[0].GetJob())
	require.Equal(t, int32(3), accessStub.replacedAccess.GetJobs()[0].GetMinimumGrade())
	require.Equal(t, int32(4), accessStub.replacedAccess.GetJobs()[0].GetAccess())
	require.True(t, accessStub.replacedAccess.GetJobs()[0].GetRequired())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateGroupDeniedByAccess(t *testing.T) {
	t.Parallel()

	srv, mock := newJobsGroupACLTestServer(t, &stubGroupAccess{allowed: false})
	now := time.Date(2026, time.August, 12, 13, 0, 0, 0, time.UTC)

	ctx := grpcauth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		UserId: 7,
		Job:    "police",
	})

	mock.ExpectQuery(regexp.QuoteMeta("FROM fivenet_job_groups AS `group` LEFT JOIN fivenet_files AS logo_file")).
		WillReturnRows(expectGroupGetRows(now, 42, "K9 Unit", 7))

	resp, err := srv.UpdateGroup(ctx, &pbjobs.UpdateGroupRequest{
		Id:          42,
		Name:        new("K9 Unit"),
		Description: new("Certified handlers and support staff."),
		ShortName:   new("K9"),
		Color:       new("#123456"),
		SortRank:    new("0|zzzzzz:"),
	})
	require.ErrorIs(t, err, errorsjobs.ErrNotFoundOrNoPerms)
	require.Nil(t, resp)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateGroupUsesActiveOnlyLookup(t *testing.T) {
	t.Parallel()

	srv, mock := newJobsGroupACLTestServer(t, &stubGroupAccess{allowed: true})

	ctx := grpcauth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		UserId: 7,
		Job:    "police",
	})

	mock.ExpectQuery(regexp.QuoteMeta("FROM fivenet_job_groups AS `group` LEFT JOIN fivenet_files AS logo_file")).
		WithArgs(
			"police",
			int64(42),
			int32(jobsgroups.GroupState_GROUP_STATE_ARCHIVED),
			int64(1),
		).
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
		}))

	resp, err := srv.UpdateGroup(ctx, &pbjobs.UpdateGroupRequest{
		Id:    42,
		Name:  new("Archived K9"),
		State: new(jobsgroups.GroupState_GROUP_STATE_ACTIVE),
	})
	require.ErrorIs(t, err, errorsjobs.ErrNotFoundOrNoPerms)
	require.Nil(t, resp)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestArchiveGroupDeniedByAccess(t *testing.T) {
	t.Parallel()

	srv, mock := newJobsGroupACLTestServer(t, &stubGroupAccess{allowed: false})

	ctx := grpcauth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		UserId: 7,
		Job:    "police",
	})

	resp, err := srv.ArchiveGroup(ctx, &pbjobs.ArchiveGroupRequest{Id: 42})
	require.ErrorIs(t, err, errorsjobs.ErrNotFoundOrNoPerms)
	require.Nil(t, resp)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestArchiveGroupAlreadyArchivedIsNoop(t *testing.T) {
	t.Parallel()

	srv, mock := newJobsGroupACLTestServer(t, &stubGroupAccess{allowed: true})
	now := time.Date(2026, time.August, 13, 10, 0, 0, 0, time.UTC)

	ctx := grpcauth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		UserId: 7,
		Job:    "police",
	})

	mock.ExpectQuery(regexp.QuoteMeta("FROM fivenet_job_groups AS `group` LEFT JOIN fivenet_files AS logo_file")).
		WithArgs("police", int64(42), int64(1)).
		WillReturnRows(expectGroupGetRowsWithStateAndDeletedAt(
			now,
			42,
			"K9 Unit",
			7,
			jobsgroups.GroupState_GROUP_STATE_ARCHIVED,
			now,
		))

	resp, err := srv.ArchiveGroup(ctx, &pbjobs.ArchiveGroupRequest{
		Id:     42,
		Reason: new("duplicate archive"),
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.GetGroup())
	require.Equal(t, jobsgroups.GroupState_GROUP_STATE_ARCHIVED, resp.GetGroup().GetState())
	require.NotNil(t, resp.GetGroup().GetDeletedAt())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRestoreGroupDeniedByAccess(t *testing.T) {
	t.Parallel()

	srv, mock := newJobsGroupACLTestServer(t, &stubGroupAccess{allowed: false})

	ctx := grpcauth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		UserId: 7,
		Job:    "police",
	})

	resp, err := srv.RestoreGroup(ctx, &pbjobs.RestoreGroupRequest{Id: 42})
	require.ErrorIs(t, err, errorsjobs.ErrNotFoundOrNoPerms)
	require.Nil(t, resp)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRestoreGroupActiveGroupIsNoop(t *testing.T) {
	t.Parallel()

	srv, mock := newJobsGroupACLTestServer(t, &stubGroupAccess{allowed: true})
	now := time.Date(2026, time.August, 13, 10, 0, 0, 0, time.UTC)

	ctx := grpcauth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		UserId: 7,
		Job:    "police",
	})

	mock.ExpectQuery(regexp.QuoteMeta("FROM fivenet_job_groups AS `group` LEFT JOIN fivenet_files AS logo_file")).
		WithArgs("police", int64(42), int64(1)).
		WillReturnRows(expectGroupGetRowsWithStateAndDeletedAt(
			now,
			42,
			"K9 Unit",
			7,
			jobsgroups.GroupState_GROUP_STATE_ACTIVE,
			nil,
		))

	resp, err := srv.RestoreGroup(ctx, &pbjobs.RestoreGroupRequest{
		Id:     42,
		Reason: new("duplicate restore"),
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.GetGroup())
	require.Equal(t, jobsgroups.GroupState_GROUP_STATE_ACTIVE, resp.GetGroup().GetState())
	require.Nil(t, resp.GetGroup().GetDeletedAt())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestArchiveThenRestoreGroup(t *testing.T) {
	t.Parallel()

	accessStub := &stubGroupAccess{allowed: true}
	srv, mock := newJobsGroupACLTestServer(t, accessStub)
	now := time.Date(2026, time.August, 13, 10, 0, 0, 0, time.UTC)

	ctx := grpcauth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		UserId: 7,
		Job:    "police",
	})

	mock.ExpectQuery(regexp.QuoteMeta("FROM fivenet_job_groups AS `group` LEFT JOIN fivenet_files AS logo_file")).
		WithArgs("police", int64(42), int64(1)).
		WillReturnRows(expectGroupGetRowsWithStateAndDeletedAt(
			now,
			42,
			"K9 Unit",
			7,
			jobsgroups.GroupState_GROUP_STATE_ACTIVE,
			nil,
		))
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE fivenet_job_groups SET`)).
		WithArgs(
			int32(jobsgroups.GroupState_GROUP_STATE_ARCHIVED),
			int64(7),
			int64(42),
			"police",
			int64(1),
		).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO fivenet_job_group_activity`)).
		WithArgs(
			"police",
			int64(42),
			int32(jobsgroups.GroupActivityType_GROUP_ACTIVITY_TYPE_ARCHIVED),
			int32(7),
			"review",
			nil,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(regexp.QuoteMeta("FROM fivenet_job_groups AS `group` LEFT JOIN fivenet_files AS logo_file")).
		WithArgs("police", int64(42), int64(1)).
		WillReturnRows(expectGroupGetRowsWithStateAndDeletedAt(
			now,
			42,
			"K9 Unit",
			7,
			jobsgroups.GroupState_GROUP_STATE_ARCHIVED,
			now,
		))
	mock.ExpectCommit()

	archiveResp, err := srv.ArchiveGroup(ctx, &pbjobs.ArchiveGroupRequest{
		Id:     42,
		Reason: new("review"),
	})
	require.NoError(t, err)
	require.NotNil(t, archiveResp)
	require.NotNil(t, archiveResp.GetGroup())
	require.Equal(t, jobsgroups.GroupState_GROUP_STATE_ARCHIVED, archiveResp.GetGroup().GetState())
	require.NotNil(t, archiveResp.GetGroup().GetDeletedAt())

	mock.ExpectQuery(regexp.QuoteMeta("FROM fivenet_job_groups AS `group` LEFT JOIN fivenet_files AS logo_file")).
		WithArgs("police", int64(42), int64(1)).
		WillReturnRows(expectGroupGetRowsWithStateAndDeletedAt(
			now,
			42,
			"K9 Unit",
			7,
			jobsgroups.GroupState_GROUP_STATE_ARCHIVED,
			now,
		))
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE fivenet_job_groups SET`)).
		WithArgs(
			int32(jobsgroups.GroupState_GROUP_STATE_ACTIVE),
			int64(7),
			int64(42),
			"police",
			int64(1),
		).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO fivenet_job_group_activity`)).
		WithArgs(
			"police",
			int64(42),
			int32(jobsgroups.GroupActivityType_GROUP_ACTIVITY_TYPE_RESTORED),
			int32(7),
			nil,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectQuery(regexp.QuoteMeta("FROM fivenet_job_groups AS `group` LEFT JOIN fivenet_files AS logo_file")).
		WithArgs("police", int64(42), int64(1)).
		WillReturnRows(expectGroupGetRowsWithStateAndDeletedAt(
			now,
			42,
			"K9 Unit",
			7,
			jobsgroups.GroupState_GROUP_STATE_ACTIVE,
			nil,
		))
	mock.ExpectCommit()

	restoreResp, err := srv.RestoreGroup(ctx, &pbjobs.RestoreGroupRequest{Id: 42})
	require.NoError(t, err)
	require.NotNil(t, restoreResp)
	require.NotNil(t, restoreResp.GetGroup())
	require.Equal(t, jobsgroups.GroupState_GROUP_STATE_ACTIVE, restoreResp.GetGroup().GetState())
	require.Nil(t, restoreResp.GetGroup().GetDeletedAt())
	require.NoError(t, mock.ExpectationsWereMet())
}

func expectGroupGetRowsWithStateAndDeletedAt(
	now time.Time,
	id int64,
	name string,
	updatedByUserID int64,
	state jobsgroups.GroupState,
	deletedAt any,
) *sqlmock.Rows {
	return sqlmock.NewRows([]string{
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
		id,
		"police",
		name,
		"Certified handlers and support staff.",
		"K9",
		nil,
		"#123456",
		int32(1),
		int32(state),
		int32(1),
		"0|zzzzzz:",
		int32(0),
		int32(0),
		int32(0),
		int32(0),
		int64(7),
		updatedByUserID,
		now,
		now,
		deletedAt,
		nil,
		nil,
		nil,
		nil,
		nil,
	)
}

func TestAddGroupLeaderDeniedByAccess(t *testing.T) {
	t.Parallel()

	srv, mock := newJobsGroupACLTestServer(t, &stubGroupAccess{allowed: false})

	ctx := grpcauth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		UserId: 7,
		Job:    "police",
	})

	resp, err := srv.AddGroupLeader(ctx, &pbjobs.AddGroupLeaderRequest{
		GroupId: 42,
		UserId:  99,
	})
	require.ErrorIs(t, err, errorsjobs.ErrNotFoundOrNoPerms)
	require.Nil(t, resp)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRemoveGroupLeaderDeniedByAccess(t *testing.T) {
	t.Parallel()

	srv, mock := newJobsGroupACLTestServer(t, &stubGroupAccess{allowed: false})

	ctx := grpcauth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		UserId: 7,
		Job:    "police",
	})

	resp, err := srv.RemoveGroupLeader(ctx, &pbjobs.RemoveGroupLeaderRequest{
		GroupId: 42,
		UserId:  99,
	})
	require.ErrorIs(t, err, errorsjobs.ErrNotFoundOrNoPerms)
	require.Nil(t, resp)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestEnsureHighestJobGradeAccessUsesSharedNormalization(t *testing.T) {
	t.Parallel()

	srv := &Server{enricher: mstlystcdata.NewDummyUserAwareEnricher()}
	in := &resourcesaccess.Access{
		Jobs: []*resourcesaccess.JobAccess{
			{
				Id:             11,
				TargetId:       77,
				Job:            "police",
				MinimumGrade:   1,
				Access:         4,
				Required:       nil,
				RequiredAccess: nil,
			},
		},
	}

	out, err := srv.ensureHighestJobGradeAccess("police", in)
	require.NoError(t, err)

	require.NotSame(t, in, out)
	require.Len(t, out.GetJobs(), 2)
	require.Equal(t, int32(1), out.GetJobs()[0].GetMinimumGrade())
	require.Equal(t, int32(3), out.GetJobs()[1].GetMinimumGrade())
	require.Equal(t, int32(4), out.GetJobs()[1].GetAccess())
	require.True(t, out.GetJobs()[1].GetRequired())
	require.Equal(t, int32(4), out.GetJobs()[1].GetRequiredAccess())
	require.Equal(t, int64(0), out.GetJobs()[1].GetId())
	require.Equal(t, int64(0), out.GetJobs()[1].GetTargetId())
	require.Len(t, in.GetJobs(), 1)
	require.Equal(t, int32(1), in.GetJobs()[0].GetMinimumGrade())
}

func TestEnsureHighestJobGradeAccessKeepsExistingHighestEntry(t *testing.T) {
	t.Parallel()

	srv := &Server{enricher: mstlystcdata.NewDummyUserAwareEnricher()}
	in := &resourcesaccess.Access{
		Jobs: []*resourcesaccess.JobAccess{
			{
				Id:             11,
				TargetId:       77,
				Job:            "police",
				MinimumGrade:   1,
				Access:         2,
				Required:       nil,
				RequiredAccess: nil,
			},
			{
				Id:             12,
				TargetId:       77,
				Job:            "police",
				MinimumGrade:   3,
				Access:         5,
				Required:       nil,
				RequiredAccess: nil,
			},
		},
	}

	out, err := srv.ensureHighestJobGradeAccess("police", in)
	require.NoError(t, err)

	require.NotSame(t, in, out)
	require.Len(t, out.GetJobs(), 2)
	require.Equal(t, int32(1), out.GetJobs()[0].GetMinimumGrade())
	require.Equal(t, int32(3), out.GetJobs()[1].GetMinimumGrade())
	require.Equal(t, int32(5), out.GetJobs()[1].GetAccess())
	require.True(t, out.GetJobs()[1].GetRequired())
	require.Equal(t, int32(5), out.GetJobs()[1].GetRequiredAccess())
}
