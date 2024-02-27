<script lang="ts" setup>
import { KeyIcon } from 'mdi-vue3';
import GenericContainerPanel from '~/components/partials/GenericContainerPanel.vue';
import GenericContainerPanelEntry from '~/components/partials/GenericContainerPanelEntry.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { useAuthStore } from '~/store/auth';
import { useClipboardStore } from '~/store/clipboard';
import { useConfigStore } from '~/store/config';
import { useSettingsStore } from '~/store/settings';

const clipboardStore = useClipboardStore();

const config = useConfigStore();

const settings = useSettingsStore();

const authStore = useAuthStore();
const { activeChar, permissions, getAccessTokenExpiration } = storeToRefs(authStore);
const { clearAuthInfo } = authStore;

async function resetLocalStorage(): Promise<void> {
    clearAuthInfo();

    window.localStorage.clear();

    await navigateTo({ name: 'index' });
}
</script>

<template>
    <div class="mx-auto max-w-5xl py-2">
        <GenericContainerPanel>
            <template #title>
                {{ $t('components.debug_info.title') }}
            </template>
            <template #description>
                {{ $t('components.debug_info.subtitle') }}
            </template>
            <template #default>
                <GenericContainerPanelEntry>
                    <template #title>
                        {{ $t('components.debug_info.version') }}
                    </template>
                    <template #default>
                        {{ settings.version }}
                    </template>
                </GenericContainerPanelEntry>
                <GenericContainerPanelEntry v-if="activeChar">
                    <template #title>
                        {{ $t('components.debug_info.active_char_id') }}
                    </template>
                    <template #default>
                        {{ activeChar.userId }}
                    </template>
                </GenericContainerPanelEntry>
                <GenericContainerPanelEntry v-if="activeChar">
                    <template #title>
                        {{ $t('common.job') }}
                    </template>
                    <template #default> {{ activeChar.job }} ({{ $t('common.rank') }}: {{ activeChar.jobGrade }}) </template>
                </GenericContainerPanelEntry>
                <GenericContainerPanelEntry v-if="getAccessTokenExpiration">
                    <template #title>
                        {{ $t('components.debug_info.access_token_expiration') }}
                    </template>
                    <template #default>
                        <GenericTime :value="getAccessTokenExpiration" :ago="true" />
                        (<GenericTime :value="getAccessTokenExpiration" type="long" />)
                    </template>
                </GenericContainerPanelEntry>
                <GenericContainerPanelEntry>
                    <template #title>
                        {{ $t('components.debug_info.nui_info') }}
                    </template>
                    <template #default>
                        {{ config.nuiEnabled ? $t('common.enabled') : $t('common.disabled') }}:
                        {{ config.nuiResourceName ?? $t('common.na') }}
                    </template>
                </GenericContainerPanelEntry>
                <GenericContainerPanelEntry>
                    <template #title>
                        {{ $t('components.debug_info.debug_functions') }}
                    </template>
                    <template #default>
                        <div class="isolate inline-flex rounded-md shadow-sm">
                            <button
                                type="button"
                                class="inline-flex w-full items-center rounded-md bg-base-500 px-3.5 py-2.5 text-sm font-semibold text-neutral hover:bg-base-400"
                                @click="clipboardStore.clear()"
                            >
                                {{ $t('components.debug_info.reset_clipboard') }}
                            </button>
                            <button
                                type="button"
                                class="ml-2 inline-flex w-full items-center rounded-md bg-base-500 px-3.5 py-2.5 text-sm font-semibold text-neutral hover:bg-base-400"
                                @click="resetLocalStorage()"
                            >
                                {{ $t('components.debug_info.reset_local_storage') }}
                            </button>
                            <NuxtLink
                                :external="true"
                                to="/api/clear-site-data"
                                class="ml-2 inline-flex w-full items-center rounded-md bg-error-800 px-3.5 py-2.5 text-center text-sm font-semibold text-neutral hover:bg-error-600"
                            >
                                {{ $t('components.debug_info.factory_reset') }}
                            </NuxtLink>
                        </div>
                    </template>
                </GenericContainerPanelEntry>
                <GenericContainerPanelEntry v-if="permissions.length > 0">
                    <template #title>
                        {{ $t('components.debug_info.perms') }}
                    </template>
                    <template #default>
                        <ul role="list" class="divide-y divide-gray-100 rounded-md border border-gray-200">
                            <li
                                v-for="perm in permissions"
                                :key="perm"
                                class="flex items-center justify-between py-4 pl-4 pr-5 text-sm leading-6"
                            >
                                <KeyIcon class="h-5 w-5 flex-shrink-0 text-gray-400" aria-hidden="true" />
                                <div class="ml-4 flex min-w-0 flex-1 gap-2">
                                    <span class="truncate font-medium">
                                        {{ perm }}
                                    </span>
                                </div>
                            </li>
                        </ul>
                    </template>
                </GenericContainerPanelEntry>
            </template>
        </GenericContainerPanel>
    </div>
</template>
