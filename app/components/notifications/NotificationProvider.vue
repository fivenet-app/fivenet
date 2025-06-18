<script lang="ts" setup>
import { v4 as uuidv4 } from 'uuid';
import { notificationTypeToColor, notificationTypeToIcon } from '~/components/notifications/helpers';
import { useCalendarStore } from '~/stores/calendar';
import { useMailerStore } from '~/stores/mailer';
import { useSettingsStore } from '~/stores/settings';
import type { Notification } from '~/utils/notifications';

const { t } = useI18n();

const { timeouts } = useAppConfig();

const { can, username, activeChar } = useAuth();

const notificationsStore = useNotificationsStore();
const { startStream, stopStream } = notificationsStore;
const { notifications } = storeToRefs(notificationsStore);

const settingsStore = useSettingsStore();
const { calendar } = storeToRefs(settingsStore);

const calendarStore = useCalendarStore();
const mailerStore = useMailerStore();

async function checkAppointments(): Promise<void> {
    await calendarStore.checkAppointments();
}

const { pause, resume } = useIntervalFn(
    async () => {
        pause();

        if (calendar.value.reminderTimes.length > 0) {
            checkAppointments();
        }

        resume();
    },
    (60 + randomNumber(2, 7)) * 1000,
    { immediate: false },
);

const { start, stop } = useTimeoutFn(
    async () => {
        if (!activeChar.value) {
            stop();
            pause();
            return;
        }

        if (can('mailer.MailerService/ListEmails').value) {
            await mailerStore.checkEmails();
        }

        if (calendar.value.reminderTimes.length > 0) {
            await checkAppointments();
        }

        resume();
    },
    randomNumber(1, 7) * 1000,
    {
        immediate: false,
    },
);

async function toggleStream(): Promise<void> {
    // Only stream notifications when a user is logged in and has a character selected
    if (username.value !== null && activeChar.value !== null) {
        start();

        try {
            startStream();
        } catch (e) {
            logger.error('exception during notification stream', e);
        }
    } else {
        pause();
        stop();
        await stopStream();

        notificationsStore.resetData();
    }
}

watch(username, async () => toggleStream());
watch(activeChar, async () => toggleStream());

onMounted(async () => await toggleStream());

onUnmounted(async () => await stopStream());

const toast = useToast();

function createNotifications(notifications: Notification[]): void {
    notifications.forEach((notification) => {
        toast.add({
            id: notification.id?.toString() ?? uuidv4(),
            title: t(notification.title.key, notification.title.parameters ?? {}),
            description: t(notification.description.key, notification.description.parameters ?? {}),
            icon: notificationTypeToIcon(notification.type),
            color: notificationTypeToColor(notification.type),
            timeout: notification.timeout ?? timeouts.notification,
            actions: notification.actions?.map((action) => ({
                ...action,
                label: t(action.label.key, action.label.parameters ?? {}),
            })),
            callback: () => {
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

watchArray(notifications, (_, _0, added) => createNotifications(added), { deep: true });
createNotifications(notifications.value);
</script>

<template>
    <div></div>
</template>
