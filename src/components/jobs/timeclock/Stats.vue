<script lang="ts" setup>
import { TimeclockStats } from '~~/gen/ts/resources/jobs/timeclock';

const props = defineProps<{
    stats?: TimeclockStats;
}>();

const entries = ref<{ name: string; value: number }[]>([]);

async function updateStats(): Promise<void> {
    if (!props.stats) {
        return;
    }

    entries.value.length = 0;
    entries.value.push({
        name: 'components.jobs.timeclock.Stats.sum',
        value: parseFloat(((Math.round(props.stats.spentTimeSum * 100) / 100) * 60 * 60).toPrecision(2)),
    });
    entries.value.push({
        name: 'components.jobs.timeclock.Stats.avg',
        value: parseFloat(((Math.round(props.stats.spentTimeAvg * 100) / 100) * 60 * 60).toPrecision(2)),
    });
    entries.value.push({
        name: 'components.jobs.timeclock.Stats.max',
        value: parseFloat(((Math.round(props.stats.spentTimeMax * 100) / 100) * 60 * 60).toPrecision(2)),
    });
}

watch(props, async () => updateStats());

onBeforeMount(async () => updateStats());
</script>

<template>
    <div class="mx-auto max-w-7xl">
        <div class="grid grid-cols-1 gap-px bg-neutral/5 sm:grid-cols-2 lg:grid-cols-3">
            <div v-for="stat in entries" class="bg-gray-900 px-4 py-4 sm:px-4 lg:px-4">
                <p class="text-sm font-medium leading-6 text-gray-400">{{ $t(stat.name) }}</p>
                <p class="mt-2 flex items-baseline gap-x-2">
                    <span class="text-3xl font-semibold tracking-tight text-neutral">{{
                        fromSecondsToFormattedDuration(stat.value, { seconds: false })
                    }}</span>
                </p>
            </div>
        </div>
    </div>
</template>
