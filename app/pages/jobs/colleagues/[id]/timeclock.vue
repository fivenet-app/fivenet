<script lang="ts" setup>
import type { TypedRouteFromName } from '@typed-router';
import TimeclockList from '~/components/jobs/timeclock/TimeclockList.vue';

useHead({
    title: 'pages.jobs.colleagues.single.timeclock',
});

definePageMeta({
    title: 'pages.jobs.colleagues.single.timeclock',
    requiresAuth: true,
    permission: 'jobs.TimeclockService/ListTimeclock',
    validate: async (route) => {
        route = route as TypedRouteFromName<'jobs-colleagues-id-timeclock'>;
        // Check if the id is made up of digits
        if (typeof route.params.id !== 'string') {
            return false;
        }
        return !!(route.params.id && !isNaN(Number(route.params.id))) && Number(route.params.id) > -1;
    },
});

const route = useRoute('jobs-colleagues-id-timeclock');
</script>

<template>
    <TimeclockList
        :user-id="parseInt(route.params.id as string)"
        :show-stats="false"
        force-historic-view
        :historic-sub-days="14"
        hide-daily
    />
</template>
