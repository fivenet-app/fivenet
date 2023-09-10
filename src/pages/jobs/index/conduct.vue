<script lang="ts" setup>
useHead({
    title: 'pages.jobs.conduct.title',
});
definePageMeta({
    title: 'pages.jobs.conduct.title',
    requiresAuth: true,
    permission: 'Jobs.ConductListEntries',
});
</script>

<template>
    <h1 class="text-xl text-white">Conduct</h1>
</template>
