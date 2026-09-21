<script lang="ts" setup>
import ColleagueInfoPopover from '~/components/jobs/colleagues/ColleagueInfoPopover.vue';
import { jobsUserActivityTypeBGColor, jobsUserActivityTypeIcon } from '~/components/jobs/colleagues/info/helpers';
import ActivityFeedItem from '~/components/partials/data/ActivityFeedItem.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { ColleagueActivityType, type ColleagueActivity } from '~~/gen/ts/resources/jobs/colleagues/activity/activity';

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
    <ActivityFeedItem
        :icon="jobsUserActivityTypeIcon(activity.activityType)"
        :icon-class="jobsUserActivityTypeBGColor(activity.activityType)"
        inline
    >
        <template #title>
            {{ $t(`enums.jobs.ColleagueActivityType.${ColleagueActivityType[activity.activityType]}`) }}
            <template v-if="activity.data?.data.oneofKind !== undefined">
                {{ '&nbsp;-&nbsp;' }}
                <template v-if="activity.data?.data.oneofKind === 'absenceDate'">
                    <span
                        v-if="activity.data?.data.absenceDate.absenceBegin && activity.data?.data.absenceDate.absenceEnd"
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
                    <span class="inline-flex flex-wrap gap-1">
                        <UBadge
                            v-for="label in activity.data.data.labelsChange?.removed"
                            :key="label.name"
                            class="justify-between gap-2 line-through"
                            :class="isColorBright(hexToRgb(label.color, rgbBlack)!) ? 'text-black!' : 'text-white!'"
                            :style="{ backgroundColor: label.color }"
                            size="md"
                            :label="label.name"
                        />

                        <UBadge
                            v-for="label in activity.data.data.labelsChange?.added"
                            :key="label.name"
                            class="justify-between gap-2"
                            :class="isColorBright(hexToRgb(label.color, rgbBlack)!) ? 'text-black!' : 'text-white!'"
                            :style="{ backgroundColor: label.color }"
                            size="md"
                            :label="label.name"
                        />
                    </span>
                </template>

                <template v-else-if="activity.data?.data.oneofKind === 'nameChange'">
                    <span class="inline-flex flex-wrap gap-1">
                        <span
                            >{{ $t('common.prefix') }}:
                            <span class="font-mono">{{ activity.data.data.nameChange.prefix ?? $t('common.na') }}</span></span
                        >
                        <span
                            >{{ $t('common.suffix') }}:
                            <span class="font-mono">{{ activity.data.data.nameChange.suffix ?? $t('common.na') }}</span></span
                        >
                    </span>
                </template>
            </template>
        </template>

        <template #timestamp>
            <GenericTime :value="activity.createdAt" type="long" />
        </template>

        <div class="flex flex-col gap-2 sm:flex-row sm:items-start sm:justify-between">
            <p class="flex min-w-0 flex-col gap-1 text-sm break-words">
                <template v-if="activity.reason">
                    <span class="inline-flex gap-1">
                        <span class="font-semibold">{{ $t('common.reason') }}:</span>
                        <span>{{ activity.reason }}</span>
                    </span>
                </template>
                <template v-if="showTargetUser">
                    <span class="inline-flex items-center gap-1 text-sm">
                        <span class="font-semibold">{{ $t('common.colleague') }}:</span>
                        <ColleagueInfoPopover :user="activity.targetUser" />
                    </span>
                </template>
            </p>

            <p class="inline-flex shrink-0 items-center gap-1 text-sm">
                <span>{{ $t('common.created_by') }}</span>
                <ColleagueInfoPopover :user="activity.sourceUser" />
            </p>
        </div>
    </ActivityFeedItem>
</template>
