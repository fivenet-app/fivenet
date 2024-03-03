<script lang="ts" setup>
import CitizenInfoPopover from '~/components/partials/citizens/CitizenInfoPopover.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { jobsUserActivityTypeBGColor, jobsUserActivityTypeIcon } from '~/components/jobs/colleagues/info/helpers';
import { JobsUserActivityType, type JobsUserActivity } from '~~/gen/ts/resources/jobs/colleagues';

defineProps<{
    activity: JobsUserActivity;
}>();
</script>

<template>
    <div class="flex space-x-3">
        <div class="my-auto flex h-10 w-10 items-center justify-center rounded-full">
            <component
                :is="jobsUserActivityTypeIcon(activity.activityType)"
                :class="[jobsUserActivityTypeBGColor(activity.activityType), 'h-full w-full']"
                aria-hidden="true"
            />
        </div>
        <div class="flex-1 space-y-1">
            <div class="flex items-center justify-between">
                <h3 class="text-sm font-medium text-neutral">
                    {{ $t(`enums.jobs.JobsUserActivityType.${JobsUserActivityType[activity.activityType]}`) }}
                    <template v-if="activity.data?.data.oneofKind !== undefined">
                        {{ ' - ' }}
                        <template v-if="activity.data?.data.oneofKind === 'absenceDate'">
                            <span class="inline-flex gap-1">
                                <GenericTime :value="activity.data?.data.absenceDate.absenceBegin" type="date" />
                                <span>{{ $t('common.to') }}</span>
                                <GenericTime :value="activity.data?.data.absenceDate.absenceEnd" type="date" />
                            </span>
                        </template>
                        <template v-else-if="activity.data?.data.oneofKind === 'gradeChange'">
                            {{ activity.data?.data.gradeChange.gradeLabel }} ({{ activity.data?.data.gradeChange.grade }})
                        </template>
                    </template>
                </h3>
                <p class="text-sm text-gray-400">
                    <GenericTime :value="activity.createdAt" type="long" />
                </p>
            </div>
            <div class="flex items-center justify-between">
                <p class="inline-flex gap-1 text-sm text-gray-300">
                    <template v-if="activity.reason">
                        <span class="font-semibold">{{ $t('common.reason') }}:</span>
                        <span>
                            {{ activity.reason }}
                        </span>
                    </template>
                </p>
                <p class="inline-flex text-sm text-gray-300">
                    {{ $t('common.created_by') }}
                    <CitizenInfoPopover class="ml-1" text-class="underline" :user="activity.sourceUser" />
                </p>
            </div>
        </div>
    </div>
</template>
