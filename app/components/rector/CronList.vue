<script lang="ts" setup>
import { CronjobState } from '~~/gen/ts/resources/common/cron/cron';
import type { ListCronjobsResponse } from '~~/gen/ts/services/rector/cron';
import DataErrorBlock from '../partials/data/DataErrorBlock.vue';
import GenericTime from '../partials/elements/GenericTime.vue';

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
        :loading="loading"
        :columns="columns"
        :rows="cronjobs?.jobs"
        :empty-state="{ icon: 'i-mdi-calendar-task', label: $t('common.not_found', [$t('common.file', 2)]) }"
        class="flex-1"
    >
        <template #name-data="{ row }">
            <span class="text-gray-900 dark:text-white">
                {{ row.name }}
            </span>
        </template>

        <template #schedule-data="{ row }">
            <UKbd size="md">{{ row.schedule }}</UKbd>
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
</template>
