<script lang="ts" setup>
import AppConfig from '~/components/rector/AppConfig.vue';

useHead({
    title: 'pages.rector.settings.title',
});
definePageMeta({
    title: 'pages.rector.settings.title',
    requiresAuth: true,
    permission: 'SuperUser',
    showQuickButtons: false,
});
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('pages.rector.settings.title')"> </UDashboardNavbar>

            <AppConfig />
        </UDashboardPanel>
    </UDashboardPage>
</template>
