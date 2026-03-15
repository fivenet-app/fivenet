<script setup lang="ts">
import { VisAxis, VisCrosshair, VisGroupedBar, VisLine, VisTooltip, VisXYContainer } from '@unovis/vue';
import type { StatsPeriod } from '~~/gen/ts/resources/stats/stats';
import { buildChartData, type ChartStats, type DataRecord, type Range } from '../../../documents/stats/helpers';

const props = defineProps<{
    stats: ChartStats;
    isPenalties: boolean;
    period: StatsPeriod;
    range: Range;
}>();

const cardRef = useTemplateRef<HTMLElement | null>('cardRef');

const { width } = useElementSize(cardRef);

const isPenalties = toRef(props, 'isPenalties');

const data = computed(() => {
    return buildChartData({
        stats: props.stats,
        isPenalties: props.isPenalties,
        period: props.period,
        range: props.range,
    });
});

const x = (_: DataRecord, i: number) => i;
const y = computed(() =>
    props.isPenalties
        ? [(d: DataRecord) => d.fine, (d: DataRecord) => d.detention, (d: DataRecord) => d.points]
        : [(d: DataRecord) => d.amount],
);
const barColors = computed(() =>
    props.isPenalties ? ['var(--ui-primary)', 'var(--ui-warning)', 'var(--ui-success)'] : ['var(--ui-primary)'],
);

const total = computed(() => {
    const calculated = data.value?.reduce((acc: number, { amount }) => acc + amount, 0) ?? 0;
    if (props.isPenalties) {
        return calculated;
    }

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
    <UCard v-if="!isPenalties" ref="cardRef" :ui="{ root: 'overflow-visible', body: '!px-0 !pt-0 !pb-3' }">
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
            <VisGroupedBar :x="x" :y="y" :color="barColors" />
            <VisLine :x="x" :y="(d: DataRecord) => d.vacation" color="var(--ui-warning)" />

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
