<script lang="ts" setup>
import CentrumSidebar from '~/components/centrum/livemap/CentrumSidebar.vue';

useHead({
    title: 'common.livemap',
});

definePageMeta({
    title: 'common.livemap',
    requiresAuth: true,
    permission: 'livemap.LivemapService/Stream',
});
</script>

<template>
    <ClientOnly>
        <CentrumSidebar />
    </ClientOnly>
</template>
