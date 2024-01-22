<script lang="ts" setup>
import { LoadingIcon } from 'mdi-vue3';
import { Bar, Chart, Grid, Responsive, Tooltip } from 'vue3-charts';
import type { ChartAxis } from 'vue3-charts/dist/types';
import GenericContainer from '~/components/partials/elements/GenericContainer.vue';
import { TimeclockStats, TimeclockWeeklyStats } from '~~/gen/ts/resources/jobs/timeclock';

const props = defineProps<{
    stats?: TimeclockStats | null;
    weekly?: TimeclockWeeklyStats[];
    hideHeader?: boolean;
}>();

const data = ref<Record<string, { name: string; value?: number }>>({
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

    data.value.sum.value = parseFloat(((Math.round(props.stats.spentTimeSum * 100) / 100) * 60 * 60).toPrecision(2));
    data.value.avg.value = parseFloat(((Math.round(props.stats.spentTimeAvg * 100) / 100) * 60 * 60).toPrecision(2));
    data.value.max.value = parseFloat(((Math.round(props.stats.spentTimeMax * 100) / 100) * 60 * 60).toPrecision(2));
}

watch(props, async () => updateStats());

onBeforeMount(async () => updateStats());

console.log('weekly', props.weekly);

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
            <h2 v-if="!hideHeader" class="text-center text-2xl font-semibold text-neutral">
                {{ $t('common.timeclock') }}
            </h2>
            <div class="flex flex-col gap-4 sm:flex-row">
                <div class="flex-0">
                    <h3 class="mb-2 text-lg font-semibold text-neutral">
                        {{ $t('components.jobs.timeclock.StatsBlock.7_days') }}
                    </h3>
                    <div class="grid grid-cols-1 gap-2 sm:grid-rows-2 lg:grid-rows-3">
                        <GenericContainer v-for="stat in data" :key="stat.name" class="bg-gray-900">
                            <p class="text-sm font-medium leading-6 text-gray-400">{{ $t(stat.name) }}</p>
                            <p class="mt-2 flex w-full items-center gap-x-2 text-2xl font-semibold tracking-tight text-neutral">
                                <template v-if="stat.value === undefined">
                                    <LoadingIcon class="h-5 w-5 animate-spin" />
                                </template>
                                <template v-else>
                                    {{ fromSecondsToFormattedDuration(stat.value, { seconds: false }) }}
                                </template>
                            </p>
                        </GenericContainer>
                    </div>
                </div>

                <div class="w-full flex-1">
                    <h3 class="mb-2 text-lg font-semibold text-neutral">
                        {{ $t('components.jobs.timeclock.StatsBlock.weekly') }}
                    </h3>

                    <Responsive class="w-full">
                        <template #main="{ width }">
                            <Chart
                                :size="{ width, height: 350 }"
                                :data="weekly"
                                :margin="margin"
                                direction="horizontal"
                                :axis="axis"
                                class="text-neutral"
                            >
                                <template #layers>
                                    <Grid stroke-dasharray="2,2" />
                                    <Bar :data-keys="['date', 'sum']" :bar-style="{ fill: '#443a8f' }" :gap="8" />
                                    <Bar :data-keys="['date', 'avg']" :bar-style="{ fill: '#1f236e' }" :gap="8" />
                                    <Bar :data-keys="['date', 'max']" :bar-style="{ fill: '#8d81f2' }" :gap="8" />
                                </template>
                                <template #widgets>
                                    <Tooltip
                                        border-color="#48CAE4"
                                        :config="{
                                            date: { label: $t('common.date'), color: '#2b2d34' },
                                            sum: { label: $t('components.jobs.timeclock.StatsBlock.sum'), color: '#443a8f' },
                                            avg: { label: $t('components.jobs.timeclock.StatsBlock.avg'), color: '#1f236e' },
                                            max: { label: $t('components.jobs.timeclock.StatsBlock.max'), color: '#8d81f2' },
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
