<script lang="ts" setup>
import { AlertCircleIcon, LoadingIcon } from 'mdi-vue3';
import { Bar, Chart, Grid, Responsive, Tooltip } from 'vue3-charts';
import type { ChartAxis } from 'vue3-charts/dist/types';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import GenericContainer from '~/components/partials/elements/GenericContainer.vue';
import { TimeclockStats, TimeclockWeeklyStats } from '~~/gen/ts/resources/jobs/timeclock';

const props = withDefaults(
    defineProps<{
        stats?: TimeclockStats | null;
        weekly?: TimeclockWeeklyStats[];
        hideHeader?: boolean;
        failed?: boolean;
    }>(),
    {
        stats: null,
        weekly: undefined,
        hideHeader: false,
    },
);

defineEmits<{
    (e: 'refresh'): void;
}>();

const statsData = ref<Record<string, { name: string; value?: number }>>({
    sum: {
        name: 'components.jobs.timeclock.Stats.sum',
    },
    avg: {
        name: 'components.jobs.timeclock.Stats.avg',
    },
    max: {
        name: 'components.jobs.timeclock.Stats.max',
    },
});

async function updateStats(): Promise<void> {
    if (!props.stats) {
        return;
    }
    console.log(props.stats);

    statsData.value.sum.value = parseFloat(((Math.round(props.stats.spentTimeSum * 100) / 100) * 60 * 60).toPrecision(2));
    statsData.value.avg.value = parseFloat(((Math.round(props.stats.spentTimeAvg * 100) / 100) * 60 * 60).toPrecision(2));
    statsData.value.max.value = parseFloat(((Math.round(props.stats.spentTimeMax * 100) / 100) * 60 * 60).toPrecision(2));
}

watch(props, async () => updateStats());

onBeforeMount(async () => updateStats());

const margin = ref({
    left: 5,
    top: 5,
    right: 5,
    bottom: 5,
});

const axis = ref<ChartAxis>({
    primary: {
        type: 'band',
        domain: ['data', 'dataMax'],
    },
    secondary: {
        domain: [0, 'dataMax + 2'],
        type: 'linear',
        ticks: 8,
        format: (v) => `${v}h`,
    },
});
</script>

<template>
    <div class="mx-auto max-w-7xl">
        <GenericContainer>
            <h2 v-if="!hideHeader" class="text-2xl font-semibold text-neutral">
                {{ $t('common.timeclock') }}
            </h2>
            <div class="flex flex-col gap-4 lg:flex-row">
                <div class="flex-0">
                    <h3 class="mb-2 ml-0.5 text-lg font-bold text-neutral">
                        {{ $t('components.jobs.timeclock.StatsBlock.7_days') }}
                    </h3>
                    <div class="grid grid-cols-1 gap-2">
                        <GenericContainer v-for="stat in statsData" :key="stat.name" class="bg-primary-900">
                            <p class="text-sm font-medium leading-6 text-gray-300">{{ $t(stat.name) }}</p>
                            <p class="mt-2 flex w-full items-center gap-x-2 text-2xl font-semibold tracking-tight text-neutral">
                                <template v-if="stat.value === undefined">
                                    <LoadingIcon class="h-5 w-5 animate-spin" aria-hidden="true" />
                                </template>
                                <template v-else-if="failed">
                                    <AlertCircleIcon class="h-5 w-5" aria-hidden="true" />
                                </template>
                                <template v-else>
                                    {{
                                        fromSecondsToFormattedDuration(stat.value, {
                                            seconds: false,
                                            emptyText: $t('common.none'),
                                        })
                                    }}
                                </template>
                            </p>
                        </GenericContainer>
                    </div>
                </div>

                <div class="w-full flex-1">
                    <h3 class="mb-2 text-lg font-bold text-neutral">
                        {{ $t('components.jobs.timeclock.StatsBlock.weekly') }}
                    </h3>

                    <DataErrorBlock v-if="failed" :retry="async () => $emit('refresh')" />
                    <DataNoDataBlock v-else-if="weekly === undefined" />
                    <Responsive v-else class="w-ful">
                        <template #main="{ width }">
                            <!-- @vue-ignore our own data format works fine.. but the package type is "off" -->
                            <Chart
                                :size="{ width: width as number, height: 375 }"
                                :data="weekly"
                                :margin="margin"
                                direction="horizontal"
                                :axis="axis"
                                class="text-neutral"
                            >
                                <template #layers>
                                    <Grid stroke-dasharray="2,2" />
                                    <Bar :data-keys="['date', 'sum']" :bar-style="{ class: 'fill-primary-600' }" :gap="12" />
                                    <Bar :data-keys="['date', 'avg']" :bar-style="{ class: 'fill-primary-800' }" :gap="12" />
                                    <Bar :data-keys="['date', 'max']" :bar-style="{ class: 'fill-primary-400' }" :gap="12" />
                                </template>
                                <template #widgets>
                                    <Tooltip
                                        border-color="#48CAE4"
                                        :config="{
                                            date: { label: $t('common.date'), color: '#2b2d34' },
                                            sum: {
                                                label: $t('components.jobs.timeclock.StatsBlock.sum'),
                                            },
                                            avg: {
                                                label: $t('components.jobs.timeclock.StatsBlock.avg'),
                                            },
                                            max: {
                                                label: $t('components.jobs.timeclock.StatsBlock.max'),
                                            },
                                        }"
                                    />
                                </template>
                            </Chart>
                        </template>
                    </Responsive>
                </div>
            </div>
        </GenericContainer>
    </div>
</template>
