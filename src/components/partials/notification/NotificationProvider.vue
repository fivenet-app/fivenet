<script lang="ts" setup>
import { useAuthStore } from '~/store/auth';
import { useNotificatorStore } from '~/store/notificator';
import { notificationTypeToIcon } from './helpers';

const authStore = useAuthStore();
const { accessToken, activeChar } = storeToRefs(authStore);

const notifications = useNotificatorStore();
const { getNotifications } = storeToRefs(notifications);
const { startStream, stopStream } = notifications;

async function toggleStream(): Promise<void> {
    // Only stream notifications when a user is logged in and has a character selected
    if (accessToken.value !== null && activeChar.value !== null) {
        return startStream();
    } else {
        await stopStream();
        notifications.$reset();
    }
}

watch(accessToken, async () => toggleStream());
watch(activeChar, async () => toggleStream());

onBeforeMount(async () => toggleStream());

onBeforeUnmount(async () => stopStream());
</script>

<template>
    <div class="hidden">
        <UNotification
            v-for="notification in getNotifications.filter((n) => n.position === undefined || n.position === 'top-right')"
            :id="notification.id"
            :key="notification.id"
            :title="$t(notification.title.key, notification.title.parameters ?? {})"
            :description="$t(notification.content.key, notification.content.parameters ?? {})"
            :icon="notificationTypeToIcon(notification.type)"
            :timeout="3500"
        />
    </div>
</template>
