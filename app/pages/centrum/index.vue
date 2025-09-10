<script lang="ts" setup>
import DispatchCenter from '~/components/centrum/DispatchCenter.vue';

useHead({
    title: 'common.dispatch_center',
});

definePageMeta({
    title: 'common.dispatch_center',
    requiresAuth: true,
    permission: 'centrum.CentrumService/TakeControl',
});
</script>

<template>
    <DispatchCenter />
</template>
