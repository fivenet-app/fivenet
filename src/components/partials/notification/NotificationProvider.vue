<script lang="ts" setup>
import { useNotificatorStore } from '~/store/notificator';
import NotificationItem from '~/components/partials/notification/NotificationItem.vue';

const notifications = useNotificatorStore();
const { getNotifications } = storeToRefs(notifications);
</script>

<template>
    <div>
        <slot />

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
    </div>
</template>
