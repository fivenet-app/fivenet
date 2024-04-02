<script lang="ts" setup>
import UserSettingsPanel from '~/components/auth/account/UserSettingsPanel.vue';

useHead({
    title: 'components.auth.settings_panel.title',
});
definePageMeta({
    title: 'components.auth.settings_panel.title',
    requiresAuth: true,
    authOnlyToken: true,

    showCookieOptions: true,
});
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('components.auth.settings_panel.title')"> </UDashboardNavbar>

            <UserSettingsPanel />
        </UDashboardPanel>
    </UDashboardPage>
</template>
