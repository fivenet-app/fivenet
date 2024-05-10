<script lang="ts" setup>
import type { TypedRouteFromName } from '@typed-router';
import TimeclockOverviewBlock from '~/components/jobs/timeclock/TimeclockOverviewBlock.vue';

useHead({
    title: 'pages.jobs.colleagues.single.timeclock',
});
definePageMeta({
    title: 'pages.jobs.colleagues.single.timeclock',
    requiresAuth: true,
    permission: 'JobsTimeclockService.ListTimeclock',
    validate: async (route) => {
        route = route as TypedRouteFromName<'jobs-colleagues-id-timeclock'>;
        // Check if the id is made up of digits
        return /^\d+$/.test(route.params.id);
    },
});

const route = useRoute('jobs-colleagues-id-timeclock');
</script>

<template>
    <TimeclockOverviewBlock :user-id="parseInt(route.params.id)" />
</template>
