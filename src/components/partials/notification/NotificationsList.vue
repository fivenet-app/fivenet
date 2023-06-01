<script lang="ts" setup>
import { BellSlashIcon, ChevronRightIcon } from '@heroicons/vue/24/solid';
import { RpcError } from '@protobuf-ts/runtime-rpc/build/types';
import { PaginationResponse } from '~~/gen/ts/resources/common/database/database';
import { Notification } from '~~/gen/ts/resources/notifications/notifications';
import DataErrorBlock from '../DataErrorBlock.vue';
import DataPendingBlock from '../DataPendingBlock.vue';
import TablePagination from '../TablePagination.vue';

const { $grpc } = useNuxtApp();

const pagination = ref<PaginationResponse>();
const offset = ref(BigInt(0));

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
</script>

<template>
    <div class="py-2">
        <div class="px-2 sm:px-6 lg:px-8">
            <div class="sm:flex sm:items-center">
                <div class="sm:flex-auto">
                    <form @submit.prevent="refresh()">
                        <label for="search" class="block mb-2 text-sm font-medium leading-6 text-neutral">
                            {{ $t('common.search') }}
                        </label>
                        <div class="flex flex-row items-center gap-2 sm:mx-auto">
                            <div class="flex-1 form-control">
                                <!-- TODO Include read notification-->
                            </div>
                            <div class="flex-initial">
                                <button
                                    disabled
                                    class="inline-flex px-3 py-2 text-sm font-semibold rounded-md bg-primary-500 text-neutral hover:bg-primary-400 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                                >
                                    Read all {{ $t('common.notification', 2) }}
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
                        <button
                            v-else-if="notifications && notifications.length === 0"
                            type="button"
                            class="relative block w-full p-12 text-center border-2 border-gray-300 border-dashed rounded-lg hover:border-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
                        >
                            <BellSlashIcon class="w-12 h-12 mx-auto text-neutral" />
                            <span class="block mt-2 text-sm font-semibold text-gray-300">
                                {{ $t('common.not_found', [$t('common.notification', 2)]) }}
                                {{ $t('components.documents.document_list.no_documents_hint') }}
                            </span>
                        </button>
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
                                            <p class="mt-1 flex text-xs leading-5 text-gray-500">
                                                {{ $t(not.content!.key, not.content?.parameters ?? []) }}
                                            </p>
                                        </div>
                                    </div>
                                    <div class="flex items-center gap-x-4">
                                        <div class="hidden sm:flex sm:flex-col sm:items-end">
                                            <p class="mt-1 text-xs leading-5 text-gray-500">
                                                Received at
                                                <time :datetime="toDate(not.createdAt)?.toLocaleDateString()">
                                                    {{ useLocaleTimeAgo(toDate(not.createdAt)!).value }}
                                                </time>
                                            </p>
                                            <div v-if="!not.readAt" class="mt-1 flex items-center gap-x-1.5">
                                                <div class="flex-none rounded-full bg-emerald-500/20 p-1">
                                                    <div class="h-1.5 w-1.5 rounded-full bg-emerald-500" />
                                                </div>
                                                <p class="text-xs leading-5 text-gray-500">Unread</p>
                                            </div>
                                        </div>
                                        <ChevronRightIcon
                                            v-if="not.data && not.data.link"
                                            class="h-5 w-5 flex-none text-gray-400"
                                            aria-hidden="true"
                                        />
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
