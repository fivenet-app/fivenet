import { GroupMembershipMode, GroupState, GroupType } from '~~/gen/ts/resources/jobs/groups/group';

export const groupStateItems: { label: string; value: GroupState }[] = [
    { label: 'Active', value: GroupState.ACTIVE },
    { label: 'Inactive', value: GroupState.INACTIVE },
    { label: 'Archived', value: GroupState.ARCHIVED },
];

export const groupStateFilterItems: { label: string; value: 'active' | 'inactive' | 'archived' | 'all' }[] = [
    { label: 'Active', value: 'active' },
    { label: 'Inactive', value: 'inactive' },
    { label: 'Archived', value: 'archived' },
    { label: 'All', value: 'all' },
];

export const groupTypeItems: { label: string; value: GroupType }[] = [
    { label: 'Manual', value: GroupType.MANUAL },
    { label: 'Smart', value: GroupType.SMART },
    { label: 'Mixed', value: GroupType.MIXED },
];

export const groupTypeFilterItems: { label: string; value: 'all' | 'manual' | 'smart' | 'mixed' }[] = [
    { label: 'All types', value: 'all' },
    { label: 'Manual', value: 'manual' },
    { label: 'Smart', value: 'smart' },
    { label: 'Mixed', value: 'mixed' },
];

export const groupMembershipModeItems: { label: string; value: GroupMembershipMode }[] = [
    { label: 'Flexible', value: GroupMembershipMode.FLEXIBLE },
    { label: 'Strict', value: GroupMembershipMode.STRICT },
];

export function groupStateLabel(state: GroupState): string {
    switch (state) {
        case GroupState.ACTIVE:
            return 'Active';
        case GroupState.INACTIVE:
            return 'Inactive';
        case GroupState.ARCHIVED:
            return 'Archived';
        default:
            return 'Unknown';
    }
}

export function groupStateColor(state: GroupState): 'success' | 'warning' | 'neutral' {
    switch (state) {
        case GroupState.ACTIVE:
            return 'success';
        case GroupState.INACTIVE:
            return 'warning';
        case GroupState.ARCHIVED:
        default:
            return 'neutral';
    }
}

export function groupTypeLabel(type: GroupType): string {
    switch (type) {
        case GroupType.MANUAL:
            return 'Manual';
        case GroupType.SMART:
            return 'Smart';
        case GroupType.MIXED:
            return 'Mixed';
        default:
            return 'Unknown';
    }
}

export function groupTypeColor(type: GroupType): 'neutral' | 'primary' | 'success' {
    switch (type) {
        case GroupType.MANUAL:
            return 'neutral';
        case GroupType.SMART:
            return 'primary';
        case GroupType.MIXED:
            return 'success';
        default:
            return 'neutral';
    }
}

export function groupMembershipModeLabel(mode: GroupMembershipMode): string {
    switch (mode) {
        case GroupMembershipMode.FLEXIBLE:
            return 'Flexible';
        case GroupMembershipMode.STRICT:
            return 'Strict';
        default:
            return 'Unknown';
    }
}

export function groupMembershipModeColor(mode: GroupMembershipMode): 'primary' | 'warning' | 'neutral' {
    switch (mode) {
        case GroupMembershipMode.FLEXIBLE:
            return 'primary';
        case GroupMembershipMode.STRICT:
            return 'warning';
        default:
            return 'neutral';
    }
}
