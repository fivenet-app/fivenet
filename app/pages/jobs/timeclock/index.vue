<script lang="ts" setup>
import TimeclockList from '~/components/jobs/timeclock/TimeclockList.vue';

useHead({
    title: 'pages.jobs.timeclock.title',
});

definePageMeta({
    title: 'pages.jobs.timeclock.title',
    requiresAuth: true,
    permission: 'jobs.TimeclockService/ListTimeclock',
});
</script>

<template>
    <TimeclockList />
</template>
