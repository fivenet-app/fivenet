<script lang="ts" setup>
import { UButton } from '#components';
import type { TableColumn } from '@nuxt/ui';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import Pagination from '~/components/partials/Pagination.vue';
import type { Job } from '~~/gen/ts/resources/jobs/jobs';

const { t } = useI18n();

const appConfig = useAppConfig();

const settingsStore = useSettingsStore();
const { streamerMode } = storeToRefs(settingsStore);

const completorStore = useCompletorStore();
const { completeJobs } = completorStore;

const { data: jobs, error, status, refresh } = useLazyAsyncData(`settings-appconfig-jobs`, () => completeJobs({}));

const columns = computed(
    () =>
        [
            {
                id: 'expand',
                cell: ({ row }) =>
                    h(UButton, {
                        color: 'neutral',
                        variant: 'ghost',
                        icon: 'i-mdi-chevron-down',
                        square: true,
                        'aria-label': 'Expand',
                        ui: {
                            leadingIcon: ['transition-transform', row.getIsExpanded() ? 'duration-200 rotate-180' : ''],
                        },
                        onClick: () => row.toggleExpanded(),
                    }),
            },
            {
                accessorKey: 'label',
                header: ({ column }) => {
                    const isSorted = column.getIsSorted();

                    return h(UButton, {
                        color: 'neutral',
                        variant: 'ghost',
                        label: t('common.label', 1),
                        icon: isSorted
                            ? isSorted === 'asc'
                                ? appConfig.custom.icons.sortAsc
                                : appConfig.custom.icons.sortDesc
                            : appConfig.custom.icons.sort,
                        class: '-mx-2.5',
                        onClick: () => column.toggleSorting(column.getIsSorted() === 'asc'),
                    });
                },
                meta: {
                    class: {
                        td: 'text-highlighted',
                    },
                },
            },
            {
                accessorKey: 'name',
                header: ({ column }) => {
                    const isSorted = column.getIsSorted();

                    return h(UButton, {
                        color: 'neutral',
                        variant: 'ghost',
                        label: t('common.name'),
                        icon: isSorted
                            ? isSorted === 'asc'
                                ? appConfig.custom.icons.sortAsc
                                : appConfig.custom.icons.sortDesc
                            : appConfig.custom.icons.sort,
                        class: '-mx-2.5',
                        onClick: () => column.toggleSorting(column.getIsSorted() === 'asc'),
                    });
                },
            },
            {
                accessorKey: 'grades',
                header: t('common.rank', 2),
                cell: ({ row }) => row.original.grades.length,
            },
        ] as TableColumn<Job>[],
);

const sorting = ref([
    {
        id: 'label',
        desc: false,
    },
]);

const expanded = ref({});
</script>

<template>
    <UDashboardPanel :ui="{ body: 'p-0 sm:p-0 gap-0 sm:gap-0' }">
        <template #header>
            <UDashboardNavbar :title="$t('pages.settings.jobs.title')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <PartialsBackButton fallback-to="/settings" />
                </template>
            </UDashboardNavbar>
        </template>

        <template #body>
            <StreamerModeAlert v-if="streamerMode" />
            <template v-else>
                <DataErrorBlock
                    v-if="error"
                    :title="$t('common.unable_to_load', [$t('common.job', 2)])"
                    :error="error"
                    :retry="refresh"
                />
                <DataPendingBlock v-else-if="isRequestPending(status)" :message="$t('common.loading', [$t('common.job', 2)])" />

                <UTable
                    v-else
                    v-model:expanded="expanded"
                    v-model:sorting="sorting"
                    class="flex-1"
                    :loading="isRequestPending(status)"
                    :columns="columns"
                    :data="jobs"
                    :empty="$t('common.not_found', [$t('common.job', 2)])"
                    :pagination-options="{ manualPagination: true }"
                    sticky
                    :ui="{ tr: 'data-[expanded=true]:bg-elevated/50' }"
                >
                    <template #expanded="{ row }">
                        <DataNoDataBlock v-if="row.original.grades.length === 0" />
                        <div v-else class="p-4">
                            <div class="mb-2 text-sm text-muted">
                                {{ $t('common.total_count') }}: {{ $t('common.ranks', row.original.grades.length) }}
                            </div>

                            <div class="flex flex-col gap-2">
                                <div
                                    v-for="grade in row.original.grades"
                                    :key="grade.grade"
                                    class="flex items-center gap-2 rounded bg-elevated px-4 py-2"
                                >
                                    <span>{{ grade.label }} ({{ grade.grade }})</span>
                                </div>
                            </div>
                        </div>
                    </template>
                </UTable>

                <Pagination :status="status" :refresh="refresh" hide-buttons hide-text />
            </template>
        </template>
    </UDashboardPanel>
</template>
