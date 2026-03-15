<script lang="ts" setup>
import { differenceInDays, sub } from 'date-fns';
import { z } from 'zod';
import DateRangePicker from '~/components/documents/stats/DateRangePicker.vue';
import ChartClient from '~/components/jobs/colleagues/stats/Chart.client.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import RefreshButton from '~/components/partials/RefreshButton.vue';
import { getJobsStatsClient } from '~~/gen/ts/clients';
import { StatsCategory, StatsPeriod } from '~~/gen/ts/resources/stats/stats';

useHead({
    title: 'pages.jobs.colleagues.stats.title',
});

definePageMeta({
    title: 'pages.jobs.colleagues.stats.title',
    requiresAuth: true,
    permission: 'jobs.StatsService/GetStats',
});

const schema = z.object({
    range: z
        .object({
            start: z.date().default(() => sub(new Date(), { days: 14 })),
            end: z.date().default(() => new Date()),
        })
        .default(() => ({
            start: sub(new Date(), { days: 14 }),
            end: new Date(),
        })),
    category: z.enum(StatsCategory).default(StatsCategory.EMPLOYEE_COUNT_OVER_TIME),
});

const query = useSearchForm('jobs-colleagues-stats', schema);

const selectedPeriod = computed(() => {
    const days = Math.max(differenceInDays(query.range.end, query.range.start), 1);
    if (days > 365) {
        return StatsPeriod.MONTHLY;
    }

    return days > 60 ? StatsPeriod.WEEKLY : StatsPeriod.DAILY;
});

const jobsStatsClient = await getJobsStatsClient();

const {
    data: response,
    status,
    error,
    refresh,
} = useLazyAsyncData('jobs-stats-chart-jobs-by-category', async () => {
    const call = jobsStatsClient.getStats({
        start: toUtcDateTimestamp(query.range.start),
        end: toUtcDateTimestamp(query.range.end),

        period: selectedPeriod.value,
        category: query.category,
    });
    const { response } = await call;

    return response;
});

watch(
    () => [query.range.start.getTime(), query.range.end.getTime()],
    async () => await refresh(),
);
</script>

<template>
    <UDashboardPanel :ui="{ root: 'min-h-0' }">
        <template #header>
            <UDashboardToolbar>
                <template #left>
                    <DateRangePicker v-model="query.range" />
                </template>

                <template #right>
                    <RefreshButton @click="() => refresh()" />
                    <UButton
                        to="/jobs/colleagues"
                        icon="i-mdi-arrow-left"
                        variant="subtle"
                        :label="$t('pages.jobs.colleagues.title')"
                    />
                </template>
            </UDashboardToolbar>
        </template>

        <template #body>
            <DataPendingBlock v-if="isRequestPending(status)" :message="$t('common.loading', [$t('common.stats')])" />
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.unable_to_load', [$t('common.stats')])"
                :error="error"
                :retry="refresh"
            />
            <DataNoDataBlock
                v-else-if="!response"
                icon="i-mdi-file-search"
                :message="$t('common.not_found', [$t('common.stats')])"
            />

            <template v-else>
                <ChartClient :stats="response" :is-penalties="false" :period="selectedPeriod" :range="query.range" />
            </template>
        </template>
    </UDashboardPanel>
</template>
