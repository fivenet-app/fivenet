<script lang="ts" setup>
import CopyToClipboardButton from '~/components/partials/CopyToClipboardButton.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { useGRPCWebsocketTransport } from '~/composables/grpcws';
import { useAuthStore } from '~/store/auth';
import { useClipboardStore } from '~/store/clipboard';
import { useNotificatorStore } from '~/store/notificator';
import { useSettingsStore } from '~/store/settings';
import { NotificationType } from '~~/gen/ts/resources/notifications/notifications';

const clipboardStore = useClipboardStore();

const settings = useSettingsStore();

const authStore = useAuthStore();
const { activeChar, permissions, getAccessTokenExpiration } = storeToRefs(authStore);
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
    NotificationTypes.forEach((notificationType, index) => {
        notifications.add({
            title: { key: 'notifications.system.test_notification.title', parameters: { index: (index + 1).toString() } },
            description: {
                key: 'notifications.system.test_notification.content',
                parameters: { type: NotificationType[notificationType] },
            },
            type: notificationType,
            onClick: () => alert('Test was successful!'),
        });
    });
}

const version = APP_VERSION;
</script>

<template>
    <UDashboardPanelContent class="pb-24">
        <UDashboardSection :title="$t('components.debug_info.title')" :description="$t('components.debug_info.subtitle')">
            <UFormGroup
                name="version"
                :label="$t('components.debug_info.version')"
                class="grid grid-cols-2 items-center gap-2"
                :ui="{ container: '' }"
            >
                <div class="inline-flex w-full justify-between">
                    <span> {{ version }}/ {{ settings.version }} </span>
                    <CopyToClipboardButton :value="`${version}/ ${settings.version}`" />
                </div>
            </UFormGroup>

            <UFormGroup
                v-if="activeChar"
                name="activeCharId"
                :label="$t('components.debug_info.active_char_id')"
                class="grid grid-cols-2 items-center gap-2"
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
                name="activeCharJob"
                :label="$t('common.job')"
                class="grid grid-cols-2 items-center gap-2"
                :ui="{ container: '' }"
            >
                <div class="inline-flex w-full justify-between">
                    <span>{{ activeChar.job }} ({{ $t('common.rank') }}: {{ activeChar.jobGrade }})</span>
                    <CopyToClipboardButton :value="`${activeChar.job} (${$t('common.rank')}: ${activeChar.jobGrade})`" />
                </div>
            </UFormGroup>

            <UFormGroup
                v-if="getAccessTokenExpiration"
                name="accessTokenExpiration"
                :label="$t('components.debug_info.access_token_expiration')"
                class="grid grid-cols-2 items-center gap-2"
                :ui="{ container: '' }"
            >
                <GenericTime :value="getAccessTokenExpiration" :ago="true" />
                (<GenericTime :value="getAccessTokenExpiration" type="long" />)
            </UFormGroup>

            <UFormGroup
                name="nuiInfo"
                :label="$t('components.debug_info.nui_info')"
                class="grid grid-cols-2 items-center gap-2"
                :ui="{ container: '' }"
            >
                {{ settings.nuiEnabled ? $t('common.enabled') : $t('common.disabled') }}:
                {{ settings.nuiResourceName ?? $t('common.na') }}
            </UFormGroup>

            <UFormGroup
                name="status"
                :label="$t('common.status')"
                class="grid grid-cols-2 items-center gap-2"
                :ui="{ container: '' }"
            >
                {{ $t('common.active') }}: {{ wsInitiated ? $t('common.yes') : $t('common.no') }}<br />
                {{ $t('common.status') }}: <code>{{ webSocket.status.value }}</code>
            </UFormGroup>

            <UFormGroup
                name="debugFunctions"
                :label="$t('components.debug_info.debug_functions')"
                class="grid grid-cols-2 items-center gap-2"
                :ui="{ container: '' }"
            >
                <UButtonGroup class="flex w-full break-words" orientation="vertical">
                    <UButton @click="clipboardStore.clear()">
                        {{ $t('components.debug_info.reset_clipboard') }}
                    </UButton>
                    <UButton @click="resetLocalStorage()">
                        {{ $t('components.debug_info.reset_local_storage') }}
                    </UButton>
                    <UButton color="red" :external="true" to="/api/clear-site-data">
                        {{ $t('components.debug_info.factory_reset') }}
                    </UButton>
                    <UButton color="gray" @click="sendTestNotifications">
                        {{ $t('components.debug_info.test_notifications') }}
                    </UButton>
                </UButtonGroup>
            </UFormGroup>

            <UFormGroup
                name="permissions"
                :label="$t('components.debug_info.perms')"
                class="grid grid-cols-2 items-center gap-2"
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
                        <ul v-else role="list" class="divide-y divide-gray-100 dark:divide-gray-800">
                            <li
                                v-for="perm in permissions"
                                :key="perm"
                                class="flex items-center justify-between py-4 pl-4 pr-5 text-sm leading-6"
                            >
                                <UIcon name="i-mdi-key" class="size-5 shrink-0 text-gray-400" />
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
