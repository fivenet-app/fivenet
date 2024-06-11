<script lang="ts" setup>
import TimeclockList from '~/components/jobs/timeclock/TimeclockList.vue';
import PagesJobsLayout from '~/components/jobs/PagesJobsLayout.vue';

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
