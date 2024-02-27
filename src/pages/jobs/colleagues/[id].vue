<script lang="ts" setup>
import ColleagueInfo from '~/components/jobs/colleagues/ColleagueInfo.vue';

useHead({
    title: 'pages.jobs.colleagues.single.title',
});
definePageMeta({
    title: 'pages.jobs.colleagues.single.title',
    requiresAuth: true,
    permission: 'JobsService.ListColleagues',
});
</script>

<template>
    <div>
        <ColleagueInfo />
    </div>
</template>
