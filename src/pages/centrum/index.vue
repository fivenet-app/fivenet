<script lang="ts" setup>
import DispatchCenter from '~/components/centrum/DispatchCenter.vue';

useHead({
    title: 'common.dispatch_center',
});
definePageMeta({
    title: 'common.dispatch_center',
    requiresAuth: true,
    permission: 'CentrumService.TakeControl',
    showQuickButtons: false,
});
</script>

<template>
    <UDashboardPage>
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('common.dispatch_center')"> </UDashboardNavbar>

            <DispatchCenter />
        </UDashboardPanel>
    </UDashboardPage>
</template>
