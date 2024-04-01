<script lang="ts" setup>
import LivemapHolder from '~/components/livemap/LivemapHolder.vue';

useHead({
    title: 'common.livemap',
});
definePageMeta({
    title: 'common.livemap',
    requiresAuth: true,
    permission: 'LivemapperService.Stream',
    showQuickButtons: false,
});
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('common.livemap')"> </UDashboardNavbar>

            <LivemapHolder class="h-dscreen w-full" />
        </UDashboardPanel>
    </UDashboardPage>
</template>
