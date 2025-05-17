<script lang="ts" setup>
import AttrJobList from '~/components/rector/attrs/AttrJobList.vue';

useHead({
    title: 'pages.rector.limiter.title',
});

definePageMeta({
    title: 'pages.rector.limiter.title',
    requiresAuth: true,
    permission: 'SuperUser',
});
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('pages.rector.limiter.title')">
                <template #right>
                    <PartialsBackButton fallback-to="/rector" />
                </template>
            </UDashboardNavbar>

            <AttrJobList />
        </UDashboardPanel>
    </UDashboardPage>
</template>
