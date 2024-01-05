<script lang="ts" setup>
import ColleaguesList from '~/components/jobs/colleagues/ColleaguesList.vue';

useHead({
    title: 'pages.jobs.colleagues.title',
});
definePageMeta({
    title: 'pages.jobs.colleagues.title',
    requiresAuth: true,
    permission: 'JobsService.ColleaguesList',
});
</script>

<template>
    <div>
        <ColleaguesList />
    </div>
</template>
