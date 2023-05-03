<script lang="ts" setup>
import { KeyIcon } from '@heroicons/vue/20/solid';
import { useAuthStore } from '~/store/auth';

const authStore = useAuthStore();

const accessTokenExpiration = computed(() => authStore.getAccessTokenExpiration);
const activeChar = computed(() => authStore.getActiveChar);
const perms = computed(() => authStore.getPermissions);
</script>

<template>
    <div class="overflow-hidden bg-base-800 shadow sm:rounded-lg text-neutral mt-3">
        <div class="px-4 py-5 sm:px-6">
            <h3 class="text-base font-semibold leading-6">
                {{ $t('components.debug_info.title') }}</h3>
            <p class="mt-1 max-w-2xl text-sm">
                {{ $t('components.debug_info.subtitle') }}
            </p>
        </div>
        <div class="border-t border-base-400 px-4 py-5 sm:p-0">
            <dl class="sm:divide-y sm:divide-base-400">
                <div v-if="activeChar" class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6 sm:py-5">
                    <dt class="text-sm font-medium">
                        {{ $t('components.debug_info.active_char_id') }}
                    </dt>
                    <dd class="mt-1 text-sm sm:col-span-2 sm:mt-0">
                        {{ activeChar.getUserId() }}
                    </dd>
                </div>
                <div v-if="accessTokenExpiration" class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6 sm:py-5">
                    <dt class="text-sm font-medium">
                        {{ $t('components.debug_info.access_token_expiration') }}
                    </dt>
                    <dd class="mt-1 text-sm sm:col-span-2 sm:mt-0">
                        <time :datetime="accessTokenExpiration.toDateString()">
                            {{ useLocaleTimeAgo(accessTokenExpiration).value }} ({{ $d(accessTokenExpiration, 'long') }})
                        </time>
                    </dd>
                </div>
                <div v-if="perms.length > 0" class="py-4 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6 sm:py-5">
                    <dt class="text-sm font-medium">
                        {{ $t('components.debug_info.perms') }}
                    </dt>
                    <dd class="mt-1 text-sm sm:col-span-2 sm:mt-0">
                        <ul role="list" class="divide-y divide-gray-100 rounded-md border border-gray-200">
                            <li v-for="perm in perms"
                                class="flex items-center justify-between py-4 pl-4 pr-5 text-sm leading-6">
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
