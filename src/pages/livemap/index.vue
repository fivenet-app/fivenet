<script lang="ts" setup>
import Livemap from '~/components/livemap/Livemap.vue';
import Sidebar from '~/components/partials/Sidebar.vue';

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
