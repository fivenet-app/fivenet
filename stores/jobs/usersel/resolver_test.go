package usersel

import (
	"context"
	"database/sql"
	"slices"
	"sort"
	"testing"

	resourcesaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/access"
	jobs "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs"
	jobsgroups "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/groups"
	pbuserinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	jobsstore "github.com/fivenet-app/fivenet/v2026/stores/jobs"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/stretchr/testify/require"
)

type resolverTestStore struct {
	groups         map[int64]*jobsgroups.Group
	manualMembers  map[int64][]*jobsgroups.GroupManualMember
	ruleMatches    map[int64][]*jobsstore.GroupRuleMemberMatch
	exclusions     map[int64][]*jobsgroups.GroupMemberExclusion
	leaders        map[int64][]*jobsgroups.GroupLeader
	getGroupCalls  int
	manualCalls    int
	ruleMatchCalls int
	exclusionCalls int
	leaderCalls    int
}

type resolverTestAccessChecker struct {
	allowed []int64
	err     error
	calls   int
}

func (c *resolverTestAccessChecker) CanUserAccessTarget(
	_ context.Context,
	targetID int64,
	_ *pbuserinfo.UserInfo,
	_ int32,
) (bool, error) {
	if slices.Contains(c.allowed, targetID) {
		return true, c.err
	}
	return false, c.err
}

func (c *resolverTestAccessChecker) CanUserAccessTargetIncludingDeleted(
	_ context.Context,
	targetID int64,
	_ *pbuserinfo.UserInfo,
	_ int32,
) (bool, error) {
	if slices.Contains(c.allowed, targetID) {
		return true, c.err
	}
	return false, c.err
}

func (c *resolverTestAccessChecker) CanUserAccessTargetIDs(
	_ context.Context,
	_ *pbuserinfo.UserInfo,
	_ int32,
	targetIDs ...int64,
) ([]int64, error) {
	c.calls++
	if c.err != nil {
		return nil, c.err
	}
	return slices.Clone(c.allowed), nil
}

func (c *resolverTestAccessChecker) ListTargetAccess(
	_ context.Context,
	_ qrm.DB,
	_ int64,
	_ access.SubjectAccessOptions,
) (*resourcesaccess.Access, error) {
	return &resourcesaccess.Access{}, c.err
}

func (c *resolverTestAccessChecker) ReplaceTargetAccess(
	_ context.Context,
	_ qrm.DB,
	_ *access.SubjectResolver,
	_ int64,
	_ *resourcesaccess.Access,
	_ access.SubjectAccessOptions,
) (*access.SubjectAccessChanges, error) {
	return nil, c.err
}

func (c *resolverTestAccessChecker) VisibleIDsByConditionQuery(
	_ *pbuserinfo.UserInfo,
	_ int32,
	_ bool,
	_ mysql.BoolExpression,
) access.VisibilityQuery {
	return access.VisibilityQuery{}
}

func (s *resolverTestStore) GetGroup(
	_ context.Context,
	_ qrm.DB,
	_ jobsstore.GroupQuery,
	id int64,
) (*jobsgroups.Group, error) {
	s.getGroupCalls++
	if group, ok := s.groups[id]; ok {
		return group, nil
	}
	return nil, sql.ErrNoRows
}

func (s *resolverTestStore) ListGroupManualMembers(
	_ context.Context,
	_ qrm.DB,
	groupID int64,
	_ string,
) ([]*jobsgroups.GroupManualMember, error) {
	s.manualCalls++
	return slices.Clone(s.manualMembers[groupID]), nil
}

func (s *resolverTestStore) ListGroupRuleMemberMatches(
	_ context.Context,
	_ qrm.DB,
	group *jobsgroups.Group,
	_ string,
) ([]*jobsstore.GroupRuleMemberMatch, error) {
	s.ruleMatchCalls++
	return slices.Clone(s.ruleMatches[group.GetId()]), nil
}

