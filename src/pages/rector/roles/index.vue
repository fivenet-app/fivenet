<script lang="ts" setup>
import RolesList from '~/components/rector/roles/RolesList.vue';

useHead({
    title: 'pages.rector.roles.title',
});
definePageMeta({
    title: 'pages.rector.roles.title',
    requiresAuth: true,
    permission: 'RectorService.GetRoles',
});
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('pages.rector.roles.title')" />

            <RolesList />
        </UDashboardPanel>
    </UDashboardPage>
</template>
