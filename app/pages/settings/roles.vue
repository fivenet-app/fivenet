<script lang="ts" setup>
import RoleList from '~/components/settings/roles/RoleList.vue';

useHead({
    title: 'pages.settings.roles.title',
});

definePageMeta({
    title: 'pages.settings.roles.title',
    requiresAuth: true,
    permission: 'settings.SettingsService/GetRoles',
});
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('pages.settings.roles.title')">
                <template #right>
                    <PartialsBackButton fallback-to="/settings" />
                </template>
            </UDashboardNavbar>

            <RoleList />
        </UDashboardPanel>
    </UDashboardPage>
</template>
