<script lang="ts" setup>
import { type TypedRouteFromName } from '@typed-router';
import ColleagueInfo from '~/components/jobs/colleagues/info/ColleagueInfo.vue';

useHead({
    title: 'pages.jobs.colleagues.single.title',
});
definePageMeta({
    title: 'pages.jobs.colleagues.single.title',
    requiresAuth: true,
    permission: 'JobsService.GetColleague',
    validate: async (route) => {
        route = route as TypedRouteFromName<'jobs-colleagues-id'>;
        // Check if the id is made up of digits
        return /^\d+$/.test(route.params.id);
    },
});

const route = useRoute('jobs-colleagues-id');
</script>

<template>
    <div>
        <ColleagueInfo :user-id="parseInt(route.params.id)" />
    </div>
</template>
