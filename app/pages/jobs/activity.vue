<script lang="ts" setup>
import ActivityFeed from '~/components/jobs/colleagues/info/ActivityFeed.vue';

useHead({
    title: 'pages.jobs.activity.title',
});

definePageMeta({
    title: 'pages.jobs.activity.title',
    requiresAuth: true,
    permission: 'jobs.ColleaguesService/ListColleagueActivity',
});
</script>

<template>
    <ActivityFeed show-target-user />
</template>
