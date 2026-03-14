<script lang="ts" setup>
import type { StatsPeriod } from '~~/gen/ts/resources/stats/stats';
import { buildChartData, type ChartStats, type Range } from './helpers';
import PenaltySeriesCardClient from './PenaltySeriesCard.client.vue';

const props = defineProps<{
    stats: ChartStats;
    period: StatsPeriod;
    range: Range;
}>();

const fineColor = 'var(--ui-primary)';
const detentionColor = 'var(--ui-warning)';
const pointsColor = 'var(--ui-success)';

const data = computed(() => {
    return buildChartData({
        stats: props.stats,
        isPenalties: true,
        period: props.period,
        range: props.range,
    });
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
</script>

<template>
    <UPageGrid class="gap-4 sm:grid-cols-1 lg:grid-cols-1 xl:grid-cols-2">
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
    </UPageGrid>
</template>
