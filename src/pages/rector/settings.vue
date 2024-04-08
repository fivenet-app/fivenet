<script lang="ts" setup>
import AppConfig from '~/components/rector/AppConfig.vue';

useHead({
    title: 'pages.rector.settings.title',
});
definePageMeta({
    title: 'pages.rector.settings.title',
    requiresAuth: true,
    permission: 'SuperUser',
});
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <AppConfig />
        </UDashboardPanel>
    </UDashboardPage>
</template>
