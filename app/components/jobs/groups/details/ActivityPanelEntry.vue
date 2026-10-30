<script lang="ts" setup>
import ColleagueInfoPopover from '~/components/jobs/colleagues/ColleagueInfoPopover.vue';
import QualificationBadge from '~/components/partials/qualifications/QualificationBadge.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { type GroupActivity, GroupActivityType } from '~~/gen/ts/resources/jobs/groups/activity';
import type { QualificationShort } from '~~/gen/ts/resources/qualifications/qualifications';
import { groupRuleLabel } from '../helpers';

const props = defineProps<{
    activity: GroupActivity;
}>();

const { t } = useI18n();

const ruleLabel = computed(() => {
    if (props.activity.data?.data.oneofKind === 'rule') {
        const rule = props.activity.data.data.rule;
        const ruleId = rule.id || props.activity.ruleId;
        return ruleId ? `#${ruleId} - ${groupRuleLabel(rule, t)}` : groupRuleLabel(rule, t);
    }

    return props.activity.ruleId ? `#${props.activity.ruleId}` : undefined;
});

function ruleQualification(qualificationId: number): QualificationShort | undefined {
    if (props.activity.data?.data.oneofKind !== 'rule' || props.activity.data.data.rule.rule.oneofKind !== 'qualification') {
        return undefined;
    }

    return props.activity.data.data.rule.rule.qualification.qualifications.find(
        (qualification) => qualification.id === qualificationId,
    );
}
</script>

<template>
    <div class="flex min-w-0 flex-auto justify-between gap-2 text-xs leading-5 text-muted">
        <div class="flex min-w-0 flex-col gap-0.5">
            <span class="text-sm font-medium text-highlighted">
                {{ $t(`enums.jobs.groups.GroupActivityType.${GroupActivityType[activity.type] ?? 'UNSPECIFIED'}`) }}
            </span>

            <div class="flex min-w-0 flex-col gap-0.5 text-sm">
                <template v-if="activity.reason">
                    <span class="inline-flex gap-1">
                        <span class="font-semibold">{{ $t('common.reason') }}:</span>
                        <span>{{ activity.reason }}</span>
                    </span>
                </template>

                <template v-if="activity.targetUserId">
                    <span class="inline-flex items-center gap-1">
                        <span class="font-semibold">{{ $t('common.target') }}:</span>
                        <ColleagueInfoPopover :user="activity.targetUser" :user-id="activity.targetUserId" />
                    </span>
                </template>

                <template v-if="ruleLabel">
                    <span class="inline-flex gap-1">
                        <span class="font-semibold">{{ $t('components.jobs.groups.details.rule') }}:</span>
                        <span>{{ ruleLabel }}</span>
                    </span>
                </template>

                <div
                    v-if="
                        activity.data?.data.oneofKind === 'rule' &&
                        activity.data.data.rule.rule.oneofKind === 'qualification' &&
                        activity.data.data.rule.rule.qualification.qualificationIds.length > 0
                    "
                    class="flex flex-wrap gap-1"
                >
                    <QualificationBadge
                        v-for="qualificationId in activity.data.data.rule.rule.qualification.qualificationIds"
                        :key="qualificationId"
                        :qualification-id="qualificationId"
                        :qualification="ruleQualification(qualificationId)"
                    />
                </div>
            </div>
        </div>

        <div class="flex flex-none flex-col items-end gap-0.5 text-xs">
            <span class="text-xs/5 text-dimmed"><GenericTime :value="activity.createdAt" type="compact" /></span>

            <span class="inline-flex items-center gap-1">
                <span class="text-dimmed">{{ $t('common.created_by') }}:</span>
                <ColleagueInfoPopover :user="activity.actorUser" :user-id="activity.actorUserId" />
            </span>
        </div>
    </div>
</template>
