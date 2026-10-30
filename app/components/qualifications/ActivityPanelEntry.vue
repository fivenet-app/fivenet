<script lang="ts" setup>
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { type QualificationActivity, QualificationActivityType } from '~~/gen/ts/resources/qualifications/activity/activity';
import { RequestStatus, ResultStatus } from '~~/gen/ts/resources/qualifications/qualifications';
import { requestStatusToBadgeColor, resultStatusToBadgeColor } from './helpers';

defineProps<{
    activity: QualificationActivity;
}>();
</script>

<template>
    <div class="flex min-w-0 flex-auto justify-between gap-2 text-xs leading-5 text-muted">
        <div class="flex min-w-0 flex-col gap-0.5">
            <span class="inline-flex flex-wrap items-center gap-x-2 gap-y-1 text-sm font-medium text-highlighted">
                <span>{{
                    $t(`enums.qualifications.QualificationActivityType.${QualificationActivityType[activity.type]}`)
                }}</span>
                <UBadge
                    v-if="activity.data?.requestStatus !== undefined"
                    :color="requestStatusToBadgeColor(activity.data.requestStatus)"
                    :label="$t(`enums.qualifications.RequestStatus.${RequestStatus[activity.data.requestStatus]}`)"
                />
                <UBadge
                    v-if="activity.data?.resultStatus !== undefined"
                    :color="resultStatusToBadgeColor(activity.data.resultStatus)"
                    :label="$t(`enums.qualifications.ResultStatus.${ResultStatus[activity.data.resultStatus]}`)"
                />
                <UBadge v-if="activity.data?.score !== undefined" color="neutral" :label="`${activity.data.score}`" />
            </span>

            <span v-if="activity.targetUser" class="inline-flex min-w-0 items-center gap-1 text-sm">
                <span class="font-semibold">{{ $t('common.citizen', 1) }}:</span>
                <CitizenInfoPopover :user="activity.targetUser" show-avatar />
            </span>
        </div>

        <div class="flex flex-none flex-col items-end gap-0.5 text-xs">
            <span class="text-xs/5 text-dimmed"><GenericTime :value="activity.createdAt" type="compact" /></span>

            <span class="inline-flex items-center gap-1">
                <span class="text-dimmed">{{ $t('common.created_by') }}:</span>
                <CitizenInfoPopover v-if="activity.actorUser" :user="activity.actorUser" show-avatar />
                <span v-else>{{ $t('common.system') }}</span>
            </span>
        </div>
    </div>
</template>
