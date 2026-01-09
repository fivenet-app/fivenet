<script lang="ts" setup>
import { getSettingsSystemClient } from '~~/gen/ts/clients';
import DataErrorBlock from '../partials/data/DataErrorBlock.vue';
import DataPendingBlock from '../partials/data/DataPendingBlock.vue';
import GenericTime from '../partials/elements/GenericTime.vue';

const settingsSystemClient = await getSettingsSystemClient();

const { data, error, status, refresh } = useLazyAsyncData('settings-system-status', () => getStatus());

async function getStatus() {
    try {
        const call = settingsSystemClient.getStatus({});
        const { response } = await call;

        return response.status;
    } catch (err) {
        console.error('Failed to fetch system status:', err);
        throw err;
    }
}
</script>

<template>
    <UCard :title="$t('components.settings.system_status.title')" icon="i-mdi-server" :ui="{ body: 'p-2 sm:p-2' }">
        <template #header>
            <div class="flex items-center gap-2">
                <UIcon name="i-mdi-server" class="size-5 text-primary" />
                <h3 class="text-md flex-1 font-semibold">{{ $t('components.settings.system_status.title') }}</h3>

                <UTooltip
                    v-if="data?.version?.newVersion?.version"
                    :text="$d(toDate(data?.version?.newVersion.releaseDate), 'long')"
                    placement="bottom"
                >
                    <UButton
                        :label="`${$t('components.settings.system_status.new_version_available')} ${data?.version?.newVersion?.version}`"
                        icon="i-mdi-cellphone-system-update"
                        :to="data?.version?.newVersion?.url"
                        external
                        variant="subtle"
                    />
                </UTooltip>
            </div>
        </template>

        <template #default>
            <DataPendingBlock v-if="isRequestPending(status)" :message="$t('common.loading', [$t('common.status')])" />
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.not_found', [$t('common.status')])"
                :error="error"
                :retry="refresh"
            />

            <div v-else-if="data" class="flex flex-wrap gap-2">
                <UPopover class="flex-1">
                    <UButton
                        variant="link"
                        size="xl"
                        :color="data.database?.connected ? 'success' : 'error'"
                        icon="i-simple-icons-mysql"
                        :label="$t('components.settings.system_status.database.title')"
                        block
                        :ui="{ leadingIcon: 'size-10' }"
                    />

                    <template #content>
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
                                    <span>{{ data.database?.migrationDirty ? $t('common.yes') : $t('common.no') }}</span>
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
                                    <strong>{{ $t('components.settings.system_status.database.tables_mismatch') }}:</strong>
                                    <span>{{ !data.database?.tablesOk ? $t('common.yes') : $t('common.no') }}</span>
                                </li>
                            </ul>
                        </div>
                    </template>
                </UPopover>

                <UPopover class="flex-1">
                    <UButton
                        variant="link"
                        size="xl"
                        :color="data.nats?.connected ? 'success' : 'error'"
                        icon="i-simple-icons-natsdotio"
                        :label="$t('components.settings.system_status.nats.title')"
                        block
                        :ui="{ leadingIcon: 'size-10' }"
                    />

                    <template #content>
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
                    <UButton
                        variant="link"
                        size="xl"
                        icon="i-mdi-database-sync"
                        :label="$t('components.settings.system_status.db_sync.title')"
                        block
                        :ui="{ leadingIcon: 'size-10' }"
                    />

                    <template #content>
                        <div class="p-4">
                            <ul class="flex flex-col gap-1">
                                <li class="inline-flex items-center gap-1">
                                    <strong>{{ $t('components.settings.system_status.db_sync.last_data_received') }}:</strong>
                                    <GenericTime v-if="data.dbsync?.lastSyncedData" :value="data.dbsync?.lastSyncedData" />
                                    <span v-else>{{ $t('common.na') }}</span>
                                </li>

                                <li class="inline-flex items-center gap-1">
                                    <strong
                                        >{{ $t('components.settings.system_status.db_sync.last_activity_received') }}:</strong
                                    >
                                    <GenericTime
                                        v-if="data.dbsync?.lastSyncedActivity"
                                        :value="data.dbsync?.lastSyncedActivity"
                                    />
                                    <span v-else>{{ $t('common.na') }}</span>
                                </li>

                                <li class="inline-flex items-center gap-1">
                                    <strong>{{ $t('components.settings.system_status.db_sync.last_dbsync_version') }}:</strong>
                                    <span v-if="data.dbsync?.lastDbsyncVersion">{{
                                        data.dbsync?.lastDbsyncVersion ?? ''
                                    }}</span>
                                    <span v-else>{{ $t('common.na') }}</span>
                                </li>
                            </ul>
                        </div>
                    </template>
                </UPopover>
            </div>
        </template>
    </UCard>
</template>
