<script lang="ts" setup>
import { NotificationProvider } from '~/components/partials/notification';
import Sidebar from '~/components/partials/Sidebar.vue';
import Livemap from '~/components/livemap/Livemap.vue';

useHead({
    title: 'Livemap',
});
definePageMeta({
    title: 'Livemap',
    requiresAuth: true,
    permission: 'LivemapperService.Stream',
    layout: 'blank',
});
</script>

<template>
    <NotificationProvider>
        <Sidebar>
            <Livemap />
        </Sidebar>
    </NotificationProvider>
</template>
