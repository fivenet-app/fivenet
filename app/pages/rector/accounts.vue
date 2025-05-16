<script lang="ts" setup>
import AccountsList from '~/components/rector/accounts/AccountsList.vue';

useHead({
    title: 'pages.rector.accounts.title',
});

definePageMeta({
    title: 'pages.rector.accounts.title',
    requiresAuth: true,
    permission: 'SuperUser',
});
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <AccountsList />
        </UDashboardPanel>
    </UDashboardPage>
</template>
