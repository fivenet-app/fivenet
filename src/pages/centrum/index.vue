<script lang="ts" setup>
import DispatchCenter from '~/components/centrum/DispatchCenter.vue';
import DisponentsInfo from '~/components/centrum/disponents/DisponentsInfo.vue';

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
    <UDashboardPage class="size-full">
        <UDashboardPanel grow>
            <UDashboardNavbar :title="$t('common.dispatch_center')">
                <template #right>
                    <DisponentsInfo />
                </template>
            </UDashboardNavbar>

            <DispatchCenter class="min-h-[calc(100vh-var(--header-height))]" />
        </UDashboardPanel>
    </UDashboardPage>
</template>
