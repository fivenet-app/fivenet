<script lang="ts" setup>
import DispatchCenter from '~/components/centrum/DispatchCenter.vue';

useHead({
    title: 'common.dispatch_center',
});

definePageMeta({
    title: 'common.dispatch_center',
    requiresAuth: true,
    permission: 'CentrumService.TakeControl',
});
</script>

<template>
    <UDashboardPage>
        <DispatchCenter />
    </UDashboardPage>
</template>
