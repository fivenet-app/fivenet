<script lang="ts" setup>
useHead({
    title: 'components.auth.user_settings.preferences',
});

definePageMeta({
    title: 'components.auth.user_settings.preferences',
    requiresAuth: true,
    authTokenOnly: true,
});

const { t } = useI18n();

const route = useRoute();

const items = computed(() => [
    {
        label: t('components.auth.user_settings.preferences'),
        icon: 'i-mdi-cog',
        to: '/user-settings',
        exact: true,
    },
    {
        label: t('common.notification', 2),
        icon: 'i-mdi-notification-settings',
        to: '/user-settings/notifications',
        active: isRoute(route.path, '/user-settings/notifications'),
    },
    {
        label: t('components.auth.user_settings.livemap_dispatch.title'),
        icon: 'i-mdi-map-marker-path',
        to: '/user-settings/livemap-dispatch',
        active: isRoute(route.path, '/user-settings/livemap-dispatch'),
    },
]);
</script>

<template>
    <UDashboardPanel :ui="{ root: 'pb-(--page-content-bottom-offset)' }">
        <template #header>
            <UDashboardNavbar :title="$t('components.auth.user_settings.preferences')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <PartialsBackButton fallback-to="/overview" />
                </template>
            </UDashboardNavbar>

            <UDashboardToolbar>
                <UNavigationMenu class="-mx-1 flex-1" :items="items" highlight />
            </UDashboardToolbar>
        </template>

        <template #body>
            <div class="space-y-4">
                <UAlert
                    color="neutral"
                    variant="subtle"
                    icon="i-mdi-content-save-check-outline"
                    :title="$t('components.auth.user_settings.autosave.title')"
                    :description="$t('components.auth.user_settings.autosave.description')"
                />

                <NuxtPage />
            </div>
        </template>
    </UDashboardPanel>
</template>
