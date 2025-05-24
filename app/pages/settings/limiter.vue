<script lang="ts" setup>
import AttrJobList from '~/components/settings/attrs/AttrJobList.vue';

useHead({
    title: 'pages.settings.limiter.title',
});

definePageMeta({
    title: 'pages.settings.limiter.title',
    requiresAuth: true,
    permission: 'Superuser',
});
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('pages.settings.limiter.title')">
                <template #right>
                    <PartialsBackButton fallback-to="/settings" />
                </template>
            </UDashboardNavbar>

            <AttrJobList />
        </UDashboardPanel>
    </UDashboardPage>
</template>
