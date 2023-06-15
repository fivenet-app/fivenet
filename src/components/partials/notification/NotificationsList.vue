<script lang="ts" setup>
import { Switch } from '@headlessui/vue';
import SvgIcon from '@jamescoyle/vue-icon';
import { mdiBell, mdiChevronRight } from '@mdi/js';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { watchDebounced } from '@vueuse/core';
import { PaginationResponse } from '~~/gen/ts/resources/common/database/database';
import { Notification } from '~~/gen/ts/resources/notifications/notifications';
import DataErrorBlock from '../DataErrorBlock.vue';
import DataNoDataBlock from '../DataNoDataBlock.vue';
import DataPendingBlock from '../DataPendingBlock.vue';
import TablePagination from '../TablePagination.vue';

const { $grpc } = useNuxtApp();

const pagination = ref<PaginationResponse>();
const offset = ref(0n);

const includeRead = ref(false);

const {
    data: notifications,
    pending,
    refresh,
    error,
} = useLazyAsyncData(`notifications-${offset.value}`, () => getNotifications());

async function getNotifications(): Promise<Array<Notification>> {
    return new Promise(async (res, rej) => {
        try {
            const call = $grpc.getNotificatorClient().getNotifications({
                pagination: {
                    offset: offset.value,
                },
                includeRead: includeRead.value,
            });

            const { response } = await call;

            offset.value = response.pagination?.offset!;
            pagination.value = response.pagination;

            return res(response.notifications);
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

async function markAllRead(): Promise<void> {
    return new Promise(async (res, rej) => {
        const now = {
            timestamp: undefined,
        };
        try {
            await $grpc.getNotificatorClient().readNotifications({
                ids: [],
                all: true,
            });

            notifications.value?.forEach((v) => (v.readAt = now));

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

async function markRead(ids: bigint[]): Promise<void> {
    return new Promise(async (res, rej) => {
        try {
            await $grpc.getNotificatorClient().readNotifications({
                ids: ids,
            });

            return res();
        } catch (e) {
            $grpc.handleError(e as RpcError);
            return rej(e as RpcError);
        }
    });
}

watchDebounced(includeRead, async () => refresh(), { debounce: 500, maxWait: 1500 });
</script>

<template>
    <div class="py-2">
        <div class="px-2 sm:px-6 lg:px-8">
            <div class="sm:flex sm:items-center">
                <div class="sm:flex-auto">
                    <form @submit.prevent="refresh()">
                        <div class="flex flex-row items-center gap-2 sm:mx-auto">
                            <div class="flex-1 form-control">
                                <label for="search" class="block text-sm font-medium leading-6 text-neutral"
                                    >{{ $t('pages.notifications.include_read') }}
                                </label>
                                <div class="relative flex items-center mt-3">
                                    <Switch
                                        v-model="includeRead"
                                        :class="[
                                            includeRead ? 'bg-error-500' : 'bg-base-700',
                                            'relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-neutral focus:ring-offset-2',
                                        ]"
                                    >
                                        <span class="sr-only">>{{ $t('pages.notifications.include_read') }}</span>
                                        <span
                                            aria-hidden="true"
                                            :class="[
                                                includeRead ? 'translate-x-5' : 'translate-x-0',
                                                'pointer-events-none inline-block h-5 w-5 transform rounded-full bg-neutral ring-0 transition duration-200 ease-in-out',
                                            ]"
                                        />
                                    </Switch>
                                </div>
                            </div>
                            <div class="flex-initial">
                                <button
                                    type="button"
                                    :disabled="!notifications || notifications.length <= 0"
                                    @click="markAllRead"
                                    class="inline-flex px-3 py-2 text-sm font-semibold rounded-md bg-primary-500 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                                >
                                    {{ $t('pages.notifications.mark_all_read') }}
                                </button>
                            </div>
                        </div>
                    </form>
                </div>
            </div>
            <div class="flow-root mt-2">
                <div class="-mx-4 -my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
                    <div class="inline-block min-w-full py-2 align-middle sm:px-6 lg:px-8">
                        <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.notification', 2)])" />
                        <DataErrorBlock
                            v-else-if="error"
                            :title="$t('common.unable_to_load', [$t('common.notification', 2)])"
                            :retry="refresh"
                        />
                        <DataNoDataBlock
                            v-else-if="notifications && notifications.length === 0"
                            :type="$t('common.notification', 2)"
                            :icon="mdiBell"
                        />
                        <div v-else>
                            <ul class="flex flex-col">
                                <li
                                    v-for="not in notifications"
                                    :key="not.id.toString()"
                                    class="relative flex justify-between my-1 gap-x-6 px-4 py-5 hover:bg-base-800 bg-base-850 sm:px-6 rounded-lg"
                                >
                                    <div class="flex gap-x-4">
                                        <div class="min-w-0 flex-auto">
                                            <p class="py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-0">
                                                <span v-if="not.data && not.data.link">
                                                    <NuxtLink :to="not.data?.link?.to">
                                                        <span class="absolute inset-x-0 -top-px bottom-0" />
                                                        {{ $t(not.title!.key, not.title?.parameters ?? []) }}
                                                    </NuxtLink>
                                                </span>
                                                <span v-else>
                                                    {{ $t(not.title!.key, not.title?.parameters ?? []) }}
                                                </span>
                                            </p>
                                            <p class="mt-1 flex text-xs leading-5 text-gray-200">
                                                {{ $t(not.content!.key, not.content?.parameters ?? []) }}
                                            </p>
                                        </div>
                                    </div>
                                    <div class="flex items-center gap-x-4">
                                        <div class="hidden sm:flex sm:flex-col sm:items-end">
                                            <p class="mt-1 text-xs leading-5 text-gray-500">
                                                {{ $t('common.received') }}
                                                <time :datetime="toDate(not.createdAt)?.toLocaleDateString()">
                                                    {{ useLocaleTimeAgo(toDate(not.createdAt)!).value }}
                                                </time>
                                            </p>
                                            <div v-if="!not.readAt" class="mt-1 flex items-center gap-x-1.5">
                                                <div class="flex-none rounded-full bg-green-500/20 p-1">
                                                    <div class="h-1.5 w-1.5 rounded-full bg-green-500" />
                                                </div>
                                                <p class="text-xs leading-5 text-gray-500">
                                                    {{ $t('pages.notifications.unread') }}
                                                </p>
                                            </div>
                                        </div>
                                        <SvgIcon
                                            v-if="not.data && not.data.link"
                                            class="h-5 w-5 flex-none text-gray-400"
                                            aria-hidden="true"
                                            type="mdi"
                                            :path="mdiChevronRight"
                                        />
                                        <span class="h-5 w-5" v-else></span>
                                    </div>
                                </li>
                            </ul>

                            <TablePagination :pagination="pagination" @offset-change="offset = $event" class="mt-2" />
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
