<script lang="ts" setup>
import ColleagueInfoPopover from '~/components/jobs/colleagues/ColleagueInfoPopover.vue';
import LabelBadge from '~/components/jobs/colleagues/labels/LabelBadge.vue';
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
    <div class="flex min-w-0 flex-auto justify-between gap-2 text-xs leading-5 text-muted">
        <div class="flex min-w-0 flex-col gap-0.5">
            <span class="text-sm font-medium text-highlighted">
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
                            <LabelBadge
                                v-for="label in activity.data.data.labelsChange?.removed"
                                :key="label.id"
                                :label="label"
                                removed
                            />

                            <LabelBadge
                                v-for="label in activity.data.data.labelsChange?.added"
                                :key="label.id"
                                :label="label"
                            />

                            <LabelBadge
                                v-for="labelId in activity.data.data.labelsChange?.removedIds"
                                :id="labelId"
                                :key="labelId"
                                removed
                            />

                            <LabelBadge
                                v-for="labelId in activity.data.data.labelsChange?.addedIds"
                                :id="labelId"
                                :key="labelId"
                            />
                        </span>
                    </template>

                    <template v-else-if="activity.data?.data.oneofKind === 'nameChange'">
                        <span class="inline-flex flex-wrap gap-1">
                            <span
                                >{{ $t('common.prefix') }}:
                                <span class="font-mono">{{
                                    activity.data.data.nameChange.prefix ?? $t('common.na')
                                }}</span></span
                            >
                            <span
                                >{{ $t('common.suffix') }}:
                                <span class="font-mono">{{
                                    activity.data.data.nameChange.suffix ?? $t('common.na')
                                }}</span></span
                            >
                        </span>
                    </template>
                </template>
            </span>

            <div class="flex min-w-0 flex-col gap-0.5 text-sm break-words">
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
            </div>
        </div>

        <div class="flex flex-none flex-col items-end gap-0.5 text-xs">
            <span class="text-xs/5 text-dimmed"><GenericTime :value="activity.createdAt" type="short" /></span>

            <span class="inline-flex items-center gap-1">
                <span class="text-dimmed">{{ $t('common.created_by') }}:</span>
                <ColleagueInfoPopover :user="activity.sourceUser" />
            </span>
        </div>
    </div>
</template>
