package jobs

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	jobs "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs"
	jobsgroups "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/groups"
	permissionsattributes "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/attributes"
	pbuserinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pbjobs "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/jobs"
	"github.com/fivenet-app/fivenet/v2026/internal/tests/permsstub"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	grpcauth "github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	pkgperms "github.com/fivenet-app/fivenet/v2026/pkg/perms"
	jobsstore "github.com/fivenet-app/fivenet/v2026/stores/jobs"
	"github.com/fivenet-app/fivenet/v2026/stores/jobs/usersel"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/stretchr/testify/require"
)

type timeclockTestPerms struct {
	permsstub.Permissions

	access []string
}

func (p *timeclockTestPerms) AttrStringList(
	_ *pbuserinfo.UserInfo,
	_ pkgperms.AttrRef[pkgperms.StringListAttr],
) (*permissionsattributes.StringList, error) {
	return &permissionsattributes.StringList{Strings: append([]string(nil), p.access...)}, nil
}

type timeclockTestGroupStore struct {
	groups    map[int64]*jobsgroups.Group
	manual    map[int64][]*jobsgroups.GroupManualMember
	ruleMatch map[int64][]*jobsstore.GroupRuleMemberMatch
	excl      map[int64][]*jobsgroups.GroupMemberExclusion
	leaders   map[int64][]*jobsgroups.GroupLeader
}

func (s *timeclockTestGroupStore) GetGroup(
	_ context.Context,
	_ qrm.DB,
	_ jobsstore.GroupQuery,
	id int64,
) (*jobsgroups.Group, error) {
	if s == nil || s.groups == nil {
		return nil, sql.ErrNoRows
	}
	group, ok := s.groups[id]
	if !ok {
		return nil, sql.ErrNoRows
	}
	return group, nil
}

func (s *timeclockTestGroupStore) ListGroupManualMembers(
	_ context.Context,
	_ qrm.DB,
	q jobsstore.GroupItemsQuery,
) ([]*jobsgroups.GroupManualMember, error) {
	if s == nil {
		return nil, nil
	}
	return append([]*jobsgroups.GroupManualMember(nil), s.manual[q.GroupID]...), nil
}

func (s *timeclockTestGroupStore) ListGroupRuleMemberMatches(
	_ context.Context,
	_ qrm.DB,
	group *jobsgroups.Group,
	_ string,
) ([]*jobsstore.GroupRuleMemberMatch, error) {
	if s == nil {
		return nil, nil
	}
	return append([]*jobsstore.GroupRuleMemberMatch(nil), s.ruleMatch[group.GetId()]...), nil
}

func (s *timeclockTestGroupStore) ListGroupMemberExclusions(
	_ context.Context,
	_ qrm.DB,
	q jobsstore.GroupItemsQuery,
) ([]*jobsgroups.GroupMemberExclusion, error) {
	if s == nil {
		return nil, nil
	}
	return append([]*jobsgroups.GroupMemberExclusion(nil), s.excl[q.GroupID]...), nil
}

func (s *timeclockTestGroupStore) ListGroupLeaders(
	_ context.Context,
	_ qrm.DB,
	q jobsstore.GroupItemsQuery,
) ([]*jobsgroups.GroupLeader, error) {
	if s == nil {
		return nil, nil
	}
	return append([]*jobsgroups.GroupLeader(nil), s.leaders[q.GroupID]...), nil
}

func newTimeclockStatsTestServer(
	t *testing.T,
	accessPs []string,
	groupStore jobsstore.IGroupsQuery,
) (*Server, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = db.Close()
	})

	if groupStore == nil {
		groupStore = &timeclockTestGroupStore{}
	}

	return &Server{
		db:    db,
		perms: &timeclockTestPerms{access: accessPs},
		store: jobsstore.New(
			db,
			&config.CustomDB{},
			access.NewJobGroupsSubjectObjectAccess(db),
		).Store,
		userSel: usersel.NewWithAccess(
			groupStore,
			&stubGroupAccess{allowed: true},
		),
	}, mock
}

func expectTimeclockStatsQueries(mock sqlmock.Sqlmock, userIDs []int32) {
	tail := strings.Join(func() []string {
		parts := make([]string, len(userIDs))
		for i := range userIDs {
			parts[i] = `\?`
		}
		return parts
	}(), ", ")

	filter := fmt.Sprintf(`timeclock_entry\.user_id IN \(%s\)`, tail)

	args := make([]driver.Value, 0, len(userIDs)+1)
	args = append(args, "police")
	for i := range userIDs {
		args = append(args, userIDs[i])
	}

	mock.ExpectQuery(
		`(?s)SELECT .*FROM .*WHERE .*` + filter + `.*GROUP BY .*`,
	).
		WithArgs(args...).
		WillReturnRows(sqlmock.NewRows([]string{
			"timeclock_stats.job",
			"timeclock_stats.spent_time_sum",
			"timeclock_stats.spent_time_avg",
			"timeclock_stats.spent_time_max",
		}).AddRow("police", 12.5, 6.25, 8.0))

	mock.ExpectQuery(
		`(?s)SELECT .*FROM .*WHERE .*` + filter + `.*GROUP BY .*ORDER BY .*LIMIT \?.*`,
	).
		WithArgs(append(append([]driver.Value(nil), args...), int64(10))...).
		WillReturnRows(sqlmock.NewRows([]string{
			"timeclock_weekly_stats.year",
			"timeclock_weekly_stats.calendar_week",
			"timeclock_weekly_stats.sum",
			"timeclock_weekly_stats.avg",
			"timeclock_weekly_stats.max",
		}).AddRow(int32(2026), int32(32), 12.5, 6.25, 8.0))
}

