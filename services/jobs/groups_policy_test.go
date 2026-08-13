package jobs

import (
	"context"
	"testing"

	jobsgroups "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/groups"
	errorsjobs "github.com/fivenet-app/fivenet/v2026/services/jobs/errors"
	jobsstore "github.com/fivenet-app/fivenet/v2026/stores/jobs"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/stretchr/testify/require"
)

type groupPolicyStoreStub struct {
	jobsstore.Store

	manualMembers []*jobsgroups.GroupManualMember
	rules         []*jobsgroups.GroupRule
	exclusions    []*jobsgroups.GroupMemberExclusion
	ruleMatches   []*jobsstore.GroupRuleMemberMatch
}

func (s *groupPolicyStoreStub) ListGroupManualMembers(
	_ context.Context,
	_ qrm.DB,
	_ int64,
	_ string,
) ([]*jobsgroups.GroupManualMember, error) {
	return s.manualMembers, nil
}

func (s *groupPolicyStoreStub) ListGroupRules(
	_ context.Context,
	_ qrm.DB,
	_ int64,
) ([]*jobsgroups.GroupRule, error) {
	return s.rules, nil
}

func (s *groupPolicyStoreStub) ListGroupMemberExclusions(
	_ context.Context,
	_ qrm.DB,
	_ int64,
	_ string,
) ([]*jobsgroups.GroupMemberExclusion, error) {
	return s.exclusions, nil
}

func (s *groupPolicyStoreStub) ListGroupRuleMemberMatches(
	_ context.Context,
	_ qrm.DB,
	_ *jobsgroups.Group,
	_ string,
) ([]*jobsstore.GroupRuleMemberMatch, error) {
	return s.ruleMatches, nil
}

func TestValidateGroupPolicyAgainstExistingDataRejectsUnmatchedManualMembers(t *testing.T) {
	t.Parallel()

	server := &Server{
		store: &groupPolicyStoreStub{
			manualMembers: []*jobsgroups.GroupManualMember{
				{UserId: 7},
			},
			rules:       []*jobsgroups.GroupRule{},
			exclusions:  []*jobsgroups.GroupMemberExclusion{},
			ruleMatches: []*jobsstore.GroupRuleMemberMatch{},
		},
	}

	err := server.validateGroupPolicyAgainstExistingData(t.Context(), nil, &jobsgroups.Group{
		Id:             42,
		Type:           jobsgroups.GroupType_GROUP_TYPE_MIXED,
		MembershipMode: jobsgroups.GroupMembershipMode_GROUP_MEMBERSHIP_MODE_STRICT,
	})
	require.ErrorIs(t, err, errorsjobs.ErrGroupMemberRulesRequired)
}