func (s *resolverTestStore) ListGroupMemberExclusions(
	_ context.Context,
	_ qrm.DB,
	groupID int64,
	_ string,
) ([]*jobsgroups.GroupMemberExclusion, error) {
	s.exclusionCalls++
	return slices.Clone(s.exclusions[groupID]), nil
}

func (s *resolverTestStore) ListGroupLeaders(
	_ context.Context,
	_ qrm.DB,
	groupID int64,
	_ string,
) ([]*jobsgroups.GroupLeader, error) {
	s.leaderCalls++
	return slices.Clone(s.leaders[groupID]), nil
}

func TestResolveUserIDsReturnsExplicitUsersOnly(t *testing.T) {
	t.Parallel()

	store := &resolverTestStore{}
	resolver := &Resolver{store: store}

	resolved, err := resolver.Resolve(
		t.Context(),
		nil,
		&pbuserinfo.UserInfo{Job: "police"},
		&jobs.UserSelector{UserIds: []int32{3, 1, 3, 2}},
		ResolveOpts{},
	)
	require.NoError(t, err)

	slices.Sort(resolved)
	require.Equal(t, []int32{1, 2, 3}, resolved)
	require.Zero(t, store.getGroupCalls)
	require.Zero(t, store.manualCalls)
	require.Zero(t, store.ruleMatchCalls)
	require.Zero(t, store.exclusionCalls)
	require.Zero(t, store.leaderCalls)
}

func TestResolveUserIDsIgnoresInaccessibleGroups(t *testing.T) {
	t.Parallel()

	checker := &resolverTestAccessChecker{allowed: []int64{10}}
	store := &resolverTestStore{
		groups: map[int64]*jobsgroups.Group{
			10: {
				Id:             10,
				Job:            "police",
				MembershipMode: jobsgroups.GroupMembershipMode_GROUP_MEMBERSHIP_MODE_FLEXIBLE,
				Name:           "Traffic",
				ShortName:      new("TRF"),
			},
		},
		manualMembers: map[int64][]*jobsgroups.GroupManualMember{
			10: {
				{GroupId: 10, UserId: 1},
				{GroupId: 10, UserId: 2},
			},
		},
		ruleMatches: map[int64][]*jobsstore.GroupRuleMemberMatch{
			10: {
				{GroupID: 10, UserID: 1, RuleID: 100},
				{GroupID: 10, UserID: 2, RuleID: 101},
			},
		},
		exclusions: map[int64][]*jobsgroups.GroupMemberExclusion{
			10: {
				{GroupId: 10, UserId: 2},
			},
		},
		leaders: map[int64][]*jobsgroups.GroupLeader{
			10: {
				{GroupId: 10, UserId: 9},
			},
		},
	}
	resolver := &Resolver{store: store, groupAccess: checker}

	resolved, err := resolver.Resolve(
		t.Context(),
		nil,
		&pbuserinfo.UserInfo{Job: "police"},
		&jobs.UserSelector{
			Groups: &jobs.GroupUserSelector{
				GroupIds:        []int64{10, 20},
				IncludeLeaders:  true,
				IncludeExcluded: false,
			},
		},
		ResolveOpts{},
	)
	require.NoError(t, err)

	sort.Slice(resolved, func(i, j int) bool { return resolved[i] < resolved[j] })
	require.Equal(t, []int32{1, 9}, resolved)
	require.Equal(t, 1, checker.calls)
	require.Equal(t, 1, store.getGroupCalls)
	require.Equal(t, 1, store.manualCalls)
	require.Equal(t, 1, store.ruleMatchCalls)
	require.Equal(t, 1, store.exclusionCalls)
	require.Equal(t, 1, store.leaderCalls)
}

