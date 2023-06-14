<script lang="ts" setup>
import { useNotificationsStore } from '~/store/notifications';
import NotificationItem from './NotificationItem.vue';

const notificationsStore = useNotificationsStore();
const { getNotifications } = storeToRefs(notificationsStore);
</script>

<template>
    <div>
        <!-- Global notification live region, render this permanently at the end of the document -->
        <div
            aria-live="assertive"
            class="pointer-events-none fixed inset-0 flex items-end px-4 py-6 sm:items-start sm:p-6 z-50"
        >
            <div class="flex w-full flex-col items-center space-y-4 sm:items-end">
                <NotificationItem
                    :key="notification.id?.toString()"
                    :notification="notification"
                    :class="idx > 0 ? 'mt-4' : ''"
                    v-for="(notification, idx) in getNotifications"
                />
            </div>
        </div>
        <slot />
    </div>
</template>
