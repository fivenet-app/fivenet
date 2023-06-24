<script lang="ts" setup>
import Livemap from '~/components/livemap/Livemap.vue';
import Notificator from '~/components/partials/Notificator.vue';
import Sidebar from '~/components/partials/Sidebar.vue';
import NotificationProvider from '~/components/partials/notification/NotificationProvider.vue';

useHead({
    title: 'common.dispatch_center',
});
definePageMeta({
    title: 'common.dispatch_center',
    requiresAuth: true,
    permission: 'CentrumService.Stream',
    showQuickButtons: false,
});
</script>

<template>
    <NotificationProvider>
        <Sidebar>
            <Livemap />
        </Sidebar>
        <Notificator />
    </NotificationProvider>
</template>
