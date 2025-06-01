<script lang="ts" setup>
import { listEnumValues } from '@protobuf-ts/runtime';
import { LogLevels } from 'consola';
import CopyToClipboardButton from '~/components/partials/CopyToClipboardButton.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { useGRPCWebsocketTransport } from '~/composables/grpc/grpcws';
import { useAuthStore } from '~/stores/auth';
import { useClipboardStore } from '~/stores/clipboard';
import { useNotificatorStore } from '~/stores/notificator';
import { useSettingsStore } from '~/stores/settings';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const clipboardStore = useClipboardStore();

const settings = useSettingsStore();

const authStore = useAuthStore();
const { activeChar, permissions, accessTokenExpiration, isSuperuser } = storeToRefs(authStore);
const { clearAuthInfo } = authStore;

const notifications = useNotificatorStore();

const { webSocket, wsInitiated } = useGRPCWebsocketTransport();

async function resetLocalStorage(): Promise<void> {
    clearAuthInfo();

    if (import.meta.client) {
        window.localStorage.clear();
    }

    await navigateTo({ name: 'index' });
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
                        click: () => alert('Test was successful!'),
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

function setLogLevel(): void {
    setDefaultLogLevel(getDefaultLogLevel() !== LogLevels.debug ? LogLevels.debug : LogLevels.warn);
    console.warn('Log Level set to', getDefaultLogLevel() === 4 ? 'DEBUG' : 'WARN');
}

const isDevEnv = import.meta.dev;

const version = APP_VERSION;
</script>

<template>
    <UDashboardPanelContent>
        <UDashboardSection :title="$t('components.debug_info.title')" :description="$t('components.debug_info.subtitle')">
            <UFormGroup
                class="grid grid-cols-2 items-center gap-2"
                name="version"
                :label="$t('components.debug_info.version')"
                :ui="{ container: '' }"
            >
                <div class="inline-flex w-full justify-between">
                    <span> {{ version }}/ {{ settings.version }} </span>
                    <CopyToClipboardButton :value="`${version}/ ${settings.version}`" />
                </div>
            </UFormGroup>

            <UFormGroup
                v-if="activeChar"
                class="grid grid-cols-2 items-center gap-2"
                name="activeCharId"
                :label="$t('components.debug_info.active_char_id')"
                :ui="{ container: '' }"
            >
                <div class="inline-flex w-full justify-between">
                    <span>
                        {{ activeChar.userId }}
                    </span>
                    <CopyToClipboardButton :value="activeChar.userId" />
                </div>
            </UFormGroup>

            <UFormGroup
                v-if="activeChar"
                class="grid grid-cols-2 items-center gap-2"
                name="activeCharJob"
                :label="$t('common.job')"
                :ui="{ container: '' }"
            >
                <div class="inline-flex w-full justify-between">
                    <span>{{ activeChar.job }} ({{ $t('common.rank') }}: {{ activeChar.jobGrade }})</span>
                    <CopyToClipboardButton :value="`${activeChar.job} (${$t('common.rank')}: ${activeChar.jobGrade})`" />
                </div>
            </UFormGroup>

            <UFormGroup
                v-if="accessTokenExpiration"
                class="grid grid-cols-2 items-center gap-2"
                name="accessTokenExpiration"
                :label="$t('components.debug_info.access_token_expiration')"
                :ui="{ container: '' }"
            >
                <GenericTime :value="accessTokenExpiration" ago />
                (<GenericTime :value="accessTokenExpiration" type="long" />)
            </UFormGroup>

            <UFormGroup
                class="grid grid-cols-2 items-center gap-2"
                name="nuiInfo"
                :label="$t('components.debug_info.nui_info')"
                :ui="{ container: '' }"
            >
                {{ settings.nuiEnabled ? $t('common.enabled') : $t('common.disabled') }}:
                {{ settings.nuiResourceName ?? $t('common.na') }}
            </UFormGroup>

            <UFormGroup
                class="grid grid-cols-2 items-center gap-2"
                name="status"
                :label="$t('common.status')"
                :ui="{ container: '' }"
            >
                {{ $t('common.active') }}: {{ wsInitiated ? $t('common.yes') : $t('common.no') }}<br />
                {{ $t('common.status') }}: <code>{{ webSocket.status.value }}</code>
            </UFormGroup>

            <UFormGroup
                class="grid grid-cols-2 items-center gap-2"
                name="debugFunctions"
                :label="$t('components.debug_info.debug_functions')"
                :ui="{ container: '' }"
            >
                <UButtonGroup class="flex w-full break-words" orientation="vertical">
                    <UButton block @click="clipboardStore.clear()">
                        <span>{{ $t('components.debug_info.reset_clipboard') }}</span>
                    </UButton>
                    <UButton block @click="resetLocalStorage()">
                        <span>{{ $t('components.debug_info.reset_local_storage') }}</span>
                    </UButton>
                    <UButton block color="error" :external="true" to="/api/clear-site-data">
                        <span>{{ $t('components.debug_info.factory_reset') }}</span>
                    </UButton>
                    <UButton block color="gray" @click="sendTestNotifications">
                        <span>{{ $t('components.debug_info.test_notifications') }}</span>
                    </UButton>
                    <UButton block color="gray" @click="triggerBannerMessage">
                        <span>{{ $t('components.debug_info.trigger_banner_message') }}</span>
                    </UButton>
                    <UButton block color="gray" @click="triggerErrorPage">
                        <span>{{ $t('components.debug_info.trigger_error') }}</span>
                    </UButton>
                    <UButton v-if="isDevEnv || isSuperuser" block color="white" @click="setLogLevel">
                        <span>{{ $t('components.debug_info.toggle_log_level') }}</span>
                    </UButton>
                </UButtonGroup>
            </UFormGroup>

            <UFormGroup
                class="grid grid-cols-2 items-center gap-2"
                name="permissions"
                :label="$t('components.debug_info.perms')"
                :ui="{ container: '' }"
            >
                <UAccordion
                    variant="soft"
                    :items="[{ label: $t('components.debug_info.perms'), slot: 'perms', icon: 'i-mdi-key' }]"
                    :ui="{ wrapper: 'flex flex-col w-full' }"
                >
                    <template #perms>
                        <p v-if="!activeChar">
                            {{ $t('components.debug_info.no_char_selected') }}
                        </p>
                        <ul v-else class="divide-y divide-gray-100 dark:divide-gray-800" role="list">
                            <li
                                v-for="perm in permissions"
                                :key="perm"
                                class="flex items-center justify-between py-4 pl-4 pr-5 text-sm leading-6"
                            >
                                <UIcon class="size-5 shrink-0 text-gray-400" name="i-mdi-key" />
                                <div class="ml-4 flex min-w-0 flex-1 gap-2">
                                    <span class="truncate font-medium">
                                        {{ perm }}
                                    </span>
                                </div>
                            </li>
                        </ul>
                    </template>
                </UAccordion>
            </UFormGroup>
        </UDashboardSection>

        <UDivider class="mb-4" />
    </UDashboardPanelContent>
</template>
