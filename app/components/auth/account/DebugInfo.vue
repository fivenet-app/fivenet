<script lang="ts" setup>
import { listEnumValues } from '@protobuf-ts/runtime';
import { LogLevels } from 'consola';
import CopyToClipboardButton from '~/components/partials/CopyToClipboardButton.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import PermList from '~/components/settings/roles/PermList.vue';
import { useGRPCWebsocketTransport } from '~/composables/grpc/grpcws';
import { useAuthStore } from '~/stores/auth';
import { useClipboardStore } from '~/stores/clipboard';
import { useSettingsStore } from '~/stores/settings';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const clipboardStore = useClipboardStore();

const searchesStore = useSearchesStore();

const settingsStore = useSettingsStore();

const authStore = useAuthStore();
const { activeChar, sessionExpiration, attributes, permissions } = storeToRefs(authStore);
const { clearAuthInfo } = authStore;

const notifications = useNotificationsStore();

const { webSocket } = useGRPCWebsocketTransport();

async function resetLocalStorage(): Promise<void> {
    clearAuthInfo();

    if (import.meta.client) {
        window.localStorage.clear();
    }

    await navigateTo('/');
}

async function sendTestNotifications(): Promise<void> {
    listEnumValues(NotificationType)
        .filter((t) => t.number !== 0)
        .forEach((notificationType, index) => {
            notifications.add({
                title: { key: 'notifications.system.test_notification.title', parameters: { index: (index + 1).toString() } },
                description: {
                    key: 'notifications.system.test_notification.content',
                    parameters: { type: notificationType.name },
                },
                type: notificationType.number,
                actions: [
                    {
                        label: { key: 'common.click_here' },
                        onClick: () => alert('Test was successful!'),
                    },
                ],
            });
        });
}

function triggerBannerMessage(): void {
    const { system } = useAppConfig();
    system.bannerMessageEnabled = true;
    system.bannerMessage = {
        id: 'test-' + new Date().getTime().toString(),
        title: 'Test Banner: Insert cool message here',
    };
}

function triggerErrorPage(): void {
    showError(new Error('You pressed the trigger error page button'));
}

function toggleLogLevel(): void {
    setDefaultLogLevel(getDefaultLogLevel() !== LogLevels.debug ? LogLevels.debug : LogLevels.warn);
    console.warn('Setting log level to', getDefaultLogLevel() === 4 ? 'DEBUG' : 'WARN');
}

const version = APP_VERSION;
</script>

<template>
    <UPageCard :title="$t('components.debug_info.title')" :description="$t('components.debug_info.subtitle')">
        <UFormField class="grid grid-cols-2 items-center gap-2" name="version" :label="$t('components.debug_info.version')">
            <div class="inline-flex w-full justify-between">
                <span>
                    <code>{{ version }}</code> / <code>{{ settingsStore.version }}</code>
                </span>
                <CopyToClipboardButton :value="`${version}/ ${settingsStore.version}`" />
            </div>
        </UFormField>

        <UFormField
            v-if="activeChar"
            class="grid grid-cols-2 items-center gap-2"
            name="activeCharId"
            :label="$t('components.debug_info.active_char_id')"
        >
            <div class="inline-flex w-full justify-between">
                <span>
                    {{ activeChar.userId }}
                </span>
                <CopyToClipboardButton :value="activeChar.userId" />
            </div>
        </UFormField>

        <UFormField
            v-if="activeChar"
            class="grid grid-cols-2 items-center gap-2"
            name="activeCharJob"
            :label="$t('common.job')"
        >
            <div class="inline-flex w-full justify-between">
                <span>{{ activeChar.job }} ({{ $t('common.rank') }}: {{ activeChar.jobGrade }})</span>
                <CopyToClipboardButton :value="`${activeChar.job} (${$t('common.rank')}: ${activeChar.jobGrade})`" />
            </div>
        </UFormField>

        <UFormField
            v-if="sessionExpiration"
            class="grid grid-cols-2 items-center gap-2"
            name="sessionExpiration"
            :label="$t('components.debug_info.access_token_expiration')"
        >
            <GenericTime :value="sessionExpiration" ago />
            (<GenericTime :value="sessionExpiration" type="long" />)
        </UFormField>

        <UFormField class="grid grid-cols-2 items-center gap-2" name="status" :label="$t('common.status')">
            {{ $t('common.status') }}: <code>{{ webSocket.status.value }}</code>
        </UFormField>

        <UFormField class="grid grid-cols-2 items-center gap-2" name="nuiInfo" :label="$t('components.debug_info.nui_info')">
            {{ settingsStore.nuiEnabled ? $t('common.enabled') : $t('common.disabled') }}:
            {{ settingsStore.nuiResourceName ?? $t('common.na') }}
        </UFormField>

        <UFormField
            class="grid grid-cols-2 items-center gap-2"
            name="debugFunctions"
            :label="$t('components.debug_info.debug_functions')"
        >
            <UFieldGroup class="flex w-full break-words" orientation="vertical">
                <UButton
                    block
                    :label="$t('components.debug_info.reset_clipboard')"
                    @click="
                        clipboardStore.clear();
                        searchesStore.clear();
                    "
                />
                <UButton block :label="$t('components.debug_info.reset_local_storage')" @click="() => resetLocalStorage()" />
                <UButton
                    block
                    color="error"
                    external
                    to="/api/clear-site-data"
                    :label="$t('components.debug_info.factory_reset')"
                />
                <UButton
                    block
                    color="neutral"
                    :label="$t('components.debug_info.test_notifications')"
                    @click="() => sendTestNotifications()"
                />
                <UButton
                    block
                    color="neutral"
                    :label="$t('components.debug_info.trigger_banner_message')"
                    @click="() => triggerBannerMessage()"
                />
                <UButton
                    block
                    color="neutral"
                    :label="$t('components.debug_info.trigger_error')"
                    @click="() => triggerErrorPage()"
                />
                <UButton
                    block
                    color="neutral"
                    :label="$t('components.debug_info.toggle_log_level')"
                    @click="() => toggleLogLevel()"
                />
            </UFieldGroup>
        </UFormField>

        <UFormField class="grid grid-cols-2 items-center gap-2" name="permissions" :label="$t('components.debug_info.perms')">
            <p v-if="!activeChar">
                {{ $t('components.debug_info.no_char_selected') }}
            </p>
            <UCollapsible v-else>
                <UButton
                    class="group"
                    variant="soft"
                    :label="$t('components.debug_info.perms')"
                    icon="i-mdi-key"
                    trailing-icon="i-mdi-chevron-down"
                    block
                    :ui="{
                        trailingIcon: 'group-data-[state=open]:rotate-180 transition-transform duration-200',
                    }"
                />

                <template #content>
                    <PermList :permissions="permissions" :attributes="attributes" disabled class="w-full" />
                </template>
            </UCollapsible>
        </UFormField>
    </UPageCard>
</template>
