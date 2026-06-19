import { GroupActivityType } from '~~/gen/ts/resources/jobs/groups/activity';
import {
    GroupGradeRuleType,
    GroupMembershipMode,
    GroupQualificationRuleType,
    GroupState,
    GroupType,
    type GroupRule,
} from '~~/gen/ts/resources/jobs/groups/group';

type SelectItem<T> = {
    labelKey: string;
    value: T;
    icon?: string;
};

export const groupStateItems: { labelKey: string; value: GroupState }[] = [
    { labelKey: 'enums.jobs.groups.GroupState.ACTIVE', value: GroupState.ACTIVE },
    { labelKey: 'enums.jobs.groups.GroupState.INACTIVE', value: GroupState.INACTIVE },
    { labelKey: 'enums.jobs.groups.GroupState.ARCHIVED', value: GroupState.ARCHIVED },
];

export const groupStateFilterItems: { labelKey: string; value: 'active' | 'inactive' | 'archived' | 'all' }[] = [
    { labelKey: 'enums.jobs.groups.GroupState.ACTIVE', value: 'active' },
    { labelKey: 'enums.jobs.groups.GroupState.INACTIVE', value: 'inactive' },
    { labelKey: 'enums.jobs.groups.GroupState.ARCHIVED', value: 'archived' },
    { labelKey: 'common.all', value: 'all' },
];

export const groupTypeItems: SelectItem<GroupType>[] = [
    { labelKey: 'enums.jobs.groups.GroupType.MANUAL', value: GroupType.MANUAL, icon: 'i-mdi-account-edit' },
    { labelKey: 'enums.jobs.groups.GroupType.SMART', value: GroupType.SMART, icon: 'i-mdi-account-cog' },
    { labelKey: 'enums.jobs.groups.GroupType.MIXED', value: GroupType.MIXED, icon: 'i-mdi-account-switch' },
];

export const groupTypeFilterItems: SelectItem<'all' | 'manual' | 'smart' | 'mixed'>[] = [
    { labelKey: 'components.jobs.groups.types.all', value: 'all', icon: 'i-mdi-account-group' },
    { labelKey: 'enums.jobs.groups.GroupType.MANUAL', value: 'manual', icon: 'i-mdi-account-edit' },
    { labelKey: 'enums.jobs.groups.GroupType.SMART', value: 'smart', icon: 'i-mdi-account-cog' },
    { labelKey: 'enums.jobs.groups.GroupType.MIXED', value: 'mixed', icon: 'i-mdi-account-switch' },
];

export const groupMembershipModeItems: SelectItem<GroupMembershipMode>[] = [
    {
        labelKey: 'enums.jobs.groups.GroupMembershipMode.FLEXIBLE',
        value: GroupMembershipMode.FLEXIBLE,
        icon: 'i-mdi-account-plus',
    },
    {
        labelKey: 'enums.jobs.groups.GroupMembershipMode.STRICT',
        value: GroupMembershipMode.STRICT,
        icon: 'i-mdi-shield-account',
    },
];

