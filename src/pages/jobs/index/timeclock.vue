<script lang="ts" setup>
useHead({
    title: 'pages.jobs.timeclock.title',
});
definePageMeta({
    title: 'pages.jobs.timeclock.title',
    requiresAuth: true,
    permission: 'Jobs.TimeclockList',
});
</script>

<template>
    <h1 class="text-xl text-white">Time Clock</h1>
</template>
