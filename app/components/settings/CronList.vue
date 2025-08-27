<script lang="ts" setup>
import { UBadge, UButton, UKbd } from '#components';
import type { TableColumn } from '@nuxt/ui';
import { h } from 'vue';
import { getSettingsCronClient } from '~~/gen/ts/clients';
import { Any } from '~~/gen/ts/google/protobuf/any';
import { type Cronjob, CronjobState, GenericCronData } from '~~/gen/ts/resources/common/cron/cron';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import type { ListCronjobsResponse } from '~~/gen/ts/services/settings/cron';
import DataErrorBlock from '../partials/data/DataErrorBlock.vue';
import GenericTime from '../partials/elements/GenericTime.vue';
import Pagination from '../partials/Pagination.vue';

const settingsCronClient = await getSettingsCronClient();

const { data: cronjobs, status, refresh, error } = useLazyAsyncData(`settings-cronjobs`, () => listCronjobs());

async function listCronjobs(): Promise<ListCronjobsResponse> {
    try {
        const { response } = settingsCronClient.listCronjobs({});

        start();

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

const { t } = useI18n();

const notifications = useNotificationsStore();

function copyLinkToClipboard(text: string): void {
    copyToClipboardWrapper(text);

    notifications.add({
        title: { key: 'notifications.clipboard.link_copied.title', parameters: {} },
        description: { key: 'notifications.clipboard.link_copied.content', parameters: {} },
        duration: 3250,
        type: NotificationType.INFO,
    });
}

const uiState = useUIStateStore();
const { windowFocus } = storeToRefs(uiState);

// Auto refresh the list every minute (if window is active)
const { remaining, start, pause, resume } = useCountdown(60, {
    onComplete() {
        refresh();
    },
});

onBeforeMount(() => start());

watchDebounced(windowFocus, () => {
    if (!windowFocus) {
        pause();
    } else {
        resume();
    }
});

const columns = computed(
    () =>
        [
            {
                accessorKey: 'name',
                header: t('common.name'),
                cell: ({ row }) => h('span', { class: 'text-highlighted' }, row.original.name),
            },
            {
                accessorKey: 'schedule',
                header: t('common.schedule'),
                cell: ({ row }) => h(UKbd, { size: 'md' }, row.original.schedule),
            },
            {
                accessorKey: 'status',
                header: t('common.status'),
                cell: ({ row }) =>
                    row.original.lastCompletedEvent?.success
                        ? h(UBadge, { icon: 'i-mdi-check-bold', color: 'success' })
                        : h(UBadge, { icon: 'i-mdi-exclamation-thick', color: 'error' }),
            },
            {
                accessorKey: 'state',
                header: t('common.state'),
                cell: ({ row }) => t(`enums.settings.CronjobState.${CronjobState[row.original.state]}`),
            },
            {
                accessorKey: 'nextScheduleTime',
                header: t('common.next_schedule_time'),
                cell: ({ row }) =>
                    row.original.nextScheduleTime ? h(GenericTime, { value: row.original.nextScheduleTime }) : null,
            },
            {
                accessorKey: 'lastAttemptTime',
                header: t('common.last_attempt_time'),
                cell: ({ row }) =>
                    row.original.lastAttemptTime ? h(GenericTime, { value: row.original.lastAttemptTime }) : null,
            },
            {
                accessorKey: 'startedTime',
                header: t('common.started_time'),
                cell: ({ row }) => (row.original.startedTime ? h(GenericTime, { value: row.original.startedTime }) : null),
            },
        ] as TableColumn<Cronjob>[],
);
</script>

<template>
    <UDashboardNavbar :title="$t('pages.settings.cron.title')">
        <template #right>
            <PartialsBackButton fallback-to="/settings" />
        </template>
    </UDashboardNavbar>

    <DataErrorBlock
        v-if="error"
        :title="$t('common.unable_to_load', [$t('common.cronjob', 2)])"
        :error="error"
        :retry="refresh"
    />

    <UTable
        v-else
        class="flex-1"
        :loading="isRequestPending(status)"
        :columns="columns"
        :data="cronjobs?.jobs"
        :pagination-options="{ manualPagination: true }"
        :sorting-options="{ manualSorting: true }"
        :empty="$t('common.not_found', [$t('pages.settings.cron.title', 2)])"
        :ui="{ tr: 'data-[expanded=true]:bg-elevated/50' }"
        sticky
    >
        <template #expanded="{ row }">
            <div class="p-2">
                <pre v-if="!row.original.lastCompletedEvent">{{ $t('common.unknown') }}</pre>
                <UCard v-else>
                    <template #header>
                        <div class="flex items-center justify-between gap-2">
                            <div class="flex gap-2">
                                <UBadge
                                    v-if="row.original.lastCompletedEvent.success"
                                    icon="i-mdi-check-bold"
                                    color="success"
                                />
                                <UBadge v-else icon="i-mdi-exclamation-thick" color="error" />

                                <div class="font-semibold">
                                    {{ $t('common.end_date') }}:
                                    <GenericTime :value="row.original.lastCompletedEvent.endDate" /> ({{
                                        $t('common.duration')
                                    }}: {{ fromDuration(row.original.lastCompletedEvent.elapsed) }}s)
                                </div>
                            </div>

                            <UButton
                                v-if="row.original.lastCompletedEvent.data?.data"
                                variant="link"
                                icon="i-mdi-share"
                                @click="
                                    copyLinkToClipboard(
                                        row.original.lastCompletedEvent.data?.data?.typeUrl.includes(
                                            '/resources.common.cron.GenericCronData',
                                        )
                                            ? Any.unpack(row.original.lastCompletedEvent.data.data, GenericCronData)
                                            : row.original.lastCompletedEvent.data,
                                    )
                                "
                            />
                        </div>
                    </template>

                    <pre
                        class="line-clamp-9 hover:line-clamp-none"
                        v-text="
                            row.original.lastCompletedEvent.data?.data?.typeUrl.includes(
                                '/resources.common.cron.GenericCronData',
                            )
                                ? Any.unpack(row.original.lastCompletedEvent.data.data, GenericCronData)
                                : row.original.lastCompletedEvent.data
                        "
                    />

                    <template #footer>
                        <div>
                            <span class="font-semibold">{{ $t('pages.error.error_message') }}</span>

                            <pre
                                class="line-clamp-4 whitespace-break-spaces hover:line-clamp-none"
                                v-text="
                                    row.original.lastCompletedEvent.errorMessage
                                        ? row.original.lastCompletedEvent.errorMessage
                                        : $t('common.none')
                                "
                            />
                        </div>
                    </template>
                </UCard>
            </div>
        </template>
    </UTable>

    <Pagination :status="status" :refresh="refresh" hide-text hide-buttons>
        <p>
            {{ $t('common.refresh_in_x', { d: remaining, unit: $t('common.time_ago.second', remaining) }) }}
        </p>
    </Pagination>
</template>
