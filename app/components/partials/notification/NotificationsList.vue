<script lang="ts" setup>
import { z } from 'zod';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import GenericTime from '~/components/partials/elements/GenericTime.vue';
import Pagination from '~/components/partials/Pagination.vue';
import { useNotificatorStore } from '~/store/notificator';
import { NotificationCategory } from '~~/gen/ts/resources/notifications/notifications';
import { GetNotificationsResponse } from '~~/gen/ts/services/notificator/notificator';
import { notificationCategoryToIcon } from './helpers';

defineEmits<{
    (e: 'clicked'): void;
}>();

const { t } = useI18n();

const notificator = useNotificatorStore();

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

const page = ref(1);
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
    await notificator.markNotifications({
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
                    <UFormGroup name="includeRead" :label="$t('components.notifications.include_read')" class="flex-initial">
                        <UToggle v-model="query.includeRead">
                            <span class="sr-only">{{ $t('components.notifications.include_read') }}</span>
                        </UToggle>
                    </UFormGroup>

                    <UFormGroup name="categories" :label="$t('common.category', 2)" class="flex-1">
                        <USelectMenu
                            v-model="query.categories"
                            multiple
                            name="categories"
                            :options="categories"
                            option-attribute="label"
                            value-attribute="chip"
                            :searchable-placeholder="$t('common.search_field')"
                            @focusin="focusTablet(true)"
                            @focusout="focusTablet(false)"
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
                    </UFormGroup>

                    <UFormGroup label="&nbsp;" class="flex-initial">
                        <UButton
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

    <UDashboardPanelContent class="p-0">
        <DataPendingBlock v-if="loading" :message="$t('common.loading', [$t('common.notification', 2)])" />
        <DataErrorBlock
            v-else-if="error"
            :title="$t('common.unable_to_load', [$t('common.notification', 2)])"
            :retry="refresh"
        />
        <DataNoDataBlock v-else-if="data?.notifications.length === 0" :type="$t('common.notification', 2)" icon="i-mdi-bell" />

        <div v-else class="relative overflow-x-auto">
            <ul role="list" class="my-1 flex flex-col">
                <li
                    v-for="not in data?.notifications"
                    :key="not.id"
                    class="relative my-1 flex justify-between gap-x-6 px-4 py-5 sm:px-6"
                >
                    <div class="flex flex-1 gap-x-2">
                        <div class="min-w-0 flex-auto">
                            <p class="py-2 pr-3 text-sm font-medium">
                                <template v-if="not.data && not.data.link">
                                    <!-- @vue-ignore the route should be valid... at least in most cases -->
                                    <UButton
                                        variant="link"
                                        :to="not.data.link.to"
                                        :icon="notificationCategoryToIcon(not.category)"
                                        trailing-icon="i-mdi-link-variant"
                                        class="inline-flex items-center gap-1"
                                        @click="
                                            markRead(not.id);
                                            $emit('clicked');
                                        "
                                    >
                                        {{ $t(not.title!.key, not.title?.parameters ?? {}) }}
                                    </UButton>
                                </template>
                                <template v-else>
                                    {{ $t(not.title!.key, not.title?.parameters ?? {}) }}
                                </template>
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
        </div>

        <Pagination v-model="page" :pagination="data?.pagination" :loading="loading" :refresh="refresh" />
    </UDashboardPanelContent>
</template>
