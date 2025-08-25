<script lang="ts" setup>
import ActivityFeed from '~/components/jobs/colleagues/info/ActivityFeed.vue';

useHead({
    title: 'pages.jobs.colleagues.title',
});

definePageMeta({
    title: 'pages.jobs.colleagues.title',
    requiresAuth: true,
    permission: 'jobs.JobsService/ListColleagueActivity',
});
</script>

<template>
    <ActivityFeed :show-target-user="true" />
</template>
