<script lang="ts" setup>
import AccountList from '~/components/settings/accounts/AccountList.vue';

useHead({
    title: 'pages.settings.accounts.title',
});

definePageMeta({
    title: 'pages.settings.accounts.title',
    requiresAuth: true,
    permission: 'Superuser/Superuser',
});
</script>

<template>
    <UDashboardPanel>
        <AccountList />
    </UDashboardPanel>
</template>
