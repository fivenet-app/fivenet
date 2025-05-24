<script lang="ts" setup>
import LivemapHolder from '~/components/livemap/LivemapHolder.vue';

useHead({
    title: 'common.livemap',
});

definePageMeta({
    title: 'common.livemap',
    requiresAuth: true,
    permission: 'livemap.LivemapService.Stream',
});
</script>

<template>
    <UDashboardPage>
        <LivemapHolder />
    </UDashboardPage>
</template>