func TestGetTimeclockStatsDefaultsToSelfWithoutAccessAll(t *testing.T) {
	t.Parallel()

	srv, mock := newTimeclockStatsTestServer(t, nil, nil)
	ctx := grpcauth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		UserId: 42,
		Job:    "police",
	})

	expectTimeclockStatsQueries(mock, []int32{42})

	resp, err := srv.GetTimeclockStats(ctx, &pbjobs.GetTimeclockStatsRequest{
		Users: &jobs.UserSelector{
			UserIds: []int32{99, 100},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "police", resp.GetStats().GetJob())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTimeclockStatsAllowsExplicitUsersWithAccessAll(t *testing.T) {
	t.Parallel()

	srv, mock := newTimeclockStatsTestServer(t, []string{"All"}, nil)
	ctx := grpcauth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		UserId: 42,
		Job:    "police",
	})

	expectTimeclockStatsQueries(mock, []int32{7, 8})

	resp, err := srv.GetTimeclockStats(ctx, &pbjobs.GetTimeclockStatsRequest{
		Users: &jobs.UserSelector{
			UserIds: []int32{7, 8},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.InDelta(t, float32(12.5), resp.GetStats().GetSpentTimeSum(), 0.000001)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTimeclockStatsAllowsGroupsWithAccessAll(t *testing.T) {
	t.Parallel()

	srv, mock := newTimeclockStatsTestServer(t, []string{"All"}, &timeclockTestGroupStore{
		groups: map[int64]*jobsgroups.Group{
			10: {
				Id:             10,
				Job:            "police",
				MembershipMode: jobsgroups.GroupMembershipMode_GROUP_MEMBERSHIP_MODE_FLEXIBLE,
				Name:           "Traffic",
			},
		},
		manual: map[int64][]*jobsgroups.GroupManualMember{
			10: {
				{GroupId: 10, UserId: 2},
				{GroupId: 10, UserId: 3},
			},
		},
		excl: map[int64][]*jobsgroups.GroupMemberExclusion{
			10: {
				{GroupId: 10, UserId: 3},
			},
		},
		leaders: map[int64][]*jobsgroups.GroupLeader{
			10: {
				{GroupId: 10, UserId: 9},
			},
		},
	})
	ctx := grpcauth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		UserId: 42,
		Job:    "police",
	})

	expectTimeclockStatsQueries(mock, []int32{2, 9})

	resp, err := srv.GetTimeclockStats(ctx, &pbjobs.GetTimeclockStatsRequest{
		Users: &jobs.UserSelector{
			Groups: &jobs.GroupUserSelector{
				GroupIds:        []int64{10},
				IncludeLeaders:  true,
				IncludeExcluded: false,
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, "police", resp.GetStats().GetJob())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTimeclockStatsAllowsMixedUsersAndGroupsWithAccessAll(t *testing.T) {
	t.Parallel()

	srv, mock := newTimeclockStatsTestServer(t, []string{"All"}, &timeclockTestGroupStore{
		groups: map[int64]*jobsgroups.Group{
			10: {
				Id:             10,
				Job:            "police",
				MembershipMode: jobsgroups.GroupMembershipMode_GROUP_MEMBERSHIP_MODE_FLEXIBLE,
				Name:           "Traffic",
			},
		},
		manual: map[int64][]*jobsgroups.GroupManualMember{
			10: {
				{GroupId: 10, UserId: 3},
			},
		},
	})
	ctx := grpcauth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		UserId: 42,
		Job:    "police",
	})

	expectTimeclockStatsQueries(mock, []int32{3, 4})

	resp, err := srv.GetTimeclockStats(ctx, &pbjobs.GetTimeclockStatsRequest{
		Users: &jobs.UserSelector{
			UserIds: []int32{4},
			Groups: &jobs.GroupUserSelector{
				GroupIds: []int64{10},
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.InDelta(t, float32(8.0), resp.GetStats().GetSpentTimeMax(), 0.000001)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTimeclockStatsReturnsEmptyForZeroResolvedUsers(t *testing.T) {
	t.Parallel()

	srv, mock := newTimeclockStatsTestServer(t, []string{"All"}, &timeclockTestGroupStore{})
	ctx := grpcauth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		UserId: 42,
		Job:    "police",
	})

	resp, err := srv.GetTimeclockStats(ctx, &pbjobs.GetTimeclockStatsRequest{
		Users: &jobs.UserSelector{
			Groups: &jobs.GroupUserSelector{
				GroupIds: []int64{999},
			},
		},
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Empty(t, resp.GetStats().GetJob())
	require.Zero(t, resp.GetStats().GetSpentTimeSum())
	require.Empty(t, resp.GetWeekly())
	require.NoError(t, mock.ExpectationsWereMet())
}
