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
    if (i === 0 || i === data.value.length - 1 || !data.value[i]) {
        return '';
    }

    return formatDate(data.value[i]?.date ?? new Date());
};

const template = (d: DataRecord) =>
    `<span class="font-semibold">${formatDate(d.date)}</span><br />
${t('components.jobs.timeclock.Stats.sum')}: ${n(d.sum, 'decimal')} h<br />
${t('components.jobs.timeclock.Stats.avg')}: ${n(d.avg, 'decimal')} h<br />
${t('components.jobs.timeclock.Stats.max')}: ${n(d.max, 'decimal')} h`;
</script>

<template>
    <UCard ref="cardRef" :ui="{ root: 'overflow-visible', body: '!px-0 !pt-0 !pb-3' }">
        <template #header>
            <div>
                <p class="mb-1.5 text-xs text-muted uppercase">
                    {{ $t('components.jobs.timeclock.Stats.sum') }}
                </p>
                <USkeleton v-if="loading" class="h-9 w-[275px]" />
                <p v-else class="text-3xl font-semibold text-highlighted">
                    {{ fromSecondsToFormattedDuration(Math.ceil(total * 60 * 60), { seconds: false }) }}
                </p>
            </div>
        </template>

        <VisXYContainer class="h-90" :data="data" :padding="{ top: 10 }" :width="width">
            <VisLine :x="x" :y="y" color="var(--ui-primary)" />
            <VisArea :x="x" :y="y" color="var(--ui-primary)" :opacity="0.1" />

            <VisLine :x="x" :y="(d: DataRecord) => d.avg" color="var(--color-gray-500)" />
            <VisLine :x="x" :y="(d: DataRecord) => d.max" color="orange" />

            <VisAxis type="x" :x="x" :tick-format="xTicks" />

            <VisCrosshair color="var(--ui-primary)" :template="template" />

            <VisTooltip />
        </VisXYContainer>
    </UCard>
</template>

<style scoped>
.unovis-xy-container {
    --vis-crosshair-line-stroke-color: var(--ui-primary);
    --vis-crosshair-circle-stroke-color: var(--ui-bg);

    --vis-axis-grid-color: var(--ui-border);
    --vis-axis-tick-color: var(--ui-border);
    --vis-axis-tick-label-color: var(--ui-text-dimmed);

    --vis-tooltip-background-color: var(--ui-bg);
    --vis-tooltip-border-color: var(--ui-border);
    --vis-tooltip-text-color: var(--ui-text-highlighted);
}
</style>
