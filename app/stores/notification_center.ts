import { defineStore } from 'pinia';
import { getNotificationsNotificationsClient } from '~~/gen/ts/clients';
import type { MarkNotificationsRequest, UpdateNotificationStateRequest } from '~~/gen/ts/services/notifications/notifications';

/**
 * Durable inbox state. Transient action feedback remains in the legacy
 * notifications store until its callers can be migrated to a toast-specific
 * composable.
 */
export const useNotificationCenterStore = defineStore('notification-center', () => {
    const unreadCount = ref(0);
    const newNotificationsAvailable = ref(false);
    const inboxRevision = ref(0);

    function setUnreadCount(count: number): void {
        unreadCount.value = Math.max(0, count);
    }

    function reset(): void {
        unreadCount.value = 0;
        newNotificationsAvailable.value = false;
        inboxRevision.value = 0;
    }

    function notifyNewNotification(): void {
        inboxRevision.value++;
        newNotificationsAvailable.value = true;
    }

    function acknowledgeNewNotifications(revision: number): void {
        if (revision === inboxRevision.value) {
            newNotificationsAvailable.value = false;
        }
    }

    async function updateState(req: UpdateNotificationStateRequest): Promise<void> {
        try {
            const client = await getNotificationsNotificationsClient();
            const { response } = await client.updateNotificationState(req);
            setUnreadCount(response.unreadCount);
        } catch (e) {
            handleGRPCError(e as RpcError);
            throw e;
        }
    }

    async function markAllRead(): Promise<void> {
        const req: MarkNotificationsRequest = { unread: false, ids: [], all: true };
        try {
            const client = await getNotificationsNotificationsClient();
            const { response } = await client.markNotifications(req);
            setUnreadCount(response.unreadCount);
        } catch (e) {
            handleGRPCError(e as RpcError);
            throw e;
        }
    }

    return {
        unreadCount,
        newNotificationsAvailable,
        inboxRevision,
        setUnreadCount,
        reset,
        notifyNewNotification,
        acknowledgeNewNotifications,
        updateState,
        markAllRead,
    };
});
