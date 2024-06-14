<script lang="ts" setup>
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import { TimeclockStats, TimeclockWeeklyStats } from '~~/gen/ts/resources/jobs/timeclock';
import TimeclockStatsChart from '~/components/jobs/timeclock/TimeclockStatsChart.vue';

const props = withDefaults(
    defineProps<{
        stats?: TimeclockStats | null;
        weekly?: TimeclockWeeklyStats[];
        hideHeader?: boolean;
        failed?: boolean;
        loading?: boolean;
    }>(),
    {
        stats: null,
        weekly: undefined,
        hideHeader: false,
        failed: false,
        loading: false,
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

    statsData.value.sum.value = parseFloat(((Math.round(props.stats.spentTimeSum * 100) / 100) * 60 * 60).toPrecision(2));
    statsData.value.avg.value = parseFloat(((Math.round(props.stats.spentTimeAvg * 100) / 100) * 60 * 60).toPrecision(2));
    statsData.value.max.value = parseFloat(((Math.round(props.stats.spentTimeMax * 100) / 100) * 60 * 60).toPrecision(2));
}

watch(props, async () => updateStats());

onBeforeMount(async () => updateStats());

const loadingState = ref(false);
watch(
    () => props.loading,
    () => {
        if (props.loading) {
            loadingState.value = true;
        }
    },
);
watchDebounced(
    () => props.loading,
    () => {
        if (!props.loading) {
            loadingState.value = false;
        }
    },
    {
        debounce: 750,
        maxWait: 1250,
    },
);
</script>

<template>
    <UCard>
        <template v-if="!hideHeader" #header>
            <h2 class="inline-flex w-full items-center justify-between text-lg font-semibold">
                {{ $t('common.timeclock') }}

                <UButton
                    variant="link"
                    icon="i-mdi-refresh"
                    :title="$t('common.refresh')"
                    :disabled="loading || loadingState"
                    :loading="loading || loadingState"
                    @click="$emit('refresh')"
                >
                    {{ $t('common.refresh') }}
                </UButton>
            </h2>
        </template>

        <div class="flex flex-col gap-4 lg:flex-row">
            <div class="flex-none">
                <h3 class="mb-2 ml-0.5 text-lg font-bold text-gray-900 dark:text-white">
                    {{ $t('components.jobs.timeclock.Stats.7_days') }}
                </h3>
                <div class="grid grid-cols-1 gap-2">
                    <UCard v-for="stat in statsData" :key="stat.name">
                        <p class="text-sm font-medium leading-6 text-gray-500 dark:text-gray-400">{{ $t(stat.name) }}</p>
                        <p
                            class="mt-2 flex w-full items-center gap-x-2 text-2xl font-semibold tracking-tight text-gray-900 dark:text-white"
                        >
                            <UIcon v-if="failed" name="i-mdi-alert-circle" class="size-5" />
                            <USkeleton v-else-if="stat.value === undefined" class="h-8 w-[175px]" />
                            <template v-else>
                                {{
                                    fromSecondsToFormattedDuration(stat.value, {
                                        seconds: false,
                                        emptyText: 'common.none',
                                    })
                                }}
                            </template>
                        </p>
                    </UCard>
                </div>
            </div>

            <div class="flex-1">
                <h3 class="mb-2 text-lg font-bold text-gray-900 dark:text-white">
                    {{ $t('components.jobs.timeclock.Stats.weekly') }}
                </h3>

                <DataErrorBlock v-if="failed" :retry="async () => $emit('refresh')" />
                <DataNoDataBlock v-else-if="weekly === undefined" />
                <ClientOnly v-else>
                    <TimeclockStatsChart :stats="weekly" :loading="loading" />
                </ClientOnly>
            </div>
        </div>
    </UCard>
</template>
