<script lang="ts" setup>
import { snakeCase } from 'scule';
import type { DBSyncStatus } from '~~/gen/ts/resources/settings/status';
import DataNoDataBlock from '../partials/data/DataNoDataBlock.vue';
import GenericTime from '../partials/elements/GenericTime.vue';
import RefreshButton from '~/components/partials/RefreshButton.vue';
import type { DBSyncTableSyncState } from '~~/gen/ts/resources/dbsync/state';

const props = defineProps<{
    dbsync?: DBSyncStatus | null;
    disabled?: boolean;
}>();

const emits = defineEmits<{
    (e: 'refresh'): void;
}>();

const isOpen = ref<boolean>(false);

const { t } = useI18n();

const uiState = useUIStateStore();
const { windowFocus } = storeToRefs(uiState);

type DBSyncCardState = DBSyncTableSyncState & { label: string };

const dbSyncTables = computed<DBSyncCardState[]>(() =>
    [...(props.dbsync?.syncState?.tables ?? [])]
        .map((table) => ({
            ...table,
            label: t(`components.settings.system_status.db_sync.tables.${snakeCase(table.table)}`),
        }))
        .sort((a, b) => a.label.localeCompare(b.label))
        .sort((a, b) => Number(b.enabled) - Number(a.enabled)),
);

const dbSyncStreamConnected = computed(() => props.dbsync?.streamConnected ?? false);

const dbSyncStreamConnectionState = computed(() => (dbSyncStreamConnected.value ? 'connected' : 'disconnected'));

const dbSyncStreamConnectionColor = computed(() => (dbSyncStreamConnected.value ? 'success' : 'warning'));

function getTableState(table: DBSyncCardState) {
    if (table.lastError) {
        return 'error';
    }

    if (table.lastSyncedAt || table.lastAttemptAt || table.checkpoint) {
        return 'healthy';
    }

    return 'idle';
}

function getTableStateLabel(table: DBSyncCardState) {
    return t(`components.settings.system_status.db_sync.state.${getTableState(table)}`);
}

function getTableStateColor(table: DBSyncCardState) {
    switch (getTableState(table)) {
        case 'error':
            return 'error';
        case 'healthy':
            return 'success';
        default:
            return 'neutral';
    }
}

// Auto refresh the list every minute (if window is active)
const { remaining, start, pause, resume } = useCountdown(60, {
    onComplete: refresh,
});

onBeforeMount(() => start());

watchDebounced(windowFocus, () => {
    if (!windowFocus.value) {
        pause();
    } else {
        resume();
    }
});

watch(isOpen, () => {
    if (isOpen.value) {
        start();
    } else {
        pause();
    }
});

function refresh() {
    emits('refresh');
    start();
}
</script>

