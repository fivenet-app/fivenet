<script lang="ts" setup>
useHead({
    title: 'pages.jobs.colleagues.title',
});
definePageMeta({
    title: 'pages.jobs.colleagues.title',
    requiresAuth: true,
    permission: 'Jobs.View',
});
</script>

<template>
    <h1 class="text-xl text-white">Colleagues</h1>
</template>
