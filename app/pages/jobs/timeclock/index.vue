<script lang="ts" setup>
import PagesJobsLayout from '~/components/jobs/PagesJobsLayout.vue';
import TimeclockList from '~/components/jobs/timeclock/TimeclockList.vue';

useHead({
    title: 'pages.jobs.timeclock.title',
});
definePageMeta({
    title: 'pages.jobs.timeclock.title',
    requiresAuth: true,
    permission: 'JobsTimeclockService.ListTimeclock',
});
</script>

<template>
    <PagesJobsLayout>
        <template #default>
            <TimeclockList />
        </template>
    </PagesJobsLayout>
</template>
