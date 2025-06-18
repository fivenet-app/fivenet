<script lang="ts" setup>
import ColleaguesList from '~/components/jobs/colleagues/ColleaguesList.vue';
import PagesJobsLayout from '~/components/jobs/PagesJobsLayout.vue';

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
    <PagesJobsLayout>
        <template #default>
            <ColleaguesList />
        </template>
    </PagesJobsLayout>
</template>
