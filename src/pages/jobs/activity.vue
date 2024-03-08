<script lang="ts" setup>
import ColleagueActivityFeed from '~/components/jobs/colleagues/info/ColleagueActivityFeed.vue';

useHead({
    title: 'pages.jobs.colleagues.title',
});
definePageMeta({
    title: 'pages.jobs.colleagues.title',
    requiresAuth: true,
    permission: 'JobsService.ListColleagueActivity',
});
</script>

<template>
    <div>
        <ColleagueActivityFeed :show-target-user="true" />
    </div>
</template>
