<script lang="ts" setup>
import ColleagueList from '~/components/jobs/colleagues/ColleagueList.vue';

useHead({
    title: 'pages.jobs.colleagues.title',
});

definePageMeta({
    title: 'pages.jobs.colleagues.title',
    requiresAuth: true,
    permission: 'jobs.JobsService/ListColleagues',
});
</script>

<template>
    <ColleagueList />
</template>
