<script lang="ts" setup>
useHead({
    title: 'components.auth.user_settings_panel.title',
});

definePageMeta({
    title: 'components.auth.user_settings_panel.title',
    requiresAuth: true,
    authTokenOnly: true,
    showCookieOptions: true,
});

const { t } = useI18n();

const items = computed(() => [
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
]);
</script>

<template>
    <UDashboardPanel>
        <template #header>
            <UDashboardNavbar :title="$t('components.auth.user_settings_panel.title')">
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
