<script lang="ts" setup>
import ColleagueList from '~/components/jobs/colleagues/ColleagueList.vue';
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
            <ColleagueList />
        </template>
    </PagesJobsLayout>
</template>
