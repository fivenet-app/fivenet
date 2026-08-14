import { describe, expect, it } from 'vitest';
import { GroupMembershipMode, GroupType } from '~~/gen/ts/resources/jobs/groups/group';
import {
    groupTypeAllowsExclusions,
    groupTypeAllowsManualMembers,
    groupTypeAllowsRules,
    groupTypeAllowsStrictMembershipMode,
    isLegacyGroupPolicyState,
    isValidGroupTypeMembershipMode,
    normalizeGroupMembershipMode,
} from './policy';

describe('job group policy', () => {
    it('allows manual members only for manual and mixed groups', () => {
        expect(groupTypeAllowsManualMembers(GroupType.MANUAL)).toBe(true);
        expect(groupTypeAllowsManualMembers(GroupType.SMART)).toBe(false);
        expect(groupTypeAllowsManualMembers(GroupType.MIXED)).toBe(true);
    });

    it('allows rules only for smart and mixed groups', () => {
        expect(groupTypeAllowsRules(GroupType.MANUAL)).toBe(false);
        expect(groupTypeAllowsRules(GroupType.SMART)).toBe(true);
        expect(groupTypeAllowsRules(GroupType.MIXED)).toBe(true);
    });

    it('allows exclusions only for mixed groups', () => {
        expect(groupTypeAllowsExclusions(GroupType.MANUAL)).toBe(false);
        expect(groupTypeAllowsExclusions(GroupType.SMART)).toBe(false);
        expect(groupTypeAllowsExclusions(GroupType.MIXED)).toBe(true);
    });

    it('allows strict membership only for mixed groups', () => {
        expect(groupTypeAllowsStrictMembershipMode(GroupType.MANUAL)).toBe(false);
        expect(groupTypeAllowsStrictMembershipMode(GroupType.SMART)).toBe(false);
        expect(groupTypeAllowsStrictMembershipMode(GroupType.MIXED)).toBe(true);
    });

    it('validates group type and membership mode combinations', () => {
        expect(isValidGroupTypeMembershipMode(GroupType.UNSPECIFIED, GroupMembershipMode.FLEXIBLE)).toBe(false);
        expect(isValidGroupTypeMembershipMode(GroupType.MANUAL, GroupMembershipMode.UNSPECIFIED)).toBe(false);
        expect(isValidGroupTypeMembershipMode(GroupType.MANUAL, GroupMembershipMode.FLEXIBLE)).toBe(true);
        expect(isValidGroupTypeMembershipMode(GroupType.SMART, GroupMembershipMode.FLEXIBLE)).toBe(true);
        expect(isValidGroupTypeMembershipMode(GroupType.MIXED, GroupMembershipMode.STRICT)).toBe(true);
        expect(isValidGroupTypeMembershipMode(GroupType.MANUAL, GroupMembershipMode.STRICT)).toBe(false);
        expect(isValidGroupTypeMembershipMode(GroupType.SMART, GroupMembershipMode.STRICT)).toBe(false);
    });

    it('normalizes membership mode to the canonical mode for each group type', () => {
        expect(normalizeGroupMembershipMode(GroupType.MANUAL, GroupMembershipMode.STRICT)).toBe(GroupMembershipMode.FLEXIBLE);
        expect(normalizeGroupMembershipMode(GroupType.SMART, GroupMembershipMode.STRICT)).toBe(GroupMembershipMode.FLEXIBLE);
        expect(normalizeGroupMembershipMode(GroupType.MIXED, GroupMembershipMode.FLEXIBLE)).toBe(GroupMembershipMode.STRICT);
        expect(normalizeGroupMembershipMode(GroupType.MIXED, GroupMembershipMode.STRICT)).toBe(GroupMembershipMode.STRICT);
    });

    it('treats invalid legacy group policy states as legacy', () => {
        expect(
            isLegacyGroupPolicyState({
                type: GroupType.MANUAL,
                membershipMode: GroupMembershipMode.STRICT,
            } as never),
        ).toBe(true);

        expect(
            isLegacyGroupPolicyState({
                type: GroupType.MIXED,
                membershipMode: GroupMembershipMode.FLEXIBLE,
            } as never),
        ).toBe(false);
    });
});
