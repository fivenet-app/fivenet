<script lang="ts" setup>
import TimeclockStatsChart from '~/components/jobs/timeclock/TimeclockStatsChart.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import type { TimeclockStats, TimeclockWeeklyStats } from '~~/gen/ts/resources/jobs/timeclock';

const props = withDefaults(
    defineProps<{
        stats?: TimeclockStats | null;
        weekly?: TimeclockWeeklyStats[];
        hideHeader?: boolean;
        error?: Error;
        loading?: boolean;
    }>(),
    {
        stats: null,
        weekly: undefined,
        hideHeader: false,
        error: undefined,
        loading: false,
    },
);

defineEmits<{
    (e: 'refresh'): void;
}>();

type Stats = { name: string; value?: number };

const statsData = ref<{
    sum: Stats;
    avg: Stats;
    max: Stats;
}>({
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

                <UTooltip :text="$t('common.refresh')">
                    <UButton
                        variant="link"
                        icon="i-mdi-refresh"
                        :disabled="loading || loadingState"
                        :loading="loading || loadingState"
                        @click="$emit('refresh')"
                    >
                        {{ $t('common.refresh') }}
                    </UButton>
                </UTooltip>
            </h2>
        </template>

        <div class="flex flex-col gap-4 lg:flex-row">
            <div class="flex-none">
                <h3 class="text-highlighted mb-2 ml-0.5 text-lg font-bold">
                    {{ $t('components.jobs.timeclock.Stats.7_days') }}
                </h3>
                <div class="grid grid-cols-1 gap-2">
                    <UCard v-for="stat in statsData" :key="stat.name">
                        <p class="text-muted text-sm font-medium leading-6">{{ $t(stat.name) }}</p>
                        <p class="text-highlighted mt-2 flex w-full items-center gap-x-2 text-2xl font-semibold tracking-tight">
                            <UIcon v-if="error" class="size-5" name="i-mdi-alert-circle" />
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
                <h3 class="text-highlighted mb-2 text-lg font-bold">
                    {{ $t('components.jobs.timeclock.Stats.weekly') }}
                </h3>

                <DataErrorBlock v-if="error" :error="error" :retry="async () => $emit('refresh')" />
                <DataNoDataBlock v-else-if="weekly === undefined" />

                <ClientOnly v-else>
                    <TimeclockStatsChart :stats="weekly" :loading="loading" />
                </ClientOnly>
            </div>
        </div>
    </UCard>
</template>
