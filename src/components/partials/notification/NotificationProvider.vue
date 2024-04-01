<script lang="ts" setup>
import { useAuthStore } from '~/store/auth';
import { useNotificatorStore } from '~/store/notificator';
import NotificationItem from '~/components/partials/notification/NotificationItem.vue';

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
    <!-- Global notification live region, render this permanently at the end of the document -->
    <div aria-live="assertive" class="pointer-events-none fixed inset-0 z-50 flex items-end px-4 py-16">
        <div class="flex w-full flex-col items-center space-y-4">
            <NotificationItem
                v-for="(notification, idx) in getNotifications.filter(
                    (n) => n.position === undefined || n.position === 'top-right',
                )"
                :key="notification.id"
                :notification="notification"
                :class="idx > 0 ? 'mb-6' : ''"
            />
        </div>
    </div>
</template>
