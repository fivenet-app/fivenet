<script lang="ts" setup>
useHead({
    title: 'pages.jobs.requests.title',
});
definePageMeta({
    title: 'pages.jobs.requests.title',
    requiresAuth: true,
    permission: 'Jobs.RequestsList',
});
</script>

<template>
    <h1 class="text-xl text-white">Requests</h1>
</template>