export function groupStateLabelKey(state: GroupState): string {
    switch (state) {
        case GroupState.ACTIVE:
            return 'enums.jobs.groups.GroupState.ACTIVE';
        case GroupState.INACTIVE:
            return 'enums.jobs.groups.GroupState.INACTIVE';
        case GroupState.ARCHIVED:
            return 'enums.jobs.groups.GroupState.ARCHIVED';
        default:
            return 'common.unknown';
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

export function groupTypeLabelKey(type: GroupType): string {
    switch (type) {
        case GroupType.MANUAL:
            return 'enums.jobs.groups.GroupType.MANUAL';
        case GroupType.SMART:
            return 'enums.jobs.groups.GroupType.SMART';
        case GroupType.MIXED:
            return 'enums.jobs.groups.GroupType.MIXED';
        default:
            return 'common.unknown';
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

export function groupTypeIcon(type: GroupType): string {
    switch (type) {
        case GroupType.MANUAL:
            return 'i-mdi-account-edit';
        case GroupType.SMART:
            return 'i-mdi-account-cog';
        case GroupType.MIXED:
            return 'i-mdi-account-switch';
        default:
            return 'i-mdi-account-question';
    }
}

export function groupMembershipModeLabelKey(mode: GroupMembershipMode): string {
    switch (mode) {
        case GroupMembershipMode.FLEXIBLE:
            return 'enums.jobs.groups.GroupMembershipMode.FLEXIBLE';
        case GroupMembershipMode.STRICT:
            return 'enums.jobs.groups.GroupMembershipMode.STRICT';
        default:
            return 'common.unknown';
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

export function groupMembershipModeIcon(mode: GroupMembershipMode): string {
    switch (mode) {
        case GroupMembershipMode.FLEXIBLE:
            return 'i-mdi-account-plus';
        case GroupMembershipMode.STRICT:
            return 'i-mdi-shield-account';
        default:
            return 'i-mdi-account-question';
    }
}

export function groupActivityTypeColor(type: GroupActivityType | undefined): string {
    switch (type) {
        case GroupActivityType.CREATED:
            return 'text-green-500!';
        case GroupActivityType.UPDATED:
            return 'text-blue-500!';
        case GroupActivityType.ARCHIVED:
            return 'text-neutral-500!';
        case GroupActivityType.RESTORED:
            return 'text-teal-500!';
        case GroupActivityType.MEMBER_ADDED:
        case GroupActivityType.LEADER_ADDED:
            return 'text-green-500!';
        case GroupActivityType.MEMBER_REMOVED:
        case GroupActivityType.LEADER_REMOVED:
            return 'text-red-500!';
        case GroupActivityType.MEMBER_EXCLUDED:
            return 'text-orange-500!';
        case GroupActivityType.MEMBER_EXCLUSION_REMOVED:
            return 'text-emerald-500!';
        case GroupActivityType.RULE_ADDED:
            return 'text-cyan-500!';
        case GroupActivityType.RULE_UPDATED:
            return 'text-blue-500!';
        case GroupActivityType.RULE_REMOVED:
            return 'text-red-500!';
        case GroupActivityType.LOGO_UPDATED:
            return 'text-purple-500!';
        case GroupActivityType.UNSPECIFIED:
        default:
            return '!text-info-500';
    }
}

export function groupActivityTypeIcon(type: GroupActivityType | undefined): string {
    switch (type) {
        case GroupActivityType.CREATED:
            return 'i-mdi-plus-circle';
        case GroupActivityType.UPDATED:
            return 'i-mdi-pencil';
        case GroupActivityType.ARCHIVED:
            return 'i-mdi-archive';
        case GroupActivityType.RESTORED:
            return 'i-mdi-restore';
        case GroupActivityType.MEMBER_ADDED:
            return 'i-mdi-account-plus';
        case GroupActivityType.MEMBER_REMOVED:
            return 'i-mdi-account-minus';
        case GroupActivityType.MEMBER_EXCLUDED:
            return 'i-mdi-account-cancel';
        case GroupActivityType.MEMBER_EXCLUSION_REMOVED:
            return 'i-mdi-account-check';
        case GroupActivityType.LEADER_ADDED:
            return 'i-mdi-account-star';
        case GroupActivityType.LEADER_REMOVED:
            return 'i-mdi-account-star-outline';
        case GroupActivityType.RULE_ADDED:
            return 'i-mdi-filter-plus';
        case GroupActivityType.RULE_UPDATED:
            return 'i-mdi-filter-cog';
        case GroupActivityType.RULE_REMOVED:
            return 'i-mdi-filter-remove';
        case GroupActivityType.LOGO_UPDATED:
            return 'i-mdi-image-edit';
        case GroupActivityType.UNSPECIFIED:
        default:
            return 'i-mdi-help';
    }
}

type TranslateFn = (key: string, params?: Record<string, unknown>) => string;

function ruleGradeLabel(label: string | undefined, grade: number | undefined): string {
    return label || `${grade ?? 0}`;
}

export function groupRuleLabel(rule: GroupRule, t: TranslateFn): string {
    if (rule.rule.oneofKind === 'grade') {
        const gradeRule = rule.rule.grade;
        switch (gradeRule.type) {
            case GroupGradeRuleType.MINIMUM:
                return t('components.jobs.groups.details.rules.grade_minimum', {
                    grade: ruleGradeLabel(gradeRule.gradeLabel, gradeRule.grade),
                });
            case GroupGradeRuleType.EXACT:
                return t('components.jobs.groups.details.rules.grade_exact', {
                    grade: ruleGradeLabel(gradeRule.gradeLabel, gradeRule.grade),
                });
            case GroupGradeRuleType.RANGE:
                return t('components.jobs.groups.details.rules.grade_range', {
                    min: ruleGradeLabel(gradeRule.minGradeLabel, gradeRule.minGrade),
                    max: ruleGradeLabel(gradeRule.maxGradeLabel, gradeRule.maxGrade),
                });
            default:
                return t('enums.jobs.groups.GroupGradeRuleType.UNSPECIFIED');
        }
    }

    if (rule.rule.oneofKind === 'qualification') {
        const qualificationRule = rule.rule.qualification;
        const typeKey = `enums.jobs.groups.GroupQualificationRuleType.${
            GroupQualificationRuleType[qualificationRule.type] ?? 'UNSPECIFIED'
        }`;
        return t('components.jobs.groups.details.rules.qualifications', {
            mode: t(typeKey),
            count: qualificationRule.qualificationIds.length,
        });
    }

    return t('components.jobs.groups.details.rules.unknown');
}