func TestResolveUserIDsMergesExplicitUsersAndGroups(t *testing.T) {
	t.Parallel()

	checker := &resolverTestAccessChecker{allowed: []int64{10}}
	store := &resolverTestStore{
		groups: map[int64]*jobsgroups.Group{
			10: {
				Id:             10,
				Job:            "police",
				MembershipMode: jobsgroups.GroupMembershipMode_GROUP_MEMBERSHIP_MODE_FLEXIBLE,
				Name:           "Traffic",
			},
		},
		manualMembers: map[int64][]*jobsgroups.GroupManualMember{
			10: {
				{GroupId: 10, UserId: 2},
			},
		},
		ruleMatches: map[int64][]*jobsstore.GroupRuleMemberMatch{
			10: {
				{GroupID: 10, UserID: 2, RuleID: 101},
			},
		},
	}
	resolver := &Resolver{store: store, groupAccess: checker}

	resolved, err := resolver.Resolve(
		t.Context(),
		nil,
		&pbuserinfo.UserInfo{Job: "police"},
		&jobs.UserSelector{
			UserIds: []int32{2, 3},
			Groups: &jobs.GroupUserSelector{
				GroupIds: []int64{10},
			},
		},
		ResolveOpts{},
	)
	require.NoError(t, err)

	sort.Slice(resolved, func(i, j int) bool { return resolved[i] < resolved[j] })
	require.Equal(t, []int32{2, 3}, resolved)
}

func TestResolveUserIDsHonorsMaxResolvedUsers(t *testing.T) {
	t.Parallel()

	checker := &resolverTestAccessChecker{allowed: []int64{10}}
	store := &resolverTestStore{
		groups: map[int64]*jobsgroups.Group{
			10: {
				Id:             10,
				Job:            "police",
				MembershipMode: jobsgroups.GroupMembershipMode_GROUP_MEMBERSHIP_MODE_FLEXIBLE,
				Name:           "Traffic",
			},
		},
		manualMembers: map[int64][]*jobsgroups.GroupManualMember{
			10: {
				{GroupId: 10, UserId: 1},
				{GroupId: 10, UserId: 2},
			},
		},
		ruleMatches: map[int64][]*jobsstore.GroupRuleMemberMatch{
			10: {
				{GroupID: 10, UserID: 1, RuleID: 100},
				{GroupID: 10, UserID: 2, RuleID: 101},
			},
		},
	}
	resolver := &Resolver{store: store, groupAccess: checker}

	_, err := resolver.Resolve(
		t.Context(),
		nil,
		&pbuserinfo.UserInfo{Job: "police"},
		&jobs.UserSelector{
			Groups: &jobs.GroupUserSelector{
				GroupIds: []int64{10},
			},
		},
		ResolveOpts{MaxResolvedUsers: 1},
	)
	require.Error(t, err)
	require.ErrorContains(t, err, "too many resolved users")
}

func TestResolveUserIDsIgnoresMissingGroups(t *testing.T) {
	t.Parallel()

	store := &resolverTestStore{}
	resolver := &Resolver{
		store:       store,
		groupAccess: &resolverTestAccessChecker{allowed: []int64{999}},
	}

	resolved, err := resolver.Resolve(
		t.Context(),
		nil,
		&pbuserinfo.UserInfo{Job: "police"},
		&jobs.UserSelector{
			UserIds: []int32{7},
			Groups: &jobs.GroupUserSelector{
				GroupIds: []int64{999},
			},
		},
		ResolveOpts{},
	)
	require.NoError(t, err)
	require.Equal(t, []int32{7}, resolved)
	require.Equal(t, 1, store.getGroupCalls)
}

func TestGroupsOnlyStripsExplicitUsers(t *testing.T) {
	t.Parallel()

	selector := GroupsOnly(&jobs.UserSelector{
		UserIds: []int32{1, 2},
		Groups: &jobs.GroupUserSelector{
			GroupIds:        []int64{10, 11},
			IncludeLeaders:  true,
			IncludeExcluded: false,
		},
	})

	require.Empty(t, selector.GetUserIds())
	require.NotNil(t, selector.GetGroups())
	require.Equal(t, []int64{10, 11}, selector.GetGroups().GetGroupIds())
	require.True(t, selector.GetGroups().GetIncludeLeaders())
	require.False(t, selector.GetGroups().GetIncludeExcluded())
}
