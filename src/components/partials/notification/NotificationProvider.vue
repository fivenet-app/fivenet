<script lang="ts" setup>
import { notificationTypeToColor, notificationTypeToIcon } from '~/components/partials/notification/helpers';
import { useAuthStore } from '~/store/auth';
import { useNotificatorStore } from '~/store/notificator';

const { t } = useI18n();

const { timeouts } = useAppConfig();

const authStore = useAuthStore();
const { username, activeChar } = storeToRefs(authStore);

const notificatorStore = useNotificatorStore();
const { abort, notifications } = storeToRefs(notificatorStore);
const { startStream, stopStream } = notificatorStore;

async function toggleStream(): Promise<void> {
    // Only stream notifications when a user is logged in and has a character selected
    if (username.value !== null && activeChar.value !== null) {
        startStream();
    } else if (abort.value !== undefined) {
        await stopStream();
        notificatorStore.$reset();
    }
}

watch(username, async () => toggleStream());
watch(activeChar, async () => toggleStream());

onMounted(async () => toggleStream());

onBeforeUnmount(async () => stopStream());

const toast = useToast();

watchArray(
    notifications,
    (_, _0, added) => {
        added.forEach((notification) => {
            toast.add({
                id: notification.id,
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
    },
    { deep: true },
);
</script>

<template>
    <div></div>
</template>
