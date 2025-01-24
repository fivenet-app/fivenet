<script lang="ts" setup>
import { notificationTypeToColor, notificationTypeToIcon } from '~/components/partials/notification/helpers';
import type { Notification } from '~/composables/notifications';
import { useCalendarStore } from '~/store/calendar';
import { useMailerStore } from '~/store/mailer';
import { useNotificatorStore } from '~/store/notificator';
import { useSettingsStore } from '~/store/settings';

const { t } = useI18n();

const { timeouts } = useAppConfig();

const { can, username, activeChar } = useAuth();

const notificatorStore = useNotificatorStore();
const { notifications } = storeToRefs(notificatorStore);
const { startStream, stopStream } = notificatorStore;

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

        if (can('MailerService.ListEmails').value) {
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

        notificatorStore.resetData();
    }
}

watch(username, async () => toggleStream());
watch(activeChar, async () => toggleStream());

onMounted(async () => toggleStream());

onBeforeUnmount(async () => stopStream());

const toast = useToast();

function createNotifications(notifications: Notification[]): void {
    notifications.forEach((notification) => {
        toast.add({
            id: notification.id?.toString(),
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
                    notificatorStore.remove(notification.id);
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
