<script lang="ts" setup>
import DispatchCenter from '~/components/dispatch/DispatchCenter.vue';

useHead({
    title: 'common.dispatch_center',
});

definePageMeta({
    title: 'common.dispatch_center',
    requiresAuth: true,
    permission: 'centrum.CentrumService/TakeControl',
    alias: '/centrum',
});
</script>

<template>
    <DispatchCenter />
</template>
