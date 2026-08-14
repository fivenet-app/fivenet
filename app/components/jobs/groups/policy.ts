import { GroupMembershipMode, GroupType, type Group } from '~~/gen/ts/resources/jobs/groups/group';

export function groupTypeAllowsManualMembers(groupType: GroupType): boolean {
    return groupType === GroupType.MANUAL || groupType === GroupType.MIXED;
}

export function groupTypeAllowsRules(groupType: GroupType): boolean {
    return groupType === GroupType.SMART || groupType === GroupType.MIXED;
}

export function groupTypeAllowsExclusions(groupType: GroupType): boolean {
    return groupType === GroupType.MIXED;
}

export function groupTypeAllowsStrictMembershipMode(groupType: GroupType): boolean {
    return groupType === GroupType.MIXED;
}

export function isValidGroupTypeMembershipMode(groupType: GroupType, membershipMode: GroupMembershipMode): boolean {
    if (groupType === GroupType.UNSPECIFIED || membershipMode === GroupMembershipMode.UNSPECIFIED) {
        return false;
    }

    if (membershipMode === GroupMembershipMode.STRICT && !groupTypeAllowsStrictMembershipMode(groupType)) {
        return false;
    }

    return true;
}

export function normalizeGroupMembershipMode(groupType: GroupType, _membershipMode: GroupMembershipMode): GroupMembershipMode {
    if (groupType === GroupType.MIXED) {
        return GroupMembershipMode.STRICT;
    }

    return GroupMembershipMode.FLEXIBLE;
}

export function isLegacyGroupPolicyState(group: Group | undefined): boolean {
    if (!group) return false;

    return !isValidGroupTypeMembershipMode(group.type, group.membershipMode);
}
