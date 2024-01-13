<script lang="ts" setup>
import { LoadingIcon } from 'mdi-vue3';
import { TimeclockStats } from '~~/gen/ts/resources/jobs/timeclock';

const props = defineProps<{
    stats?: TimeclockStats | null;
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
</script>

<template>
    <div class="mx-auto max-w-7xl">
        <div class="grid grid-cols-1 gap-px bg-neutral/5 sm:grid-cols-2 lg:grid-cols-3">
            <div v-for="stat in data" :key="stat.name" class="bg-gray-900 px-4 py-4 sm:px-4 lg:px-4">
                <p class="text-sm font-medium leading-6 text-gray-400">{{ $t(stat.name) }}</p>
                <p class="mt-2 w-full flex items-center gap-x-2 text-2xl font-semibold tracking-tight text-neutral">
                    <template v-if="stat.value === undefined">
                        <LoadingIcon class="h-5 w-5 animate-spin" />
                    </template>
                    <template v-else>
                        {{ fromSecondsToFormattedDuration(stat.value, { seconds: false }) }}
                    </template>
                </p>
            </div>
        </div>
    </div>
</template>
