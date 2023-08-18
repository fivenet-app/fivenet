<script lang="ts" setup>
import LivemapHolder from '~/components/livemap/LivemapHolder.vue';
import CommandPalette from '~/components/partials/CommandPalette.vue';
import NotificationProvider from '~/components/partials/notification/NotificationProvider.vue';
import Notificator from '~/components/partials/notification/Notificator.vue';
import Sidebar from '~/components/partials/sidebar/Sidebar.vue';
import { useAuthStore } from '~/store/auth';

useHead({
    title: 'common.livemap',
});
definePageMeta({
    title: 'common.livemap',
    layout: 'blank',
    requiresAuth: true,
    permission: 'LivemapperService.Stream',
    showQuickButtons: false,
});

const authStore = useAuthStore();
const { activeChar } = storeToRefs(authStore);
</script>

<template>
    <NotificationProvider>
        <Sidebar>
            <div class="h-full">
                <LivemapHolder />
            </div>
        </Sidebar>
        <CommandPalette v-if="activeChar" />
        <Notificator />
    </NotificationProvider>
</template>
