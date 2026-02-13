<script lang="ts" setup>
import type { TimeclockStats, TimeclockWeeklyStats } from '~~/gen/ts/resources/jobs/timeclock/timeclock';

withDefaults(
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

const isOpen = ref(false);
</script>

<template>
    <UDrawer
        v-model:open="isOpen"
        :ui="{ container: 'flex-1', content: 'min-h-[50%]', title: 'flex flex-row gap-2', body: 'h-full' }"
    >
        <UButton
            :label="$t('common.stats')"
            color="neutral"
            variant="subtle"
            icon="i-mdi-graph-line"
            trailing-icon="i-mdi-chevron-up"
            :ui="{
                leadingIcon: 'hidden sm:block',
                trailingIcon: 'group-data-[state=open]:rotate-180 transition-transform duration-200',
            }"
        />

        <template #title>
            <span class="flex-1">{{ $t('common.stats') }}</span>
            <UButton icon="i-mdi-close" color="neutral" variant="link" size="sm" @click="isOpen = false" />
        </template>

        <template #body>
            <LazyJobsTimeclockStatsBlock
                :weekly="weekly"
                :stats="stats"
                :hide-header="hideHeader"
                :failed="!!error"
                :loading="loading"
            />
        </template>
    </UDrawer>
</template>
