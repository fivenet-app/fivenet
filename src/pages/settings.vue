<script lang="ts" setup>
import UserSettingsPanel from '~/components/auth/account/UserSettingsPanel.vue';

useHead({
    title: 'components.auth.UserSettingsPanel.title',
});
definePageMeta({
    title: 'components.auth.UserSettingsPanel.title',
    requiresAuth: true,
    authOnlyToken: true,

    showCookieOptions: true,
});
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('components.auth.UserSettingsPanel.title')" />

            <UserSettingsPanel />
        </UDashboardPanel>
    </UDashboardPage>
</template>
