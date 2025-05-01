<script lang="ts" setup>
import { CronjobState } from '~~/gen/ts/resources/common/cron/cron';
import type { ListCronjobsResponse } from '~~/gen/ts/services/rector/cron';
import DataErrorBlock from '../partials/data/DataErrorBlock.vue';
import GenericTime from '../partials/elements/GenericTime.vue';
import Pagination from '../partials/Pagination.vue';

const { $grpc } = useNuxtApp();

const { data: cronjobs, pending: loading, refresh, error } = useLazyAsyncData(`rector-cronjobs`, () => listCronjobs());

async function listCronjobs(): Promise<ListCronjobsResponse> {
    try {
        const { response } = $grpc.rector.rectorCron.listCronjobs({});

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const { t } = useI18n();

const columns = [
    {
        key: 'name',
        label: t('common.name'),
    },
    {
        key: 'schedule',
        label: t('common.schedule'),
    },
    {
        key: 'status',
        label: t('common.status'),
    },
    {
        key: 'state',
        label: t('common.state'),
    },
    {
        key: 'nextScheduleTime',
        label: t('common.next_schedule_time'),
    },
    {
        key: 'lastAttemptTime',
        label: t('common.last_attempt_time'),
    },
    {
        key: 'startedTime',
        label: t('common.started_time'),
    },
];

const expand = ref({
    openedRows: [],
    row: {},
});
</script>

<template>
    <UDashboardNavbar :title="$t('pages.rector.cron.title')">
        <template #right>
            <PartialsBackButton fallback-to="/rector" />
        </template>
    </UDashboardNavbar>

    <DataErrorBlock v-if="error" :title="$t('common.unable_to_load', [$t('common.file', 2)])" :error="error" :retry="refresh" />

    <UTable
        v-else
        v-model:expand="expand"
        :loading="loading"
        :columns="columns"
        :rows="cronjobs?.jobs"
        :empty-state="{ icon: 'i-mdi-calendar-task', label: $t('common.not_found', [$t('common.file', 2)]) }"
        class="flex-1"
    >
        <template #expand="{ row }">
            <div class="p-2">
                <pre v-if="!row.lastCompletedEvent">{{ $t('common.na') }}</pre>
                <div v-else class="flex items-center gap-2">
                    <UBadge v-if="row.lastCompletedEvent.success" icon="i-mdi-check-bold" color="success" />
                    <UBadge v-else icon="i-mdi-exclamation-thick" color="error" />

                    <div class="font-semibold">
                        {{ $t('common.end_date') }}: <GenericTime :value="row.lastCompletedEvent.endDate" /> ({{
                            $t('common.duration')
                        }}: {{ fromDuration(row.lastCompletedEvent.elapsed) }}s)
                    </div>
                </div>
            </div>
        </template>

        <template #name-data="{ row }">
            <span class="text-gray-900 dark:text-white">
                <pre>{{ row.name }}</pre>
            </span>
        </template>

        <template #schedule-data="{ row }">
            <UKbd size="md">{{ row.schedule }}</UKbd>
        </template>

        <template #status-data="{ row }">
            <UBadge v-if="row.lastCompletedEvent?.success" icon="i-mdi-check-bold" color="success" />
            <UBadge v-else icon="i-mdi-exclamation-thick" color="error" />
        </template>

        <template #state-data="{ row }">
            {{ $t(`enums.rector.CronjobState.${CronjobState[row.state]}`) }}
        </template>

        <template #nextScheduleTime-data="{ row }">
            <GenericTime v-if="row.nextScheduleTime" :value="row.nextScheduleTime" />
        </template>

        <template #lastAttemptTime-data="{ row }">
            <GenericTime v-if="row.lastAttemptTime" :value="row.lastAttemptTime" />
        </template>

        <template #startedTime-data="{ row }">
            <GenericTime v-if="row.startedTime" :value="row.startedTime" />
        </template>
    </UTable>

    <Pagination :loading="loading" :refresh="refresh" />
</template>
