package groupspolicy

import jobsgroups "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/groups"

type Mutation int

const (
	MutationManualMemberAdd Mutation = iota
	MutationManualMemberRemove
	MutationRuleAdd
	MutationRuleUpdate
	MutationRuleDelete
	MutationExclusionAdd
	MutationExclusionRemove
)

func AllowsManualMembers(groupType jobsgroups.GroupType) bool {
	switch groupType {
	case jobsgroups.GroupType_GROUP_TYPE_MANUAL, jobsgroups.GroupType_GROUP_TYPE_MIXED:
		return true
	default:
		return false
	}
}

func AllowsRules(groupType jobsgroups.GroupType) bool {
	switch groupType {
	case jobsgroups.GroupType_GROUP_TYPE_SMART, jobsgroups.GroupType_GROUP_TYPE_MIXED:
		return true
	default:
		return false
	}
}

func AllowsExclusions(groupType jobsgroups.GroupType) bool {
	return groupType == jobsgroups.GroupType_GROUP_TYPE_MIXED
}

func RequiresManualMembersMatchRules(
	groupType jobsgroups.GroupType,
	membershipMode jobsgroups.GroupMembershipMode,
) bool {
	return groupType == jobsgroups.GroupType_GROUP_TYPE_MIXED &&
		membershipMode == jobsgroups.GroupMembershipMode_GROUP_MEMBERSHIP_MODE_STRICT
}

func AllowsMutation(groupType jobsgroups.GroupType, mutation Mutation) bool {
	switch mutation {
	case MutationManualMemberAdd:
		return AllowsManualMembers(groupType)
	case MutationManualMemberRemove:
		return true
	case MutationRuleAdd, MutationRuleUpdate:
		return AllowsRules(groupType)
	case MutationRuleDelete:
		return true
	case MutationExclusionAdd:
		return AllowsExclusions(groupType)
	case MutationExclusionRemove:
		return true
	default:
		return false
	}
}

func ValidTypeAndMembershipMode(
	groupType jobsgroups.GroupType,
	membershipMode jobsgroups.GroupMembershipMode,
) bool {
	if groupType == jobsgroups.GroupType_GROUP_TYPE_UNSPECIFIED ||
		membershipMode == jobsgroups.GroupMembershipMode_GROUP_MEMBERSHIP_MODE_UNSPECIFIED {
		return false
	}

	if membershipMode == jobsgroups.GroupMembershipMode_GROUP_MEMBERSHIP_MODE_STRICT &&
		groupType != jobsgroups.GroupType_GROUP_TYPE_MIXED {
		return false
	}

	return true
}
