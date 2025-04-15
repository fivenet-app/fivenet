<script lang="ts" setup>
import CronList from '~/components/rector/CronList.vue';

useHead({
    title: 'pages.rector.cron.title',
});

definePageMeta({
    title: 'pages.rector.cron.title',
    requiresAuth: true,
    permission: 'SuperUser',
});
</script>

<template>

        <UDashboardPanel grow>
            <CronList />
        </UDashboardPanel>

</template>
