<script lang="ts" setup>
import { VisArea, VisAxis, VisCrosshair, VisLine, VisTooltip, VisXYContainer } from '@unovis/vue';
import { format, parse } from 'date-fns';
import type { TimeclockWeeklyStats } from '~~/gen/ts/resources/jobs/timeclock';

const { n, t } = useI18n();

const props = defineProps<{
    stats: TimeclockWeeklyStats[];
    loading?: boolean;
}>();

type DataRecord = TimeclockWeeklyStats & {
    date: Date;
};

const data = computed(() =>
    props.stats.map((s) => {
        const date = parse(s.calendarWeek.toString(), 'I', new Date());
        date.setFullYear(s.year);

        return {
            date: date,
            sum: s.sum,
            avg: s.avg,
            max: s.max,
        };
    }),
);

const cardRef = useTemplateRef('cardRef');
const { width } = useElementSize(cardRef);

const x = (_: DataRecord, i: number) => i;
const y = (d: DataRecord) => d.sum;

const total = computed(() => data.value.reduce((acc: number, { sum }) => acc + sum, 0));

const formatDate = (date: Date): string => format(date, `yyyy '${t('common.calendar_week')}' w`);

const xTicks = (i: number) => {
    return formatDate(data.value[i]?.date ?? new Date());
};

const template = (d: DataRecord) =>
    `<span class="font-semibold">${formatDate(d.date)}</span><br />
${t('components.jobs.timeclock.Stats.sum')}: ${n(d.sum, 'decimal')} h<br />
${t('components.jobs.timeclock.Stats.avg')}: ${n(d.avg, 'decimal')} h<br />
${t('components.jobs.timeclock.Stats.max')}: ${n(d.max, 'decimal')} h`;
</script>

<template>
    <UDashboardCard ref="cardRef" :ui="{ body: { padding: '!pb-3 !px-0' } as any }">
        <template #header>
            <div>
                <p class="mb-1 text-sm font-medium text-gray-500 dark:text-gray-400">
                    {{ $t('components.jobs.timeclock.Stats.sum') }}
                </p>
                <USkeleton v-if="loading" class="h-9 w-[275px]" />
                <p v-else class="text-3xl font-semibold text-gray-900 dark:text-white">
                    {{ fromSecondsToFormattedDuration(Math.ceil(total * 60 * 60), { seconds: false }) }}
                </p>
            </div>
        </template>

        <VisXYContainer class="h-80" :data="data" :padding="{ top: 10, left: 2, right: 2 }" :width="width">
            <VisLine :x="x" :y="y" color="rgb(var(--color-primary-DEFAULT))" />
            <VisArea :x="x" :y="y" color="rgb(var(--color-primary-DEFAULT))" :opacity="0.1" />

            <VisLine :x="x" :y="(d: DataRecord) => d.avg" color="rgb(var(--color-gray-500))" />
            <VisLine :x="x" :y="(d: DataRecord) => d.max" color="orange" />

            <VisAxis type="x" :x="x" :tick-format="xTicks" />

            <VisCrosshair color="rgb(var(--color-primary-600))" :template="template" />

            <VisTooltip />
        </VisXYContainer>
    </UDashboardCard>
</template>

<style scoped>
.unovis-xy-container {
    --vis-crosshair-line-stroke-color: rgb(var(--color-primary-500));
    --vis-crosshair-circle-stroke-color: #fff;

    --vis-axis-grid-color: rgb(var(--color-gray-200));
    --vis-axis-tick-color: rgb(var(--color-gray-200));
    --vis-axis-tick-label-color: rgb(var(--color-gray-400));

    --vis-tooltip-background-color: #fff;
    --vis-tooltip-border-color: rgb(var(--color-gray-200));
    --vis-tooltip-text-color: rgb(var(--color-gray-900));
}

.dark {
    .unovis-xy-container {
        --vis-crosshair-line-stroke-color: rgb(var(--color-primary-400));
        --vis-crosshair-circle-stroke-color: rgb(var(--color-gray-900));

        --vis-axis-grid-color: rgb(var(--color-gray-800));
        --vis-axis-tick-color: rgb(var(--color-gray-800));
        --vis-axis-tick-label-color: rgb(var(--color-gray-400));

        --vis-tooltip-background-color: rgb(var(--color-gray-900));
        --vis-tooltip-border-color: rgb(var(--color-gray-800));
        --vis-tooltip-text-color: #fff;
    }
}
</style>
