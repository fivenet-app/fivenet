<script lang="ts" setup>
import DebugInfo from '~/components/auth/account/DebugInfo.vue';
import StreamerModeAlert from '~/components/partials/StreamerModeAlert.vue';

useHead({
    title: 'components.auth.Debug.title',
});

definePageMeta({
    title: 'components.auth.Debug.title',
    requiresAuth: true,
    authTokenOnly: true,
    showCookieOptions: true,
});

const { streamerMode } = storeToRefs(useSettingsStore());
</script>

<template>
    <UDashboardPanel :ui="{ root: 'pb-(--page-content-bottom-offset)' }">
        <template #header>
            <UDashboardNavbar :title="$t('components.auth.Debug.title')">
                <template #leading>
                    <UDashboardSidebarCollapse />
                </template>

                <template #right>
                    <PartialsBackButton fallback-to="/overview" />
                </template>
            </UDashboardNavbar>
        </template>

        <template #body>
            <StreamerModeAlert v-if="streamerMode" />
            <DebugInfo v-else />
        </template>
    </UDashboardPanel>
</template>
