<script lang="ts" setup>
import { snakeCase } from 'scule';
import type { DBSyncStatus } from '~~/gen/ts/resources/settings/status';
import DataNoDataBlock from '../partials/data/DataNoDataBlock.vue';
import GenericTime from '../partials/elements/GenericTime.vue';

const props = defineProps<{
    dbsync?: DBSyncStatus | null;
}>();

const isOpen = ref<boolean>(false);

const { t } = useI18n();

const dbSyncTables = computed(() =>
    [...(props.dbsync?.tables ?? [])]
        .map((table) => ({
            ...table,
            label: t(`components.settings.system_status.db_sync.tables.${snakeCase(table.table)}`),
        }))
        .sort((a, b) => a.label.localeCompare(b.label)),
);
</script>

<template>
    <UDrawer
        v-model:open="isOpen"
        :title="$t('components.settings.system_status.db_sync.drawer_title')"
        handle-only
        :ui="{ body: 'p-0 sm:mx-auto sm:max-w-7xl sm:w-full', title: 'flex flex-row gap-2' }"
    >
        <UButton
            variant="link"
            size="xl"
            icon="i-mdi-database-sync"
            :label="$t('components.settings.system_status.db_sync.title')"
            block
            :ui="{ leadingIcon: 'size-10' }"
        />

        <template #title>
            <span class="flex-1">{{ $t('components.settings.system_status.db_sync.drawer_title') }}</span>

            <UButton icon="i-mdi-close" color="neutral" variant="ghost" size="md" @click="isOpen = false" />
        </template>

        <template #body>
            <div class="space-y-6 p-4">
                <UCard variant="subtle" :ui="{ body: 'p-4' }">
                    <template #title>
                        <div class="flex items-center gap-2">
                            <UIcon name="i-mdi-database-sync" class="size-5 text-primary" />
                            <span>{{ $t('components.settings.system_status.db_sync.summary') }}</span>
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
                            <span v-if="props.dbsync?.lastDbsyncVersion" class="text-sm">
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

                    <div v-else class="grid gap-4 sm:grid-cols-2 xl:grid-cols-3">
                        <UCard v-for="table in dbSyncTables" :key="table.table" variant="subtle" :ui="{ body: 'space-y-3' }">
                            <template #title>
                                <div class="flex items-center gap-2">
                                    <UIcon name="i-mdi-table" class="size-5 text-primary" />
                                    <span>{{ table.label }}</span>
                                </div>
                            </template>

                            <dl class="space-y-2 text-sm">
                                <div class="grid grid-cols-[1fr_auto] gap-3">
                                    <dt class="font-medium text-muted">
                                        {{ $t('components.settings.system_status.db_sync.last_check') }}
                                    </dt>
                                    <dd class="text-right">
                                        <GenericTime v-if="table.lastCheck" :value="table.lastCheck" />
                                        <span v-else class="text-muted">{{ $t('common.na') }}</span>
                                    </dd>
                                </div>

                                <div class="grid grid-cols-[1fr_auto] gap-3">
                                    <dt class="font-medium text-muted">
                                        {{ $t('components.settings.system_status.db_sync.last_id') }}
                                    </dt>
                                    <dd class="text-right">
                                        <code v-if="table.lastId">{{ table.lastId }}</code>
                                        <span v-else class="text-muted">{{ $t('common.na') }}</span>
                                    </dd>
                                </div>
                            </dl>
                        </UCard>
                    </div>
                </div>
            </div>
        </template>
    </UDrawer>
</template>
