<script lang="ts" setup>
import { TimeclockStats } from '~~/gen/ts/resources/jobs/timeclock';

const props = defineProps<{
    stats: TimeclockStats;
}>();

const entries = ref<{ name: string; value: number }[]>([]);

function updateStats(): void {
    entries.value.length = 0;
    entries.value.push({
        name: 'components.jobs.timeclock.Stats.sum',
        value: props.stats.spentTimeSum * 60 * 60,
    });
    entries.value.push({
        name: 'components.jobs.timeclock.Stats.avg',
        value: props.stats.spentTimeAvg * 60 * 60,
    });
    entries.value.push({
        name: 'components.jobs.timeclock.Stats.max',
        value: props.stats.spentTimeMax * 60 * 60,
    });
}

watch(props, () => updateStats());

onBeforeMount(() => updateStats());
</script>

<template>
    <div class="mx-auto max-w-7xl">
        <div class="grid grid-cols-1 gap-px bg-white/5 sm:grid-cols-2 lg:grid-cols-3">
            <div v-for="stat in entries" class="bg-gray-900 px-4 py-4 sm:px-4 lg:px-4">
                <p class="text-sm font-medium leading-6 text-gray-400">{{ $t(stat.name) }}</p>
                <p class="mt-2 flex items-baseline gap-x-2">
                    <span class="text-3xl font-semibold tracking-tight text-white">{{
                        fromSecondsToFormattedDuration(stat.value, { seconds: false })
                    }}</span>
                </p>
            </div>
        </div>
    </div>
</template>
