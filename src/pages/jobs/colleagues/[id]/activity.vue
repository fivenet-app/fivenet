<script lang="ts" setup>
import type { TypedRouteFromName } from '@typed-router';
import ColleagueActivityFeed from '~/components/jobs/colleagues/info/ColleagueActivityFeed.vue';

useHead({
    title: 'pages.citizens.id.title',
});
definePageMeta({
    title: 'pages.citizens.id.title',
    requiresAuth: true,
    permission: 'JobsService.GetColleague',
    validate: async (route) => {
        route = route as TypedRouteFromName<'jobs-colleagues-id-activity'>;
        // Check if the id is made up of digits
        return /^\d+$/.test(route.params.id);
    },
});

const route = useRoute('jobs-colleagues-id-activity');
</script>

<template>
    <div>
        <ColleagueActivityFeed :user-id="parseInt(route.params.id)" />
    </div>
</template>
