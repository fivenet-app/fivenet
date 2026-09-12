<script lang="ts" setup>
import { z } from 'zod';
import DataErrorBlock from '~/components/partials/data/DataErrorBlock.vue';
import DataNoDataBlock from '~/components/partials/data/DataNoDataBlock.vue';
import DataPendingBlock from '~/components/partials/data/DataPendingBlock.vue';
import Pagination from '~/components/partials/Pagination.vue';
import type { Form } from '@nuxt/ui';
import { getNotificationsNotificationsClient } from '~~/gen/ts/clients';
import { NotificationCategory, type Notification } from '~~/gen/ts/resources/notifications/notifications';
import type { GetNotificationsResponse } from '~~/gen/ts/services/notifications/notifications';
import DNBToggle from './DNBToggle.vue';
import { notificationEntryComponent } from './entries/registry';

defineProps<{
    hideHeader?: boolean;
    hideFooter?: boolean;
    scrollable?: boolean;
}>();

const emit = defineEmits<{
    (e: 'clicked'): void;
}>();

const notificationCenter = useNotificationCenterStore();

const notificationsNotificationsClient = await getNotificationsNotificationsClient();

const categories: { mode: NotificationCategory }[] = [
    { mode: NotificationCategory.GENERAL },
    { mode: NotificationCategory.DOCUMENT },
    { mode: NotificationCategory.CALENDAR },
    { mode: NotificationCategory.JOBS },
    { mode: NotificationCategory.QUALIFICATIONS },
    { mode: NotificationCategory.MAILER },
    { mode: NotificationCategory.SYSTEM },
];

const scopes = ['all', 'unread', 'starred', 'archived'] as const;
type Scope = (typeof scopes)[number];

const schema = z.object({
    categories: z.enum(NotificationCategory).array().max(categories.length).default([]),
    page: pageNumberSchema,
});

type Schema = z.output<typeof schema>;

const query = reactive<Schema>(schema.parse({}));
const scope = ref<Scope>('all');
const formRef = useTemplateRef<Form<typeof schema>>('formRef');
const { validatedQuery, commitValidatedQuery } = useFormSearchValidation<typeof schema>(query, formRef);

const { data, status, refresh, error } = useAuthedLazyAsyncData(
    'userState',
    () => `notifications-${scope.value}-${validatedQuery.value.page}-${[...validatedQuery.value.categories].sort().join('-')}`,
    ({ signal }) => getNotifications(validatedQuery.value, signal),
);

async function getNotifications(values: Schema, signal: AbortSignal): Promise<GetNotificationsResponse> {
    try {
        const call = notificationsNotificationsClient.getNotifications(
            {
                pagination: {
                    offset: calculateOffset(values.page, data.value?.pagination),
                },
                includeRead: scope.value !== 'unread',
                categories: values.categories,
                kinds: [],
                includeArchived: scope.value === 'archived',
                archivedOnly: scope.value === 'archived',
                starredOnly: scope.value === 'starred',
            },
            { abort: signal },
        );

        const { response } = await call;

        return response;
    } catch (e) {
        handleGRPCError(e as RpcError);
        throw e;
    }
}

async function markAllRead(): Promise<void> {
    await notificationCenter.markAllRead();

    const now = toTimestamp(new Date());
    data.value?.notifications.forEach((v) => {
        if (!v.readAt) v.readAt = now;
    });
}

function setScope(value: Scope): void {
    scope.value = value;
    query.page = 1;
    commitValidatedQuery();
}

function navigate(notification: Notification): void {
    if (!notification.readAt) {
        void notificationCenter.updateState({ ids: [notification.id], unread: false });
    }
    emit('clicked');
}

function handleStateChanged(): void {
    void refresh();
}

defineShortcuts({
    shift_r: () =>
        scope.value !== 'archived' &&
        canSubmit.value &&
        data.value?.notifications !== undefined &&
        data.value?.notifications.length > 0 &&
        markAllRead().finally(timeoutFn),
});

