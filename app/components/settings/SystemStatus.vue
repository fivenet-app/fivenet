<script lang="ts" setup>
import { getSettingsSystemClient } from '~~/gen/ts/clients';
import DataErrorBlock from '../partials/data/DataErrorBlock.vue';
import SystemStatusDBSyncDrawer from '~/components/settings/SystemStatusDBSyncDrawer.vue';
import RefreshButton from '~/components/partials/RefreshButton.vue';
import { isRequestPending } from '~/utils/data';

const settingsSystemClient = await getSettingsSystemClient();

const { data, error, status, refresh } = useLazyAsyncData('settings-system-status', () => getStatus());
const isStatusLoading = computed(() => isRequestPending(status.value));

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

const version = APP_VERSION;
async function copyVersionToClipboard() {
    copyToClipboardWrapper(`${$t('common.version')}: ${version}`);
}
</script>

<template>
    <UCard
        :title="$t('components.settings.system_status.title')"
        icon="i-mdi-server"
        variant="subtle"
        :ui="{ body: 'p-2 sm:p-2' }"
    >
        <template #header>
            <div class="flex items-center gap-2">
                <UIcon class="size-5 text-primary" name="i-mdi-server" />
                <h3 class="text-md flex-1 font-semibold">{{ $t('components.settings.system_status.title') }}</h3>

                <UTooltip :text="$t('common.copy')">
                    <UButton variant="soft" size="xs" @click="copyVersionToClipboard">
                        <span class="hidden truncate sm:block">{{ $t('common.version') }}:</span>
                        <code class="font-mono">{{ version }}</code>
                    </UButton>
                </UTooltip>

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

                <RefreshButton icon-only :loading="isStatusLoading" :disabled="isStatusLoading" @click="() => refresh()" />
            </div>
        </template>

        <template #default>
            <DataErrorBlock
                v-if="error"
                :title="$t('common.not_found', [$t('common.status')])"
                :error="error"
                :retry="refresh"
            />

            <div v-else-if="isStatusLoading && !data" class="flex flex-wrap gap-2">
                <USkeleton class="h-16 min-w-[16rem] flex-1" />
                <USkeleton class="h-16 min-w-[16rem] flex-1" />
                <USkeleton class="h-16 min-w-[16rem] flex-1" />
            </div>

            <div v-else class="flex flex-row justify-around gap-2">
                <UPopover>
                    <UChip :color="data?.database?.connected ? 'success' : 'error'">
                        <UButton
                            icon="i-simple-icons-mysql"
                            :label="$t('components.settings.system_status.database.title')"
                            size="xl"
                            variant="link"
                            :disabled="isStatusLoading || !data"
                            :ui="{ leadingIcon: 'size-8' }"
                        />
                    </UChip>

                    <template #content>
                        <div class="p-4">
                            <ul class="flex flex-col gap-1">
                                <li class="inline-flex items-center gap-1">
                                    <strong>{{ $t('common.version') }}:</strong>
                                    <code class="font-mono">{{ data?.database?.version }}</code>
                                </li>

                                <li class="inline-flex items-center gap-1">
                                    <strong>{{ $t('components.settings.system_status.database.migration_version') }}:</strong>
                                    <code class="font-mono">{{ data?.database?.migrationVersion }}</code>
                                </li>

                                <li class="inline-flex items-center gap-1">
                                    <strong>{{ $t('components.settings.system_status.database.migration_dirty') }}:</strong>
                                    <span>{{ data?.database?.migrationDirty ? $t('common.yes') : $t('common.no') }}</span>
                                </li>

                                <li class="inline-flex items-center gap-1">
                                    <strong>{{ $t('components.settings.system_status.database.db_charset') }}:</strong>
                                    <code class="font-mono">{{ data?.database?.dbCharset }}</code>
                                </li>

                                <li class="inline-flex items-center gap-1">
                                    <strong>{{ $t('components.settings.system_status.database.db_collation') }}:</strong>
                                    <code class="font-mono">{{ data?.database?.dbCollation }}</code>
                                </li>

                                <li class="inline-flex items-center gap-1">
                                    <strong>{{ $t('components.settings.system_status.database.tables_mismatch') }}:</strong>
                                    <span>{{ !data?.database?.tablesOk ? $t('common.yes') : $t('common.no') }}</span>
                                </li>
                            </ul>
                        </div>
                    </template>
                </UPopover>

                <UPopover>
                    <UChip :color="data?.nats?.connected ? 'success' : 'error'">
                        <UButton
                            variant="link"
                            size="xl"
                            icon="i-simple-icons-natsdotio"
                            :label="$t('components.settings.system_status.nats.title')"
                            :disabled="isStatusLoading || !data"
                            :ui="{ leadingIcon: 'size-8' }"
                        />
                    </UChip>

                    <template #content>
                        <div class="p-4">
                            <ul class="flex flex-col gap-1">
                                <li class="inline-flex items-center gap-1">
                                    <strong>{{ $t('common.version') }}:</strong>
                                    <code class="font-mono">{{ data?.nats?.version }}</code>
                                </li>

                                <li class="inline-flex items-center gap-1">
                                    <strong>{{ $t('components.settings.system_status.database.migration_version') }}:</strong>
                                    <code class="font-mono">{{ data?.nats?.migrationVersion }}</code>
                                </li>
                            </ul>
                        </div>
                    </template>
                </UPopover>

                <SystemStatusDBSyncDrawer
                    v-if="data?.dbsync?.enabled"
                    :dbsync="data.dbsync"
                    :disabled="isStatusLoading"
                    @refresh="() => refresh()"
                />
            </div>
        </template>
    </UCard>
</template>
