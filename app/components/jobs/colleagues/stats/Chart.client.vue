<script setup lang="ts">
import { VisArea, VisAxis, VisCrosshair, VisGroupedBar, VisLine, VisTooltip, VisXYContainer } from '@unovis/vue';
import type { StatsPeriod } from '~~/gen/ts/resources/stats/stats';
import { buildChartData, type ChartStats, type DataRecord, type Range } from '../../../documents/stats/helpers';

const props = defineProps<{
    stats: ChartStats;
    period: StatsPeriod;
    range: Range;
}>();

const cardRef = useTemplateRef<HTMLElement | null>('cardRef');

const { width } = useElementSize(cardRef);

const data = computed(() => {
    return buildChartData({
        stats: props.stats,
        period: props.period,
        range: props.range,
        isPenalties: false,
    });
});

const x = (_: DataRecord, i: number) => i;
const y = computed(() => [(d: DataRecord) => d.amount]);

const total = computed(() => {
    const calculated = data.value?.reduce((acc: number, { amount }) => acc + amount, 0) ?? 0;
    return props.stats?.totalValue ?? calculated;
});
const averageVacation = computed(() => {
    const points = data.value ?? [];
    if (points.length === 0) {
        return 0;
    }

    const totalVacation = points.reduce((acc: number, point) => acc + point.vacation, 0);
    return totalVacation / points.length;
});

const { format: formatNumber } = useIntlNumberFormatWithOptions({
    style: 'decimal',
    currency: undefined,
    maximumFractionDigits: 0,
});
const { format: formatDate } = useDateFormatterWithOptions('short');

const xTicks = (i: number) => {
    if (!data.value?.[i]) {
        return '';
    }

    return formatDate(data.value[i].date);
};

const template = (d?: DataRecord) => {
    if (!d || !(d.date instanceof Date)) {
        return '';
    }

    return `${formatDate(d.date)}: ${formatNumber(d.amount)} (${formatNumber(d.vacation)})`;
};
</script>

<template>
    <UCard ref="cardRef" :ui="{ root: 'overflow-visible', body: '!px-0 !pt-0 !pb-3' }">
        <template #header>
            <div>
                <p class="mb-1.5 text-xs text-muted uppercase">
                    {{ stats.averageValue ? $t('common.avg') : $t('common.count') }}: {{ $t('common.colleague', 2) }}

                    ({{ $t('common.absent') }})
                </p>
                <p class="text-3xl font-semibold text-highlighted">
                    {{ formatNumber(stats.averageValue ? stats.averageValue : total) }}
                    ({{ formatNumber(averageVacation) }})
                </p>
            </div>
        </template>

        <VisXYContainer :data="data ?? []" :padding="{ top: 40 }" class="h-96" :width="width">
            <VisGroupedBar :x="x" :y="y" color="var(--ui-primary)" />

            <VisLine :x="x" :y="(d: DataRecord) => d.vacation" color="var(--ui-warning)" />
            <VisArea :x="x" :y="(d: DataRecord) => d.vacation" color="var(--ui-warning)" :opacity="0.1" />

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
