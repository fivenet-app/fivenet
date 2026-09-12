<script lang="ts" setup>
import { notificationEntryComponent } from './entries/registry';
import { getNotificationsNotificationsClient } from '~~/gen/ts/clients';
import type { Notification } from '~~/gen/ts/resources/notifications/notifications';

const notificationCenter = useNotificationCenterStore();
const { unreadCount, newNotificationsAvailable } = storeToRefs(notificationCenter);

const notificationsStore = useNotificationsStore();
const { doNotDisturb } = storeToRefs(notificationsStore);

const open = ref(false);
const loading = ref(false);
const notifications = ref<Notification[]>([]);

async function refresh(): Promise<void> {
    loading.value = true;
    const revision = notificationCenter.inboxRevision;
    try {
        const client = await getNotificationsNotificationsClient();
        const { response } = await client.getNotifications({
            pagination: { offset: 0, pageSize: 6 },
            includeRead: true,
            categories: [],
            kinds: [],
        });
        notifications.value = response.notifications;
        notificationCenter.acknowledgeNewNotifications(revision);
    } catch (e) {
        handleGRPCError(e as RpcError);
    } finally {
        loading.value = false;
    }
}

watch(open, (value) => {
    if (value) void refresh();
});

async function markAllRead(): Promise<void> {
    await notificationCenter.markAllRead();
    const now = toTimestamp(new Date());
    notifications.value.forEach((notification) => {
        if (!notification.readAt) notification.readAt = now;
    });
}

function navigate(notification: Notification): void {
    if (!notification.readAt) {
        void notificationCenter.updateState({ ids: [notification.id], unread: false });
    }
    open.value = false;
}
</script>

<template>
    <UPopover
        v-model:open="open"
        :content="{ align: 'end', sideOffset: 8 }"
        :ui="{ content: 'w-[min(28rem,calc(100vw-1rem))] p-0' }"
    >
        <UChip :show="unreadCount > 0" color="error" inset :text="unreadCount <= 99 ? unreadCount : '99+'" size="xl">
            <UTooltip :text="$t('components.partials.sidebar_notifications')" :kbds="['B']">
                <UButton
                    color="neutral"
                    variant="ghost"
                    square
                    :icon="doNotDisturb ? 'i-mdi-notifications-paused' : 'i-mdi-bell-outline'"
                />
            </UTooltip>
        </UChip>

        <template #content>
            <section class="flex max-h-[min(34rem,calc(100vh-5rem))] flex-col">
                <header class="flex items-center gap-3 border-b border-default px-4 py-3">
                    <div class="min-w-0 flex-1">
                        <h2 class="font-semibold text-highlighted">{{ $t('components.notifications.title') }}</h2>
                        <p class="text-xs text-muted">{{ unreadCount }} {{ $t('components.notifications.unread') }}</p>
                    </div>
                    <UTooltip :text="$t('components.notifications.mark_all_read')">
                        <UButton
                            color="neutral"
                            variant="ghost"
                            size="sm"
                            icon="i-mdi-notification-clear-all"
                            :disabled="unreadCount === 0"
                            @click="markAllRead"
                        />
                    </UTooltip>
                </header>

                <div v-if="newNotificationsAvailable && !loading" class="border-b border-default p-2">
                    <UButton
                        block
                        color="primary"
                        variant="soft"
                        icon="i-mdi-refresh"
                        :label="$t('components.notifications.new_notifications')"
                        @click="refresh"
                    />
                </div>

                <div v-if="loading" class="p-6 text-center text-sm text-muted">{{ $t('common.loading') }}</div>
                <div v-else-if="notifications.length === 0" class="p-6 text-center text-sm text-muted">
                    {{ $t('common.empty') }}
                </div>
                <ul v-else class="divide-y divide-default overflow-y-auto" role="list">
                    <component
                        :is="notificationEntryComponent(notification)"
                        v-for="notification in notifications"
                        :key="notification.id"
                        :notification="notification"
                        @navigate="navigate"
                        @state-changed="refresh"
                    />
                </ul>

                <footer class="border-t border-default p-2">
                    <UButton
                        block
                        color="neutral"
                        variant="ghost"
                        :label="$t('components.partials.sidebar_notifications')"
                        @click="
                            navigateTo('/notifications');
                            open = false;
                        "
                    />
                </footer>
            </section>
        </template>
    </UPopover>
</template>
