<script lang="ts" setup>
import RolesList from '~/components/settings/roles/RolesList.vue';

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

            <RolesList />
        </UDashboardPanel>
    </UDashboardPage>
</template>
