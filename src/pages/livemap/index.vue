<script lang="ts" setup>
import Livemap from '~/components/livemap/Livemap.vue';

useHead({
    title: 'Livemap',
});
definePageMeta({
    title: 'Livemap',
    requiresAuth: true,
    permission: 'LivemapperService.Stream',
});
</script>

<template>
    <Livemap />
</template>
