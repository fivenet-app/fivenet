<script lang="ts" setup>
useHead({
    title: 'pages.jobs.title',
});
definePageMeta({
    title: 'pages.jobs.title',
    requiresAuth: true,
    permission: 'Jobs.View',
    alias: '/jobs',
});
</script>

<template>
    <h1 class="text-xl text-white">Overview</h1>
</template>