<template>
    <UDrawer
        v-model:open="isOpen"
        :title="$t('components.settings.system_status.db_sync.drawer_title')"
        handle-only
        :ui="{ body: 'p-0 sm:mx-auto sm:max-w-7xl sm:w-full', title: 'flex flex-row gap-2' }"
    >
        <UChip :show="true" :color="dbSyncStreamConnectionColor" position="top-right" size="xl">
            <UButton
                variant="link"
                size="xl"
                icon="i-mdi-database-sync"
                :label="$t('components.settings.system_status.db_sync.title')"
                block
                :disabled="props.disabled"
                :ui="{ leadingIcon: 'size-8' }"
            />
        </UChip>

        <template #title>
            <p class="flex-1">
                {{ $t('components.settings.system_status.db_sync.drawer_title') }}
                <span class="text-sm text-secondary"
                    >{{ $t('common.refresh_in_x', { d: remaining, unit: $t('common.time_ago.second', remaining) }) }}
                </span>
            </p>

            <UTooltip :text="$t('common.close', 1)">
                <UButton icon="i-mdi-close" color="neutral" variant="ghost" size="md" @click="isOpen = false" />
            </UTooltip>
        </template>

        <template #actions>
            <RefreshButton icon-only :disabled="props.disabled" @click="refresh" />
        </template>

        <template #body>
            <div class="space-y-6 p-4">
                <UCard variant="subtle" :ui="{ body: 'p-4' }">
                    <template #title>
                        <div class="flex flex-1 items-center justify-between gap-2">
                            <div class="flex items-center gap-2">
                                <UIcon name="i-mdi-database-sync" class="size-5 text-primary" />
                                <span>{{ $t('components.settings.system_status.db_sync.summary') }}</span>
                            </div>

                            <UBadge
                                :color="dbSyncStreamConnectionColor"
                                variant="soft"
                                :label="
                                    $t(
                                        `components.settings.system_status.db_sync.stream_connected.${dbSyncStreamConnectionState}`,
                                    )
                                "
                            />
                        </div>
                    </template>

                    <div class="grid gap-4 md:grid-cols-3">
                        <div class="space-y-1">
                            <p class="text-xs font-medium tracking-wide text-muted uppercase">
                                {{ $t('components.settings.system_status.db_sync.last_data_received') }}
                            </p>
                            <GenericTime v-if="props.dbsync?.lastSyncedData" :value="props.dbsync.lastSyncedData" />
                            <span v-else class="text-sm text-muted">{{ $t('common.na') }}</span>
                        </div>

                        <div class="space-y-1">
                            <p class="text-xs font-medium tracking-wide text-muted uppercase">
                                {{ $t('components.settings.system_status.db_sync.last_activity_received') }}
                            </p>
                            <GenericTime v-if="props.dbsync?.lastSyncedActivity" :value="props.dbsync.lastSyncedActivity" />
                            <span v-else class="text-sm text-muted">{{ $t('common.na') }}</span>
                        </div>

                        <div class="space-y-1">
                            <p class="text-xs font-medium tracking-wide text-muted uppercase">
                                {{ $t('components.settings.system_status.db_sync.last_dbsync_version') }}
                            </p>
                            <span v-if="props.dbsync?.lastDbsyncVersion" class="font-mono text-sm">
                                {{ props.dbsync.lastDbsyncVersion }}
                            </span>
                            <span v-else class="text-sm text-muted">{{ $t('common.na') }}</span>
                        </div>
                    </div>
                </UCard>

                <div class="space-y-3">
                    <div class="flex items-center justify-between gap-3">
                        <h4 class="text-sm font-semibold tracking-wide text-muted uppercase">
                            {{ $t('components.settings.system_status.db_sync.per_table') }}
                        </h4>

                        <UBadge color="neutral" variant="soft" :label="dbSyncTables.length" />
                    </div>

                    <DataNoDataBlock
                        v-if="dbSyncTables.length === 0"
                        :message="$t('components.settings.system_status.db_sync.no_tables')"
                        :title="$t('components.settings.system_status.db_sync.per_table')"
                        icon="i-mdi-table"
                        :padded="false"
                    />

                    <UPageGrid v-else class="gap-4">
                        <UCard
                            v-for="table in dbSyncTables"
                            :key="table.table"
                            :variant="table.enabled ? 'subtle' : 'outline'"
                            :ui="{ body: 'space-y-3' }"
                        >
                            <template #title>
                                <div class="flex items-center justify-between gap-2">
                                    <div class="flex items-center gap-2">
                                        <UIcon name="i-mdi-table" class="size-5 text-primary" />
                                        <span>{{ table.label }}</span>
                                    </div>

                                    <UBadge
                                        :color="getTableStateColor(table)"
                                        variant="soft"
                                        :label="getTableStateLabel(table)"
                                    />
                                </div>
                            </template>

                            <dl class="space-y-3 text-sm">
                                <div class="grid grid-cols-[1fr_auto] gap-3">
                                    <dt class="font-medium text-muted">
                                        {{ $t('components.settings.system_status.db_sync.last_synced_at') }}
                                    </dt>
                                    <dd class="text-right">
                                        <GenericTime v-if="table.lastSyncedAt" :value="table.lastSyncedAt" />
                                        <span v-else class="text-muted">{{ $t('common.na') }}</span>
                                    </dd>
                                </div>

                                <div class="grid grid-cols-[1fr_auto] gap-3">
                                    <dt class="font-medium text-muted">
                                        {{ $t('components.settings.system_status.db_sync.last_attempt_at') }}
                                    </dt>
                                    <dd class="text-right">
                                        <GenericTime v-if="table.lastAttemptAt" :value="table.lastAttemptAt" />
                                        <span v-else class="text-muted">{{ $t('common.na') }}</span>
                                    </dd>
                                </div>

                                <div class="grid grid-cols-[1fr_auto] gap-3">
                                    <dt class="font-medium text-muted">
                                        {{ $t('components.settings.system_status.db_sync.checkpoint') }}
                                    </dt>
                                    <dd class="text-right">
                                        <div v-if="table.checkpoint" class="space-x-2">
                                            <GenericTime
                                                v-if="table.checkpoint.lastCheck"
                                                :value="table.checkpoint.lastCheck"
                                            />
                                            <code v-if="table.checkpoint.lastId">{{ table.checkpoint.lastId }}</code>
                                            <span
                                                v-if="!table.checkpoint.lastCheck && !table.checkpoint.lastId"
                                                class="text-muted"
                                            >
                                                {{ $t('common.na') }}
                                            </span>
                                        </div>
                                        <span v-else class="text-muted">{{ $t('common.na') }}</span>
                                    </dd>
                                </div>

                                <div v-if="table.lastError" class="grid grid-cols-[1fr_auto] gap-3">
                                    <dt class="font-medium text-muted">
                                        {{ $t('components.settings.system_status.db_sync.last_error') }}
                                    </dt>
                                    <dd class="max-w-56 text-right break-words text-error">
                                        {{ table.lastError }}
                                    </dd>
                                </div>
                            </dl>
                        </UCard>
                    </UPageGrid>
                </div>
            </div>
        </template>
    </UDrawer>
</template>
