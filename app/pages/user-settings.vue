<script lang="ts" setup>
useHead({
    title: 'components.auth.UserSettingsPanel.title',
});

definePageMeta({
    title: 'components.auth.UserSettingsPanel.title',
    requiresAuth: true,
    authTokenOnly: true,
    showCookieOptions: true,
});

const { t } = useI18n();

const items = [
    {
        label: t('common.settings'),
        icon: 'i-mdi-cog',
        to: '/user-settings',
        exact: true,
    },
    {
        label: t('common.notification', 2),
        icon: 'i-mdi-notification-settings',
        to: '/user-settings/notifications',
    },
];
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

            <UDashboardToolbar>
                <UNavigationMenu :items="items" highlight class="-mx-1 flex-1" />
            </UDashboardToolbar>
        </template>

        <template #body>
            <NuxtPage />
        </template>
    </UDashboardPanel>
</template>
