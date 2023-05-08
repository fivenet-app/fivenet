<script lang="ts" setup>
import NotificationProvider from '~/components/partials/notification/NotificationProvider.vue';
import Sidebar from '~/components/partials/Sidebar.vue';
import Livemap from '~/components/livemap/Livemap.vue';
import Notificator from '~/components/partials/Notificator.vue';

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
</script>

<template>
    <NotificationProvider>
        <Sidebar>
            <Livemap />
        </Sidebar>
        <Notificator />
    </NotificationProvider>
</template>
