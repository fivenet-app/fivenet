<script lang="ts" setup>
import { Switch } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { watchDebounced } from '@vueuse/core';
import { BellIcon, ChevronRightIcon } from 'mdi-vue3';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import TablePagination from '~/components/partials/elements/TablePagination.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { useNotificatorStore } from '~/store/notificator';
import { GetNotificationsResponse } from '~~/gen/ts/services/notificator/notificator';

const { $grpc } = useNuxtApp();

const notificator = useNotificatorStore();
const offset = ref(0n);

const includeRead = ref(false);

const { data, pending, refresh, error } = useLazyAsyncData(`notifications-${offset.value}`, () => getNotifications());

async function getNotifications(): Promise<GetNotificationsResponse> {
    try {
        const call = $grpc.getNotificatorClient().getNotifications({
            pagination: {
                offset: offset.value,
            },
            includeRead: includeRead.value,
        });

        const { response } = await call;

        return response;
    } catch (e) {
        $grpc.handleError(e as RpcError);
        throw e;
    }
}

async function markAllRead(): Promise<void> {
    await notificator.markNotifications({
        ids: [],
        all: true,
    });

    const now = {
        timestamp: undefined,
    };
    data.value?.notifications.forEach((v) => (v.readAt = now));
}

watch(offset, async () => refresh());
watchDebounced(includeRead, async () => refresh(), { debounce: 500, maxWait: 1500 });
</script>

<template>
    <div class="py-2 pb-14">
        <div class="px-1 sm:px-2 lg:px-4">
            <div class="sm:flex sm:items-center">
                <div class="sm:flex-auto">
                    <form @submit.prevent="refresh()">
                        <div class="flex flex-row items-center gap-2 sm:mx-auto">
                            <div class="flex-1 form-control">
                                <label for="search" class="block text-sm font-medium leading-6 text-neutral"
                                    >{{ $t('components.notifications.include_read') }}
                                </label>
                                <div class="relative flex items-center mt-3">
                                    <Switch
                                        v-model="includeRead"
                                        :class="[
                                            includeRead ? 'bg-error-500' : 'bg-base-700',
                                            'relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-neutral focus:ring-offset-2',
                                        ]"
                                    >
                                        <span class="sr-only">{{ $t('components.notifications.include_read') }}</span>
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
                                    :disabled="data?.notifications === undefined || data?.notifications.length === 0"
                                    class="inline-flex px-3 py-2 text-sm font-semibold rounded-md text-neutral focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                                    :class="
                                        data?.notifications === undefined || data?.notifications.length === 0
                                            ? 'bg-primary-500 hover:bg-primary-400'
                                            : 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                    "
                                    @click="markAllRead()"
                                >
                                    {{ $t('components.notifications.mark_all_read') }}
                                </button>
                            </div>
                        </div>
                    </form>
                </div>
            </div>
            <div class="flow-root mt-2">
                <div class="mx-0 -my-2 overflow-x-auto">
                    <div class="inline-block min-w-full py-2 align-middle px-1">
                        <DataPendingBlock v-if="pending" :message="$t('common.loading', [$t('common.notification', 2)])" />
                        <DataErrorBlock
                            v-else-if="error"
                            :title="$t('common.unable_to_load', [$t('common.notification', 2)])"
                            :retry="refresh"
                        />
                        <DataNoDataBlock
                            v-else-if="data?.notifications.length === 0"
                            :type="$t('common.notification', 2)"
                            :icon="BellIcon"
                        />
                        <div v-else>
                            <ul class="flex flex-col">
                                <li
                                    v-for="not in data?.notifications"
                                    :key="not.id"
                                    class="relative flex justify-between my-1 gap-x-6 px-4 py-5 hover:bg-base-700 bg-base-800 sm:px-6 rounded-lg"
                                >
                                    <div class="flex gap-x-4">
                                        <div class="min-w-0 flex-auto">
                                            <p class="py-2 pl-4 pr-3 text-sm font-medium text-neutral sm:pl-1">
                                                <span v-if="not.data && not.data.link">
                                                    <!-- @vue-expect-error the route should be valid... at least in most cases -->
                                                    <NuxtLink :to="not.data.link.to">
                                                        <span class="absolute inset-x-0 -top-px bottom-0" />
                                                        {{ $t(not.title!.key, not.title?.parameters ?? {}) }}
                                                    </NuxtLink>
                                                </span>
                                                <span v-else>
                                                    {{ $t(not.title!.key, not.title?.parameters ?? {}) }}
                                                </span>
                                            </p>
                                            <p class="mt-1 flex text-xs leading-5 text-gray-200">
                                                {{ $t(not.content!.key, not.content?.parameters ?? {}) }}
                                            </p>
                                        </div>
                                    </div>
                                    <div class="flex items-center gap-x-4">
                                        <div class="hidden sm:flex sm:flex-col sm:items-end">
                                            <p class="mt-1 text-xs leading-5 text-gray-500">
                                                {{ $t('common.received') }}
                                                <GenericTime :value="not.createdAt" :ago="true" />
                                            </p>
                                            <div v-if="!not.readAt" class="mt-1 flex items-center gap-x-1.5">
                                                <div class="flex-none rounded-full bg-success-500/20 p-1">
                                                    <div class="h-1.5 w-1.5 rounded-full bg-success-500" />
                                                </div>
                                                <p class="text-xs leading-5 text-gray-500">
                                                    {{ $t('components.notifications.unread') }}
                                                </p>
                                            </div>
                                        </div>
                                        <ChevronRightIcon
                                            v-if="not.data && not.data.link"
                                            class="h-5 w-5 flex-none text-gray-400"
                                            aria-hidden="true"
                                        />
                                        <span v-else class="h-5 w-5"></span>
                                    </div>
                                </li>
                            </ul>

                            <TablePagination
                                class="mt-2"
                                :pagination="data?.pagination"
                                :refresh="refresh"
                                @offset-change="offset = $event"
                            />
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
