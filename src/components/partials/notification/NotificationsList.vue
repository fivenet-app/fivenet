<script lang="ts" setup>
import { z } from 'zod';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';

import GenericTime from '~/components/partials/elements/GenericTime.vue';
import { useNotificatorStore } from '~/store/notificator';
import { GetNotificationsResponse } from '~~/gen/ts/services/notificator/notificator';
import Pagination from '../Pagination.vue';

defineEmits<{
    (e: 'clicked'): void;
}>();

const { $grpc } = useNuxtApp();

const notificator = useNotificatorStore();

const schema = z.object({
    includeRead: z.boolean(),
});

type Schema = z.output<typeof schema>;

const query = reactive<Schema>({
    includeRead: false,
});

const page = ref(1);
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const { data, pending, refresh, error } = useLazyAsyncData(`notifications-${page.value}-${query.includeRead}`, () =>
    getNotifications(),
);

async function getNotifications(): Promise<GetNotificationsResponse> {
    try {
        const call = $grpc.getNotificatorClient().getNotifications({
            pagination: {
                offset: offset.value,
            },
            includeRead: query.includeRead,
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
    data.value?.notifications.forEach((v) => {
        if (ids.includes(v.id)) {
            v.readAt = now;
        }
    });
}

watch(offset, async () => refresh());
watchDebounced(query, async () => refresh(), { debounce: 500, maxWait: 1500 });

const { start: timeoutFn } = useTimeoutFn(() => (canSubmit.value = true), 400, { immediate: false });
const canSubmit = ref(true);
</script>

<template>
    <div class="sm:flex sm:items-center">
        <div class="sm:flex-auto">
            <UForm :schema="schema" :state="query" @submit="refresh()">
                <div class="flex flex-row items-center gap-2">
                    <UFormGroup name="search" :label="$t('components.notifications.include_read')" class="flex-1">
                        <UToggle v-model="query.includeRead">
                            <span class="sr-only">{{ $t('components.notifications.include_read') }}</span>
                        </UToggle>
                    </UFormGroup>

                    <UFormGroup class="flex-initial">
                        <UButton
                            :disabled="!canSubmit || data?.notifications === undefined || data?.notifications.length === 0"
                            @click="markAllRead().finally(timeoutFn)"
                        >
                            {{ $t('components.notifications.mark_all_read') }}
                        </UButton>
                    </UFormGroup>
                </div>
            </UForm>
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
                    icon="i-mdi-bell"
                />
                <template v-else>
                    <ul class="flex flex-col">
                        <li
                            v-for="not in data?.notifications"
                            :key="not.id"
                            class="hover:bg-background relative my-1 flex justify-between gap-x-6 rounded-lg bg-base-700 px-4 py-5 sm:px-6"
                        >
                            <div class="flex flex-1 gap-x-2">
                                <div class="min-w-0 flex-auto">
                                    <p class="py-2 pr-3 text-sm font-medium">
                                        <template v-if="not.data && not.data.link">
                                            <!-- @vue-ignore the route should be valid... at least in most cases -->
                                            <UButton
                                                variant="link"
                                                :to="not.data.link.to"
                                                trailing-icon="i-mdi-link-variant"
                                                class="inline-flex items-center gap-1"
                                                @click="$emit('clicked')"
                                            >
                                                {{ $t(not.title!.key, not.title?.parameters ?? {}) }}
                                            </UButton>
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

                            <div class="flex items-center gap-x-2">
                                <div class="hidden sm:flex sm:flex-col sm:items-end">
                                    <p class="mt-1 text-xs leading-5">
                                        {{ $t('common.received') }}
                                        <GenericTime :value="not.createdAt" :ago="true" />
                                    </p>
                                    <div class="mt-1 flex items-center gap-x-1.5">
                                        <template v-if="not.readAt">
                                            <p class="text-xs leading-5">
                                                {{ $t('common.read') }}
                                                <GenericTime :value="not.readAt" :ago="true" />
                                            </p>
                                        </template>
                                        <template v-else>
                                            <div class="flex-none rounded-full bg-error-500/20 p-1">
                                                <div class="size-1.5 rounded-full bg-error-500" />
                                            </div>
                                            <p class="text-xs leading-5">
                                                {{ $t('components.notifications.unread') }}
                                            </p>
                                        </template>
                                    </div>
                                </div>
                            </div>

                            <div class="-my-5 -mr-6 flex">
                                <UButton
                                    v-if="!not.readAt"
                                    class="flex shrink items-center rounded-r-md p-1 text-sm font-semibold focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2"
                                    :disabled="!canSubmit"
                                    icon="i-mdi-check"
                                    @click="markRead(not.id).finally(timeoutFn)"
                                >
                                    <span class="sr-only">{{ $t('components.notifications.mark_read') }}</span>
                                </UButton>
                                <span v-else class="size-5"></span>
                            </div>
                        </li>
                    </ul>
                </template>

                <Pagination v-model="page" :pagination="data?.pagination" />
            </div>
        </div>
    </div>
</template>
