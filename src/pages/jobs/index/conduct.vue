<script lang="ts" setup>
useHead({
    title: 'pages.jobs.title',
});
definePageMeta({
    title: 'pages.jobs.title',
    requiresAuth: true,
    permission: 'Jobs.View',
});
</script>

<template>
    <h1 class="text-xl text-white">Conduct</h1>
</template>
