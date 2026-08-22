package jobs

import (
	"context"

	jobsgroups "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/groups"
	errorsjobs "github.com/fivenet-app/fivenet/v2026/services/jobs/errors"
	jobsstore "github.com/fivenet-app/fivenet/v2026/stores/jobs"
	groupspolicy "github.com/fivenet-app/fivenet/v2026/stores/jobs/groupspolicy"
	"github.com/go-jet/jet/v2/qrm"
)

func normalizeGroupPolicyDefaults(group *jobsgroups.Group) {
	if group == nil {
		return
	}
	if group.GetType() == jobsgroups.GroupType_GROUP_TYPE_UNSPECIFIED {
		group.Type = jobsgroups.GroupType_GROUP_TYPE_MANUAL
	}
	if group.GetMembershipMode() == jobsgroups.GroupMembershipMode_GROUP_MEMBERSHIP_MODE_UNSPECIFIED {
		group.MembershipMode = jobsgroups.GroupMembershipMode_GROUP_MEMBERSHIP_MODE_FLEXIBLE
	}
}

func validateGroupPolicyCombination(group *jobsgroups.Group) error {
	if group == nil {
		return errorsjobs.ErrGroupPolicyViolation
	}
	if !groupspolicy.ValidTypeAndMembershipMode(group.GetType(), group.GetMembershipMode()) {
		return errorsjobs.ErrGroupPolicyViolation
	}
	return nil
}

func validateGroupPolicyAllowedMutation(
	group *jobsgroups.Group,
	mutation groupspolicy.Mutation,
) error {
	if group == nil {
		return errorsjobs.ErrGroupPolicyViolation
	}
	if !groupspolicy.AllowsMutation(group.GetType(), mutation) {
		return errorsjobs.ErrGroupPolicyViolation
	}
	return nil
}

func (s *Server) validateGroupPolicyAgainstExistingData(
	ctx context.Context,
	db qrm.DB,
	group *jobsgroups.Group,
) error {
	if err := validateGroupPolicyCombination(group); err != nil {
		return err
	}

	manualMembers, err := s.store.ListGroupManualMembers(
		ctx,
		db,
		jobsstore.GroupItemsQuery{GroupID: group.GetId()},
	)
	if err != nil {
		return err
	}
	rules, err := s.store.ListGroupRules(ctx, db, jobsstore.GroupItemsQuery{GroupID: group.GetId()})
	if err != nil {
		return err
	}
	exclusions, err := s.store.ListGroupMemberExclusions(
		ctx,
		db,
		jobsstore.GroupItemsQuery{GroupID: group.GetId()},
	)
	if err != nil {
		return err
	}

	if !groupspolicy.AllowsManualMembers(group.GetType()) && len(manualMembers) > 0 {
		return errorsjobs.ErrGroupPolicyViolation
	}
	if !groupspolicy.AllowsRules(group.GetType()) && len(rules) > 0 {
		return errorsjobs.ErrGroupPolicyViolation
	}
	if !groupspolicy.AllowsExclusions(group.GetType()) && len(exclusions) > 0 {
		return errorsjobs.ErrGroupPolicyViolation
	}
	if groupspolicy.RequiresManualMembersMatchRules(group.GetType(), group.GetMembershipMode()) {
		ruleMatches, err := s.store.ListGroupRuleMemberMatches(ctx, db, group, "")
		if err != nil {
			return err
		}
		ruleMemberIDs := map[int32]struct{}{}
		for _, match := range ruleMatches {
			ruleMemberIDs[match.UserID] = struct{}{}
		}
		for _, member := range manualMembers {
			if _, ok := ruleMemberIDs[member.GetUserId()]; !ok {
				return errorsjobs.ErrGroupMemberRulesRequired
			}
		}
	}

	return nil
}
