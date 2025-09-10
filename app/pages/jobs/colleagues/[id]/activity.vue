<script lang="ts" setup>
import type { TypedRouteFromName } from '@typed-router';
import ActivityFeed from '~/components/jobs/colleagues/info/ActivityFeed.vue';

useHead({
    title: 'pages.jobs.colleagues.id.activity',
});

definePageMeta({
    title: 'pages.jobs.colleagues.id.activity',
    requiresAuth: true,
    permission: 'jobs.JobsService/GetColleague',
    validate: async (route) => {
        route = route as TypedRouteFromName<'jobs-colleagues-id-activity'>;
        // Check if the id is made up of digits
        if (typeof route.params.id !== 'string') {
            return false;
        }
        return !!(route.params.id && !isNaN(Number(route.params.id))) && Number(route.params.id) > -1;
    },
});

const route = useRoute('jobs-colleagues-id-activity');
</script>

<template>
    <ActivityFeed :user-id="parseInt(route.params.id as string)" />
</template>
