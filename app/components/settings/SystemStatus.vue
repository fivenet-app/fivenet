<script lang="ts" setup>
import DataErrorBlock from '../partials/data/DataErrorBlock.vue';
import DataPendingBlock from '../partials/data/DataPendingBlock.vue';
import GenericTime from '../partials/elements/GenericTime.vue';

const { $grpc } = useNuxtApp();

const { data, error, status, refresh } = useLazyAsyncData('settings-system-status', () => getStatus());

async function getStatus() {
    try {
        const call = $grpc.settings.system.getStatus({});
        const { response } = await call;

        return response;
    } catch (err) {
        console.error('Failed to fetch system status:', err);
        throw err;
    }
}
</script>

<template>
    <UCard>
        <template #header>
            <h2 class="text-lg font-medium">{{ $t('components.settings.system_status.title') }}</h2>
        </template>

        <DataPendingBlock v-if="isRequestPending(status)" :message="$t('common.loading', [$t('common.status')])" />
        <DataErrorBlock
            v-else-if="error"
            :title="$t('common.not_found', [$t('common.status')])"
            :error="error"
            :retry="refresh"
        />
        <div v-if="data" class="flex flex-wrap gap-4">
            <UPopover class="flex-1">
                <UButton variant="link" size="xl" :padded="false" :color="data.database?.connected ? 'success' : 'error'">
                    <UIcon name="i-simple-icons-mysql" class="size-10" />

                    {{ $t('components.settings.system_status.database.title') }}
                </UButton>

                <template #panel>
                    <div class="p-4">
                        <ul class="flex flex-col gap-1">
                            <li class="inline-flex items-center gap-1">
                                <strong>{{ $t('common.version') }}:</strong> <code>{{ data.database?.version }}</code>
                            </li>
                            <li class="inline-flex items-center gap-1">
                                <strong>{{ $t('components.settings.system_status.database.migration_version') }}:</strong>
                                <code>{{ data.database?.migrationVersion }}</code>
                            </li>
                            <li class="inline-flex items-center gap-1">
                                <strong>{{ $t('components.settings.system_status.database.migration_dirty') }}:</strong>
                                {{ data.database?.migrationDirty ? $t('common.yes') : $t('common.no') }}
                            </li>
                            <li class="inline-flex items-center gap-1">
                                <strong>{{ $t('components.settings.system_status.database.db_charset') }}:</strong>
                                <code>{{ data.database?.dbCharset }}</code>
                            </li>
                            <li class="inline-flex items-center gap-1">
                                <strong>{{ $t('components.settings.system_status.database.db_collation') }}:</strong>
                                <code>{{ data.database?.dbCollation }}</code>
                            </li>
                            <li class="inline-flex items-center gap-1">
                                <strong>{{ $t('components.settings.system_status.database.tables_ok') }}:</strong>
                                <code>{{ data.database?.tablesOk ? $t('common.yes') : $t('common.no') }}</code>
                            </li>
                        </ul>
                    </div>
                </template>
            </UPopover>

            <UPopover class="flex-1">
                <UButton variant="link" size="xl" :padded="false" :color="data.nats?.connected ? 'success' : 'error'">
                    <UIcon name="i-simple-icons-natsdotio" class="size-10" />

                    {{ $t('components.settings.system_status.nats.title') }}
                </UButton>

                <template #panel>
                    <div class="p-4">
                        <ul class="flex flex-col gap-1">
                            <li class="inline-flex items-center gap-1">
                                <strong>{{ $t('common.version') }}:</strong> <code>{{ data.nats?.version }}</code>
                            </li>
                        </ul>
                    </div>
                </template>
            </UPopover>

            <UPopover v-if="data.dbsync?.enabled" class="flex-1">
                <UButton variant="link" size="xl" :padded="false">
                    <UIcon name="i-mdi-database-sync" class="size-10" />

                    {{ $t('components.settings.system_status.db_sync.title') }}
                </UButton>

                <template #panel>
                    <div class="p-4">
                        <ul class="flex flex-col gap-1">
                            <li class="inline-flex items-center gap-1">
                                <strong>{{ $t('components.settings.system_status.db_sync.last_data_received') }}:</strong>
                                <GenericTime v-if="data.dbsync?.lastSyncedData" :value="data.dbsync?.lastSyncedData" />
                                <span v-else>{{ $t('common.na') }}</span>
                            </li>
                            <li class="inline-flex items-center gap-1">
                                <strong>{{ $t('components.settings.system_status.db_sync.last_activity_received') }}:</strong>
                                <GenericTime v-if="data.dbsync?.lastSyncedActivity" :value="data.dbsync?.lastSyncedActivity" />
                                <span v-else>{{ $t('common.na') }}</span>
                            </li>
                            <li class="inline-flex items-center gap-1">
                                <strong>{{ $t('components.settings.system_status.db_sync.last_dbsync_version') }}:</strong>
                                <span v-if="data.dbsync?.lastDbsyncVersion">{{ data.dbsync?.lastDbsyncVersion ?? '' }}</span>
                                <span v-else>{{ $t('common.na') }}</span>
                            </li>
                        </ul>
                    </div>
                </template>
            </UPopover>
        </div>
    </UCard>
</template>
