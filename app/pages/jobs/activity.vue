<script lang="ts" setup>
import ColleagueActivityFeed from '~/components/jobs/colleagues/info/ColleagueActivityFeed.vue';
import PagesJobsLayout from '~/components/jobs/PagesJobsLayout.vue';

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
    <PagesJobsLayout>
        <template #default>
            <ColleagueActivityFeed :show-target-user="true" />
        </template>
    </PagesJobsLayout>
</template>
