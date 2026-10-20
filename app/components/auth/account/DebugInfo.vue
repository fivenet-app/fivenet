<script lang="ts" setup>
import { listEnumValues } from '@protobuf-ts/runtime';
import { LogLevels } from 'consola';
import CopyToClipboardButton from '~/components/partials/CopyToClipboardButton.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import PermList from '~/components/settings/roles/PermList.vue';
import { useGRPCWebsocketTransport } from '~/composables/grpcws';
import { useAuthStore } from '~/stores/auth';
import { useClipboardStore } from '~/stores/clipboard';
import { useSettingsStore } from '~/stores/settings';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';
import { clearClientErrors, getClientErrors, type ClientError } from '~/utils/client-errors';

const clipboardStore = useClipboardStore();

const searchesStore = useSearchesStore();

const settingsStore = useSettingsStore();

const authStore = useAuthStore();
const { activeChar, attributes, permissions } = storeToRefs(authStore);
const { clearAuthInfo } = authStore;

const authSessionStore = useAuthSessionStore();
const { userInfo } = storeToRefs(authSessionStore);

const notifications = useNotificationsStore();

const { webSocket } = useGRPCWebsocketTransport();

async function resetLocalStorage(): Promise<void> {
    clearAuthInfo();

    if (import.meta.client) {
        window.localStorage.clear();
        window.sessionStorage.clear();
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
    console.warn('Setting log level to', getDefaultLogLevel() >= 4 ? 'DEBUG' : 'WARN');
}

const version = APP_VERSION;

const { name: browserName, platform: browserPlatform } = getBrowserNameAndPlatform();

// The session buffer is intentionally loaded once when this page mounts; users reach this page after an error and can reload it to refresh the list.
const clientErrors = ref<ClientError[]>(getClientErrors());

function clearCollectedErrors(): void {
    clearClientErrors();
    clientErrors.value = [];
}
</script>

<template>
    <UPageCard :description="$t('components.debug_info.subtitle')" :ui="{ body: 'w-full' }">
        <template #title>
            <div class="flex flex-1 items-center gap-2">
                <span class="flex-1">{{ $t('components.debug_info.title') }}</span>

                <CopyToClipboardButton :label="$t('common.copy')" :value="collectDebugInfo" show-text />
            </div>
        </template>

        <UFormField class="grid grid-cols-2 items-center gap-2" name="version" :label="$t('components.debug_info.version')">
            <div class="inline-flex w-full justify-between">
                <span>
                    <code>{{ version }}</code> / <code>{{ settingsStore.version }}</code>
                </span>

                <CopyToClipboardButton :value="`${version}/ ${settingsStore.version}`" />
            </div>
        </UFormField>

        <UFormField class="grid grid-cols-2 items-center gap-2" name="version" :label="$t('common.browser')">
            <div class="inline-flex w-full justify-between">
                <span>
                    <code>{{ browserName }}</code> (<code>{{ browserPlatform }}</code
                    >)
                </span>

                <CopyToClipboardButton :value="browserName" />
            </div>
        </UFormField>

        <USeparator :label="$t('components.debug_info.context_title')" class="my-4" />

        <UAlert
            v-if="!activeChar"
            color="info"
            variant="subtle"
            icon="i-mdi-account-question-outline"
            :title="$t('components.debug_info.no_char_selected_title')"
            :description="$t('components.debug_info.no_char_selected_description')"
            :actions="[
                {
                    label: $t('components.debug_info.choose_character'),
                    icon: 'i-mdi-account-switch-outline',
                    to: '/auth/character-selector',
                },
            ]"
        />

        <UFormField
            v-if="activeChar"
            class="grid grid-cols-2 items-center gap-2"
            name="activeCharId"
            :label="$t('components.debug_info.active_char_id')"
        >
            <div class="inline-flex w-full justify-between">
                <code>{{ activeChar.userId }}</code>

                <CopyToClipboardButton :value="activeChar.userId" />
            </div>
        </UFormField>

        <UFormField
            v-if="activeChar"
            class="grid grid-cols-2 items-center gap-2"
            name="activeCharJob"
            :label="$t('common.job')"
        >
            <div class="flex w-full items-center justify-between">
                <div class="flex flex-col gap-1">
                    <template v-if="userInfo?.originalJob">
                        <div>
                            {{ userInfo?.originalJob.job }} ({{ $t('common.rank') }}:
                            {{ userInfo?.originalJob.jobGrade ?? '?' }})
                        </div>

                        <div>
                            ({{
                                $t('common.impersonation_active', {
                                    job: `${userInfo.job} (${userInfo.jobGrade ?? '?'})`,
                                })
                            }})
                        </div>
                    </template>
                    <div v-else>
                        <code>{{ activeChar.job }}</code> ({{ $t('common.rank') }}: <code>{{ activeChar.jobGrade }}</code
                        >)
                    </div>
                </div>

                <div>
                    <CopyToClipboardButton :value="`${activeChar.job} (${$t('common.rank')}: ${activeChar.jobGrade})`" />
                </div>
            </div>
        </UFormField>

        <UFormField
            v-if="userInfo"
            class="grid grid-cols-2 items-center gap-2"
            name="sessionExpiration"
            :label="$t('components.debug_info.access_token_expiration')"
        >
            <span v-if="!userInfo.expiration">{{ $t('common.na') }}</span>
            <template v-else>
                <GenericTime :value="userInfo.expiration" ago />
                (<GenericTime :value="userInfo.expiration" type="long" />)
            </template>
        </UFormField>

        <UFormField class="grid grid-cols-2 items-center gap-2" name="status" :label="$t('common.websocket')">
            <code>{{ webSocket.status.value }}</code>
        </UFormField>

        <UFormField class="grid grid-cols-2 items-center gap-2" name="nuiInfo" :label="$t('components.debug_info.nui_info')">
            {{ settingsStore.nuiEnabled ? $t('common.enabled') : $t('common.disabled') }}:
            <code>{{ settingsStore.nuiResourceName ?? $t('common.na') }}</code>
        </UFormField>

        <USeparator :label="$t('components.debug_info.client_errors')" class="my-4" />

        <UAlert
            v-if="clientErrors.length === 0"
            color="success"
            variant="subtle"
            icon="i-mdi-check-circle-outline"
            :title="$t('components.debug_info.no_client_errors')"
        />
        <UCollapsible v-else class="w-full" :ui="{ content: 'mt-2' }">
            <UButton
                class="group"
                block
                color="warning"
                variant="soft"
                icon="i-mdi-alert-circle-outline"
                :label="$t('components.debug_info.client_errors_count', { count: clientErrors.length })"
                trailing-icon="i-mdi-chevron-down"
                :ui="{ trailingIcon: 'group-data-[state=open]:rotate-180 transition-transform' }"
            />

            <template #content>
                <UFieldGroup class="w-full" orientation="vertical">
                    <UCard v-for="(clientError, index) in clientErrors" :key="index" class="text-sm">
                        <div class="flex flex-col gap-1">
                            <div class="flex justify-between gap-2 font-semibold">
                                <span>{{ clientError.source }}</span>
                                <GenericTime :value="new Date(clientError.timestamp)" type="long" />
                            </div>
                            <code class="break-words">{{ clientError.message }}</code>
                            <span class="break-words text-muted">{{ clientError.route }}</span>
                            <pre
                                v-if="clientError.stack !== 'N/A'"
                                class="max-h-32 overflow-auto text-xs whitespace-pre-wrap"
                                >{{ clientError.stack }}</pre>
                        </div>
                    </UCard>

                    <UButton
                        color="neutral"
                        variant="soft"
                        icon="i-mdi-delete"
                        :label="$t('components.debug_info.clear_client_errors')"
                        @click="clearCollectedErrors()"
                    />
                </UFieldGroup>
            </template>
        </UCollapsible>

        <USeparator :label="$t('components.debug_info.tools_title')" class="my-4" />

        <UFormField
            class="grid grid-cols-2 items-center gap-2"
            name="debugFunctions"
            :label="$t('components.debug_info.debug_functions')"
        >
            <UFieldGroup class="flex w-full break-words" orientation="vertical">
                <UButton
                    block
                    color="warning"
                    :label="$t('components.debug_info.factory_reset')"
                    external
                    to="/api/clear-site-data"
                    trailing-icon="i-mdi-restart-alert"
                    variant="soft"
                />

                <UButton
                    block
                    color="neutral"
                    :label="$t('components.debug_info.reset_clipboard')"
                    trailing-icon="i-mdi-clipboard-remove"
                    variant="soft"
                    @click="
                        clipboardStore.clear();
                        searchesStore.clear();
                    "
                />

                <UButton
                    block
                    color="error"
                    :label="$t('components.debug_info.reset_local_storage')"
                    trailing-icon="i-mdi-delete"
                    variant="soft"
                    @click="() => resetLocalStorage()"
                />
            </UFieldGroup>

            <UCollapsible class="mt-2" :ui="{ content: 'mt-2' }">
                <UButton
                    class="group"
                    variant="ghost"
                    :label="$t('components.debug_info.advanced_debug')"
                    icon="i-mdi-bug-check-outline"
                    trailing-icon="i-mdi-chevron-down"
                    block
                    :ui="{
                        trailingIcon: 'group-data-[state=open]:rotate-180 transition-transform duration-200',
                    }"
                />

                <template #content>
                    <UFieldGroup orientation="vertical">
                        <UButton
                            block
                            color="neutral"
                            :label="$t('components.debug_info.test_notifications')"
                            trailing-icon="i-mdi-bell-notification-outline"
                            variant="soft"
                            @click="() => sendTestNotifications()"
                        />

                        <UButton
                            block
                            color="neutral"
                            :label="$t('components.debug_info.trigger_banner_message')"
                            trailing-icon="i-mdi-message-cog-outline"
                            variant="soft"
                            @click="() => triggerBannerMessage()"
                        />

                        <UButton
                            block
                            color="warning"
                            :label="$t('components.debug_info.toggle_log_level')"
                            trailing-icon="i-mdi-bug-outline"
                            variant="soft"
                            @click="() => toggleLogLevel()"
                        />

                        <UButton
                            block
                            color="error"
                            variant="soft"
                            trailing-icon="i-mdi-error-outline"
                            :label="$t('components.debug_info.trigger_error')"
                            @click="() => triggerErrorPage()"
                        />
                    </UFieldGroup>
                </template>
            </UCollapsible>
        </UFormField>

        <UFormField
            v-if="activeChar"
            class="grid grid-cols-2 items-center gap-2"
            name="permissions"
            :label="$t('components.debug_info.perms')"
        >
            <UCollapsible>
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
                    <PermList class="w-full" :permissions="permissions" :attributes="attributes" disabled />
                </template>
            </UCollapsible>
        </UFormField>
    </UPageCard>
</template>
