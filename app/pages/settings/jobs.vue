<script lang="ts" setup>
import JobList from '~/components/settings/jobs/JobList.vue';

useHead({
    title: 'pages.settings.jobs.title',
});

definePageMeta({
    title: 'pages.settings.jobs.title',
    requiresAuth: true,
    permission: 'Superuser/Superuser',
});
</script>

<template>
    <JobList />
</template>
