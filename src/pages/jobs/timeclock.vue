<script lang="ts" setup>
import List from '~/components/jobs/timeclock/List.vue';

useHead({
    title: 'pages.jobs.timeclock.title',
});
definePageMeta({
    title: 'pages.jobs.timeclock.title',
    requiresAuth: true,
    permission: 'JobsService.TimeclockListEntries',
});
</script>

<template>
    <div>
        <List />
    </div>
</template>
