<script lang="ts" setup>
import Livemap from '~/components/livemap/Livemap.vue';

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
    <Livemap />
</template>
