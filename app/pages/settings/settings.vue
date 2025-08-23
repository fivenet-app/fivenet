<script lang="ts" setup>
import AppConfig from '~/components/settings/AppConfig.vue';

useHead({
    title: 'pages.settings.settings.title',
});

definePageMeta({
    title: 'pages.settings.settings.title',
    requiresAuth: true,
    permission: 'Superuser/Superuser',
});
</script>

<template>
    <UDashboardPanel>
        <AppConfig />
    </UDashboardPanel>
</template>
