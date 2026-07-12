<script lang="ts" setup>
import { v4 as uuidv4 } from 'uuid';
import { notificationTypeToColor, notificationTypeToIcon } from '~/components/notifications/helpers';
import { useGRPCWebsocketTransport } from '~/composables/grpcws';
import { notificationToastEvents } from '~/composables/useNotificationToasts';
import { useCalendarStore } from '~/stores/calendar';
import { useMailerStore } from '~/stores/mailer';
import { logger } from '~/stores/notifications';
import { useSettingsStore } from '~/stores/settings';
import type { Notification } from '~/types/notifications';

const { t } = useI18n();

const { timeouts } = useAppConfig();

const { can, username, activeChar } = useAuth();

const { webSocket } = useGRPCWebsocketTransport();

const notificationsStore = useNotificationsStore();
const { startStream, stopStream } = notificationsStore;
const { notifications, abort, ready } = storeToRefs(notificationsStore);

const settingsStore = useSettingsStore();
const { calendar } = storeToRefs(settingsStore);

const calendarStore = useCalendarStore();
const mailerStore = useMailerStore();

async function checkAppointments(): Promise<void> {
    try {
        await calendarStore.checkAppointments();
    } catch (e) {
        logger.error('Exception during check appointments call', e);
    }
}

const { pause, resume } = useIntervalFn(
    async () => {
        pause();

        if (calendar.value.reminderTimes.length > 0) checkAppointments();

        resume();
    },
    (300 + randomNumber(1, 17)) * 1000,
    { immediate: false },
);

const { start, stop } = useTimeoutFn(
    async () => {
        if (!activeChar.value) {
            stop();
            pause();
            return;
        }

        if (can('mailer.MailerService/ListEmails').value) await mailerStore.checkEmails();

        if (calendar.value.reminderTimes.length > 0) checkAppointments();

        resume();
    },
    randomNumber(1, 7) * 1000,
    {
        immediate: false,
    },
);

async function syncStream(previousActiveChar: typeof activeChar.value): Promise<void> {
    const sessionReady = username.value !== null && webSocket.status.value === 'OPEN';

    if (!sessionReady) {
        pause();
        stop();
        await stopStream(true);
        notificationsStore.reset();
        return;
    }

    const shouldRestart = ready.value && previousActiveChar !== activeChar.value;
    if (shouldRestart) {
        pause();
        stop();
        await stopStream(true);
    }

    if (!ready.value && abort.value !== undefined) {
        await stopStream(true);
    }

    if (activeChar.value !== null) {
        start();
    } else {
        pause();
        stop();
    }

    if (!ready.value) {
        try {
            await startStream();
        } catch (e) {
            logger.error('exception during notification stream', e);
        }
    }
}

watch([username, activeChar, webSocket.status], async (_values, oldValues) => {
    const [, oldActiveChar] = oldValues;
    await syncStream(oldActiveChar);
});

const toast = useToast();
const seenNotificationIds = new Set<number>();

function markNotificationSeen(notification: Notification): void {
    if (notification.id !== undefined) seenNotificationIds.add(notification.id);
}

function createNotification(notification: Notification): void {
    if (notification.id !== undefined && seenNotificationIds.has(notification.id)) return;

    markNotificationSeen(notification);
    toast.add({
        id: notification.id?.toString() ?? uuidv4(),
        title:
            typeof notification.title === 'string'
                ? notification.title
                : t(notification.title.key, notification.title.parameters ?? {}),
        description:
            typeof notification.description === 'string'
                ? notification.description
                : t(notification.description.key, notification.description.parameters ?? {}),
        icon: notificationTypeToIcon(notification.type),
        color: notificationTypeToColor(notification.type),
        duration: notification.duration ?? timeouts.notification,
        actions: notification.actions?.map((action) => ({
            label: t(action.label.key, action.label.parameters ?? {}),
            to: action.to,
            external: action.external,
            onClick: action.onClick,
        })),
        'onUpdate:open': () => {
            if (notification.id) {
                notificationsStore.remove(notification.id);
                seenNotificationIds.delete(notification.id);
            }
            if (notification.callback) {
                notification.callback();
            }
        },
    });
}

function createNotifications(notifications: Notification[]): void {
    notifications.forEach((notification) => createNotification(notification));
}

const handleNotificationAdded = (notification: Notification): void => createNotification(notification);

onMounted(async () => await syncStream(activeChar.value));

onUnmounted(async () => await stopStream(true));

onMounted(() => {
    notificationToastEvents.on('add', handleNotificationAdded);
    createNotifications(notifications.value);
});

onUnmounted(() => {
    notificationToastEvents.off('add', handleNotificationAdded);
    seenNotificationIds.clear();
});
</script>

<template>
    <div></div>
</template>
