<script lang="ts" setup>
import ColleagueInfoPopover from '~/components/jobs/colleagues/ColleagueInfoPopover.vue';
import { jobsUserActivityTypeBGColor, jobsUserActivityTypeIcon } from '~/components/jobs/colleagues/info/helpers';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { JobsUserActivityType, type JobsUserActivity } from '~~/gen/ts/resources/jobs/colleagues';

withDefaults(
    defineProps<{
        activity: JobsUserActivity;
        showTargetUser?: boolean;
    }>(),
    {
        showTargetUser: false,
    },
);
</script>

<template>
    <div class="flex space-x-3">
        <div class="my-auto flex size-10 items-center justify-center rounded-full">
            <UIcon
                :name="jobsUserActivityTypeIcon(activity.activityType)"
                :class="[jobsUserActivityTypeBGColor(activity.activityType), 'size-full']"
                :inline="true"
            />
        </div>

        <div class="flex-1 space-y-1">
            <div class="flex items-center justify-between">
                <h3 class="text-sm font-medium">
                    {{ $t(`enums.jobs.JobsUserActivityType.${JobsUserActivityType[activity.activityType]}`) }}
                    <template v-if="activity.data?.data.oneofKind !== undefined">
                        {{ ' - ' }}
                        <template v-if="activity.data?.data.oneofKind === 'absenceDate'">
                            <template
                                v-if="
                                    activity.data?.data.absenceDate.absenceBegin && activity.data?.data.absenceDate.absenceEnd
                                "
                            >
                                <span class="inline-flex gap-1">
                                    <GenericTime :value="activity.data?.data.absenceDate.absenceBegin" type="date" />
                                    <span>{{ $t('common.to') }}</span>
                                    <GenericTime :value="activity.data?.data.absenceDate.absenceEnd" type="date" />
                                </span>
                            </template>
                            <template v-else>
                                <span>{{ $t('common.annul', 2) }}</span>
                            </template>
                        </template>
                        <template v-else-if="activity.data?.data.oneofKind === 'gradeChange'">
                            {{ activity.data?.data.gradeChange.gradeLabel }} ({{ activity.data?.data.gradeChange.grade }})
                        </template>
                        <template v-else-if="activity.data?.data.oneofKind === 'labelsChange'">
                            <div class="inline-flex gap-1">
                                <UBadge
                                    v-for="label in activity.data.data.labelsChange?.removed"
                                    :key="label.name"
                                    :style="{ backgroundColor: label.color }"
                                    class="justify-between gap-2 line-through"
                                    :class="isColourBright(hexToRgb(label.color, RGBBlack)!) ? '!text-black' : '!text-white'"
                                    size="xs"
                                >
                                    {{ label.name }}
                                </UBadge>

                                <UBadge
                                    v-for="label in activity.data.data.labelsChange?.added"
                                    :key="label.name"
                                    :style="{ backgroundColor: label.color }"
                                    class="justify-between gap-2"
                                    :class="isColourBright(hexToRgb(label.color, RGBBlack)!) ? '!text-black' : '!text-white'"
                                    size="xs"
                                >
                                    {{ label.name }}
                                </UBadge>
                            </div>
                        </template>
                        <template v-else-if="activity.data?.data.oneofKind === 'nameChange'">
                            <div class="inline-flex gap-1">
                                <span
                                    >{{ $t('common.prefix') }}:
                                    <span class="font-mono">{{ activity.data.data.nameChange.prefix ?? $t('common.na') }}</span>
                                </span>
                                <span
                                    >{{ $t('common.suffix') }}:
                                    <span class="font-mono">{{ activity.data.data.nameChange.suffix ?? $t('common.na') }}</span>
                                </span>
                            </div>
                        </template>
                    </template>
                </h3>
                <p class="text-sm text-gray-400">
                    <GenericTime :value="activity.createdAt" type="long" />
                </p>
            </div>

            <div class="flex items-center justify-between">
                <p class="flex flex-col gap-1 text-sm">
                    <template v-if="activity.reason">
                        <div class="inline-flex gap-1">
                            <span class="font-semibold">{{ $t('common.reason') }}:</span>
                            <span>
                                {{ activity.reason }}
                            </span>
                        </div>
                    </template>
                    <template v-if="showTargetUser">
                        <div class="inline-flex items-center gap-1 text-sm">
                            <span class="font-semibold">{{ $t('common.colleague') }}:</span>
                            <ColleagueInfoPopover :user="activity.targetUser" />
                        </div>
                    </template>
                </p>
                <p class="inline-flex items-center gap-1 text-sm">
                    <span>{{ $t('common.created_by') }}</span>
                    <ColleagueInfoPopover :user="activity.sourceUser" :hide-props="true" />
                </p>
            </div>
        </div>
    </div>
</template>
