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
const { notifications } = storeToRefs(notificationsStore);

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

async function toggleStream(): Promise<void> {
    // Only stream notifications when a user is logged in and has a character selected
    if (username.value !== null && activeChar.value !== null && webSocket.status.value === 'OPEN') {
        start();

        try {
            startStream();
        } catch (e) {
            logger.error('exception during notification stream', e);
        }
    } else {
        pause();
        stop();
        await stopStream(true);

        notificationsStore.reset();
    }
}

watch([username, activeChar, webSocket.status], async () => toggleStream());

const toast = useToast();

function createNotifications(notifications: Notification[]): void {
    notifications.forEach((notification) => {
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
                }
                if (notification.callback) {
                    notification.callback();
                }
            },
        });
    });
}

const handleNotificationAdded = (notification: Notification): void => createNotifications([notification]);

onMounted(async () => await toggleStream());

onUnmounted(async () => await stopStream(true));

onMounted(() => {
    createNotifications(notifications.value);
    notificationToastEvents.on('add', handleNotificationAdded);
});

onUnmounted(() => {
    notificationToastEvents.off('add', handleNotificationAdded);
});
</script>

<template>
    <div></div>
</template>
