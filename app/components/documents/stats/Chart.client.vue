<script setup lang="ts">
import { VisAxis, VisCrosshair, VisGroupedBar, VisTooltip, VisXYContainer } from '@unovis/vue';
import { addDays, addWeeks, differenceInDays, format, startOfDay, startOfWeek } from 'date-fns';
import { type GetStatsResponse, StatsCategory, StatsPeriod } from '~~/gen/ts/services/documents/stats';
import PenaltySeriesCardClient from './PenaltySeriesCard.client.vue';
import type { Range } from './helpers';

const cardRef = useTemplateRef<HTMLElement | null>('cardRef');

const props = defineProps<{
    stats: GetStatsResponse;
    category: StatsCategory;
    period: StatsPeriod;
    range: Range;
}>();

type DataRecord = {
    date: Date;
    amount: number;
    fine: number;
    detention: number;
    points: number;
};

const { width } = useElementSize(cardRef);

const selectedPeriod = computed(() => {
    if (props.period !== StatsPeriod.UNSPECIFIED) {
        return props.period;
    }

    const days = Math.max(differenceInDays(props.range.end, props.range.start), 1);
    if (days > 365) {
        return StatsPeriod.MONTHLY;
    }

    return days > 60 ? StatsPeriod.WEEKLY : StatsPeriod.DAILY;
});

const isPenalties = computed(() => props.category === StatsCategory.PENALTIES_OVER_TIME);

const bucketStart = (date: Date): Date => {
    if (selectedPeriod.value === StatsPeriod.MONTHLY) {
        return startOfDay(new Date(date.getFullYear(), date.getMonth(), 1));
    }

    return selectedPeriod.value === StatsPeriod.WEEKLY ? startOfWeek(date, { weekStartsOn: 1 }) : startOfDay(date);
};

const nextBucket = (date: Date): Date => {
    if (selectedPeriod.value === StatsPeriod.MONTHLY) {
        return startOfDay(new Date(date.getFullYear(), date.getMonth() + 1, 1));
    }

    return selectedPeriod.value === StatsPeriod.WEEKLY ? addWeeks(date, 1) : addDays(date, 1);
};

const data = computed(() => {
    const valueByBucket = new Map<number, DataRecord>();

    const ensureBucket = (time: number): DataRecord => {
        const existing = valueByBucket.get(time);
        if (existing) {
            return existing;
        }

        const created: DataRecord = {
            date: new Date(time),
            amount: 0,
            fine: 0,
            detention: 0,
            points: 0,
        };
        valueByBucket.set(time, created);
        return created;
    };

    if (isPenalties.value) {
        const values = props.stats?.periodSeriesValues ?? [];
        for (const item of values) {
            if (!item.day) {
                continue;
            }

            const bucket = bucketStart(toDate(item.day)).getTime();
            const target = ensureBucket(bucket);

            switch (item.key) {
                case 'fine_total':
                    target.fine += item.value;
                    break;
                case 'detention_time_total':
                    target.detention += item.value;
                    break;
                case 'stvo_points_total':
                    target.points += item.value;
                    break;
            }

            target.amount = target.fine + target.detention + target.points;
        }
    } else {
        const values = props.stats?.periodValues ?? [];
        for (const item of values) {
            if (!item.day) {
                continue;
            }

            const bucket = bucketStart(toDate(item.day)).getTime();
            const target = ensureBucket(bucket);
            target.amount += item.value;
        }
    }

    const filled: DataRecord[] = [];
    let cursor = bucketStart(props.range.start);
    const end = bucketStart(props.range.end).getTime();

    while (cursor.getTime() <= end) {
        const key = cursor.getTime();
        filled.push({
            date: new Date(key),
            amount: valueByBucket.get(key)?.amount ?? 0,
            fine: valueByBucket.get(key)?.fine ?? 0,
            detention: valueByBucket.get(key)?.detention ?? 0,
            points: valueByBucket.get(key)?.points ?? 0,
        });
        cursor = nextBucket(cursor);
    }

    return filled;
});

const x = (_: DataRecord, i: number) => i;
const y = computed(() =>
    isPenalties.value
        ? [(d: DataRecord) => d.fine, (d: DataRecord) => d.detention, (d: DataRecord) => d.points]
        : [(d: DataRecord) => d.amount],
);
const barColors = computed(() =>
    isPenalties.value ? ['var(--ui-primary)', 'var(--ui-warning)', 'var(--ui-success)'] : ['var(--ui-primary)'],
);
const fineColor = 'var(--ui-primary)';
const detentionColor = 'var(--ui-warning)';
const pointsColor = 'var(--ui-success)';

const total = computed(() => {
    const calculated = data.value?.reduce((acc: number, { amount }) => acc + amount, 0) ?? 0;
    if (isPenalties.value) {
        return calculated;
    }

    return props.stats?.totalValue ?? calculated;
});

const penaltyTotals = computed(() =>
    data.value.reduce(
        (acc, item) => {
            acc.fine += item.fine;
            acc.detention += item.detention;
            acc.points += item.points;
            return acc;
        },
        {
            fine: 0,
            detention: 0,
            points: 0,
        },
    ),
);

const fineSeries = computed(() =>
    data.value.map((item) => ({
        date: item.date,
        value: item.fine,
    })),
);

const detentionSeries = computed(() =>
    data.value.map((item) => ({
        date: item.date,
        value: item.detention,
    })),
);

const pointsSeries = computed(() =>
    data.value.map((item) => ({
        date: item.date,
        value: item.points,
    })),
);

const { format: formatNumber } = useIntlNumberFormat({
    style: 'decimal',
    currency: undefined,
});
const formatDate = (date: Date): string => format(date, 'd MMM');

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

    return `${formatDate(d.date)}: ${formatNumber(d.amount)}`;
};
</script>

<template>
    <UCard v-if="!isPenalties" ref="cardRef" :ui="{ root: 'overflow-visible', body: '!px-0 !pt-0 !pb-3' }">
        <template #header>
            <div>
                <p class="mb-1.5 text-xs text-muted uppercase">{{ $t('common.count') }}</p>
                <p class="text-3xl font-semibold text-highlighted">
                    {{ formatNumber(total) }}
                </p>
            </div>
        </template>

        <VisXYContainer :data="data ?? []" :padding="{ top: 40 }" class="h-96" :width="width">
            <VisGroupedBar :x="x" :y="y" :color="barColors" />

            <VisAxis type="x" :x="x" :tick-format="xTicks" />

            <VisCrosshair color="var(--ui-primary)" :template="template" />

            <VisTooltip />
        </VisXYContainer>
    </UCard>

    <div v-else class="grid grid-cols-1 gap-4 xl:grid-cols-3">
        <PenaltySeriesCardClient
            :title="$t('common.fine', 2)"
            :total="penaltyTotals.fine"
            :color="fineColor"
            :data="fineSeries"
            currency
        />

        <PenaltySeriesCardClient
            :title="$t('common.detention_time')"
            :total="penaltyTotals.detention"
            :color="detentionColor"
            :data="detentionSeries"
            detention
        />

        <PenaltySeriesCardClient
            :title="$t('common.traffic_infraction_points', 2)"
            :total="penaltyTotals.points"
            :color="pointsColor"
            :data="pointsSeries"
        />
    </div>
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
