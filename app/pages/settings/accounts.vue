<script lang="ts" setup>
import AccountsList from '~/components/settings/accounts/AccountsList.vue';

useHead({
    title: 'pages.settings.accounts.title',
});

definePageMeta({
    title: 'pages.settings.accounts.title',
    requiresAuth: true,
    permission: 'Superuser',
});
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <AccountsList />
        </UDashboardPanel>
    </UDashboardPage>
</template>
