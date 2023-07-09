<script lang="ts" setup>
import Centrum from '~/components/centrum/Centrum.vue';

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
    <Centrum />
</template>