const { start: timeoutFn } = useTimeoutFn(() => (canSubmit.value = true), 400, { immediate: false });
const canSubmit = ref<boolean>(true);
</script>

<template>
    <UDashboardPanel
        :ui="{
            root: 'pb-(--page-content-bottom-offset)',
            body: 'p-0 sm:p-0 gap-0 sm:gap-0' + (scrollable ? ' overflow-y-auto' : ''),
        }"
    >
        <template #header>
            <UDashboardNavbar v-if="!hideHeader" :title="$t('components.notifications.title')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <DNBToggle />
                </template>
            </UDashboardNavbar>

            <UDashboardToolbar>
                <UTabs
                    class="min-w-0 flex-1"
                    :model-value="scope"
                    :items="scopes.map((value) => ({ value, label: $t(`components.notifications.${value}`) }))"
                    @update:model-value="($event) => setScope($event as Scope)"
                />
            </UDashboardToolbar>

            <UDashboardToolbar>
                <template #default>
                    <UForm ref="formRef" class="my-2 flex-1" :schema="schema" :state="query" @submit="commitValidatedQuery">
                        <div class="flex flex-row gap-2">
                            <UFormField class="flex-1" name="categories" :label="$t('common.category', 2)">
                                <ClientOnly>
                                    <USelectMenu
                                        v-model="query.categories"
                                        class="w-full"
                                        multiple
                                        name="categories"
                                        :items="categories"
                                        value-key="mode"
                                        :search-input="{ placeholder: $t('common.search_field') }"
                                    >
                                        <template #default>
                                            {{
                                                query.categories.length === 0 || query.categories.length === categories.length
                                                    ? $t('components.notifications.all_categories')
                                                    : query.categories
                                                          .map((c) =>
                                                              $t(
                                                                  `enums.notifications.NotificationCategory.${NotificationCategory[c]}`,
                                                              ),
                                                          )
                                                          .join(', ')
                                            }}
                                        </template>
                                        <template #item-label="{ item }">
                                            {{
                                                $t(
                                                    `enums.notifications.NotificationCategory.${NotificationCategory[item.mode ?? 0]}`,
                                                )
                                            }}
                                        </template>
                                    </USelectMenu>
                                </ClientOnly>
                            </UFormField>

                            <UFormField v-if="scope !== 'archived'" class="flex-initial" label="&nbsp;">
                                <UTooltip :text="$t('components.notifications.mark_all_read')" :kbds="['shift', 'R']">
                                    <UButton
                                        icon="i-mdi-notification-clear-all"
                                        :disabled="
                                            !canSubmit || data?.notifications === undefined || data?.notifications.length === 0
                                        "
                                        :label="$t('components.notifications.mark_all_read')"
                                        @click="() => markAllRead().finally(timeoutFn)"
                                    />
                                </UTooltip>
                            </UFormField>
                        </div>
                    </UForm>
                </template>
            </UDashboardToolbar>
        </template>

        <template #body>
            <DataPendingBlock v-if="isRequestPending(status)" :message="$t('common.loading', [$t('common.notification', 2)])" />
            <DataErrorBlock
                v-else-if="error"
                :title="$t('common.unable_to_load', [$t('common.notification', 2)])"
                :error="error"
                :retry="refresh"
            />
            <DataNoDataBlock
                v-else-if="data?.notifications.length === 0"
                :type="$t('common.notification', 2)"
                icon="i-mdi-bell"
            />

            <ul v-else class="min-w-full divide-y divide-default" :class="scrollable ? 'pb-2' : ''" role="list">
                <component
                    :is="notificationEntryComponent(notification)"
                    v-for="notification in data?.notifications"
                    :key="notification.id"
                    :notification="notification"
                    full
                    @navigate="navigate"
                    @state-changed="handleStateChanged"
                />
            </ul>
        </template>

        <template #footer>
            <Pagination
                v-if="!hideFooter"
                v-model="query.page"
                :pagination="data?.pagination"
                :status="status"
                :refresh="refresh"
            />
        </template>
    </UDashboardPanel>
</template>
