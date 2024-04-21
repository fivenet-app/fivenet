<script lang="ts" setup>
useHead({
    title: 'pages.jobs.title',
});
definePageMeta({
    title: 'pages.jobs.title',
    requiresAuth: true,
    permission: 'JobsService.ListColleagues',
    redirect: { name: 'jobs-overview' },
});
</script>

<template>
    <div></div>
</template>
