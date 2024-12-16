<script lang="ts" setup>
import { z } from 'zod';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import Pagination from '~/components/partials/Pagination.vue';
import { useNotificatorStore } from '~/store/notificator';
import { NotificationCategory } from '~~/gen/ts/resources/notifications/notifications';
import type { GetNotificationsResponse } from '~~/gen/ts/services/notificator/notificator';
import { notificationCategoryToIcon } from './helpers';

defineEmits<{
    (e: 'clicked'): void;
}>();

const { t } = useI18n();

const notifications = useNotificatorStore();

const categories: { mode: NotificationCategory }[] = [
    { mode: NotificationCategory.GENERAL },
    { mode: NotificationCategory.DOCUMENT },
    { mode: NotificationCategory.CALENDAR },
];

const schema = z.object({
    includeRead: z.boolean(),
    categories: z.nativeEnum(NotificationCategory).array().max(4),
});

type Schema = z.output<typeof schema>;

const query = reactive<Schema>({
    includeRead: false,
    categories: [...categories.map((c) => c.mode)],
});

const page = useRouteQuery('page', '1', { transform: Number });
const offset = computed(() => (data.value?.pagination?.pageSize ? data.value?.pagination?.pageSize * (page.value - 1) : 0));

const {
    data,
    pending: loading,
    refresh,
    error,
} = useLazyAsyncData(`notifications-${page.value}-${query.includeRead}`, () => getNotifications());

async function getNotifications(): Promise<GetNotificationsResponse> {
    try {
        const call = getGRPCNotificatorClient().getNotifications({
            pagination: {
                offset: offset.value,
            },
            includeRead: query.includeRead,
            categories: query.categories,
        });

        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function markAllRead(): Promise<void> {
    await notifications.markNotifications({
        ids: [],
        all: true,
    });

    data.value?.notifications.forEach(
        (v) =>
            (v.readAt = {
                timestamp: undefined,
            }),
    );
}

async function markRead(...ids: string[]): Promise<void> {
    await notifications.markNotifications({
        ids: ids,
    });

    const now = toTimestamp(new Date());
    data.value?.notifications.forEach((v) => {
        if (ids.includes(v.id)) {
            v.readAt = now;
        }
    });
}

function notificationCategoriesToLabel(categories: NotificationCategory[]): string {
    return categories.map((c) => t(`enums.notifications.NotificationCategory.${NotificationCategory[c ?? 0]}`)).join(', ');
}

watch(offset, async () => refresh());
watchDebounced(query, async () => refresh(), { debounce: 500, maxWait: 1500 });

const { start: timeoutFn } = useTimeoutFn(() => (canSubmit.value = true), 400, { immediate: false });
const canSubmit = ref(true);
</script>

<template>
    <UDashboardToolbar>
        <template #default>
            <UForm :schema="schema" :state="query" class="w-full" @submit="refresh()">
                <div class="flex flex-row gap-2">
                    <UFormGroup
                        name="includeRead"
                        :label="$t('components.notifications.include_read')"
                        class="flex flex-initial flex-col"
                        :ui="{ container: 'flex-1 flex' }"
                    >
                        <div class="flex flex-1 items-center">
                            <UToggle v-model="query.includeRead">
                                <span class="sr-only">{{ $t('components.notifications.include_read') }}</span>
                            </UToggle>
                        </div>
                    </UFormGroup>

                    <UFormGroup name="categories" :label="$t('common.category', 2)" class="flex-1">
                        <ClientOnly>
                            <USelectMenu
                                v-model="query.categories"
                                multiple
                                name="categories"
                                :options="categories"
                                option-attribute="label"
                                value-attribute="chip"
                                :searchable-placeholder="$t('common.search_field')"
                            >
                                <template #label>
                                    <template v-if="query.categories">
                                        <span class="truncate">{{ notificationCategoriesToLabel(query.categories) }}</span>
                                    </template>
                                </template>
                                <template #option="{ option }">
                                    <span class="truncate">{{
                                        $t(`enums.notifications.NotificationCategory.${NotificationCategory[option.mode ?? 0]}`)
                                    }}</span>
                                </template>
                            </USelectMenu>
                        </ClientOnly>
                    </UFormGroup>

                    <UFormGroup label="&nbsp;" class="flex-initial">
                        <UButton
                            icon="i-mdi-notification-clear-all"
                            :disabled="!canSubmit || data?.notifications === undefined || data?.notifications.length === 0"
                            @click="markAllRead().finally(timeoutFn)"
                        >
                            {{ $t('components.notifications.mark_all_read') }}
                        </UButton>
                    </UFormGroup>
                </div>
            </UForm>
        </template>
    </UDashboardToolbar>

    <UDashboardPanelContent class="p-0 sm:pb-0">
        <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.notification', 2)])" />
        <DataErrorBlock
            v-else-if="error"
            :title="$t('common.unable_to_load', [$t('common.notification', 2)])"
            :error="error"
            :retry="refresh"
        />
        <DataNoDataBlock v-else-if="data?.notifications.length === 0" :type="$t('common.notification', 2)" icon="i-mdi-bell" />

        <div v-else class="relative h-full overflow-x-auto">
            <ul role="list" class="flex flex-col divide-y divide-gray-100 dark:divide-gray-800">
                <li
                    v-for="not in data?.notifications"
                    :key="not.id"
                    class="hover:border-primary-500/25 dark:hover:border-primary-400/25 hover:bg-primary-100/50 dark:hover:bg-primary-900/10 relative flex justify-between gap-x-2 border-white px-4 py-5 sm:px-6 dark:border-gray-900"
                >
                    <div class="flex flex-1 gap-x-2">
                        <div class="min-w-0 flex-auto">
                            <p class="px-2 py-2 text-sm font-semibold">
                                <UButton
                                    v-if="not.data && not.data.link"
                                    variant="link"
                                    :padded="false"
                                    :to="not.data.link.to"
                                    :icon="notificationCategoryToIcon(not.category)"
                                    trailing-icon="i-mdi-link-variant"
                                    @click="
                                        markRead(not.id);
                                        $emit('clicked');
                                    "
                                >
                                    {{ $t(not.title!.key, not.title?.parameters ?? {}) }}
                                </UButton>
                                <span v-else>
                                    {{ $t(not.title!.key, not.title?.parameters ?? {}) }}
                                </span>
                            </p>

                            <p class="flex text-sm leading-4 text-gray-200">
                                {{ $t(not.content!.key, not.content?.parameters ?? {}) }}
                            </p>
                        </div>
                    </div>

                    <div class="flex items-center gap-x-2">
                        <div class="hidden sm:flex sm:flex-col sm:items-end">
                            <p class="mt-1 text-xs leading-5">
                                {{ $t('common.received') }}
                                <GenericTime :value="not.createdAt" ago />
                            </p>
                            <div class="mt-1 flex items-center gap-x-2">
                                <template v-if="not.readAt">
                                    <p class="text-xs leading-5">
                                        {{ $t('common.read') }}
                                        <GenericTime :value="not.readAt" ago />
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
        </div>
    </UDashboardPanelContent>

    <Pagination v-model="page" :pagination="data?.pagination" :loading="loading" :refresh="refresh" />
</template>
