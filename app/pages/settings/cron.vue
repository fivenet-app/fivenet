<script lang="ts" setup>
import CronList from '~/components/settings/CronList.vue';

useHead({
    title: 'pages.settings.cron.title',
});

definePageMeta({
    title: 'pages.settings.cron.title',
    requiresAuth: true,
    permission: 'Superuser/Superuser',
});
</script>

<template>
    <UDashboardPanel>
        <CronList />
    </UDashboardPanel>
</template>
