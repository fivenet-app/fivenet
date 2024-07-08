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
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('components.auth.UserSettingsPanel.title')">
                <template #right>
                    <UButton color="black" icon="i-mdi-arrow-back" to="/overview">
                        {{ $t('common.back') }}
                    </UButton>
                </template>
            </UDashboardNavbar>

            <UserSettingsPanel />
        </UDashboardPanel>
    </UDashboardPage>
</template>
