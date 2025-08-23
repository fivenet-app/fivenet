<script lang="ts" setup>
import UserSettingsPanel from '~/components/auth/account/UserSettingsPanel.vue';

useHead({
    title: 'components.auth.UserSettingsPanel.title',
});

definePageMeta({
    title: 'components.auth.UserSettingsPanel.title',
    requiresAuth: true,
    authTokenOnly: true,
    showCookieOptions: true,
});
</script>

<template>
    <UDashboardPanel>
        <template #header>
            <UDashboardNavbar :title="$t('components.auth.UserSettingsPanel.title')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <PartialsBackButton fallback-to="/overview" />
                </template>
            </UDashboardNavbar>
        </template>

        <template #body>
            <UserSettingsPanel />
        </template>
    </UDashboardPanel>
</template>
