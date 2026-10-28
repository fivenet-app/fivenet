<script lang="ts" setup>
import { getJobsColleaguesClient } from '~~/gen/ts/clients';
import type { Label } from '~~/gen/ts/resources/jobs/labels/labels';

const props = defineProps<{
    id?: number;
    label?: Label;
    removed?: boolean;
}>();

const jobsColleaguesClient = await getJobsColleaguesClient();

const {
    data: labels,
    pending,
    error,
} = useAuthedLazyAsyncData(
    'userState',
    'jobs-colleague-labels',
    async ({ signal }) => {
        const { response } = await jobsColleaguesClient.getColleagueLabels({}, { abort: signal });
        return response?.labels ?? [];
    },
    {
        immediate: !props.label && props.id !== undefined,
        default: () => [] as Label[],
    },
);

const label = computed(() => props.label ?? labels.value.find((label) => label.id === props.id));
</script>

<template>
    <USkeleton v-if="!props.label && pending" class="h-6 w-[120px]" />
    <UBadge v-else-if="error" :label="$t('common.error')" />
    <UBadge
        v-else-if="label"
        class="justify-between gap-2"
        :class="[isColorBright(hexToRgb(label.color, rgbBlack)!) ? 'text-black!' : 'text-white!', removed && 'line-through']"
        :style="{ backgroundColor: label.color }"
        :label="label.name"
    />
    <UBadge v-else icon="i-mdi-question-mark" variant="subtle" :label="`${$t('common.unknown')} (${props.id})`" />
</template>
