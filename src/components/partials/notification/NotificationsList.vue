<script lang="ts" setup>
import { Switch } from '@headlessui/vue';
import { RpcError } from '@protobuf-ts/runtime-rpc';
import { useTimeoutFn, watchDebounced } from '@vueuse/core';
import { BellIcon, CheckIcon, LinkVariantIcon } from 'mdi-vue3';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import TablePagination from '~/components/partials/elements/TablePagination.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { useNotificatorStore } from '~/store/notificator';
import { GetNotificationsResponse } from '~~/gen/ts/services/notificator/notificator';

withDefaults(
    defineProps<{
        compact?: boolean;
    }>(),
    {
        compact: false,
    },
);

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

async function markRead(...ids: string[]): Promise<void> {
    await notificator.markNotifications({
        ids,
    });

    const now = toTimestamp(new Date());
    data.value?.notifications.forEach((v) => (v.readAt = now));
}

watch(offset, async () => refresh());
watchDebounced(includeRead, async () => refresh(), { debounce: 500, maxWait: 1500 });

const { start: timeoutFn } = useTimeoutFn(() => (canSubmit.value = true), 400, { immediate: false });
const canSubmit = ref(true);
</script>

<template>
    <div class="py-2 pb-14">
        <div class="px-1 sm:px-2 lg:px-4">
            <div class="sm:flex sm:items-center">
                <div class="sm:flex-auto">
                    <form @submit.prevent="refresh()">
                        <div class="flex flex-row items-center gap-2 sm:mx-auto">
                            <div class="form-control flex-1">
                                <label for="search" class="block text-sm font-medium leading-6 text-neutral"
                                    >{{ $t('components.notifications.include_read') }}
                                </label>
                                <div class="relative flex items-center">
                                    <Switch
                                        v-model="includeRead"
                                        :class="[
                                            includeRead ? 'bg-primary-600' : 'bg-gray-200',
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
                                    :disabled="
                                        !canSubmit || data?.notifications === undefined || data?.notifications.length === 0
                                    "
                                    class="inline-flex rounded-md px-3 py-2 text-sm font-semibold text-neutral focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                                    :class="
                                        !canSubmit || data?.notifications === undefined || data?.notifications.length === 0
                                            ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                            : 'bg-primary-500 hover:bg-primary-400'
                                    "
                                    @click="markAllRead().finally(timeoutFn)"
                                >
                                    {{ $t('components.notifications.mark_all_read') }}
                                </button>
                            </div>
                        </div>
                    </form>
                </div>
            </div>
            <div class="mt-2 flow-root">
                <div class="-my-2 mx-0 overflow-x-auto">
                    <div class="inline-block min-w-full px-1 py-2 align-middle">
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
                        <template v-else>
                            <ul class="flex flex-col">
                                <li
                                    v-for="not in data?.notifications"
                                    :key="not.id"
                                    class="relative my-1 flex justify-between gap-x-6 rounded-lg bg-base-700 px-4 py-5 hover:bg-base-600 sm:px-6"
                                >
                                    <div class="flex flex-1 gap-x-4">
                                        <div class="min-w-0 flex-auto">
                                            <p class="py-2 pr-3 text-sm font-medium text-neutral">
                                                <template v-if="not.data && not.data.link">
                                                    <!-- @vue-ignore the route should be valid... at least in most cases -->
                                                    <NuxtLink :to="not.data.link.to" class="inline-flex items-center gap-1">
                                                        {{ $t(not.title!.key, not.title?.parameters ?? {}) }}
                                                        <LinkVariantIcon class="h-5 w-5" />
                                                    </NuxtLink>
                                                </template>
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
                                            <p class="mt-1 text-xs leading-5 text-gray-300">
                                                {{ $t('common.received') }}
                                                <GenericTime :value="not.createdAt" :ago="true" />
                                            </p>
                                            <div class="mt-1 flex items-center gap-x-1.5">
                                                <template v-if="not.readAt">
                                                    <p class="text-xs leading-5 text-gray-300">
                                                        {{ $t('common.read') }}
                                                        <GenericTime :value="not.readAt" :ago="true" />
                                                    </p>
                                                </template>
                                                <template v-else>
                                                    <div class="flex-none rounded-full bg-error-500/20 p-1">
                                                        <div class="h-1.5 w-1.5 rounded-full bg-error-500" />
                                                    </div>
                                                    <p class="text-xs leading-5 text-gray-300">
                                                        {{ $t('components.notifications.unread') }}
                                                    </p>
                                                </template>
                                            </div>
                                        </div>
                                    </div>

                                    <button
                                        v-if="!not.readAt"
                                        type="button"
                                        class="-mt-5 -mr-6 -mb-5 flex shrink items-center rounded-r-md px-1 py-1 text-sm font-semibold text-neutral focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-500"
                                        :class="
                                            !canSubmit
                                                ? 'disabled bg-base-500 hover:bg-base-400 focus-visible:outline-base-500'
                                                : 'bg-primary-500 hover:bg-primary-400'
                                        "
                                        :disabled="!canSubmit"
                                        @click="markRead(not.id).finally(timeoutFn)"
                                    >
                                        <span class="sr-only">{{ $t('components.notifications.mark_read') }}</span>
                                        <CheckIcon class="h-5 w-5 text-gray-300" aria-hidden="true" />
                                    </button>
                                    <span v-else class="h-4 w-4"></span>
                                </li>
                            </ul>

                            <TablePagination
                                v-if="!compact"
                                class="mt-2"
                                :pagination="data?.pagination"
                                :refresh="refresh"
                                @offset-change="offset = $event"
                            />
                        </template>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
