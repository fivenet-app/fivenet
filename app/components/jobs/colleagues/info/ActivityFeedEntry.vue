<script lang="ts" setup>
import ColleagueInfoPopover from '~/components/jobs/colleagues/ColleagueInfoPopover.vue';
import { jobsUserActivityTypeBGColor, jobsUserActivityTypeIcon } from '~/components/jobs/colleagues/info/helpers';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { ColleagueActivityType, type ColleagueActivity } from '~~/gen/ts/resources/jobs/activity';

withDefaults(
    defineProps<{
        activity: ColleagueActivity;
        showTargetUser?: boolean;
    }>(),
    {
        showTargetUser: false,
    },
);
</script>

<template>
    <li
        class="border-default px-2 py-4 hover:border-primary-500/25 hover:bg-primary-100/50 dark:hover:border-primary-400/25 dark:hover:bg-primary-900/10"
    >
        <div class="flex space-x-3">
            <div class="my-auto flex size-10 items-center justify-center rounded-full">
                <UIcon
                    :class="[jobsUserActivityTypeBGColor(activity.activityType), 'size-full']"
                    :name="jobsUserActivityTypeIcon(activity.activityType)"
                    :inline="true"
                />
            </div>

            <div class="flex-1 space-y-1">
                <div class="flex items-center justify-between">
                    <h3 class="text-sm font-medium">
                        {{ $t(`enums.jobs.ColleagueActivityType.${ColleagueActivityType[activity.activityType]}`) }}
                        <template v-if="activity.data?.data.oneofKind !== undefined">
                            {{ '&nbsp;-&nbsp;' }}
                            <template v-if="activity.data?.data.oneofKind === 'absenceDate'">
                                <span
                                    v-if="
                                        activity.data?.data.absenceDate.absenceBegin &&
                                        activity.data?.data.absenceDate.absenceEnd
                                    "
                                    class="inline-flex gap-1"
                                >
                                    <GenericTime :value="activity.data?.data.absenceDate.absenceBegin" type="date" />
                                    <span>{{ $t('common.to') }}</span>
                                    <GenericTime :value="activity.data?.data.absenceDate.absenceEnd" type="date" />
                                </span>
                                <span v-else>{{ $t('common.annul', 2) }}</span>
                            </template>

                            <template v-else-if="activity.data?.data.oneofKind === 'gradeChange'">
                                {{ activity.data?.data.gradeChange.gradeLabel }} ({{ activity.data?.data.gradeChange.grade }})
                            </template>

                            <template v-else-if="activity.data?.data.oneofKind === 'labelsChange'">
                                <div class="inline-flex gap-1">
                                    <UBadge
                                        v-for="label in activity.data.data.labelsChange?.removed"
                                        :key="label.name"
                                        class="justify-between gap-2 line-through"
                                        :class="isColorBright(hexToRgb(label.color, rgbBlack)!) ? 'text-black!' : 'text-white!'"
                                        :style="{ backgroundColor: label.color }"
                                        size="md"
                                    >
                                        {{ label.name }}
                                    </UBadge>

                                    <UBadge
                                        v-for="label in activity.data.data.labelsChange?.added"
                                        :key="label.name"
                                        class="justify-between gap-2"
                                        :class="isColorBright(hexToRgb(label.color, rgbBlack)!) ? 'text-black!' : 'text-white!'"
                                        :style="{ backgroundColor: label.color }"
                                        size="md"
                                    >
                                        {{ label.name }}
                                    </UBadge>
                                </div>
                            </template>

                            <template v-else-if="activity.data?.data.oneofKind === 'nameChange'">
                                <div class="inline-flex gap-1">
                                    <span
                                        >{{ $t('common.prefix') }}:
                                        <span class="font-mono">{{
                                            activity.data.data.nameChange.prefix ?? $t('common.na')
                                        }}</span>
                                    </span>
                                    <span
                                        >{{ $t('common.suffix') }}:
                                        <span class="font-mono">{{
                                            activity.data.data.nameChange.suffix ?? $t('common.na')
                                        }}</span>
                                    </span>
                                </div>
                            </template>
                        </template>
                    </h3>

                    <p class="text-sm text-dimmed">
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
    </li>
</template>
