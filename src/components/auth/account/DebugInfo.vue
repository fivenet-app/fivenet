<script lang="ts" setup>
import { KeyIcon } from 'mdi-vue3';
import Time from '~/components/partials/elements/Time.vue';
import { useAuthStore } from '~/store/auth';
import { useClipboardStore } from '~/store/clipboard';
import { useSettingsStore } from '~/store/settings';

const authStore = useAuthStore();
const clipboardStore = useClipboardStore();

const settings = useSettingsStore();

const { activeChar, permissions, getAccessTokenExpiration } = storeToRefs(authStore);
const { clearAuthInfo } = authStore;

async function resetLocalStorage(): Promise<void> {
    clearAuthInfo();

    window.localStorage.clear();

    await navigateTo({ name: 'index' });
}
</script>

<template>
    <div class="overflow-hidden bg-base-800 shadow sm:rounded-lg text-neutral mt-3">
        <div class="px-4 py-5 sm:px-6">
            <h3 class="text-base font-semibold leading-6">
                {{ $t('components.debug_info.title') }}
            </h3>
            <p class="mt-1 max-w-2xl text-sm">
                {{ $t('components.debug_info.subtitle') }}
            </p>
        </div>
        <div class="border-t border-base-400 px-4 py-5 sm:p-0">
            <dl class="sm:divide-y sm:divide-base-400">
                <div class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6 sm:py-5">
                    <dt class="text-sm font-medium">
                        {{ $t('components.debug_info.version') }}
                    </dt>
                    <dd class="mt-1 text-sm sm:col-span-2 sm:mt-0">
                        {{ settings.getVersion }}
                    </dd>
                </div>
                <div v-if="activeChar" class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6 sm:py-5">
                    <dt class="text-sm font-medium">
                        {{ $t('components.debug_info.active_char_id') }}
                    </dt>
                    <dd class="mt-1 text-sm sm:col-span-2 sm:mt-0">
                        {{ activeChar.userId }}
                    </dd>
                </div>
                <div v-if="activeChar" class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6 sm:py-5">
                    <dt class="text-sm font-medium">
                        {{ $t('common.job') }}
                    </dt>
                    <dd class="mt-1 text-sm sm:col-span-2 sm:mt-0">
                        {{ activeChar.job }} ({{ $t('common.rank') }}: {{ activeChar.jobGrade }})
                    </dd>
                </div>
                <div v-if="getAccessTokenExpiration" class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6 sm:py-5">
                    <dt class="text-sm font-medium">
                        {{ $t('components.debug_info.access_token_expiration') }}
                    </dt>
                    <dd class="mt-1 text-sm sm:col-span-2 sm:mt-0">
                        <Time :value="getAccessTokenExpiration" :ago="true" />
                        (<Time :value="getAccessTokenExpiration" type="long" />)
                    </dd>
                </div>
                <div v-if="clipboardStore" class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6 sm:py-5">
                    <dt class="text-sm font-medium">
                        {{ $t('components.debug_info.debug_functions') }}
                    </dt>
                    <dd class="mt-1 text-sm sm:col-span-2 sm:mt-0">
                        <span class="isolate inline-flex rounded-md shadow-sm">
                            <button
                                type="button"
                                @click="clipboardStore.clear()"
                                class="rounded-md bg-base-500 py-2.5 px-3.5 text-sm font-semibold text-neutral hover:bg-base-400"
                            >
                                {{ $t('components.debug_info.reset_clipboard') }}
                            </button>
                            <button
                                type="button"
                                @click="resetLocalStorage()"
                                class="rounded-md bg-base-500 py-2.5 px-3.5 ml-2 text-sm font-semibold text-neutral hover:bg-base-400"
                            >
                                {{ $t('components.debug_info.reset_local_storage') }}
                            </button>
                            <NuxtLink
                                :external="true"
                                to="/api/clear-site-data"
                                class="rounded-md bg-error-800 text-center py-2.5 px-3.5 ml-2 text-sm font-semibold text-neutral hover:bg-error-600"
                            >
                                {{ $t('components.debug_info.factory_reset') }}
                            </NuxtLink>
                        </span>
                    </dd>
                </div>
                <div v-if="permissions.length > 0" class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6 sm:py-5">
                    <dt class="text-sm font-medium">
                        {{ $t('components.debug_info.perms') }}
                    </dt>
                    <dd class="mt-1 text-sm sm:col-span-2 sm:mt-0">
                        <ul role="list" class="divide-y divide-gray-100 rounded-md border border-gray-200">
                            <li
                                v-for="perm in permissions"
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
                    </dd>
                </div>
            </dl>
        </div>
    </div>
</template>
