<script lang="ts" setup>
import { differenceInDays, sub } from 'date-fns';
import { z } from 'zod';
import ChartClient from '~/components/documents/stats/Chart.client.vue';
import DateRangePicker from '~/components/documents/stats/DateRangePicker.vue';
import PenaltySeriesChartClient from '~/components/documents/stats/PenaltySeriesChart.client.vue';
import Table from '~/components/documents/stats/Table.vue';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import RefreshButton from '~/components/partials/RefreshButton.vue';
import { getDocumentsStatsClient } from '~~/gen/ts/clients';
import { StatsCategory, StatsPeriod } from '~~/gen/ts/resources/stats/stats';

useHead({
    title: 'pages.documents.stats.title',
});

definePageMeta({
    title: 'pages.documents.stats.title',
    requiresAuth: true,
    permission: 'documents.StatsService/GetStats',
});

const { t } = useI18n();

const { activeChar, attrStringList, attrJobList, isSuperuser } = useAuth();

const completorStore = useCompletorStore();
const { jobs } = storeToRefs(completorStore);
const { listJobs } = completorStore;

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
    category: z.enum(StatsCategory).default(StatsCategory.DOCUMENTS_BY_CATEGORY),
    jobs: z.string().max(40).array().max(15).default([]),
});

const query = useSearchForm('documents-stats', schema);

const selectedPeriod = computed(() => {
    const days = Math.max(differenceInDays(query.range.end, query.range.start), 1);
    if (days > 365) {
        return StatsPeriod.MONTHLY;
    }

    return days > 60 ? StatsPeriod.WEEKLY : StatsPeriod.DAILY;
});

const documentsStatsClient = await getDocumentsStatsClient();

const {
    data: response,
    status,
    error,
    refresh,
} = useLazyAsyncData('documents-stats-chart-documents-by-category', async () => {
    const call = documentsStatsClient.getStats({
        start: toUtcDateTimestamp(query.range.start),
        end: toUtcDateTimestamp(query.range.end),
        period: selectedPeriod.value,
        category: query.category,
        jobs: query.jobs,
    });
    const { response } = await call;

    if (response.documentsByCategory.length > 0) {
        response.documentsByCategory.push({
            id: 0,
            job: '',
            name: t('common.categories', 0),
            color: '',
            value: response.totalValue - response.documentsByCategory.reduce((acc, { value }) => acc + value, 0),
        });
    }

    return response;
});

watch(
    () => query.jobs,
    () => {
        if (query.jobs.length > 0) return;
        if (!activeChar.value) return;

        query.jobs.push(activeChar.value.job);
    },
    { immediate: true },
);
watchDebounced(
    () => [query.category, query.range.start.getTime(), query.range.end.getTime(), query.jobs],
    async () => await refresh(),
    { debounce: 200, maxWait: 1250 },
);

const allowedJobs = computed(() => {
    // Superuser can see all jobs, how helpful that might be, is up to the superuser.. :-)
    if (isSuperuser.value) return jobs.value.map((j) => ({ name: j.name, label: j.label }));

    // Map the job names from the attribute to "actual" jobs
    return attrJobList('documents.StatsService/GetStats', 'Jobs').value.map((j) => {
        const job = jobs.value.find((js) => js.name === j);

        return { name: job?.name ?? j, label: job?.label ?? j };
    });
});

const canSeePenalties = computed(
    () =>
        isSuperuser.value ||
        attrStringList('documents.StatsService/GetStats', 'Categories').value.includes('PenaltyCalculator'),
);

const categories = computed(() =>
    [
        {
            label: t('enums.stats.StatsCategory.DOCUMENTS_BY_CATEGORY'),
            value: StatsCategory.DOCUMENTS_BY_CATEGORY,
            icon: 'i-mdi-shape',
        },
        {
            label: t('enums.stats.StatsCategory.TOP_LAWS'),
            value: StatsCategory.TOP_LAWS,
            icon: 'i-mdi-gavel',
            hidden: !canSeePenalties.value,
        },
        {
            label: t('enums.stats.StatsCategory.PENALTIES_OVER_TIME'),
            value: StatsCategory.PENALTIES_OVER_TIME,
            icon: 'i-mdi-gavel',
            hidden: !canSeePenalties.value,
        },
    ].flatMap((item) => (item.hidden ? [] : [item])),
);

const breadcrumbs = computed(() => [
    {
        label: t('pages.documents.title'),
        icon: 'i-mdi-file-document',
        to: '/documents',
    },
    {
        label: t('pages.documents.stats.title'),
        icon: 'i-mdi-graph-box-outline',
    },
]);

const isPenalties = computed(() => query.category === StatsCategory.PENALTIES_OVER_TIME);

onBeforeMount(async () => listJobs());
</script>

<template>
    <UDashboardPanel>
        <template #header>
            <UDashboardNavbar :title="$t('pages.documents.stats.title')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <PartialsBackButton fallback-to="/documents" />
                </template>
            </UDashboardNavbar>

            <UDashboardToolbar>
                <template #left>
                    <UBreadcrumb :items="breadcrumbs" />
                </template>
            </UDashboardToolbar>

            <UDashboardToolbar>
                <template #left>
                    <UFormField name="range">
                        <DateRangePicker v-model="query.range" />
                    </UFormField>

                    <UFormField name="category">
                        <USelectMenu v-model="query.category" :items="categories" value-key="value" />
                    </UFormField>
                </template>

                <template #right>
                    <RefreshButton @click="() => refresh()" />
                </template>
            </UDashboardToolbar>

            <UDashboardToolbar v-if="allowedJobs.length > 1">
                <template #default>
                    <UFormField class="w-full" name="jobs" :label="$t('common.job', 2)" :ui="{ root: 'py-2' }">
                        <USelectMenu
                            v-model="query.jobs"
                            class="w-full"
                            block
                            :items="allowedJobs"
                            label-key="label"
                            value-key="name"
                            multiple
                            :search-input="{ placeholder: $t('common.search_field') }"
                            :filter-fields="['name', 'label']"
                        >
                            <template v-if="query.jobs.length === 0" #default>
                                {{ $t('common.none_selected', [$t('common.job', 2)]) }}
                            </template>

                            <template #empty>
                                {{ $t('common.not_found', [$t('common.job', 2)]) }}
                            </template>
                        </USelectMenu>
                    </UFormField>
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
                <ChartClient
                    v-if="!isPenalties"
                    :stats="response"
                    :is-penalties="isPenalties"
                    :period="selectedPeriod"
                    :range="query.range"
                />
                <PenaltySeriesChartClient v-else :stats="response" :period="selectedPeriod" :range="query.range" />

                <Table
                    v-if="!isPenalties"
                    :category="query.category"
                    :stats="response"
                    :period="selectedPeriod"
                    :range="query.range"
                />
            </template>
        </template>
    </UDashboardPanel>
</template>
