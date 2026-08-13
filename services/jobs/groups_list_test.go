package jobs

import (
	"context"
	"testing"

	jobsgroups "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/groups"
	pbuserinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pbjobs "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/jobs"
	grpcauth "github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	jobsstore "github.com/fivenet-app/fivenet/v2026/stores/jobs"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/stretchr/testify/require"
)

type groupsQueryStoreStub struct {
	jobsstore.Store

	countQuery jobsstore.GroupsQuery
	listQuery  jobsstore.GroupsQuery
	count      int64
	groups     []*jobsgroups.Group
}

func (s *groupsQueryStoreStub) CountGroups(
	_ context.Context,
	_ qrm.DB,
	q jobsstore.GroupsQuery,
	_ *pbuserinfo.UserInfo,
) (int64, error) {
	s.countQuery = q
	return s.count, nil
}

func (s *groupsQueryStoreStub) ListGroups(
	_ context.Context,
	_ qrm.DB,
	q jobsstore.GroupsQuery,
	_ *pbuserinfo.UserInfo,
) ([]*jobsgroups.Group, error) {
	s.listQuery = q
	return s.groups, nil
}

func TestListGroupsForwardsGroupIDsToBothQueries(t *testing.T) {
	t.Parallel()

	store := &groupsQueryStoreStub{
		count:  1,
		groups: []*jobsgroups.Group{{Id: 42}},
	}
	server := &Server{
		store:             store,
		colleagueHydrator: noopColleagueHydrator{},
	}

	ctx := grpcauth.ContextWithUserInfo(t.Context(), &pbuserinfo.UserInfo{
		Job: "police",
	})

	resp, err := server.ListGroups(ctx, &pbjobs.ListGroupsRequest{
		GroupIds: []int32{42, 99},
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, []int64{42, 99}, store.countQuery.IDs)
	require.Equal(t, []int64{42, 99}, store.listQuery.IDs)
}
